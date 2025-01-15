package provider

import (
	"context"
	"fmt"
	"strings"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongTargetResource{}
	_ resource.ResourceWithConfigure   = &kongTargetResource{}
	_ resource.ResourceWithImportState = &kongTargetResource{}
)

// KongTargetResource is a helper function to simplify the provider implementation.
func KongTargetResource() resource.Resource {
	return &kongTargetResource{}
}

// kongTargetResource is the resource implementation.
type kongTargetResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_target"
}

// Schema defines the schema for the resource.
func (r *kongTargetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"kong_id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"upstream_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"target": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"weight": schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(100), PlanModifiers: []planmodifier.Int32{
				int32planmodifier.RequiresReplace(),
			}},
			"tags": schema.ListAttribute{Optional: true, ElementType: types.StringType, PlanModifiers: []planmodifier.List{
				listplanmodifier.RequiresReplace(),
			}},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kong.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kong.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *kongTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongTargetModel
	client := r.client.Targets
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromTargetModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_target", map[string]any{"body": reqBody, "model": plan})
	target, err := client.Create(ctx, plan.UpstreamId.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_target", err.Error())
		return
	}
	result := kongModels.TargetModelFromResponse(target)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongTargetModel
	diags := req.State.Get(ctx, &state)
	idParts := strings.Split(state.ID.ValueString(), "|")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	targets, err := r.client.Targets.ListAll(ctx, &idParts[0])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_target", "Could not read kong_target with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	var target *kong.Target
	for _, t := range targets {
		if idParts[1] == *t.ID {
			target = t
			break
		}
	}
	tflog.Debug(ctx, "Reading kong_target state", map[string]any{"state_file": state, "upstream": target})
	diags = resp.State.Set(ctx, kongModels.TargetModelFromResponse(target))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongTargetModel
	diags := req.Plan.Get(ctx, &plan)
	idParts := strings.Split(plan.ID.ValueString(), "|")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromTargetModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_target", map[string]any{"body": reqBody, "model": plan})
	target, err := r.client.Targets.Create(ctx, &idParts[0], reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_target", err.Error())
		return
	}
	result := kongModels.TargetModelFromResponse(target)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongTargetModel
	diags := req.State.Get(ctx, &state)
	idParts := strings.Split(state.ID.ValueString(), "|")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Targets.Delete(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_target", err.Error())
		return
	}
}

func (r *kongTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "|")
	targets, err := r.client.Targets.ListAll(ctx, &idParts[0])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_target", "Could not read kong_target with ID "+idParts[0]+": "+err.Error())
		return
	}
	var target *kong.Target
	for _, t := range targets {
		if idParts[1] == *t.ID || idParts[1] == *t.Target {
			target = t
			break
		}
	}
	id := *target.Upstream.ID + "|" + *target.ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
