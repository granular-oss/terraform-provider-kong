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
	_ datasource.DataSource              = &kongServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &kongServiceDataSource{}
)

func NewKongServiceDataSource() datasource.DataSource {
	return &kongServiceDataSource{}
}

type kongServiceDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongServiceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

// Configure adds the provider configured client to the data source.
func (d *kongServiceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongServiceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                    schema.StringAttribute{Optional: true},
			"name":                  schema.StringAttribute{Optional: true},
			"protocol":              schema.StringAttribute{Computed: true},
			"host":                  schema.StringAttribute{Computed: true},
			"port":                  schema.Int32Attribute{Computed: true},
			"path":                  schema.StringAttribute{Computed: true},
			"retries":               schema.Int32Attribute{Computed: true},
			"connect_timeout":       schema.Int32Attribute{Computed: true},
			"write_timeout":         schema.Int32Attribute{Computed: true},
			"read_timeout":          schema.Int32Attribute{Computed: true},
			"tags":                  schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"tls_verify":            schema.BoolAttribute{Computed: true},
			"tls_depth":             schema.Int32Attribute{Computed: true},
			"client_certificate_id": schema.StringAttribute{Computed: true},
			"ca_certificate_ids":    schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data kongModels.KongServiceModel
	client := d.client.Services
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
	service, err := client.Get(ctx, lookupProp)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_service", err.Error())
		return
	}

	result := kongModels.ServiceModelFromResponse(service)
	diags2 := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
