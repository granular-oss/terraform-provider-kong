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
	_ datasource.DataSource              = &kongUpstreamDataSource{}
	_ datasource.DataSourceWithConfigure = &kongUpstreamDataSource{}
)

func KongUpstreamDataSource() datasource.DataSource {
	return &kongUpstreamDataSource{}
}

type kongUpstreamDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongUpstreamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_upstream"
}

// Configure adds the provider configured client to the data source.
func (d *kongUpstreamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongUpstreamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                    schema.StringAttribute{Optional: true},
			"name":                  schema.StringAttribute{Optional: true},
			"slots":                 schema.Int32Attribute{Computed: true},
			"hash_on":               schema.StringAttribute{Computed: true},
			"hash_fallback":         schema.StringAttribute{Computed: true},
			"hash_on_header":        schema.StringAttribute{Computed: true},
			"host_header":           schema.StringAttribute{Computed: true},
			"tags":                  schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"client_certificate_id": schema.StringAttribute{Computed: true},
			"hash_fallback_header":  schema.StringAttribute{Computed: true},
			"hash_on_cookie":        schema.StringAttribute{Computed: true},
			"hash_on_cookie_path":   schema.StringAttribute{Computed: true},
			"healthchecks": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
				"active": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
					"type":                     schema.StringAttribute{Computed: true},
					"timeout":                  schema.Int32Attribute{Computed: true},
					"concurrency":              schema.Int32Attribute{Computed: true},
					"http_path":                schema.StringAttribute{Computed: true},
					"https_verify_certificate": schema.BoolAttribute{Computed: true},
					"https_sni":                schema.StringAttribute{Computed: true},
					"healthy": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
						"interval":      schema.Int32Attribute{Computed: true},
						"http_statuses": schema.ListAttribute{Computed: true, ElementType: types.Int32Type},
						"successes":     schema.Int32Attribute{Computed: true},
					}},
					"unhealthy": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
						"interval":      schema.Int32Attribute{Computed: true},
						"http_statuses": schema.ListAttribute{Computed: true, ElementType: types.Int32Type},
						"tcp_failures":  schema.Int32Attribute{Computed: true},
						"http_failures": schema.Int32Attribute{Computed: true},
						"timeouts":      schema.Int32Attribute{Computed: true},
					}},
				}},
				"passive": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{Computed: true},
					"healthy": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
						"http_statuses": schema.ListAttribute{Computed: true, ElementType: types.Int32Type},
						"successes":     schema.Int32Attribute{Computed: true},
					}},
					"unhealthy": schema.SingleNestedAttribute{Computed: true, Attributes: map[string]schema.Attribute{
						"http_statuses": schema.ListAttribute{Computed: true, ElementType: types.Int32Type},
						"tcp_failures":  schema.Int32Attribute{Computed: true},
						"http_failures": schema.Int32Attribute{Computed: true},
						"timeouts":      schema.Int32Attribute{Computed: true},
					}},
				}},
			}},
		},
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *kongUpstreamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data kongModels.KongUpstreamModel
	client := d.client.Upstreams
	diags1 := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags1...)
	var lookupProp *string
	if !data.ID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'id' or 'name' must be set.")
		return
	}
	if !data.ID.IsNull() {
		lookupProp = data.ID.ValueStringPointer()
	} else if !data.Name.IsNull() {
		lookupProp = data.Name.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'id' or 'name' must be set.")
		return
	}
	upstream, err := client.Get(ctx, lookupProp)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_upstream", err.Error())
		return
	}

	result := kongModels.UpstreamModelFromResponse(resp.Diagnostics, upstream)
	diags2 := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
