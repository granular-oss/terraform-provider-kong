package provider

// import (
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/compare"
// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-testing/statecheck"
// 	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
// )

// func TestAccKongConsumerJwtDataSource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_jwt" "test" {
// 	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
// 	id = "da8b0cd6-f1d3-4731-ba20-873254d9d474"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_jwt.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:da8b0cd6-f1d3-4731-ba20-873254d9d474"),
// 					assertStringValue("data.kong_consumer_jwt.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_jwt.test", "key", "fake_key"),
// 					assertStringValue("data.kong_consumer_jwt.test", "secret", "secret"),
// 					assertStringValue("data.kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("data.kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("data.kong_consumer_jwt.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_jwt" "test" {
// 	consumer_username = "test"
// 	id = "da8b0cd6-f1d3-4731-ba20-873254d9d474"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_jwt.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:da8b0cd6-f1d3-4731-ba20-873254d9d474"),
// 					assertStringValue("data.kong_consumer_jwt.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_jwt.test", "key", "fake_key"),
// 					assertStringValue("data.kong_consumer_jwt.test", "secret", "secret"),
// 					assertStringValue("data.kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("data.kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("data.kong_consumer_jwt.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_jwt" "test" {
// 	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
// 	key = "fake_key"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_jwt.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:da8b0cd6-f1d3-4731-ba20-873254d9d474"),
// 					assertStringValue("data.kong_consumer_jwt.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_jwt.test", "key", "fake_key"),
// 					assertStringValue("data.kong_consumer_jwt.test", "secret", "secret"),
// 					assertStringValue("data.kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("data.kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("data.kong_consumer_jwt.test", "tags"),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// data "kong_consumer_jwt" "test" {
// 	consumer_username = "test"
// 	key = "fake_key"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					assertStringValue("data.kong_consumer_jwt.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d:da8b0cd6-f1d3-4731-ba20-873254d9d474"),
// 					assertStringValue("data.kong_consumer_jwt.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
// 					assertStringValue("data.kong_consumer_jwt.test", "key", "fake_key"),
// 					assertStringValue("data.kong_consumer_jwt.test", "secret", "secret"),
// 					assertStringValue("data.kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("data.kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("data.kong_consumer_jwt.test", "tags"),
// 				},
// 			},
// 		},
// 	})
// }

// func TestAccKongConsumerJwtResource(t *testing.T) {
// 	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
// 	keyStaysSame := statecheck.CompareValue(compare.ValuesSame())
// 	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
// 	keyReplaced := statecheck.CompareValue(compare.ValuesDiffer())
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_jwt"
// }
// resource "kong_consumer_jwt" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	secret = "secret"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_jwt.test", "id"),
// 					keyStaysSame.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("key")),
// 					assertNotNull("kong_consumer_jwt.test", "key"),
// 					assertStringValue("kong_consumer_jwt.test", "secret", "secret"),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_jwt.test", "consumer_id"),
// 					assertStringValue("kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("kong_consumer_jwt.test", "tags"),
// 				},
// 			},
// 			{
// 				ResourceName:      "kong_consumer_jwt.test",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_jwt"
// }
// resource "kong_consumer_jwt" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	secret = "secret"
// 	rsa_public_key = "public_key"
// 	tags = ["test"]
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_jwt.test", "id"),
// 					keyStaysSame.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("key")),
// 					assertNotNull("kong_consumer_jwt.test", "key"),
// 					assertStringValue("kong_consumer_jwt.test", "secret", "secret"),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_jwt.test", "consumer_id"),
// 					assertStringValue("kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertStringValue("kong_consumer_jwt.test", "rsa_public_key", "public_key"),
// 					assertStringArrayValue("kong_consumer_jwt.test", "tags", []string{"test"}),
// 				},
// 			},
// 			{
// 				Config: providerConfig + `
// resource "kong_consumer" "consumer" {
// 	username = "test_jwt"
// }
// resource "kong_consumer_jwt" "test" {
// 	consumer_id = kong_consumer.consumer.id
// 	secret = "secret"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idStaysSame.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("id")),
// 					idReplaced.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_jwt.test", "id"),
// 					keyStaysSame.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("key")),
// 					keyReplaced.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("key")),
// 					assertNotNull("kong_consumer_jwt.test", "key"),
// 					assertStringValue("kong_consumer_jwt.test", "secret", "secret"),
// 					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_jwt.test", "consumer_id"),
// 					assertStringValue("kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("kong_consumer_jwt.test", "tags"),
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
// resource "kong_consumer_jwt" "test" {
// 	consumer_id = kong_consumer.consumer2.id
// 	secret = "secret"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idReplaced.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_jwt.test", "id"),
// 					keyReplaced.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("key")),
// 					assertNotNull("kong_consumer_jwt.test", "key"),
// 					assertStringValue("kong_consumer_jwt.test", "secret", "secret"),
// 					assertStateMatch("kong_consumer.consumer2", "id", "kong_consumer_jwt.test", "consumer_id"),
// 					assertStringValue("kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("kong_consumer_jwt.test", "tags"),
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
// resource "kong_consumer_jwt" "test" {
// 	key = "key"
// 	consumer_id = kong_consumer.consumer2.id
// 	secret = "secret"
// }`,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					idReplaced.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("id")),
// 					assertNotNull("kong_consumer_jwt.test", "id"),
// 					keyReplaced.AddStateValue("kong_consumer_jwt.test", tfjsonpath.New("key")),
// 					assertNotNull("kong_consumer_jwt.test", "key"),
// 					assertStringValue("kong_consumer_jwt.test", "secret", "secret"),
// 					assertStateMatch("kong_consumer.consumer2", "id", "kong_consumer_jwt.test", "consumer_id"),
// 					assertStringValue("kong_consumer_jwt.test", "algorithm", "HS256"),
// 					assertNullValue("kong_consumer_jwt.test", "rsa_public_key"),
// 					assertNullValue("kong_consumer_jwt.test", "tags"),
// 				},
// 			},
// 		},
// 	})

// }
