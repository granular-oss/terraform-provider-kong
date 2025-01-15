package provider

// import (
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/compare"
// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-testing/statecheck"
// 	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
// )

// func TestAccKongConsumerBasicAuthDataSource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_basic_auth" "test" {
// 	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
// 	id = "bd0f8189-e21c-4bac-9e4f-c2a7f2d04780"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_basic_auth.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:bd0f8189-e21c-4bac-9e4f-c2a7f2d04780"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "username", "test_basic"),
// 					assertNullValue("data.kong_consumer_basic_auth.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_basic_auth" "test" {
// 	consumer_username = "test"
// 	id = "bd0f8189-e21c-4bac-9e4f-c2a7f2d04780"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_basic_auth.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:bd0f8189-e21c-4bac-9e4f-c2a7f2d04780"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "username", "test_basic"),
// 					assertNullValue("data.kong_consumer_basic_auth.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_basic_auth" "test" {
// 	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
// 	username = "test_basic"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_basic_auth.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:bd0f8189-e21c-4bac-9e4f-c2a7f2d04780"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "username", "test_basic"),
// 					assertNullValue("data.kong_consumer_basic_auth.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_basic_auth" "test" {
// 	consumer_username = "test"
// 	username = "test_basic"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_basic_auth.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:bd0f8189-e21c-4bac-9e4f-c2a7f2d04780"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_basic_auth.test", "username", "test_basic"),
// 					assertNullValue("data.kong_consumer_basic_auth.test", "tags"),
// 				},
// 			},
// 		},
// 	})
// }

// func TestAccKongConsumerBasicAuthResource(t *testing.T) {
// 	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
// 	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_basic"
// }
// resource "kong_consumer_basic_auth" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	username = "foo"
// 	password = "bar"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_basic_auth.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_basic_auth.test", "id"),
// 					assertStringValue("kong_consumer_basic_auth.test", "username", "foo"),
// 					assertStringValue("kong_consumer_basic_auth.test", "password", "bar"),
// 					statecheck.ExpectSensitiveValue("kong_consumer_basic_auth.test", tfjsonpath.New("password")),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_basic_auth.test", "consumer_id"),
// 					assertNullValue("kong_consumer_basic_auth.test", "tags"),
// 				},
// 			},
// 			{
// 				ResourceName:            "kong_consumer_basic_auth.test",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"password"},
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_basic"
// }
// resource "kong_consumer_basic_auth" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	username = "foo1"
// 	password = "baz"
// 	tags = ["test"]
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_basic_auth.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_basic_auth.test", "id"),
// 					assertStringValue("kong_consumer_basic_auth.test", "username", "foo1"),
// 					assertStringValue("kong_consumer_basic_auth.test", "password", "baz"),
// 					statecheck.ExpectSensitiveValue("kong_consumer_basic_auth.test", tfjsonpath.New("password")),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_basic_auth.test", "consumer_id"),
// 					assertStringArrayValue("kong_consumer_basic_auth.test", "tags", []string{"test"}),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_basic"
// }
// resource "kong_consumer_basic_auth" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	username = "foo"
// 	password = "bar"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_basic_auth.test", tfjsonpath.New("id")),
// 					idReplaced.AddStateValue("kong_consumer_basic_auth.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_basic_auth.test", "id"),
// 					assertStringValue("kong_consumer_basic_auth.test", "username", "foo"),
// 					assertStringValue("kong_consumer_basic_auth.test", "password", "bar"),
// 					statecheck.ExpectSensitiveValue("kong_consumer_basic_auth.test", tfjsonpath.New("password")),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_basic_auth.test", "consumer_id"),
// 					assertNullValue("kong_consumer_basic_auth.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_basic"
// }
// resource "kong_consumer" "consumer2" {
// 	username = "test_basic2"
// }
// resource "kong_consumer_basic_auth" "test" {
// 	consumer_id = kong_consumer.consumer2.id
// 	username = "foo"
// 	password = "bar"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idReplaced.AddStateValue("kong_consumer_basic_auth.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_basic_auth.test", "id"),
// 					assertStringValue("kong_consumer_basic_auth.test", "username", "foo"),
// 					assertStringValue("kong_consumer_basic_auth.test", "password", "bar"),
// 					statecheck.ExpectSensitiveValue("kong_consumer_basic_auth.test", tfjsonpath.New("password")),
// 					assertStateMatch("kong_consumer.consumer2", "id", "kong_consumer_basic_auth.test", "consumer_id"),
// 					assertNullValue("kong_consumer_basic_auth.test", "tags"),
// 				},
// 			},
// 		},
// 	})

// }
