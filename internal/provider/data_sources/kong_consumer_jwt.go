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
	_ datasource.DataSource              = &kongConsumerJwtDataSource{}
	_ datasource.DataSourceWithConfigure = &kongConsumerJwtDataSource{}
)

func KongConsumerJwtDataSource() datasource.DataSource {
	return &kongConsumerJwtDataSource{}
}

type kongConsumerJwtDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongConsumerJwtDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consumer_jwt"
}

// Configure adds the provider configured client to the data source.
func (d *kongConsumerJwtDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongConsumerJwtDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"consumer_id":       schema.StringAttribute{Optional: true, Computed: true},
			"consumer_username": schema.StringAttribute{Optional: true},
			"id":                schema.StringAttribute{Optional: true, Computed: true},
			"key":               schema.StringAttribute{Optional: true, Computed: true},
			"algorithm":         schema.StringAttribute{Computed: true},
			"rsa_public_key":    schema.StringAttribute{Computed: true},
			"secret":            schema.StringAttribute{Computed: true, Sensitive: true},
			"tags":              schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongConsumerJwtDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type dataSourceModel struct {
		ConsumerID       types.String `tfsdk:"consumer_id"`
		ConsumerUsername types.String `tfsdk:"consumer_username"`
		ID               types.String `tfsdk:"id"`
		Key              types.String `tfsdk:"key"`
		Algorithm        types.String `tfsdk:"algorithm"`
		RsaPublicKey     types.String `tfsdk:"rsa_public_key"`
		Secret           types.String `tfsdk:"secret"`
		Tags             types.List   `tfsdk:"tags"`
	}
	var data dataSourceModel
	client := d.client.JWTAuths
	diags1 := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags1...)
	var consumerUsernameOrId *string
	if !data.ConsumerID.IsNull() && !data.ConsumerUsername.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'consumer_id' or 'consumer_username' must be set.")
	}
	if !data.ConsumerID.IsNull() {
		consumerUsernameOrId = data.ConsumerID.ValueStringPointer()
	} else if !data.ConsumerUsername.IsNull() {
		consumerUsernameOrId = data.ConsumerUsername.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'consumer_id' or 'consumer_username' must be set.")
	}
	var usernameOrId *string
	if !data.ID.IsNull() && !data.Key.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'id' or 'username' must be set.")
	}
	if !data.ID.IsNull() {
		usernameOrId = data.ID.ValueStringPointer()
	} else if !data.Key.IsNull() {
		usernameOrId = data.Key.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'id' or 'username' must be set.")
	}

	jwt, err := client.Get(ctx, consumerUsernameOrId, usernameOrId)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_consumer_jwt", err.Error())
		return
	}

	result := kongModels.ConsumerJwtModelFromResponse(jwt)
	diags2 := resp.State.Set(ctx, dataSourceModel{
		ID:               result.KongId,
		Key:              result.Key,
		ConsumerID:       result.ConsumerID,
		ConsumerUsername: data.ConsumerUsername,
		Algorithm:        result.Algorithm,
		RsaPublicKey:     result.RsaPublicKey,
		Secret:           result.Secret,
		Tags:             result.Tags,
	})
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
