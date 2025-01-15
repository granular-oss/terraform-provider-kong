package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongConsumerModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	CustomId types.String `tfsdk:"custom_id"`
	Tags     types.List   `tfsdk:"tags"`
}

func ConsumerModelFromResponse(consumer *kong.Consumer) KongConsumerModel {
	return KongConsumerModel{
		ID:       types.StringPointerValue(consumer.ID),
		Username: types.StringPointerValue(consumer.Username),
		CustomId: types.StringPointerValue(consumer.CustomID),
		Tags:     ParseResponseStringList(consumer.Tags),
	}
}

func RequestFromConsumerModel(ctx context.Context, model *KongConsumerModel) *kong.Consumer {
	return &kong.Consumer{
		ID:       TfStringToKongString(model.ID),
		Username: TfStringToKongString(model.Username),
		CustomID: TfStringToKongString(model.CustomId),
		Tags:     TFListToKongStringArray(ctx, model.Tags),
	}
}
