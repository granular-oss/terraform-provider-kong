package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongConsumerJwtModel struct {
	ID           types.String `tfsdk:"id"`
	KongId       types.String `tfsdk:"kong_id"`
	ConsumerID   types.String `tfsdk:"consumer_id"`
	Key          types.String `tfsdk:"key"`
	Algorithm    types.String `tfsdk:"algorithm"`
	RsaPublicKey types.String `tfsdk:"rsa_public_key"`
	Secret       types.String `tfsdk:"secret"`
	Tags         types.List   `tfsdk:"tags"`
}

func ConsumerJwtModelFromResponse(jwt *kong.JWTAuth) KongConsumerJwtModel {
	id := *jwt.Consumer.ID + ":" + *jwt.ID
	return KongConsumerJwtModel{
		ID:           types.StringValue(id),
		KongId:       types.StringPointerValue(jwt.ID),
		ConsumerID:   types.StringPointerValue(jwt.Consumer.ID),
		Key:          types.StringPointerValue(jwt.Key),
		Algorithm:    types.StringPointerValue(jwt.Algorithm),
		RsaPublicKey: types.StringPointerValue(jwt.RSAPublicKey),
		Secret:       types.StringPointerValue(jwt.Secret),
		Tags:         ParseResponseStringList(jwt.Tags),
	}
}

func RequestFromConsumerJwtModel(ctx context.Context, model *KongConsumerJwtModel) *kong.JWTAuth {
	return &kong.JWTAuth{
		ID:           TfStringToKongString(model.KongId),
		Key:          TfStringToKongString(model.Key),
		Algorithm:    TfStringToKongString(model.Algorithm),
		RSAPublicKey: TfStringToKongString(model.RsaPublicKey),
		Secret:       TfStringToKongString(model.Secret),
		Tags:         TFListToKongStringArray(ctx, model.Tags),
		Consumer: &kong.Consumer{
			ID: TfStringToKongString(model.ConsumerID),
		},
	}
}
