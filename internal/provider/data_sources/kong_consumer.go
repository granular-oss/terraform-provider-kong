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
	_ datasource.DataSource              = &kongConsumerDataSource{}
	_ datasource.DataSourceWithConfigure = &kongConsumerDataSource{}
)

func KongConsumerDataSource() datasource.DataSource {
	return &kongConsumerDataSource{}
}

type kongConsumerDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongConsumerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consumer"
}

// Configure adds the provider configured client to the data source.
func (d *kongConsumerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongConsumerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Optional: true},
			"username":  schema.StringAttribute{Optional: true},
			"custom_id": schema.StringAttribute{Computed: true},
			"tags":      schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongConsumerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data kongModels.KongConsumerModel
	diags1 := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags1...)

	var lookupProp *string
	if !data.Username.IsNull() && !data.ID.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'consumer_id' or 'consumer_username' must be set.")
	}
	if !data.ID.IsNull() {
		lookupProp = data.ID.ValueStringPointer()
	} else if !data.Username.IsNull() {
		lookupProp = data.Username.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'consumer_id' or 'consumer_username' must be set.")
	}

	consumer, err := d.client.Consumers.Get(ctx, lookupProp)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_consumer", err.Error())
		return
	}

	result := kongModels.ConsumerModelFromResponse(consumer)
	diags2 := resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
