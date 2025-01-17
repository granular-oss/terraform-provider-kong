package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/kong/go-kong/kong"
)

type kongHealthyModel struct {
	HttpStatuses types.List  `tfsdk:"http_statuses"`
	Successes    types.Int32 `tfsdk:"successes"`
}
type kongUnhealthyModel struct {
	HttpStatuses types.List  `tfsdk:"http_statuses"`
	TcpFailures  types.Int32 `tfsdk:"tcp_failures"`
	HttpFailures types.Int32 `tfsdk:"http_failures"`
	Timeouts     types.Int32 `tfsdk:"timeouts"`
}
type kongActiveHealthyModel struct {
	Interval types.Int32 `tfsdk:"interval"`
	kongHealthyModel
}
type kongActiveUnhealthyModel struct {
	Interval types.Int32 `tfsdk:"interval"`
	kongUnhealthyModel
}
type kongActiveHealthcheckModel struct {
	Type                   types.String `tfsdk:"type"`
	Timeout                types.Int32  `tfsdk:"timeout"`
	Concurrency            types.Int32  `tfsdk:"concurrency"`
	HttpPath               types.String `tfsdk:"http_path"`
	HttpsVerifyCertificate types.Bool   `tfsdk:"https_verify_certificate"`
	HttpsSNI               types.String `tfsdk:"https_sni"`
	Healthy                types.Object `tfsdk:"healthy"`
	UnHealthy              types.Object `tfsdk:"unhealthy"`
}
type kongPassiveHealthcheckModel struct {
	Type      types.String `tfsdk:"type"`
	Healthy   types.Object `tfsdk:"healthy"`
	UnHealthy types.Object `tfsdk:"unhealthy"`
}
type kongHealthcheckModel struct {
	Active  types.Object `tfsdk:"active"`
	Passive types.Object `tfsdk:"passive"`
}
type KongUpstreamModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Slots               types.Int32  `tfsdk:"slots"`
	HashOn              types.String `tfsdk:"hash_on"`
	HashFallback        types.String `tfsdk:"hash_fallback"`
	HashOnHeader        types.String `tfsdk:"hash_on_header"`
	HostHeader          types.String `tfsdk:"host_header"`
	Tags                types.List   `tfsdk:"tags"`
	ClientCertificateId types.String `tfsdk:"client_certificate_id"`
	HashFallbackHeader  types.String `tfsdk:"hash_fallback_header"`
	HashOnCookie        types.String `tfsdk:"hash_on_cookie"`
	HashOnCookiePath    types.String `tfsdk:"hash_on_cookie_path"`
	Healthchecks        types.Object `tfsdk:"healthchecks"`
}

var (
	kongHealthyModelAttributes = map[string]attr.Type{
		"http_statuses": types.ListType{ElemType: types.Int32Type},
		"successes":     types.Int32Type,
	}
	kongUnhealthyModelAttributes = map[string]attr.Type{
		"http_statuses": types.ListType{ElemType: types.Int32Type},
		"tcp_failures":  types.Int32Type,
		"http_failures": types.Int32Type,
		"timeouts":      types.Int32Type,
	}
	kongActiveHealthyModelAttributes = map[string]attr.Type{
		"interval":      types.Int32Type,
		"http_statuses": types.ListType{ElemType: types.Int32Type},
		"successes":     types.Int32Type,
	}
	kongActiveUnhealthyModelAttributes = map[string]attr.Type{
		"interval":      types.Int32Type,
		"http_statuses": types.ListType{ElemType: types.Int32Type},
		"tcp_failures":  types.Int32Type,
		"http_failures": types.Int32Type,
		"timeouts":      types.Int32Type,
	}
	kongActiveHealthcheckModelAttributes = map[string]attr.Type{
		"type":                     types.StringType,
		"timeout":                  types.Int32Type,
		"concurrency":              types.Int32Type,
		"http_path":                types.StringType,
		"https_verify_certificate": types.BoolType,
		"https_sni":                types.StringType,
		"healthy":                  types.ObjectType{AttrTypes: kongActiveHealthyModelAttributes},
		"unhealthy":                types.ObjectType{AttrTypes: kongActiveUnhealthyModelAttributes},
	}
	kongPassiveHealthcheckModelAttributes = map[string]attr.Type{
		"type":      types.StringType,
		"healthy":   types.ObjectType{AttrTypes: kongHealthyModelAttributes},
		"unhealthy": types.ObjectType{AttrTypes: kongUnhealthyModelAttributes},
	}
	kongHealthcheckModelAttributes = map[string]attr.Type{
		"active":  types.ObjectType{AttrTypes: kongActiveHealthcheckModelAttributes},
		"passive": types.ObjectType{AttrTypes: kongPassiveHealthcheckModelAttributes},
	}
)

