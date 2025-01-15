package provider

import (
	"context"
	"fmt"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongUpstreamResource{}
	_ resource.ResourceWithConfigure   = &kongUpstreamResource{}
	_ resource.ResourceWithImportState = &kongUpstreamResource{}
)

// KongUpstreamResource is a helper function to simplify the provider implementation.
func KongUpstreamResource() resource.Resource {
	return &kongUpstreamResource{}
}

// kongUpstreamResource is the resource implementation.
type kongUpstreamResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongUpstreamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_upstream"
}

// Schema defines the schema for the resource.
func (r *kongUpstreamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"name":                  schema.StringAttribute{Required: true},
			"slots":                 schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(10000)},
			"hash_on":               schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("none")},
			"hash_fallback":         schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("none")},
			"hash_on_header":        schema.StringAttribute{Optional: true},
			"host_header":           schema.StringAttribute{Optional: true},
			"tags":                  schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"client_certificate_id": schema.StringAttribute{Optional: true},
			"hash_fallback_header":  schema.StringAttribute{Optional: true},
			"hash_on_cookie":        schema.StringAttribute{Optional: true},
			"hash_on_cookie_path":   schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("/")},
			"healthchecks": schema.SingleNestedAttribute{Optional: true, Computed: true, Attributes: map[string]schema.Attribute{
				"active": schema.SingleNestedAttribute{Optional: true, Computed: true, Attributes: map[string]schema.Attribute{
					"type":                     schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("http")},
					"timeout":                  schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(1)},
					"concurrency":              schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(10)},
					"http_path":                schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("/")},
					"https_verify_certificate": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
					"https_sni":                schema.StringAttribute{Optional: true},
					"healthy": schema.SingleNestedAttribute{Optional: true, Computed: true, Attributes: map[string]schema.Attribute{
						"interval":      schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
						"http_statuses": schema.ListAttribute{Optional: true, ElementType: types.Int32Type, Computed: true},
						"successes":     schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
					}},
					"unhealthy": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
						"interval":      schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
						"http_statuses": schema.ListAttribute{Optional: true, ElementType: types.Int32Type, Computed: true},
						"tcp_failures":  schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
						"http_failures": schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
						"timeouts":      schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
					}},
				}},
				"passive": schema.SingleNestedAttribute{Optional: true, Computed: true, Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Optional: true, Computed: true, Default: stringdefault.StaticString("http")},
					"healthy": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
						"http_statuses": schema.ListAttribute{Optional: true, ElementType: types.Int32Type, Computed: true},
						"successes":     schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
					}},
					"unhealthy": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
						"http_statuses": schema.ListAttribute{Optional: true, ElementType: types.Int32Type, Computed: true},
						"tcp_failures":  schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
						"http_failures": schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
						"timeouts":      schema.Int32Attribute{Optional: true, Computed: true, Default: int32default.StaticInt32(0)},
					}},
				}},
			}, PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()}},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongUpstreamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongUpstreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongUpstreamModel
	client := r.client.Upstreams
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromUpstreamModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_upstream", map[string]any{"body": reqBody, "model": plan})
	upstream, err := client.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_upstream", err.Error())
		return
	}
	result := kongModels.UpstreamModelFromResponse(resp.Diagnostics, upstream)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongUpstreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongUpstreamModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	upstream, err := r.client.Upstreams.Get(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_upstream", "Could not read kong_upstream with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_upstream state", map[string]any{"state_file": state, "upstream": upstream})
	diags = resp.State.Set(ctx, kongModels.UpstreamModelFromResponse(resp.Diagnostics, upstream))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongUpstreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongUpstreamModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromUpstreamModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_upstream", map[string]any{"body": reqBody, "model": plan})
	upstream, err := r.client.Upstreams.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_upstream", err.Error())
		return
	}
	result := kongModels.UpstreamModelFromResponse(resp.Diagnostics, upstream)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongUpstreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongUpstreamModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Upstreams.Delete(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_upstream", err.Error())
		return
	}
}

func (r *kongUpstreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	upstream, err := r.client.Upstreams.Get(ctx, &req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_upstream for the given username/id", err.Error())
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), upstream.ID)...)
}
