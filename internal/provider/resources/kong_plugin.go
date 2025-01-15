package provider

import (
	"context"
	"fmt"
	"strings"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongPluginResource{}
	_ resource.ResourceWithConfigure   = &kongPluginResource{}
	_ resource.ResourceWithImportState = &kongPluginResource{}
)

// KongPluginResource is a helper function to simplify the provider implementation.
func KongPluginResource() resource.Resource {
	return &kongPluginResource{}
}

// kongPluginResource is the resource implementation.
type kongPluginResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_plugin"
}

// Schema defines the schema for the resource.
func (r *kongPluginResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"service_id": schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"route_id": schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"consumer_id": schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"enabled":         schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"config_json":     schema.StringAttribute{Required: true},
			"computed_config": schema.StringAttribute{Computed: true},
			"strict_match":    schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(false)},
			"tags":            schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongPluginModel
	client := r.client.Plugins
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromPluginModel(ctx, resp.Diagnostics, &plan)
	tflog.Debug(ctx, "Creating kong_plugin", map[string]any{"body": reqBody, "model": plan})
	plugin, err := client.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_plugin", err.Error())
		return
	}
	result := kongModels.PluginModelFromResponse(resp.Diagnostics, plugin, plan)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongPluginModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plugin, err := r.client.Plugins.Get(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_plugin", "Could not read kong_plugin with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_plugin state", map[string]any{"state_file": state, "upstream": plugin})
	diags = resp.State.Set(ctx, kongModels.PluginModelFromResponse(resp.Diagnostics, plugin, state))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongPluginModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromPluginModel(ctx, resp.Diagnostics, &plan)
	tflog.Debug(ctx, "Updating kong_plugin", map[string]any{"body": reqBody, "model": plan})
	plugin, err := r.client.Plugins.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_plugin", err.Error())
		return
	}
	result := kongModels.PluginModelFromResponse(resp.Diagnostics, plugin, plan)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongPluginModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Plugins.Delete(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_plugin", err.Error())
		return
	}
}

func (r *kongPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if !strings.Contains(req.ID, ":") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
		return
	}
	idPairs := strings.Split(req.ID, ":")
	if len(idPairs) != 3 {
		resp.Diagnostics.AddError("Invalid import ID", "You must import by plugin ID or `type:parent_id:plugin_name` where type in [`service`,`route`,`consumer`]")
		return
	}
	var plugin *kong.Plugin
	if idPairs[0] == "service" {
		plugin = kongModels.GetPluginByService(ctx, resp.Diagnostics, r.client, &idPairs[1], idPairs[2])
		tflog.Debug(ctx, "Found plugin for service", map[string]any{"plugin": plugin})
	} else if idPairs[0] == "route" {
		plugin = kongModels.GetPluginByRoute(ctx, resp.Diagnostics, r.client, &idPairs[1], idPairs[2])
		tflog.Debug(ctx, "Found plugin for route", map[string]any{"plugin": plugin})
	} else if idPairs[0] == "consumer" {
		plugin = kongModels.GetPluginByConsumer(ctx, resp.Diagnostics, r.client, &idPairs[1], idPairs[2])
		tflog.Debug(ctx, "Found plugin for consumer", map[string]any{"plugin": plugin})
	}
	if plugin == nil {
		resp.Diagnostics.AddError("Failed to find matching plugin", "Import ID is invalid or plugin does not exist")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plugin.ID)...)
}