func activeHealthyModelFromResponse(healthy *kong.Healthy) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(kongActiveHealthyModelAttributes, map[string]attr.Value{
		"interval":      types.Int32Value(int32(*healthy.Interval)),
		"http_statuses": ParseResponseIntList(healthy.HTTPStatuses),
		"successes":     types.Int32Value(int32(*healthy.Interval)),
	})
}
func activeUnhealthyModelFromResponse(unhealthy *kong.Unhealthy) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(kongActiveUnhealthyModelAttributes, map[string]attr.Value{
		"interval":      types.Int32Value(int32(*unhealthy.Interval)),
		"http_statuses": ParseResponseIntList(unhealthy.HTTPStatuses),
		"tcp_failures":  types.Int32Value(int32(*unhealthy.TCPFailures)),
		"http_failures": types.Int32Value(int32(*unhealthy.HTTPFailures)),
		"timeouts":      types.Int32Value(int32(*unhealthy.Timeouts)),
	})
}
func healthcheckActiveModelFromResponse(active *kong.ActiveHealthcheck) (basetypes.ObjectValue, diag.Diagnostics) {
	healthy, hDiags := activeHealthyModelFromResponse(active.Healthy)
	unhealthy, uDiags := activeUnhealthyModelFromResponse(active.Unhealthy)
	act, diags := types.ObjectValue(kongActiveHealthcheckModelAttributes, map[string]attr.Value{
		"type":                     types.StringPointerValue(active.Type),
		"timeout":                  types.Int32Value(int32(*active.Timeout)),
		"concurrency":              types.Int32Value(int32(*active.Concurrency)),
		"http_path":                types.StringPointerValue(active.HTTPPath),
		"https_verify_certificate": types.BoolPointerValue(active.HTTPSVerifyCertificate),
		"https_sni":                types.StringPointerValue(active.HTTPSSni),
		"healthy":                  healthy,
		"unhealthy":                unhealthy,
	})
	diags.Append(hDiags...)
	diags.Append(uDiags...)
	return act, diags
}
func passiveHealthyModelFromResponse(healthy *kong.Healthy) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(kongHealthyModelAttributes, map[string]attr.Value{
		"http_statuses": ParseResponseIntList(healthy.HTTPStatuses),
		"successes":     types.Int32Value(int32(*healthy.Successes)),
	})
}
func passiveUnhealthyModelFromResponse(unhealthy *kong.Unhealthy) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(kongUnhealthyModelAttributes, map[string]attr.Value{
		"http_statuses": ParseResponseIntList(unhealthy.HTTPStatuses),
		"tcp_failures":  types.Int32Value(int32(*unhealthy.TCPFailures)),
		"http_failures": types.Int32Value(int32(*unhealthy.HTTPFailures)),
		"timeouts":      types.Int32Value(int32(*unhealthy.Timeouts)),
	})
}
func healthcheckPassiveModelFromResponse(passive *kong.PassiveHealthcheck) (basetypes.ObjectValue, diag.Diagnostics) {
	healthy, hDiags := passiveHealthyModelFromResponse(passive.Healthy)
	unhealthy, uDiags := passiveUnhealthyModelFromResponse(passive.Unhealthy)
	val, diags := types.ObjectValue(kongPassiveHealthcheckModelAttributes, map[string]attr.Value{
		"type":      types.StringPointerValue(passive.Type),
		"healthy":   healthy,
		"unhealthy": unhealthy,
	})
	diags.Append(hDiags...)
	diags.Append(uDiags...)
	return val, diags
}
func healthcheckModelFromResponse(healthcheck *kong.Healthcheck) (basetypes.ObjectValue, diag.Diagnostics) {
	active, aDiags := healthcheckActiveModelFromResponse(healthcheck.Active)
	passive, pDiags := healthcheckPassiveModelFromResponse(healthcheck.Passive)
	obj, diags := types.ObjectValue(kongHealthcheckModelAttributes, map[string]attr.Value{
		"active":  active,
		"passive": passive,
	})
	diags.Append(aDiags...)
	diags.Append(pDiags...)
	return obj, diags
}

