package provider

// import (
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/compare"
// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-testing/statecheck"
// 	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
// )

// func TestAccKongConsumerDataSource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// data "kong_consumer" "test_id" {
// 	id = "068141a0-93c3-4286-8329-97a95b9fe22d"
// }
// data "kong_consumer" "test_username" {
// 	username = "test"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer.test_id", "id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer.test_id", "username", "test"),
// 					assertNullValue("data.kong_consumer.test_id", "custom_id"),
// 					assertNullValue("data.kong_consumer.test_id", "tags"),
// 					assertStringValue("data.kong_consumer.test_username", "id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer.test_username", "username", "test"),
// 					assertNullValue("data.kong_consumer.test_username", "custom_id"),
// 					assertNullValue("data.kong_consumer.test_username", "tags"),
// 				},
// 			},
// 		},
// 	})
// }

// func TestAccKongConsumerResource(t *testing.T) {
// 	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "test" {
// 	username = "foobar"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue(
// 						"kong_consumer.test", tfjsonpath.New("id"),
// 					),
// 					assertNotNull("kong_consumer.test", "id"),
// 					assertStringValue("kong_consumer.test", "username", "foobar"),
// 					assertNullValue("kong_consumer.test", "custom_id"),
// 					assertNullValue("kong_consumer.test", "tags"),
// 				},
// 			},
// 			{
// 				ResourceName:      "kong_consumer.test",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				ResourceName:      "kong_consumer.test",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 				ImportStateId:     "foobar",
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "test" {
// 	username = "foobaz"
// 	custom_id = "test_32"
// 	tags = ["fake"]
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue(
// 						"kong_consumer.test", tfjsonpath.New("id"),
// 					),
// 					assertStringValue("kong_consumer.test", "username", "foobaz"),
// 					assertStringValue("kong_consumer.test", "custom_id", "test_32"),
// 					assertStringArrayValue("kong_consumer.test", "tags", []string{"fake"}),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "test" {
// 	username = "foobar"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue(
// 						"kong_consumer.test", tfjsonpath.New("id"),
// 					),
// 					assertNotNull("kong_consumer.test", "id"),
// 					assertStringValue("kong_consumer.test", "username", "foobar"),
// 					assertNullValue("kong_consumer.test", "custom_id"),
// 					assertNullValue("kong_consumer.test", "tags"),
// 				},
// 			},
// 		},
// 	})

// }
