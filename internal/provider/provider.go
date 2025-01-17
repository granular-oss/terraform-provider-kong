// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	kongData "github.com/granular-oss/terraform-provider-kong/internal/provider/data_sources"
	kongResources "github.com/granular-oss/terraform-provider-kong/internal/provider/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure Provider satisfies various provider interfaces.
var _ provider.Provider = &kongProvider{}

// Provider defines the provider implementation.
type kongProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProviderModel describes the provider data model.
type kongProviderModel struct {
	KongAdminUri       types.String `tfsdk:"kong_admin_uri"`
	KongAdminUsername  types.String `tfsdk:"kong_admin_username"`
	KongAdminPassword  types.String `tfsdk:"kong_admin_password"`
	TlsSkipVerify      types.Bool   `tfsdk:"tls_skip_verify"`
	KongApiKey         types.String `tfsdk:"kong_api_key"`
	KongAdminToken     types.String `tfsdk:"kong_admin_token"`
	KongWorkspace      types.String `tfsdk:"kong_workspace"`
	StrictPluginsMatch types.Bool   `tfsdk:"strict_plugins_match"`
}

func (p *kongProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kong"
	resp.Version = p.version
}

func (p *kongProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"kong_admin_uri": schema.StringAttribute{
				MarkdownDescription: "The address of the kong admin url e.g. http://localhost:8001",
				Optional:            true,
			},
			"kong_admin_username": schema.StringAttribute{
				MarkdownDescription: "A basic auth user for kong admin",
				Optional:            true,
			},
			"kong_admin_password": schema.StringAttribute{
				MarkdownDescription: "A basic auth password for kong admin",
				Optional:            true,
				Sensitive:           true,
			},
			"tls_skip_verify": schema.BoolAttribute{
				MarkdownDescription: "Whether to skip tls verify for https kong api endpoint using self signed or untrusted certs",
				Optional:            true,
			},
			"kong_api_key": schema.StringAttribute{
				MarkdownDescription: "API key for the kong api (if you have locked it down)",
				Optional:            true,
				Sensitive:           true,
			},
			"kong_admin_token": schema.StringAttribute{
				MarkdownDescription: "API key for the kong api (Enterprise Edition)",
				Optional:            true,
				Sensitive:           true,
			},
			"kong_workspace": schema.StringAttribute{
				MarkdownDescription: "Workspace context (Enterprise Edition)",
				Optional:            true,
			},
			"strict_plugins_match": schema.BoolAttribute{
				MarkdownDescription: "Should plugins `config_json` field strictly match plugin configuration",
				Optional:            true,
			},
		},
	}
}

func configWithEnvFallback(configValue types.String, envKey string, fallback string) string {
	value, exists := os.LookupEnv(envKey)
	if !configValue.IsNull() {
		value = configValue.ValueString()
	}
	if !exists {
		value = fallback
	}
	return value
}

func (p *kongProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config kongProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	kongConfig := &Config{
		Address:            configWithEnvFallback(config.KongAdminUri, "KONG_ADMIN_ADDR", "http://localhost:8001"),
		Username:           configWithEnvFallback(config.KongAdminUsername, "KONG_ADMIN_USERNAME", ""),
		Password:           configWithEnvFallback(config.KongAdminPassword, "KONG_ADMIN_PASSWORD", ""),
		InsecureSkipVerify: false,
		APIKey:             configWithEnvFallback(config.KongApiKey, "KONG_API_KEY", ""),
		AdminToken:         configWithEnvFallback(config.KongAdminToken, "KONG_ADMIN_TOKEN", ""),
		Workspace:          config.KongWorkspace.ValueString(),
	}
	client, err := GetKongClient(*kongConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create kong client",
			"An unexpected error occurred while creating kong client.\n\n"+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *kongProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		kongData.KongCertificateDataSource,
		kongData.KongConsumerDataSource,
		kongData.KongConsumerAclDataSource,
		kongData.KongConsumerBasicAuthDataSource,
		kongData.KongConsumerJwtDataSource,
		kongData.KongConsumerKeyAuthDataSource,
		kongData.KongConsumerOauth2DataSource,
		kongData.KongPluginDataSource,
		kongData.KongRouteDataSource,
		kongData.KongServiceDataSource,
		kongData.KongTargetDataSource,
		kongData.KongUpstreamDataSource,
	}
}

func (p *kongProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		kongResources.KongCertificateResource,
		kongResources.KongConsumerResource,
		kongResources.KongConsumerAclResource,
		kongResources.KongConsumerBasicAuthResource,
		kongResources.KongConsumerJwtResource,
		kongResources.KongConsumerKeyAuthResource,
		kongResources.KongConsumerOauth2Resource,
		kongResources.KongPluginResource,
		kongResources.KongServiceResource,
		kongResources.KongRouteResource,
		kongResources.KongTargetResource,
		kongResources.KongUpstreamResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kongProvider{
			version: version,
		}
	}
}
