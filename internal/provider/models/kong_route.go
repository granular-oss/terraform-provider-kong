// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongRouteSource struct {
	IP   types.String `tfsdk:"ip"`
	Port types.Int32  `tfsdk:"port"`
}

type KongRouteHeader struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

type KongRouteModel struct {
	ID                      types.String      `tfsdk:"id"`
	Name                    types.String      `tfsdk:"name"`
	Protocols               types.List        `tfsdk:"protocols"`
	Methods                 types.List        `tfsdk:"methods"`
	Hosts                   types.List        `tfsdk:"hosts"`
	Paths                   types.List        `tfsdk:"paths"`
	StripPath               types.Bool        `tfsdk:"strip_path"`
	Source                  []KongRouteSource `tfsdk:"source"`
	Destination             []KongRouteSource `tfsdk:"destination"`
	Snis                    types.List        `tfsdk:"snis"`
	PreserveHost            types.Bool        `tfsdk:"preserve_host"`
	RegexPriority           types.Int32       `tfsdk:"regex_priority"`
	ServiceId               types.String      `tfsdk:"service_id"`
	PathHandling            types.String      `tfsdk:"path_handling"`
	HttpsRedirectStatusCode types.Int32       `tfsdk:"https_redirect_status_code"`
	RequestBuffering        types.Bool        `tfsdk:"request_buffering"`
	ResponseBuffering       types.Bool        `tfsdk:"response_buffering"`
	Tags                    types.List        `tfsdk:"tags"`
	Header                  []KongRouteHeader `tfsdk:"header"`
}

func RouteModelFromResponse(route *kong.Route) KongRouteModel {
	var sources []KongRouteSource
	if route.Sources != nil {
		sources = []KongRouteSource{}
		for _, source := range route.Sources {
			sources = append(sources, KongRouteSource{
				IP:   types.StringPointerValue(source.IP),
				Port: types.Int32Value(int32(*source.Port)),
			})
		}
	}
	var destinations []KongRouteSource
	if route.Destinations != nil {
		destinations = []KongRouteSource{}
		for _, dest := range route.Sources {
			destinations = append(sources, KongRouteSource{
				IP:   types.StringPointerValue(dest.IP),
				Port: types.Int32Value(int32(*dest.Port)),
			})
		}
	}
	var header []KongRouteHeader
	if route.Headers != nil {
		header = []KongRouteHeader{}
		for name, values := range route.Headers {
			valList := []attr.Value{}
			for _, val := range values {
				valList = append(valList, types.StringValue(val))
			}
			header = append(header, KongRouteHeader{
				Name:   types.StringValue(name),
				Values: types.ListValueMust(types.StringType, valList),
			})
		}
	}
	return KongRouteModel{
		ID:                      types.StringPointerValue(route.ID),
		Name:                    types.StringPointerValue(route.Name),
		Protocols:               ParseResponseStringList(route.Protocols),
		Methods:                 ParseResponseStringList(route.Methods),
		Hosts:                   ParseResponseStringList(route.Hosts),
		Paths:                   ParseResponseStringList(route.Paths),
		Tags:                    ParseResponseStringList(route.Tags),
		StripPath:               types.BoolPointerValue(route.StripPath),
		Source:                  sources,
		Destination:             destinations,
		Snis:                    ParseResponseStringList(route.SNIs),
		PreserveHost:            types.BoolPointerValue(route.PreserveHost),
		RegexPriority:           types.Int32Value(int32(*route.RegexPriority)),
		ServiceId:               types.StringPointerValue(route.Service.ID),
		PathHandling:            types.StringPointerValue(route.PathHandling),
		HttpsRedirectStatusCode: types.Int32Value(int32(*route.HTTPSRedirectStatusCode)),
		ResponseBuffering:       types.BoolPointerValue(route.ResponseBuffering),
		RequestBuffering:        types.BoolPointerValue(route.RequestBuffering),
		Header:                  header,
	}
}

func RequestFromRouteModel(ctx context.Context, model *KongRouteModel) *kong.Route {
	source := []*kong.CIDRPort{}
	for _, s := range model.Source {
		source = append(source, &kong.CIDRPort{
			IP:   TfStringToKongString(s.IP),
			Port: TfIntToKongInt(s.Port),
		})
	}
	destination := []*kong.CIDRPort{}
	for _, d := range model.Destination {
		destination = append(destination, &kong.CIDRPort{
			IP:   TfStringToKongString(d.IP),
			Port: TfIntToKongInt(d.Port),
		})
	}
	header := make(map[string][]string)
	for _, head := range model.Header {
		values := []string{}
		valStrings := make([]types.String, 0, len(head.Values.Elements()))
		_ = head.Values.ElementsAs(ctx, &valStrings, false)
		for _, val := range valStrings {
			values = append(values, val.ValueString())
		}
		header[head.Name.ValueString()] = values
	}
	route := &kong.Route{
		ID:                      TfStringToKongString(model.ID),
		Name:                    TfStringToKongString(model.Name),
		Protocols:               TFListToKongStringArray(ctx, model.Protocols),
		Methods:                 TFListToKongStringArray(ctx, model.Methods),
		Hosts:                   TFListToKongStringArray(ctx, model.Hosts),
		Paths:                   TFListToKongStringArray(ctx, model.Paths),
		Tags:                    TFListToKongStringArray(ctx, model.Tags),
		StripPath:               TfBoolToKongBool(model.StripPath),
		Sources:                 source,
		Destinations:            destination,
		SNIs:                    TFListToKongStringArray(ctx, model.Snis),
		PreserveHost:            TfBoolToKongBool(model.PreserveHost),
		RegexPriority:           TfIntToKongInt(model.RegexPriority),
		Service:                 &kong.Service{ID: TfStringToKongString(model.ServiceId)},
		PathHandling:            TfStringToKongString(model.PathHandling),
		HTTPSRedirectStatusCode: TfIntToKongInt(model.HttpsRedirectStatusCode),
		RequestBuffering:        TfBoolToKongBool(model.RequestBuffering),
		ResponseBuffering:       TfBoolToKongBool(model.ResponseBuffering),
		Headers:                 header,
	}
	return route
}
