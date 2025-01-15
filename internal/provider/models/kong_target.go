// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongTargetModel struct {
	ID         types.String `tfsdk:"id"`
	KongId     types.String `tfsdk:"kong_id"`
	Target     types.String `tfsdk:"target"`
	Weight     types.Int32  `tfsdk:"weight"`
	UpstreamId types.String `tfsdk:"upstream_id"`
	Tags       types.List   `tfsdk:"tags"`
}

func TargetModelFromResponse(target *kong.Target) KongTargetModel {
	id := *target.Upstream.ID + "|" + *target.ID
	return KongTargetModel{
		ID:         types.StringValue(id),
		KongId:     types.StringPointerValue(target.ID),
		Target:     types.StringPointerValue(target.Target),
		Weight:     types.Int32Value(int32(*target.Weight)),
		UpstreamId: types.StringPointerValue(target.Upstream.ID),
		Tags:       ParseResponseStringList(target.Tags),
	}
}

func RequestFromTargetModel(ctx context.Context, model *KongTargetModel) *kong.Target {
	return &kong.Target{
		ID:       TfStringToKongString(model.KongId),
		Target:   TfStringToKongString(model.Target),
		Weight:   TfIntToKongInt(model.Weight),
		Upstream: &kong.Upstream{ID: TfStringToKongString(model.UpstreamId)},
		Tags:     TFListToKongStringArray(ctx, model.Tags),
	}
}
