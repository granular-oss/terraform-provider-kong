package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

func TestAccKongUpstreamDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "kong_upstream" "test" {
	name = "test-upstream"
}`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertStringValue("data.kong_upstream.test", "id", "5133ba73-4ca3-42b2-9b31-87cfd13951b3"),
					assertStringValue("data.kong_upstream.test", "name", "test-upstream"),
					assertInt32Value("data.kong_upstream.test", "slots", 10000),
					assertStringValue("data.kong_upstream.test", "hash_on", "none"),
					assertInt32Value("data.kong_upstream.test", "healthchecks.active.unhealthy.interval", 0),
					assertInt32Value("data.kong_upstream.test", "healthchecks.passive.healthy.successes", 0),
				},
			},
		},
	})
}

func TestAccKongUpstreamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_upstream" "test" {
  name = "fake_upstream_for_test"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertNotNull("kong_upstream.test", "id"),
					assertStringValue("kong_upstream.test", "name", "fake_upstream_for_test"),
				},
			},
			{
				ResourceName:      "kong_upstream.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
