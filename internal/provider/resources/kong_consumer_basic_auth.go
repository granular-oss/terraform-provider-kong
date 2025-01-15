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
	_ resource.Resource                = &kongConsumerBasicAuthResource{}
	_ resource.ResourceWithConfigure   = &kongConsumerBasicAuthResource{}
	_ resource.ResourceWithImportState = &kongConsumerBasicAuthResource{}
)

// KongConsumerBasicAuthResource is a helper function to simplify the provider implementation.
func KongConsumerBasicAuthResource() resource.Resource {
	return &kongConsumerBasicAuthResource{}
}

// kongConsumerBasicAuthResource is the resource implementation.
type kongConsumerBasicAuthResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongConsumerBasicAuthResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_consumer_basic_auth"
}

// Schema defines the schema for the resource.
func (r *kongConsumerBasicAuthResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"kong_id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"consumer_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"username": schema.StringAttribute{Required: true},
			"password": schema.StringAttribute{Required: true, Sensitive: true},
			"tags":     schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongConsumerBasicAuthResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongConsumerBasicAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongConsumerBasicAuthModel
	client := r.client.BasicAuths
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromConsumerBasicAuthModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_consumer_basic", map[string]any{"body": reqBody, "model": plan})

	auth, err := client.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_basic_auth", err.Error())
		return
	}
	result := kongModels.ConsumerBasicAuthModelFromResponse(auth)
	result.Password = plan.Password
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongConsumerBasicAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongConsumerBasicAuthModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	auth, err := r.client.BasicAuths.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_consumer_basic_auth", "Could not read kong_consumer_basic_auth with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_consumer_basic_auth state", map[string]any{"state_file": state, "upstream": auth})
	model := kongModels.ConsumerBasicAuthModelFromResponse(auth)
	model.Password = state.Password
	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongConsumerBasicAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongConsumerBasicAuthModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromConsumerBasicAuthModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_consumer_basic_auth", map[string]any{"body": reqBody, "model": plan})
	acl, err := r.client.BasicAuths.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_basic_auth", err.Error())
		return
	}
	result := kongModels.ConsumerBasicAuthModelFromResponse(acl)
	result.Password = plan.Password
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongConsumerBasicAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongConsumerBasicAuthModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.BasicAuths.Delete(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_consumer_basic_auth", err.Error())
		return
	}
}

func (r *kongConsumerBasicAuthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := kongModels.ParseCompositeId(req.ID)
	auth, err := r.client.BasicAuths.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_consumer_basic_auths for the given name/id|username/id", err.Error())
	}
	id := kongModels.BuildCompositeId([]string{*auth.Consumer.ID, *auth.ID})
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
