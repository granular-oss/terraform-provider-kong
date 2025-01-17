// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kong/go-kong/kong"
)

func TfStringToKongString(value types.String) *string {
	if value.IsUnknown() || value.IsNull() {
		return nil
	}
	return value.ValueStringPointer()
}

func TfIntToKongInt(value types.Int32) *int {
	if value.IsUnknown() || value.IsNull() {
		return nil
	}
	return kong.Int(int(value.ValueInt32()))
}

func TfBoolToKongBool(value types.Bool) *bool {
	if value.IsUnknown() || value.IsNull() {
		return nil
	}
	return kong.Bool(value.ValueBool())
}

func TFListToKongStringArray(ctx context.Context, value types.List) []*string {
	if value.IsUnknown() || value.IsNull() {
		return nil
	}
	toRet := []*string{}
	strings := make([]types.String, 0, len(value.Elements()))
	_ = value.ElementsAs(ctx, &strings, false)
	for _, val := range strings {
		toRet = append(toRet, val.ValueStringPointer())
	}
	return toRet
}
func TfListToIntArray(ctx context.Context, value types.List) []int {
	if value.IsUnknown() || value.IsNull() {
		return nil
	}
	toRet := []int{}
	ints := make([]types.Int32, 0, len(value.Elements()))
	_ = value.ElementsAs(ctx, &ints, false)
	for _, val := range ints {
		toRet = append(toRet, int(val.ValueInt32()))
	}
	return toRet
}

func ParseResponseStringList(value []*string) types.List {
	var ret []attr.Value
	if value != nil {
		ret = []attr.Value{}
		for _, val := range value {
			ret = append(ret, types.StringPointerValue(val))
		}
		return types.ListValueMust(types.StringType, ret)
	}
	return types.ListNull(types.StringType)
}

func ParseResponseIntList(value []int) types.List {
	var ret []attr.Value
	if value != nil {
		ret = []attr.Value{}
		for _, val := range value {
			ret = append(ret, types.Int32Value(int32(val)))
		}
		return types.ListValueMust(types.Int32Type, ret)
	}
	return types.ListNull(types.Int32Type)
}

func BuildCompositeId(parts []string) string {
	return strings.Join(parts, "|")
}

func ParseCompositeId(id string) []string {
	return strings.Split(id, "|")
}