func UpstreamModelFromResponse(diags diag.Diagnostics, upstream *kong.Upstream) KongUpstreamModel {
	var clientCertId types.String
	if upstream.ClientCertificate != nil {
		clientCertId = types.StringPointerValue(upstream.ClientCertificate.ID)
	}
	health, hDiags := healthcheckModelFromResponse(upstream.Healthchecks)
	diags.Append(hDiags...)
	return KongUpstreamModel{
		ID:                  types.StringPointerValue(upstream.ID),
		Name:                types.StringPointerValue(upstream.Name),
		Slots:               types.Int32Value(int32(*upstream.Slots)),
		HashOn:              types.StringPointerValue(upstream.HashOn),
		HashFallback:        types.StringPointerValue(upstream.HashFallback),
		HashOnHeader:        types.StringPointerValue(upstream.HashOnHeader),
		HostHeader:          types.StringPointerValue(upstream.HostHeader),
		Tags:                ParseResponseStringList(upstream.Tags),
		ClientCertificateId: clientCertId,
		HashFallbackHeader:  types.StringPointerValue(upstream.HashFallbackHeader),
		HashOnCookie:        types.StringPointerValue(upstream.HashOnCookie),
		HashOnCookiePath:    types.StringPointerValue(upstream.HashOnCookiePath),
		Healthchecks:        health,
	}
}

