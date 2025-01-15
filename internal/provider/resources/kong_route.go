package provider

import (
	"context"
	"fmt"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongRouteResource{}
	_ resource.ResourceWithConfigure   = &kongRouteResource{}
	_ resource.ResourceWithImportState = &kongRouteResource{}
)

func KongRouteResource() resource.Resource {
	return &kongRouteResource{}
}

// kongRouteResource is the resource implementation.
type kongRouteResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongRouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_route"
}

// Schema defines the schema for the resource.
func (r *kongRouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"name":       schema.StringAttribute{Required: true},
			"service_id": schema.StringAttribute{Required: true},
			"protocols": schema.ListAttribute{Optional: true, Computed: true, ElementType: types.StringType, Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("http"),
				types.StringValue("https"),
			}))},
			"methods":    schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"hosts":      schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"paths":      schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"tags":       schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"strip_path": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"source": schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip":   schema.StringAttribute{Required: true},
					"port": schema.Int32Attribute{Required: true},
				},
			}},
			"destination": schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip":   schema.StringAttribute{Required: true},
					"port": schema.Int32Attribute{Required: true},
				},
			}},
			"snis":                       schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"preserve_host":              schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(false)},
			"regex_priority":             schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
			"path_handling":              schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("v0")},
			"https_redirect_status_code": schema.Int32Attribute{Optional: true, Computed: true},
			"request_buffering":          schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"response_buffering":         schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"header": schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name":   schema.StringAttribute{Required: true},
					"values": schema.ListAttribute{Required: true, ElementType: types.StringType},
				},
			}},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongRouteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongRouteModel
	client := r.client.Routes
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromRouteModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_route", map[string]any{"body": reqBody, "model": plan})
	route, err := client.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_route", err.Error())
		return
	}
	result := kongModels.RouteModelFromResponse(route)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongRouteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	route, err := r.client.Routes.Get(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_route", "Could not read kong_route with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_route state", map[string]any{"state_file": state, "upstream": route})
	diags = resp.State.Set(ctx, kongModels.RouteModelFromResponse(route))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongRouteModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromRouteModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_route", map[string]any{"body": reqBody, "model": plan})
	route, err := r.client.Routes.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_route", err.Error())
		return
	}
	result := kongModels.RouteModelFromResponse(route)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongRouteModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Routes.Delete(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_route", err.Error())
		return
	}
}

func (r *kongRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	route, err := r.client.Routes.Get(ctx, kong.String(req.ID))
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_route for the given name/id", err.Error())
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), route.ID)...)
}
