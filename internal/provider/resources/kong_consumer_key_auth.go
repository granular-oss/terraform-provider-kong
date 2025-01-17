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
	_ resource.Resource                = &kongConsumerKeyAuthResource{}
	_ resource.ResourceWithConfigure   = &kongConsumerKeyAuthResource{}
	_ resource.ResourceWithImportState = &kongConsumerKeyAuthResource{}
)

// KongConsumerKeyAuthResource is a helper function to simplify the provider implementation.
func KongConsumerKeyAuthResource() resource.Resource {
	return &kongConsumerKeyAuthResource{}
}

// kongConsumerKeyAuthResource is the resource implementation.
type kongConsumerKeyAuthResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongConsumerKeyAuthResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_consumer_key_auth"
}

// Schema defines the schema for the resource.
func (r *kongConsumerKeyAuthResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"key": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"tags": schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongConsumerKeyAuthResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongConsumerKeyAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongConsumerKeyAuthModel
	client := r.client.KeyAuths
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromConsumerKeyAuthModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_consumer_key_auth", map[string]any{"body": reqBody, "model": plan})

	key, err := client.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_key_auth", err.Error())
		return
	}
	result := kongModels.ConsumerKeyAuthModelFromResponse(key)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongConsumerKeyAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongConsumerKeyAuthModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	key, err := r.client.KeyAuths.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_consumer_key_auth", "Could not read kong_consumer_key_auth with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_consumer_key_auth state", map[string]any{"state_file": state, "upstream": key})
	diags = resp.State.Set(ctx, kongModels.ConsumerKeyAuthModelFromResponse(key))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongConsumerKeyAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongConsumerKeyAuthModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromConsumerKeyAuthModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_consumer_key_auth", map[string]any{"body": reqBody, "model": plan})
	key, err := r.client.KeyAuths.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_key_auth", err.Error())
		return
	}
	result := kongModels.ConsumerKeyAuthModelFromResponse(key)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongConsumerKeyAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongConsumerKeyAuthModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.KeyAuths.Delete(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_consumer_key_auth", err.Error())
		return
	}
}

func (r *kongConsumerKeyAuthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := kongModels.ParseCompositeId(req.ID)
	key, err := r.client.KeyAuths.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_consumer_key_auth for the given name/id|group", err.Error())
	}
	id := kongModels.BuildCompositeId([]string{*key.Consumer.ID, *key.ID})
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