func kongHealthyFromModel(ctx context.Context, model *kongActiveHealthyModel) *kong.Healthy {
	return &kong.Healthy{
		HTTPStatuses: TfListToIntArray(ctx, model.HttpStatuses),
		Interval:     TfIntToKongInt(model.Interval),
		Successes:    TfIntToKongInt(model.Successes),
	}
}
func kongUnhealthyFromModel(ctx context.Context, model *kongActiveUnhealthyModel) *kong.Unhealthy {
	return &kong.Unhealthy{
		HTTPFailures: TfIntToKongInt(model.HttpFailures),
		HTTPStatuses: TfListToIntArray(ctx, model.HttpStatuses),
		TCPFailures:  TfIntToKongInt(model.TcpFailures),
		Timeouts:     TfIntToKongInt(model.Timeouts),
		Interval:     TfIntToKongInt(model.Interval),
	}
}
func kongActiveHealthcheckFromModel(ctx context.Context, model *kongActiveHealthcheckModel) *kong.ActiveHealthcheck {
	ret := &kong.ActiveHealthcheck{
		Concurrency:            TfIntToKongInt(model.Concurrency),
		HTTPPath:               TfStringToKongString(model.HttpPath),
		HTTPSSni:               TfStringToKongString(model.HttpsSNI),
		HTTPSVerifyCertificate: TfBoolToKongBool(model.HttpsVerifyCertificate),
		Type:                   TfStringToKongString(model.Type),
		Timeout:                TfIntToKongInt(model.Timeout),
	}
	if !model.Healthy.IsNull() && !model.Healthy.IsUnknown() {
		var healthyModel kongActiveHealthyModel
		model.Healthy.As(ctx, healthyModel, basetypes.ObjectAsOptions{})
		ret.Healthy = kongHealthyFromModel(ctx, &healthyModel)
	}
	if !model.UnHealthy.IsNull() && !model.UnHealthy.IsUnknown() {
		var unhealthyModel kongActiveUnhealthyModel
		model.UnHealthy.As(ctx, unhealthyModel, basetypes.ObjectAsOptions{})
		ret.Unhealthy = kongUnhealthyFromModel(ctx, &unhealthyModel)
	}
	return ret
}
func kongPassiveHealthcheckFromModel(ctx context.Context, model *kongPassiveHealthcheckModel) *kong.PassiveHealthcheck {
	ret := &kong.PassiveHealthcheck{
		Type: TfStringToKongString(model.Type),
	}
	if !model.Healthy.IsNull() && !model.Healthy.IsUnknown() {
		var healthyModel kongActiveHealthyModel
		model.Healthy.As(ctx, healthyModel, basetypes.ObjectAsOptions{})
		ret.Healthy = kongHealthyFromModel(ctx, &healthyModel)
	}
	if !model.UnHealthy.IsNull() && !model.UnHealthy.IsUnknown() {
		var unhealthyModel kongActiveUnhealthyModel
		model.UnHealthy.As(ctx, unhealthyModel, basetypes.ObjectAsOptions{})
		ret.Unhealthy = kongUnhealthyFromModel(ctx, &unhealthyModel)
	}
	return ret
}

func kongHealthcheckFromModel(ctx context.Context, model *kongHealthcheckModel) *kong.Healthcheck {
	ret := &kong.Healthcheck{}
	if !model.Active.IsNull() && !model.Active.IsUnknown() {
		var activeModel kongActiveHealthcheckModel
		model.Active.As(ctx, &activeModel, basetypes.ObjectAsOptions{})
		ret.Active = kongActiveHealthcheckFromModel(ctx, &activeModel)
	}
	if !model.Passive.IsNull() && !model.Passive.IsUnknown() {
		var passiveModel kongPassiveHealthcheckModel
		model.Active.As(ctx, &passiveModel, basetypes.ObjectAsOptions{})
		ret.Passive = kongPassiveHealthcheckFromModel(ctx, &passiveModel)
	}
	return ret
}

func RequestFromUpstreamModel(ctx context.Context, model *KongUpstreamModel) *kong.Upstream {
	upstream := &kong.Upstream{
		ID:                 TfStringToKongString(model.ID),
		Name:               TfStringToKongString(model.Name),
		Slots:              TfIntToKongInt(model.Slots),
		HashOn:             TfStringToKongString(model.HashOn),
		HashFallback:       TfStringToKongString(model.HashFallback),
		HashOnHeader:       TfStringToKongString(model.HashOnHeader),
		HostHeader:         TfStringToKongString(model.HostHeader),
		Tags:               TFListToKongStringArray(ctx, model.Tags),
		HashFallbackHeader: TfStringToKongString(model.HashFallbackHeader),
		HashOnCookie:       TfStringToKongString(model.HashOnCookie),
		HashOnCookiePath:   TfStringToKongString(model.HashOnCookiePath),
	}
	if !model.ClientCertificateId.IsUnknown() && !model.ClientCertificateId.IsNull() {
		upstream.ClientCertificate = &kong.Certificate{ID: model.ClientCertificateId.ValueStringPointer()}
	}
	if !model.Healthchecks.IsNull() && !model.Healthchecks.IsUnknown() {
		var healthcheckModel kongHealthcheckModel
		model.Healthchecks.As(ctx, &healthcheckModel, basetypes.ObjectAsOptions{})
		upstream.Healthchecks = kongHealthcheckFromModel(ctx, &healthcheckModel)
	}
	return upstream
}
