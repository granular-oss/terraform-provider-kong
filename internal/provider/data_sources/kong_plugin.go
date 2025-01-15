package provider

import (
	"context"
	"fmt"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

var (
	_ datasource.DataSource              = &kongPluginDataSource{}
	_ datasource.DataSourceWithConfigure = &kongPluginDataSource{}
)

func KongPluginDataSource() datasource.DataSource {
	return &kongPluginDataSource{}
}

type kongPluginDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongPluginDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plugin"
}

// Configure adds the provider configured client to the data source.
func (d *kongPluginDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

// Schema defines the schema for the data source.
func (d *kongPluginDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Optional: true},
			"name":          schema.StringAttribute{Optional: true},
			"service_id":    schema.StringAttribute{Optional: true},
			"service_name":  schema.StringAttribute{Optional: true},
			"route_id":      schema.StringAttribute{Optional: true},
			"route_name":    schema.StringAttribute{Optional: true},
			"consumer_id":   schema.StringAttribute{Optional: true},
			"consumer_name": schema.StringAttribute{Optional: true},
			"enabled":       schema.BoolAttribute{Computed: true},
			"config_json":   schema.StringAttribute{Computed: true},
			"tags":          schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

type pluginDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ServiceId    types.String `tfsdk:"service_id"`
	ServiceName  types.String `tfsdk:"service_name"`
	RouteId      types.String `tfsdk:"route_id"`
	RouteName    types.String `tfsdk:"route_name"`
	ConsumerId   types.String `tfsdk:"consumer_id"`
	ConsumerName types.String `tfsdk:"consumer_name"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	ConfigJson   types.String `tfsdk:"config_json"`
	Tags         types.List   `tfsdk:"tags"`
}

// Read refreshes the Terraform state with the latest data.
func (d *kongPluginDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data pluginDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	var plugin *kong.Plugin
	var pluginIdOrName string
	if !data.ID.IsNull() {
		pluginIdOrName = data.ID.ValueString()
	} else if !data.Name.IsNull() {
		pluginIdOrName = data.Name.ValueString()
	}
	if !data.ServiceId.IsNull() {
		plugin = kongModels.GetPluginByService(ctx, resp.Diagnostics, d.client, data.ServiceId.ValueStringPointer(), pluginIdOrName)
	} else if !data.ServiceName.IsNull() {
		plugin = kongModels.GetPluginByService(ctx, resp.Diagnostics, d.client, data.ServiceName.ValueStringPointer(), pluginIdOrName)
	} else if !data.RouteId.IsNull() {
		plugin = kongModels.GetPluginByRoute(ctx, resp.Diagnostics, d.client, data.RouteId.ValueStringPointer(), pluginIdOrName)
	} else if !data.RouteName.IsNull() {
		plugin = kongModels.GetPluginByRoute(ctx, resp.Diagnostics, d.client, data.RouteName.ValueStringPointer(), pluginIdOrName)
	} else if !data.ConsumerId.IsNull() {
		plugin = kongModels.GetPluginByConsumer(ctx, resp.Diagnostics, d.client, data.ConsumerId.ValueStringPointer(), pluginIdOrName)
	} else if !data.ConsumerName.IsNull() {
		plugin = kongModels.GetPluginByConsumer(ctx, resp.Diagnostics, d.client, data.ConsumerName.ValueStringPointer(), pluginIdOrName)
	}
	if plugin == nil {
		resp.Diagnostics.AddError("Failed to find plugin.", "No plugin found for the provided information")
		return
	}
	result := kongModels.PluginModelFromResponse(resp.Diagnostics, plugin, kongModels.KongPluginModel{ConfigJson: types.StringValue("{}")})
	diags = resp.State.Set(ctx, pluginDataSourceModel{
		ID:           result.ID,
		Name:         result.Name,
		ServiceId:    result.ServiceId,
		ServiceName:  data.ServiceName,
		RouteId:      result.RouteId,
		RouteName:    data.RouteName,
		ConsumerId:   result.ConsumerId,
		ConsumerName: data.ConsumerName,
		Enabled:      result.Enabled,
		ConfigJson:   result.ComputedConfig,
		Tags:         result.Tags,
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
