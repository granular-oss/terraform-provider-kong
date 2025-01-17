package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccKongTargetDataSource(t *testing.T) {
	checks := []statecheck.StateCheck{
		assertStringValue("data.kong_target.test", "id", "921687b2-ac57-4fa5-a99f-ae08fbf30eff"),
		assertStringValue("data.kong_target.test", "target", "test.com:8000"),
		assertStringValue("data.kong_target.test", "upstream_id", "5133ba73-4ca3-42b2-9b31-87cfd13951b3"),
		assertInt32Value("data.kong_target.test", "weight", 100),
		assertNullValue("data.kong_target.test", "tags"),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "kong_target" "test" {
 	id = "921687b2-ac57-4fa5-a99f-ae08fbf30eff"
	upstream_name = "test-upstream"
}`,
				ConfigStateChecks: checks,
			},
			{
				Config: providerConfig + `
data "kong_target" "test" {
 	id = "921687b2-ac57-4fa5-a99f-ae08fbf30eff"
	upstream_id = "5133ba73-4ca3-42b2-9b31-87cfd13951b3"
}`,
				ConfigStateChecks: checks,
			},
			{
				Config: providerConfig + `
data "kong_target" "test" {
 	target = "test.com:8000"
	upstream_name = "test-upstream"
}`,
				ConfigStateChecks: checks,
			},
			{
				Config: providerConfig + `
data "kong_target" "test" {
 	target = "test.com:8000"
	upstream_id = "5133ba73-4ca3-42b2-9b31-87cfd13951b3"
}`,
				ConfigStateChecks: checks,
			},
		},
	})
}

func TestAccKongTargetResource(t *testing.T) {
	idReplaced := statecheck.CompareValue(compare.ValuesDiffer())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_upstream" "upstream" {
	name = "upstream-target-test"
}
resource "kong_target" "test" {
	target = "foo.bar:8000"
	upstream_id = kong_upstream.upstream.id
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue("kong_target.test", tfjsonpath.New("id")),
					assertNotNull("kong_target.test", "id"),
					assertInt32Value("kong_target.test", "weight", 100),
					assertStringValue("kong_target.test", "target", "foo.bar:8000"),
					assertNullValue("kong_target.test", "tags"),
				},
			},
			{
				ResourceName:      "kong_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "kong_target.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "upstream-target-test|foo.bar:8000",
			},
			{
				Config: providerConfig + `
resource "kong_upstream" "upstream" {
	name = "upstream-target-test"
}
resource "kong_target" "test" {
	target = "foo.bar:8000"
	upstream_id = kong_upstream.upstream.id
	weight = 30
	tags = ["foo"]
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue("kong_target.test", tfjsonpath.New("id")),
					assertNotNull("kong_target.test", "id"),
					assertInt32Value("kong_target.test", "weight", 30),
					assertStringValue("kong_target.test", "target", "foo.bar:8000"),
					assertStringArrayValue("kong_target.test", "tags", []string{"foo"}),
				},
			},
			{
				Config: providerConfig + `
resource "kong_upstream" "upstream" {
	name = "upstream-target-test"
}
resource "kong_target" "test" {
	target = "foo.bar:8000"
	upstream_id = kong_upstream.upstream.id
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue("kong_target.test", tfjsonpath.New("id")),
					assertNotNull("kong_target.test", "id"),
					assertInt32Value("kong_target.test", "weight", 100),
					assertStringValue("kong_target.test", "target", "foo.bar:8000"),
					assertNullValue("kong_target.test", "tags"),
				},
			},
			{
				Config: providerConfig + `
resource "kong_upstream" "upstream" {
	name = "upstream-target-test"
}
resource "kong_target" "test" {
	target = "foo.bar:8080"
	upstream_id = kong_upstream.upstream.id
}
			`,
				ConfigStateChecks: []statecheck.StateCheck{
					idReplaced.AddStateValue("kong_target.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}
