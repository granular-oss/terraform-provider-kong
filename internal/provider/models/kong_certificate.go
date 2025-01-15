package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongCertificateModel struct {
	ID   types.String `tfsdk:"id"`
	Cert types.String `tfsdk:"cert"`
	Key  types.String `tfsdk:"key"`
	SNIs types.List   `tfsdk:"snis"`
	Tags types.List   `tfsdk:"tags"`
}

func CertificateModelFromResponse(cert *kong.Certificate) KongCertificateModel {
	return KongCertificateModel{
		ID:   types.StringPointerValue(cert.ID),
		Cert: types.StringPointerValue(cert.Cert),
		Key:  types.StringPointerValue(cert.Key),
		SNIs: ParseResponseStringList(cert.SNIs),
		Tags: ParseResponseStringList(cert.Tags),
	}
}

func RequestFromCertificateModel(ctx context.Context, model *KongCertificateModel) *kong.Certificate {
	return &kong.Certificate{
		ID:   TfStringToKongString(model.ID),
		Cert: TfStringToKongString(model.Cert),
		Key:  TfStringToKongString(model.Key),
		SNIs: TFListToKongStringArray(ctx, model.SNIs),
		Tags: TFListToKongStringArray(ctx, model.Tags),
	}
}
