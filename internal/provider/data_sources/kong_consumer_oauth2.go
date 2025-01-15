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
	_ datasource.DataSource              = &kongConsumerOauth2DataSource{}
	_ datasource.DataSourceWithConfigure = &kongConsumerOauth2DataSource{}
)

func NewKongConsumerOauth2DataSource() datasource.DataSource {
	return &kongConsumerOauth2DataSource{}
}

type kongConsumerOauth2DataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongConsumerOauth2DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consumer_oauth2"
}

// Configure adds the provider configured client to the data source.
func (d *kongConsumerOauth2DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongConsumerOauth2DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                schema.StringAttribute{Optional: true},
			"consumer_id":       schema.StringAttribute{Optional: true},
			"consumer_username": schema.StringAttribute{Optional: true},
			"name":              schema.StringAttribute{Computed: true},
			"client_id":         schema.StringAttribute{Computed: true},
			"client_secret":     schema.StringAttribute{Computed: true, Sensitive: true},
			"hash_secret":       schema.BoolAttribute{Computed: true},
			"redirect_uris":     schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"tags":              schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongConsumerOauth2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type dataSourceModel struct {
		ID               types.String `tfsdk:"id"`
		ConsumerID       types.String `tfsdk:"consumer_id"`
		ConsumerUsername types.String `tfsdk:"consumer_username"`
		Name             types.String `tfsdk:"name"`
		ClientId         types.String `tfsdk:"client_id"`
		ClientSecret     types.String `tfsdk:"client_secret"`
		HashSecret       types.Bool   `tfsdk:"hash_secret"`
		RedirectURIs     types.List   `tfsdk:"redirect_uris"`
		Tags             types.List   `tfsdk:"tags"`
	}

	var data dataSourceModel
	client := d.client.Oauth2Credentials
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

	oauth, err := client.Get(ctx, lookupProp, data.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_consumer_oauth2", err.Error())
		return
	}

	result := kongModels.ConsumerOauth2ModelFromResponse(oauth)
	diags2 := resp.State.Set(ctx, dataSourceModel{
		ID:               result.ID,
		ConsumerID:       result.ConsumerID,
		ClientId:         result.ClientId,
		ClientSecret:     result.ClientSecret,
		Name:             result.Name,
		HashSecret:       result.HashSecret,
		RedirectURIs:     result.RedirectURIs,
		ConsumerUsername: data.ConsumerUsername,
		Tags:             result.Tags,
	})
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
