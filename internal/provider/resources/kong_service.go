package provider

import (
	"context"
	"fmt"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongServiceResource{}
	_ resource.ResourceWithConfigure   = &kongServiceResource{}
	_ resource.ResourceWithImportState = &kongServiceResource{}
)

// KongServiceResource is a helper function to simplify the provider implementation.
func KongServiceResource() resource.Resource {
	return &kongServiceResource{}
}

// kongServiceResource is the resource implementation.
type kongServiceResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_service"
}

// Schema defines the schema for the resource.
func (r *kongServiceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"name":                  schema.StringAttribute{Required: true},
			"host":                  schema.StringAttribute{Required: true},
			"protocol":              schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("http")},
			"port":                  schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(80)},
			"path":                  schema.StringAttribute{Optional: true, Computed: true},
			"retries":               schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(5)},
			"connect_timeout":       schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(60000)},
			"write_timeout":         schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(60000)},
			"read_timeout":          schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(60000)},
			"tags":                  schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"tls_verify":            schema.BoolAttribute{Optional: true},
			"tls_depth":             schema.Int32Attribute{Optional: true},
			"client_certificate_id": schema.StringAttribute{Optional: true, Default: nil},
			"ca_certificate_ids":    schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongServiceModel
	client := r.client.Services
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromServiceModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_service", map[string]any{"body": reqBody, "model": plan})
	service, err := client.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_service", err.Error())
		return
	}
	result := kongModels.ServiceModelFromResponse(service)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongServiceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	service, err := r.client.Services.Get(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_service", "Could not read kong_service with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_service state", map[string]any{"state_file": state, "upstream": service})
	diags = resp.State.Set(ctx, kongModels.ServiceModelFromResponse(service))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongServiceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromServiceModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_service", map[string]any{"body": reqBody, "model": plan})
	service, err := r.client.Services.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_service", err.Error())
		return
	}
	result := kongModels.ServiceModelFromResponse(service)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongServiceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Services.Delete(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_service", err.Error())
		return
	}
}

func (r *kongServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	service, err := r.client.Services.Get(ctx, kong.String(req.ID))
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_service for the given name/id", err.Error())
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), service.ID)...)
}
