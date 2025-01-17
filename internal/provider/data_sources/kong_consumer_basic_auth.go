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
	_ datasource.DataSource              = &kongConsumerBasicAuthDataSource{}
	_ datasource.DataSourceWithConfigure = &kongConsumerBasicAuthDataSource{}
)

func KongConsumerBasicAuthDataSource() datasource.DataSource {
	return &kongConsumerBasicAuthDataSource{}
}

type kongConsumerBasicAuthDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongConsumerBasicAuthDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consumer_basic_auth"
}

// Configure adds the provider configured client to the data source.
func (d *kongConsumerBasicAuthDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongConsumerBasicAuthDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"consumer_id":       schema.StringAttribute{Optional: true, Computed: true},
			"consumer_username": schema.StringAttribute{Optional: true},
			"id":                schema.StringAttribute{Optional: true, Computed: true},
			"username":          schema.StringAttribute{Optional: true, Computed: true},
			"tags":              schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongConsumerBasicAuthDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type dataSourceModel struct {
		ID               types.String `tfsdk:"id"`
		Username         types.String `tfsdk:"username"`
		ConsumerID       types.String `tfsdk:"consumer_id"`
		ConsumerUsername types.String `tfsdk:"consumer_username"`
		Tags             types.List   `tfsdk:"tags"`
	}
	var data dataSourceModel
	client := d.client.BasicAuths
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
	if !data.ID.IsNull() && !data.Username.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'id' or 'username' must be set.")
	}
	if !data.ID.IsNull() {
		usernameOrId = data.ID.ValueStringPointer()
	} else if !data.Username.IsNull() {
		usernameOrId = data.Username.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'id' or 'username' must be set.")
	}

	acl, err := client.Get(ctx, consumerUsernameOrId, usernameOrId)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_consumer_acl", err.Error())
		return
	}

	result := kongModels.ConsumerBasicAuthModelFromResponse(acl)
	diags2 := resp.State.Set(ctx, dataSourceModel{
		ID:               result.KongId,
		Username:         result.Username,
		ConsumerID:       result.ConsumerID,
		ConsumerUsername: data.ConsumerUsername,
		Tags:             result.Tags,
	})
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
