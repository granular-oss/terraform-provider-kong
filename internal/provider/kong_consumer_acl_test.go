package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccKongConsumerACLDataSource(t *testing.T) {
	stateChecks := []statecheck.StateCheck{
		assertStringValue("data.kong_consumer_acl.test", "id", "068141a0-93c3-4286-8329-97a95b9fe22d|test"),
		assertStringValue("data.kong_consumer_acl.test", "kong_id", "5dd417e8-c8c4-48cf-be08-e288e29649a6"),
		assertStringValue("data.kong_consumer_acl.test", "group", "test"),
		assertStringValue("data.kong_consumer_acl.test", "consumer_id", "068141a0-93c3-4286-8329-97a95b9fe22d"),
		assertNullValue("data.kong_consumer_acl.test", "tags"),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "kong_consumer_acl" "test" {
	consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
	group = "test"
}`,
				ConfigStateChecks: stateChecks,
			},
			{
				Config: providerConfig + `
data "kong_consumer_acl" "test" {
	consumer_username = "test"
	group = "test"
}`,
				ConfigStateChecks: stateChecks,
			},
			// TODO: Add error tests
		},
	})
}

func TestAccKongConsumerAclResource(t *testing.T) {
	idStaysSame := statecheck.CompareValue(compare.ValuesSame())
	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_acl"
}
resource "kong_consumer_acl" "test_acl" {
	consumer_id = kong_consumer.consumer.id
	group = "test"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_consumer_acl.test_acl", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_acl.test_acl", "id"),
					assertStringValue("kong_consumer_acl.test_acl", "group", "test"),
					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_acl.test_acl", "consumer_id"),
					assertNullValue("kong_consumer_acl.test_acl", "tags"),
				},
			},
			{
				ResourceName:      "kong_consumer_acl.test_acl",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "kong_consumer_acl.test_acl",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test_acl|test",
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_acl"
}
resource "kong_consumer_acl" "test_acl" {
	consumer_id = kong_consumer.consumer.id
	group = "test"
	tags = ["tag"]
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_consumer_acl.test_acl", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_acl.test_acl", "id"),
					assertStringValue("kong_consumer_acl.test_acl", "group", "test"),
					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_acl.test_acl", "consumer_id"),
					assertStringArrayValue("kong_consumer_acl.test_acl", "tags", []string{"tag"}),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_acl"
}
resource "kong_consumer_acl" "test_acl" {
	consumer_id = kong_consumer.consumer.id
	group = "test"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idStaysSame.AddStateValue("kong_consumer_acl.test_acl", tfjsonpath.New("id")),
					idReplaced.AddStateValue("kong_consumer_acl.test_acl", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_acl.test_acl", "id"),
					assertStringValue("kong_consumer_acl.test_acl", "group", "test"),
					assertStateMatch("kong_consumer.consumer", "id", "kong_consumer_acl.test_acl", "consumer_id"),
					assertNullValue("kong_consumer_acl.test_acl", "tags"),
				},
			},
			{
				Config: providerConfig + `
resource "kong_consumer" "consumer" {
	username = "test_acl"
}
resource "kong_consumer" "consumer2" {
	username = "test_acl2"
}
resource "kong_consumer_acl" "test_acl" {
	consumer_id = kong_consumer.consumer2.id
	group = "test"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue("kong_consumer_acl.test_acl", tfjsonpath.New("id")),
					assertNotNull("kong_consumer_acl.test_acl", "id"),
					assertStringValue("kong_consumer_acl.test_acl", "group", "test"),
					assertStateMatch("kong_consumer.consumer2", "id", "kong_consumer_acl.test_acl", "consumer_id"),
					assertNullValue("kong_consumer_acl.test_acl", "tags"),
				},
			},
		},
	})

}
