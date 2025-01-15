package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongConsumerKeyAuthModel struct {
	ID         types.String `tfsdk:"id"`
	KongId     types.String `tfsdk:"kong_id"`
	Key        types.String `tfsdk:"key"`
	ConsumerID types.String `tfsdk:"consumer_id"`
	Tags       types.List   `tfsdk:"tags"`
}

func ConsumerKeyAuthModelFromResponse(key *kong.KeyAuth) KongConsumerKeyAuthModel {
	id := *key.Consumer.ID + ":" + *key.ID
	return KongConsumerKeyAuthModel{
		ID:         types.StringValue(id),
		KongId:     types.StringPointerValue(key.ID),
		ConsumerID: types.StringPointerValue(key.Consumer.ID),
		Key:        types.StringPointerValue(key.Key),
		Tags:       ParseResponseStringList(key.Tags),
	}
}

func RequestFromConsumerKeyAuthModel(ctx context.Context, model *KongConsumerKeyAuthModel) *kong.KeyAuth {
	return &kong.KeyAuth{
		ID:   TfStringToKongString(model.KongId),
		Key:  TfStringToKongString(model.Key),
		Tags: TFListToKongStringArray(ctx, model.Tags),
		Consumer: &kong.Consumer{
			ID: TfStringToKongString(model.ConsumerID),
		},
	}
}
