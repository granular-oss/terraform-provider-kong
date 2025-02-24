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
	_ datasource.DataSource              = &kongRouteDataSource{}
	_ datasource.DataSourceWithConfigure = &kongRouteDataSource{}
)

func KongRouteDataSource() datasource.DataSource {
	return &kongRouteDataSource{}
}

type kongRouteDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongRouteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route"
}

// Configure adds the provider configured client to the data source.
func (d *kongRouteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongRouteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Optional: true},
			"name":       schema.StringAttribute{Optional: true},
			"protocols":  schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"methods":    schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"hosts":      schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"paths":      schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"tags":       schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"strip_path": schema.BoolAttribute{Computed: true},
			"source": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip":   schema.StringAttribute{Computed: true},
					"port": schema.Int32Attribute{Computed: true},
				},
			}},
			"destination": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip":   schema.StringAttribute{Computed: true},
					"port": schema.Int32Attribute{Computed: true},
				},
			}},
			"snis":                       schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"preserve_host":              schema.BoolAttribute{Computed: true},
			"regex_priority":             schema.Int32Attribute{Computed: true},
			"service_id":                 schema.StringAttribute{Computed: true},
			"path_handling":              schema.StringAttribute{Computed: true},
			"https_redirect_status_code": schema.Int32Attribute{Computed: true},
			"request_buffering":          schema.BoolAttribute{Computed: true},
			"response_buffering":         schema.BoolAttribute{Computed: true},
			"header": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name":  schema.StringAttribute{Computed: true},
					"value": schema.ListAttribute{Computed: true, ElementType: types.StringType},
				},
			}},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongRouteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data kongModels.KongRouteModel
	client := d.client.Routes
	diags1 := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags1...)
	var lookupProp string
	if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'id' or 'name' must be set.")
		return
	}
	if !data.ID.IsNull() {
		lookupProp = *data.ID.ValueStringPointer()
	} else if !data.Name.IsNull() {
		lookupProp = *data.Name.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'id' or 'name' must be set.")
		return
	}
	route, err := client.Get(ctx, &lookupProp)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_route", err.Error())
		return
	}

	result := kongModels.RouteModelFromResponse(route)
	diags2 := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
