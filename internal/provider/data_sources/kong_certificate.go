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
	_ datasource.DataSource              = &kongCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &kongCertificateDataSource{}
)

func NewKongCertificateDataSource() datasource.DataSource {
	return &kongCertificateDataSource{}
}

type kongCertificateDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongCertificateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

// Configure adds the provider configured client to the data source.
func (d *kongCertificateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongCertificateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Required: true},
			"cert": schema.StringAttribute{Computed: true, Sensitive: true},
			"key":  schema.StringAttribute{Computed: true, Sensitive: true},
			"snis": schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"tags": schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data kongModels.KongCertificateModel
	client := d.client.Certificates
	diags1 := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags1...)

	cert, err := client.Get(ctx, data.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_certificate", err.Error())
		return
	}

	result := kongModels.CertificateModelFromResponse(cert)
	diags2 := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
