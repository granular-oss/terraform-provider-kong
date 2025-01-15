package provider

import (
	"context"
	"fmt"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongConsumerAclResource{}
	_ resource.ResourceWithConfigure   = &kongConsumerAclResource{}
	_ resource.ResourceWithImportState = &kongConsumerAclResource{}
)

// NewkongConsumerAclResource is a helper function to simplify the provider implementation.
func KongConsumerAclResource() resource.Resource {
	return &kongConsumerAclResource{}
}

// kongConsumerAclResource is the resource implementation.
type kongConsumerAclResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongConsumerAclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_consumer_acl"
}

// Schema defines the schema for the resource.
func (r *kongConsumerAclResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"kong_id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"consumer_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"group": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"tags": schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongConsumerAclResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongConsumerAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongConsumerAclModel
	client := r.client.ACLs
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromConsumerAclModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_consumer_acl", map[string]any{"body": reqBody, "model": plan})

	acl, err := client.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_acl", err.Error())
		return
	}
	result := kongModels.ConsumerAclModelFromResponse(acl)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongConsumerAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongConsumerAclModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	acl, err := r.client.ACLs.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_consumer_acl", "Could not read kong_consumer_acl with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_consumer_acl state", map[string]any{"state_file": state, "upstream": acl})
	diags = resp.State.Set(ctx, kongModels.ConsumerAclModelFromResponse(acl))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongConsumerAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongConsumerAclModel
	diags := req.Plan.Get(ctx, &plan)
	idParts := kongModels.ParseCompositeId(plan.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromConsumerAclModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_consumer_acl", map[string]any{"body": reqBody, "model": plan})
	acl, err := r.client.ACLs.Create(ctx, &idParts[0], reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_acl", err.Error())
		return
	}
	result := kongModels.ConsumerAclModelFromResponse(acl)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongConsumerAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongConsumerAclModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.ACLs.Delete(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_consumer_acl", err.Error())
		return
	}
}

func (r *kongConsumerAclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := kongModels.ParseCompositeId(req.ID)
	acl, err := r.client.ACLs.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_consumer_acl for the given name/id:group", err.Error())
	}
	id := kongModels.BuildCompositeId([]string{*acl.Consumer.ID, *acl.Group})
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
