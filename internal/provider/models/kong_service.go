package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongServiceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Protocol            types.String `tfsdk:"protocol"`
	Host                types.String `tfsdk:"host"`
	Port                types.Int32  `tfsdk:"port"`
	Path                types.String `tfsdk:"path"`
	Retries             types.Int32  `tfsdk:"retries"`
	ConnectTimeout      types.Int32  `tfsdk:"connect_timeout"`
	WriteTimeout        types.Int32  `tfsdk:"write_timeout"`
	ReadTimeout         types.Int32  `tfsdk:"read_timeout"`
	Tags                types.List   `tfsdk:"tags"`
	TlsVerify           types.Bool   `tfsdk:"tls_verify"`
	TlsDepth            types.Int32  `tfsdk:"tls_depth"`
	ClientCertificateId types.String `tfsdk:"client_certificate_id"`
	CaCertificateIds    types.List   `tfsdk:"ca_certificate_ids"`
}

func ServiceModelFromResponse(service *kong.Service) KongServiceModel {
	var tlsVerify types.Bool
	if service.TLSVerify != nil {
		tlsVerify = types.BoolPointerValue(service.TLSVerify)
	}
	var tlsDepth types.Int32
	if service.TLSVerifyDepth != nil {
		tlsDepth = types.Int32Value(int32(*service.TLSVerifyDepth))
	}
	var clientCertId types.String
	if service.ClientCertificate != nil {
		clientCertId = types.StringPointerValue(service.ClientCertificate.ID)
	}
	return KongServiceModel{
		ID:                  types.StringPointerValue(service.ID),
		Name:                types.StringPointerValue(service.Name),
		Protocol:            types.StringPointerValue((service.Protocol)),
		Host:                types.StringPointerValue(service.Host),
		Port:                types.Int32Value(int32(*service.Port)),
		Path:                types.StringPointerValue(service.Path),
		Retries:             types.Int32Value(int32(*service.Retries)),
		ConnectTimeout:      types.Int32Value(int32(*service.ConnectTimeout)),
		WriteTimeout:        types.Int32Value(int32(*service.WriteTimeout)),
		ReadTimeout:         types.Int32Value(int32(*service.ReadTimeout)),
		TlsVerify:           tlsVerify,
		TlsDepth:            tlsDepth,
		ClientCertificateId: clientCertId,
		Tags:                ParseResponseStringList(service.Tags),
		CaCertificateIds:    ParseResponseStringList(service.CACertificates),
	}
}

func RequestFromServiceModel(ctx context.Context, model *KongServiceModel) *kong.Service {
	var clientCert *kong.Certificate
	if !model.ClientCertificateId.IsNull() {
		clientCert = &kong.Certificate{
			ID: model.ClientCertificateId.ValueStringPointer(),
		}
	}

	service := &kong.Service{
		ID:                TfStringToKongString(model.ID),
		Name:              TfStringToKongString(model.Name),
		Protocol:          TfStringToKongString(model.Protocol),
		Host:              TfStringToKongString(model.Host),
		Port:              TfIntToKongInt(model.Port),
		Path:              TfStringToKongString(model.Path),
		Retries:           TfIntToKongInt(model.Retries),
		ConnectTimeout:    TfIntToKongInt(model.ConnectTimeout),
		WriteTimeout:      TfIntToKongInt(model.WriteTimeout),
		ReadTimeout:       TfIntToKongInt(model.ReadTimeout),
		TLSVerify:         TfBoolToKongBool(model.TlsVerify),
		TLSVerifyDepth:    TfIntToKongInt(model.TlsDepth),
		Tags:              TFListToKongStringArray(ctx, model.Tags),
		ClientCertificate: clientCert,
		CACertificates:    TFListToKongStringArray(ctx, model.CaCertificateIds),
	}
	return service
}
