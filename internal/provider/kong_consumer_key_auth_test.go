package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccKongConsumerKeyAuthDataSource(t *testing.T) {
	expectedState := []statecheck.StateCheck{

		assertStringValue("data.kong_consumer_key_auth.test", "id", "58babb72-d2b2-4b9d-a80c-148dc8ead8f4"),
		assertStringValue("data.kong_consumer_key_auth.test", "key", "super_secret"),
		assertStringValue("data.kong_consumer_key_auth.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
		assertNullValue("data.kong_consumer_key_auth.test", "tags"),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "kong_consumer_key_auth" "test" {
	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
	id = "58babb72-d2b2-4b9d-a80c-148dc8ead8f4"
}`,
				ConfigStateChecks: expectedState,
			},
			{
				Config: providerConfig + `
data "kong_consumer_key_auth" "test" {
	consumer_username = "test"
	id = "58babb72-d2b2-4b9d-a80c-148dc8ead8f4"
}`,
				ConfigStateChecks: expectedState,
			},
			{
				Config: providerConfig + `
data "kong_consumer_key_auth" "test" {
	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
	key = "super_secret"
}`,
				ConfigStateChecks: expectedState,
			},
			{
				Config: providerConfig + `
data "kong_consumer_key_auth" "test" {
	consumer_username = "test"
	key = "super_secret"
}`,
				ConfigStateChecks: expectedState,
			},
		},
	})
}

func TestAccKongConsumerKeyAuthResource(t *testing.T) {
	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_key_auth"
}
resource "kong_consumer_key_auth" "test" {
	consumer_id = kong_consumer.consumer.id
	key = "test_key"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_consumer_key_auth.test", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_key_auth.test", "id"),
					assertStringValue("kong_consumer_key_auth.test", "key", "test_key"),
					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_key_auth.test", "consumer_id"),
					assertNullValue("kong_consumer_key_auth.test", "tags"),
				},
			},
			{
				ResourceName:      "kong_consumer_key_auth.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_key_auth"
}
resource "kong_consumer_key_auth" "test" {
	consumer_id = kong_consumer.consumer.id
	key = "test_key"
	tags = ["tag"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_consumer_key_auth.test", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_key_auth.test", "id"),
					assertStringValue("kong_consumer_key_auth.test", "key", "test_key"),
					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_key_auth.test", "consumer_id"),
					assertStringArrayValue("kong_consumer_key_auth.test", "tags", []string{"tag"}),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_key_auth"
}
resource "kong_consumer_key_auth" "test" {
	consumer_id = kong_consumer.consumer.id
	key = "test_key"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_consumer_key_auth.test", tfjsonpath.New("id")),
					idReplaced.AddStateValue("kong_consumer_key_auth.test", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_key_auth.test", "id"),
					assertStringValue("kong_consumer_key_auth.test", "key", "test_key"),
					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_key_auth.test", "consumer_id"),
					assertNullValue("kong_consumer_key_auth.test", "tags"),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_key_auth"
}
resource "kong_consumer" "consumer2" {
	username = "test_key_auth2"
}
resource "kong_consumer_key_auth" "test" {
	consumer_id = kong_consumer.consumer2.id
	key = "test_key"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue("kong_consumer_key_auth.test", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_key_auth.test", "id"),
					assertStringValue("kong_consumer_key_auth.test", "key", "test_key"),
					assertStateMatch("kong_consumer.consumer2", "id", "kong_consumer_key_auth.test", "consumer_id"),
					assertNullValue("kong_consumer_key_auth.test", "tags"),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_key_auth"
}
resource "kong_consumer" "consumer2" {
	username = "test_key_auth2"
}
resource "kong_consumer_key_auth" "test" {
	consumer_id = kong_consumer.consumer2.id
	key = "test_key2"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue("kong_consumer_key_auth.test", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_key_auth.test", "id"),
					assertStringValue("kong_consumer_key_auth.test", "key", "test_key2"),
					assertStateMatch("kong_consumer.consumer2", "id", "kong_consumer_key_auth.test", "consumer_id"),
					assertNullValue("kong_consumer_key_auth.test", "tags"),
				},
			},
		},
	})

}
