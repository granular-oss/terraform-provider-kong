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
	_ resource.Resource                = &kongConsumerResource{}
	_ resource.ResourceWithConfigure   = &kongConsumerResource{}
	_ resource.ResourceWithImportState = &kongConsumerResource{}
)

// KongConsumerResource is a helper function to simplify the provider implementation.
func KongConsumerResource() resource.Resource {
	return &kongConsumerResource{}
}

// kongConsumerResource is the resource implementation.
type kongConsumerResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongConsumerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_consumer"
}

// Schema defines the schema for the resource.
func (r *kongConsumerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"username":  schema.StringAttribute{Required: true},
			"custom_id": schema.StringAttribute{Optional: true},
			"tags":      schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongConsumerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongConsumerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongConsumerModel
	client := r.client.Consumers
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromConsumerModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_consumer", map[string]any{"body": reqBody, "model": plan})
	consumer, err := client.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer", err.Error())
		return
	}
	result := kongModels.ConsumerModelFromResponse(consumer)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongConsumerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongConsumerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	consumer, err := r.client.Consumers.Get(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_consumer", "Could not read kong_consumer with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_consumer state", map[string]any{"state_file": state, "upstream": consumer})
	diags = resp.State.Set(ctx, kongModels.ConsumerModelFromResponse(consumer))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongConsumerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongConsumerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromConsumerModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_consumer", map[string]any{"body": reqBody, "model": plan})
	consumer, err := r.client.Consumers.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer", err.Error())
		return
	}
	result := kongModels.ConsumerModelFromResponse(consumer)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongConsumerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongConsumerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Consumers.Delete(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_consumer", err.Error())
		return
	}
}

func (r *kongConsumerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	consumer, err := r.client.Consumers.Get(ctx, &req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_consumer for the given username/id", err.Error())
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), consumer.ID)...)
}
