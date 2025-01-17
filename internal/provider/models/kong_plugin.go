// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

type KongPluginModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ConsumerId     types.String `tfsdk:"consumer_id"`
	ServiceId      types.String `tfsdk:"service_id"`
	RouteId        types.String `tfsdk:"route_id"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	ConfigJson     types.String `tfsdk:"config_json"`
	ComputedConfig types.String `tfsdk:"computed_config"`
	StrictMatch    types.Bool   `tfsdk:"strict_match"`
	Tags           types.List   `tfsdk:"tags"`
}

func PluginModelFromResponse(diag diag.Diagnostics, plugin *kong.Plugin, plan KongPluginModel) KongPluginModel {
	b, err := json.Marshal(plugin.Config)
	if err != nil {
		diag.AddError("Failed to parse plugin config", err.Error())
	}
	configJsonString := string(b)
	var serviceId *string
	if plugin.Service != nil {
		serviceId = plugin.Service.ID
	}
	var routeId *string
	if plugin.Route != nil {
		routeId = plugin.Route.ID
	}
	var consumerId *string
	if plugin.Consumer != nil {
		consumerId = plugin.Consumer.ID
	}
	var configJson = plan.ConfigJson
	if plan.StrictMatch.ValueBool() {
		configJson = types.StringValue(configJsonString)
	}
	return KongPluginModel{
		ID:             types.StringPointerValue(plugin.ID),
		Name:           types.StringPointerValue(plugin.Name),
		ConsumerId:     types.StringPointerValue(consumerId),
		ServiceId:      types.StringPointerValue(serviceId),
		RouteId:        types.StringPointerValue(routeId),
		Enabled:        types.BoolPointerValue(plugin.Enabled),
		ConfigJson:     configJson,
		ComputedConfig: types.StringValue(configJsonString),
		StrictMatch:    plan.StrictMatch,
		Tags:           ParseResponseStringList(plugin.Tags),
	}
}

func RequestFromPluginModel(ctx context.Context, diag diag.Diagnostics, model *KongPluginModel) *kong.Plugin {
	var configJson map[string]interface{}
	err := json.Unmarshal([]byte(model.ConfigJson.ValueString()), &configJson)
	if err != nil {
		diag.AddError("Invalid config json", err.Error())
	}
	plugin := &kong.Plugin{
		ID:     TfStringToKongString(model.ID),
		Name:   TfStringToKongString(model.Name),
		Tags:   TFListToKongStringArray(ctx, model.Tags),
		Config: configJson,
	}
	if !model.ServiceId.IsUnknown() && !model.ServiceId.IsNull() {
		plugin.Service = &kong.Service{ID: model.ServiceId.ValueStringPointer()}
	} else if !model.RouteId.IsUnknown() && !model.RouteId.IsNull() {
		plugin.Route = &kong.Route{ID: model.RouteId.ValueStringPointer()}
	} else if !model.ConsumerId.IsUnknown() && !model.ConsumerId.IsNull() {
		plugin.Consumer = &kong.Consumer{ID: model.ConsumerId.ValueStringPointer()}
	}
	return plugin
}

func GetPluginByService(ctx context.Context, diag diag.Diagnostics, client *kong.Client, serviceIdOrName *string, pluginIdOrName string) *kong.Plugin {
	plugins, err := client.Plugins.ListAllForService(ctx, serviceIdOrName)
	if err != nil {
		diag.AddError("Failed to list plugins for service", err.Error())
		return nil
	}
	for _, plug := range plugins {
		if *plug.ID == pluginIdOrName || *plug.Name == pluginIdOrName {
			return plug
		}
	}
	return nil
}

func GetPluginByRoute(ctx context.Context, diag diag.Diagnostics, client *kong.Client, routeIdOrName *string, pluginIdOrName string) *kong.Plugin {
	plugins, err := client.Plugins.ListAllForRoute(ctx, routeIdOrName)
	if err != nil {
		diag.AddError("Failed to list plugins for route", err.Error())
		return nil
	}
	for _, plug := range plugins {
		if *plug.ID == pluginIdOrName || *plug.Name == pluginIdOrName {
			return plug
		}
	}
	return nil
}

func GetPluginByConsumer(ctx context.Context, diag diag.Diagnostics, client *kong.Client, consumerIdOrName *string, pluginIdOrName string) *kong.Plugin {
	plugins, err := client.Plugins.ListAllForConsumer(ctx, consumerIdOrName)
	if err != nil {
		diag.AddError("Failed to list plugins for route", err.Error())
		return nil
	}
	for _, plug := range plugins {
		if *plug.ID == pluginIdOrName || *plug.Name == pluginIdOrName {
			return plug
		}
	}
	return nil
}
