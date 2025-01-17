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
	_ resource.Resource                = &kongConsumerJwtResource{}
	_ resource.ResourceWithConfigure   = &kongConsumerJwtResource{}
	_ resource.ResourceWithImportState = &kongConsumerJwtResource{}
)

// KongConsumerJwtResource is a helper function to simplify the provider implementation.
func KongConsumerJwtResource() resource.Resource {
	return &kongConsumerJwtResource{}
}

// kongConsumerJwtResource is the resource implementation.
type kongConsumerJwtResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongConsumerJwtResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_consumer_jwt"
}

// Schema defines the schema for the resource.
func (r *kongConsumerJwtResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"key": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			}},
			"secret": schema.StringAttribute{Optional: true, Computed: true, Sensitive: true},
			"algorithm": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"rsa_public_key": schema.StringAttribute{Optional: true},
			"tags":           schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongConsumerJwtResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongConsumerJwtResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongConsumerJwtModel
	client := r.client.JWTAuths
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromConsumerJwtModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_consumer_jwt", map[string]any{"body": reqBody, "model": plan})

	jwt, err := client.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_jwt", err.Error())
		return
	}
	result := kongModels.ConsumerJwtModelFromResponse(jwt)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongConsumerJwtResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongConsumerJwtModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	jwt, err := r.client.JWTAuths.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_consumer_jwt", "Could not read kong_consumer_jwt with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_consumer_jwt state", map[string]any{"state_file": state, "upstream": jwt})
	model := kongModels.ConsumerJwtModelFromResponse(jwt)
	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongConsumerJwtResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongConsumerJwtModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromConsumerJwtModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_consumer_jwt", map[string]any{"body": reqBody, "model": plan})
	jwt, err := r.client.JWTAuths.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_jwt", err.Error())
		return
	}
	result := kongModels.ConsumerJwtModelFromResponse(jwt)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongConsumerJwtResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongConsumerJwtModel
	diags := req.State.Get(ctx, &state)
	idParts := kongModels.ParseCompositeId(state.ID.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.JWTAuths.Delete(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_consumer_jwt", err.Error())
		return
	}
}

func (r *kongConsumerJwtResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := kongModels.ParseCompositeId(req.ID)
	auth, err := r.client.JWTAuths.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_consumer_jwt for the given name/id|key/id", err.Error())
	}
	id := kongModels.BuildCompositeId([]string{*auth.Consumer.ID, *auth.ID})
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
