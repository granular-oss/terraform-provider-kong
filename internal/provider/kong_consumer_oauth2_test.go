package provider

// import (
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/compare"
// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-testing/statecheck"
// 	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
// )

// func TestAccKongConsumerOauth2DataSource(t *testing.T) {
// 	datasourceChecks := []statecheck.StateCheck{
// 		assertStringValue("data.kong_consumer_oauth2.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:1428d5e5-32e0-482a-bac0-c6a4c6de6504"),
// 		assertStringValue("data.kong_consumer_oauth2.test", "client_id", "1673b842-f396-4c9e-83df-f5fe05fcef1d"),
// 		assertStringValue("data.kong_consumer_oauth2.test", "client_secret", "super_secret"),
// 		assertStringValue("data.kong_consumer_oauth2.test", "name", "test-oauth"),
// 		assertBoolValue("data.kong_consumer_oauth2.test", "hash_secret", false),
// 		assertStringValue("data.kong_consumer_oauth2.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 		assertNullValue("data.kong_consumer_oauth2.test", "redirect_uris"),
// 		assertNullValue("data.kong_consumer_oauth2.test", "tags"),
// 	}
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_oauth2" "test" {
// 	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
// 	id = "1428d5e5-32e0-482a-bac0-c6a4c6de6504"
// }`,
// 				ConfigStateChecks: datasourceChecks,
// 			},
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_oauth2" "test" {
// 	consumer_username = "test"
// 	id = "1428d5e5-32e0-482a-bac0-c6a4c6de6504"
// }`,
// 				ConfigStateChecks: datasourceChecks,
// 			},
// 		},
// 	})
// }

// func TestAccKongConsumerOauth2Resource(t *testing.T) {
// 	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
// 	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_oauth"
// }
// resource "kong_consumer_oauth2" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	name = "testing_oauth"
// 	client_id = "e85edd74-67cc-41db-85ed-b0725a1de285"
// 	client_secret = "test_secret"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_oauth2.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_oauth2.test", "id"),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_oauth2.test", "consumer_id"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_id", "e85edd74-67cc-41db-85ed-b0725a1de285"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_secret", "test_secret"),
// 					assertStringValue("kong_consumer_oauth2.test", "name", "testing_oauth"),
// 					assertBoolValue("kong_consumer_oauth2.test", "hash_secret", false),
// 					assertNullValue("kong_consumer_oauth2.test", "redirect_uris"),
// 					assertNullValue("kong_consumer_oauth2.test", "tags"),
// 				},
// 			},
// 			{
// 				ResourceName:      "kong_consumer_oauth2.test",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_key_auth"
// }
// resource "kong_consumer_oauth2" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	name = "testing_oauth_again"
// 	client_id = "e917cddd-b6e4-4579-a0b1-bdf49110645f"
// 	client_secret = "test_secret2"
// 	hash_secret = true
// 	redirect_uris = ["http://testing.com"]
// 	tags = ["tag1"]
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_oauth2.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_oauth2.test", "id"),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_oauth2.test", "consumer_id"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_id", "e917cddd-b6e4-4579-a0b1-bdf49110645f"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_secret", "test_secret2"),
// 					assertStringValue("kong_consumer_oauth2.test", "name", "testing_oauth_again"),
// 					assertBoolValue("kong_consumer_oauth2.test", "hash_secret", true),
// 					assertStringArrayValue("kong_consumer_oauth2.test", "redirect_uris", []string{"http://testing.com"}),
// 					assertStringArrayValue("kong_consumer_oauth2.test", "tags", []string{"tag1"}),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_oauth"
// }
// resource "kong_consumer_oauth2" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	name = "testing_oauth"
// 	client_id = "e85edd74-67cc-41db-85ed-b0725a1de285"
// 	client_secret = "test_secret"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_oauth2.test", tfjsonpath.New("id")),
// 					idReplaced.AddStateValue("kong_consumer_oauth2.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_oauth2.test", "id"),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_oauth2.test", "consumer_id"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_id", "e85edd74-67cc-41db-85ed-b0725a1de285"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_secret", "test_secret"),
// 					assertStringValue("kong_consumer_oauth2.test", "name", "testing_oauth"),
// 					assertBoolValue("kong_consumer_oauth2.test", "hash_secret", false),
// 					assertNullValue("kong_consumer_oauth2.test", "redirect_uris"),
// 					assertNullValue("kong_consumer_oauth2.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_oauth"
// }
// resource "kong_consumer" "consumer2" {
// 	username = "test_oauth_2"
// }
// resource "kong_consumer_oauth2" "test" {
// 	consumer_id = kong_consumer.consumer2.id
// 	name = "testing_oauth"
// 	client_id = "e85edd74-67cc-41db-85ed-b0725a1de285"
// 	client_secret = "test_secret"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idReplaced.AddStateValue("kong_consumer_oauth2.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_oauth2.test", "id"),
// 					assertStateMatch("kong_consumer.consumer2", "id", "kong_consumer_oauth2.test", "consumer_id"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_id", "e85edd74-67cc-41db-85ed-b0725a1de285"),
// 					assertStringValue("kong_consumer_oauth2.test", "client_secret", "test_secret"),
// 					assertStringValue("kong_consumer_oauth2.test", "name", "testing_oauth"),
// 					assertBoolValue("kong_consumer_oauth2.test", "hash_secret", false),
// 					assertNullValue("kong_consumer_oauth2.test", "redirect_uris"),
// 					assertNullValue("kong_consumer_oauth2.test", "tags"),
// 				},
// 			},
// 		},
// 	})

// }
