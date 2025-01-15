package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongConsumerOauth2Model struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	KongId       types.String `tfsdk:"kong_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	HashSecret   types.Bool   `tfsdk:"hash_secret"`
	ConsumerID   types.String `tfsdk:"consumer_id"`
	RedirectURIs types.List   `tfsdk:"redirect_uris"`
	Tags         types.List   `tfsdk:"tags"`
}

func ConsumerOauth2ModelFromResponse(oauth *kong.Oauth2Credential) KongConsumerOauth2Model {
	return KongConsumerOauth2Model{
		ID:           types.StringValue(BuildCompositeId([]string{*oauth.Consumer.ID, *oauth.ID})),
		Name:         types.StringValue(*oauth.Name),
		KongId:       types.StringPointerValue(oauth.ID),
		ConsumerID:   types.StringPointerValue(oauth.Consumer.ID),
		ClientId:     types.StringPointerValue(oauth.ClientID),
		ClientSecret: types.StringPointerValue(oauth.ClientSecret),
		HashSecret:   types.BoolPointerValue(oauth.HashSecret),
		RedirectURIs: ParseResponseStringList(oauth.RedirectURIs),
		Tags:         ParseResponseStringList(oauth.Tags),
	}
}

func RequestFromConsumerOauth2Model(ctx context.Context, model *KongConsumerOauth2Model) *kong.Oauth2Credential {
	return &kong.Oauth2Credential{
		ID:           TfStringToKongString(model.KongId),
		Name:         TfStringToKongString(model.Name),
		ClientID:     TfStringToKongString(model.ClientId),
		ClientSecret: TfStringToKongString(model.ClientSecret),
		HashSecret:   TfBoolToKongBool(model.HashSecret),
		RedirectURIs: TFListToKongStringArray(ctx, model.RedirectURIs),
		Tags:         TFListToKongStringArray(ctx, model.Tags),
		Consumer: &kong.Consumer{
			ID: TfStringToKongString(model.ConsumerID),
		},
	}
}
