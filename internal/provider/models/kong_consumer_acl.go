package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongConsumerAclModel struct {
	ID         types.String `tfsdk:"id"`
	KongId     types.String `tfsdk:"kong_id"`
	Group      types.String `tfsdk:"group"`
	ConsumerID types.String `tfsdk:"consumer_id"`
	Tags       types.List   `tfsdk:"tags"`
}

func ConsumerAclModelFromResponse(acl *kong.ACLGroup) KongConsumerAclModel {
	return KongConsumerAclModel{
		ID:         types.StringValue(BuildCompositeId([]string{*acl.Consumer.ID, *acl.Group})),
		KongId:     types.StringPointerValue(acl.ID),
		ConsumerID: types.StringPointerValue(acl.Consumer.ID),
		Group:      types.StringPointerValue(acl.Group),
		Tags:       ParseResponseStringList(acl.Tags),
	}
}

func RequestFromConsumerAclModel(ctx context.Context, model *KongConsumerAclModel) *kong.ACLGroup {
	return &kong.ACLGroup{
		ID:    TfStringToKongString(model.KongId),
		Group: TfStringToKongString(model.Group),
		Tags:  TFListToKongStringArray(ctx, model.Tags),
		Consumer: &kong.Consumer{
			ID: TfStringToKongString(model.ConsumerID),
		},
	}
}
