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
	_ datasource.DataSource              = &kongConsumerKeyAuthDataSource{}
	_ datasource.DataSourceWithConfigure = &kongConsumerKeyAuthDataSource{}
)

func NewKongConsumerKeyAuthDataSource() datasource.DataSource {
	return &kongConsumerKeyAuthDataSource{}
}

type kongConsumerKeyAuthDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongConsumerKeyAuthDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consumer_key_auth"
}

// Configure adds the provider configured client to the data source.
func (d *kongConsumerKeyAuthDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongConsumerKeyAuthDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                schema.StringAttribute{Optional: true},
			"key":               schema.StringAttribute{Optional: true},
			"consumer_id":       schema.StringAttribute{Optional: true},
			"consumer_username": schema.StringAttribute{Optional: true},
			"tags":              schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongConsumerKeyAuthDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type dataSourceModel struct {
		ID               types.String `tfsdk:"id"`
		Key              types.String `tfsdk:"key"`
		ConsumerID       types.String `tfsdk:"consumer_id"`
		ConsumerUsername types.String `tfsdk:"consumer_username"`
		Tags             types.List   `tfsdk:"tags"`
	}

	var data dataSourceModel
	client := d.client.KeyAuths
	diags1 := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags1...)
	var lookupProp *string
	if !data.ConsumerID.IsNull() && !data.ConsumerUsername.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'consumer_id' or 'consumer_username' must be set.")
	}
	if !data.ConsumerID.IsNull() {
		lookupProp = data.ConsumerID.ValueStringPointer()
	} else if !data.ConsumerUsername.IsNull() {
		lookupProp = data.ConsumerUsername.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'consumer_id' or 'consumer_username' must be set.")
	}
	var keyOrId *string
	if !data.Key.IsNull() && !data.ID.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'id' or 'key' must be set.")
	}
	if !data.ID.IsNull() {
		keyOrId = data.ID.ValueStringPointer()
	} else if !data.Key.IsNull() {
		keyOrId = data.Key.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'id' or 'key' must be set.")
	}

	key, err := client.Get(ctx, lookupProp, keyOrId)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_consumer_key_auth", err.Error())
		return
	}

	result := kongModels.ConsumerKeyAuthModelFromResponse(key)
	diags2 := resp.State.Set(ctx, dataSourceModel{
		ID:               result.ID,
		ConsumerID:       result.ConsumerID,
		Key:              result.Key,
		ConsumerUsername: data.ConsumerUsername,
		Tags:             result.Tags,
	})
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
