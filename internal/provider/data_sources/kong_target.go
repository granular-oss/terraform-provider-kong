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
	_ datasource.DataSource              = &kongTargetDataSource{}
	_ datasource.DataSourceWithConfigure = &kongTargetDataSource{}
)

func KongTargetDataSource() datasource.DataSource {
	return &kongTargetDataSource{}
}

type kongTargetDataSource struct {
	client *kong.Client
}

// Metadata returns the data source type name.
func (d *kongTargetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_target"
}

// Configure adds the provider configured client to the data source.
func (d *kongTargetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *kongTargetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"upstream_id":   schema.StringAttribute{Optional: true},
			"upstream_name": schema.StringAttribute{Optional: true},
			"id":            schema.StringAttribute{Optional: true},
			"target":        schema.StringAttribute{Optional: true},
			"weight":        schema.Int32Attribute{Computed: true},

			"tags": schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *kongTargetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	type dataSourceModel struct {
		ID           types.String `tfsdk:"id"`
		Target       types.String `tfsdk:"target"`
		Weight       types.Int32  `tfsdk:"weight"`
		UpstreamId   types.String `tfsdk:"upstream_id"`
		UpstreamName types.String `tfsdk:"upstream_name"`
		Tags         types.List   `tfsdk:"tags"`
	}
	var data dataSourceModel
	diags1 := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags1...)

	var lookupProp string
	if !data.UpstreamId.IsNull() && !data.UpstreamName.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'upstream_id' or 'upstream_name' must be set.")
		return
	}
	if !data.UpstreamId.IsNull() {
		lookupProp = *data.UpstreamId.ValueStringPointer()
	} else if !data.UpstreamName.IsNull() {
		lookupProp = *data.UpstreamName.ValueStringPointer()
	} else {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'upstream_id' or 'upstream_name' must be set.")
		return
	}

	if !data.ID.IsNull() && !data.Target.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of 'id' or 'target' must be set.")
		return
	} else if data.ID.IsNull() && data.Target.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "One of 'id' or 'target' must be set.")
		return
	}

	var target *kong.Target
	targets, err := d.client.Targets.ListAll(ctx, &lookupProp)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read kong_target", err.Error())
		return
	}
	for _, t := range targets {
		if !data.ID.IsNull() && data.ID.ValueString() == *t.ID {
			target = t
			break
		} else if !data.Target.IsNull() && data.Target.ValueString() == *t.Target {
			target = t
			break
		}
	}
	if target == nil {
		resp.Diagnostics.AddError("Kong_target does not exist", "Could not find kong target for given upstream and id/target")
		return
	}

	result := kongModels.TargetModelFromResponse(target)
	diags2 := resp.State.Set(ctx, dataSourceModel{
		ID:           result.KongId,
		Target:       result.Target,
		Weight:       result.Weight,
		UpstreamId:   result.UpstreamId,
		UpstreamName: data.UpstreamName,
		Tags:         result.Tags,
	})
	resp.Diagnostics.Append(diags2...)
	if resp.Diagnostics.HasError() {
		return
	}
}
