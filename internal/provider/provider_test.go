// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const (
	providerConfig = `
provider "kong" {
	kong_admin_uri = "http://localhost:8003"
}
`
)

func buildJsonPath(path string) tfjsonpath.Path {
	parts := strings.Split(path, ".")
	var ret tfjsonpath.Path
	for i, p := range parts {
		if i == 0 {
			ret = tfjsonpath.New(p)
		} else {
			ret = ret.AtMapKey(p)
		}
	}
	return ret
}

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"kong": providerserver.NewProtocol6WithError(New("test")()),
	}
	assertNullValue = func(resource_path string, attribute string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(resource_path, buildJsonPath(attribute), knownvalue.Null())
	}
	assertStringValue = func(resource_path string, attribute string, expected string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(resource_path, buildJsonPath(attribute), knownvalue.StringExact(expected))
	}
	assertStringArrayValue = func(resource_path string, attribute string, expected []string) statecheck.StateCheck {
		compare := []knownvalue.Check{}
		for _, val := range expected {
			compare = append(compare, knownvalue.StringExact(val))
		}
		return statecheck.ExpectKnownValue(resource_path, buildJsonPath(attribute), knownvalue.ListExact(compare))
	}
	assertInt32Value = func(resource_path string, attribute string, expected int32) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(resource_path, buildJsonPath(attribute), knownvalue.Int32Exact(expected))
	}
	assertNotNull = func(resource_path string, attribute string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(resource_path, buildJsonPath(attribute), knownvalue.NotNull())
	}
	assertBoolValue = func(resource_path string, attribute string, expected bool) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(resource_path, buildJsonPath(attribute), knownvalue.Bool(expected))
	}
	assertStateMatch = func(resource_path1 string, attribute1 string, resource_path2 string, attribute2 string) statecheck.StateCheck {
		return statecheck.CompareValuePairs(resource_path1, buildJsonPath(attribute1), resource_path2, tfjsonpath.New(attribute2), compare.ValuesSame())
	}
)
