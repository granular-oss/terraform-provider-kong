package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongConsumerBasicAuthModel struct {
	ID         types.String `tfsdk:"id"`
	KongId     types.String `tfsdk:"kong_id"`
	ConsumerID types.String `tfsdk:"consumer_id"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Tags       types.List   `tfsdk:"tags"`
}

func ConsumerBasicAuthModelFromResponse(basic *kong.BasicAuth) KongConsumerBasicAuthModel {
	id := *basic.Consumer.ID + ":" + *basic.ID
	return KongConsumerBasicAuthModel{
		ID:         types.StringValue(id),
		KongId:     types.StringPointerValue(basic.ID),
		ConsumerID: types.StringPointerValue(basic.Consumer.ID),
		Username:   types.StringPointerValue(basic.Username),
		Password:   types.StringPointerValue(basic.Password),
		Tags:       ParseResponseStringList(basic.Tags),
	}
}

func RequestFromConsumerBasicAuthModel(ctx context.Context, model *KongConsumerBasicAuthModel) *kong.BasicAuth {
	return &kong.BasicAuth{
		ID:       TfStringToKongString(model.KongId),
		Username: TfStringToKongString(model.Username),
		Password: TfStringToKongString(model.Password),
		Tags:     TFListToKongStringArray(ctx, model.Tags),
		Consumer: &kong.Consumer{
			ID: TfStringToKongString(model.ConsumerID),
		},
	}
}
