package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccKongRouteDataSource(t *testing.T) {
	expectedState := []statecheck.StateCheck{
		assertStringValue("data.kong_route.test", "id", "772e8b14-b020-4164-b52d-3643f6310821"),
		assertStringValue("data.kong_route.test", "name", "test"),
		assertStringArrayValue("data.kong_route.test", "protocols", []string{"http", "https"}),
		assertNullValue("data.kong_route.test", "methods"),
		assertNullValue("data.kong_route.test", "hosts"),
		assertStringArrayValue("data.kong_route.test", "paths", []string{"/foobar"}),
		assertBoolValue("data.kong_route.test", "strip_path", true),
		assertNullValue("data.kong_route.test", "source"),
		assertNullValue("data.kong_route.test", "destination"),
		assertNullValue("data.kong_route.test", "snis"),
		assertBoolValue("data.kong_route.test", "preserve_host", false),
		assertInt32Value("data.kong_route.test", "regex_priority", 0),
		assertStringValue("data.kong_route.test", "service_id", "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"),
		assertStringValue("data.kong_route.test", "path_handling", "v0"),
		assertInt32Value("data.kong_route.test", "https_redirect_status_code", 426),
		assertBoolValue("data.kong_route.test", "request_buffering", true),
		assertBoolValue("data.kong_route.test", "response_buffering", true),
		assertNullValue("data.kong_route.test", "tags"),
		assertNullValue("data.kong_route.test", "header"),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "kong_route" "test" {
	name = "test"
}`,
				ConfigStateChecks: expectedState,
			},
			{
				Config: providerConfig + `
data "kong_route" "test" {
	id = "772e8b14-b020-4164-b52d-3643f6310821"
}`,
				ConfigStateChecks: expectedState,
			},
		},
	})
}

func TestAccKongRouteResource(t *testing.T) {
	expectedState := []statecheck.StateCheck{
		assertNotNull("kong_route.test", "id"),
		assertStringValue("kong_route.test", "name", "test2"),
		assertStringArrayValue("kong_route.test", "protocols", []string{"http", "https"}),
		assertNullValue("kong_route.test", "methods"),
		assertNullValue("kong_route.test", "hosts"),
		assertStringArrayValue("kong_route.test", "paths", []string{"/foobar"}),
		assertBoolValue("kong_route.test", "strip_path", true),
		assertNullValue("kong_route.test", "source"),
		assertNullValue("kong_route.test", "destination"),
		assertNullValue("kong_route.test", "snis"),
		assertBoolValue("kong_route.test", "preserve_host", false),
		assertInt32Value("kong_route.test", "regex_priority", 0),
		assertStateMatch("kong_route.test", "service_id", "kong_service.test", "id"),
		assertStringValue("kong_route.test", "path_handling", "v0"),
		assertInt32Value("kong_route.test", "https_redirect_status_code", 426),
		assertBoolValue("kong_route.test", "request_buffering", true),
		assertBoolValue("kong_route.test", "response_buffering", true),
		assertNullValue("kong_route.test", "tags"),
		statecheck.ExpectKnownValue("kong_route.test", buildJsonPath("header"), knownvalue.ListExact([]knownvalue.Check{})),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_service" "test" {
  name = "test2"
  host = "example.com"
}

resource "kong_route" "test" {
  name  = "test2"
  paths = ["/foobar"]
  service_id = kong_service.test.id
}
`,
				ConfigStateChecks: expectedState,
			},
			{
				ResourceName:      "kong_route.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
resource "kong_service" "test" {
  name = "test2"
  host = "example.com"
}

resource "kong_route" "test" {
  name       = "test2"
  paths      = ["/foobaz"]
  service_id = kong_service.test.id
  protocols  = ["http"]
  methods    = ["GET"]
  hosts      = ["test.example.com"]
  tags       = ["test", "app:fake"]
  strip_path = false
  preserve_host              = true
  regex_priority             = 3
  path_handling              = "v1"
  https_redirect_status_code = 301
  request_buffering          = false
  response_buffering         = false

  header {
    name  = "header-test"
    values = ["1", "2"]
  }

  header {
    name  = "header-test2"
    values = ["1", "2"]
  }
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertNotNull("kong_route.test", "id"),
					assertStringValue("kong_route.test", "name", "test2"),
					assertStringArrayValue("kong_route.test", "protocols", []string{"http"}),
					assertStringArrayValue("kong_route.test", "methods", []string{"GET"}),
					assertStringArrayValue("kong_route.test", "hosts", []string{"test.example.com"}),
					assertStringArrayValue("kong_route.test", "paths", []string{"/foobaz"}),
					assertBoolValue("kong_route.test", "strip_path", false),
					assertNullValue("kong_route.test", "source"),
					assertNullValue("kong_route.test", "destination"),
					assertNullValue("kong_route.test", "snis"),
					assertBoolValue("kong_route.test", "preserve_host", true),
					assertInt32Value("kong_route.test", "regex_priority", 3),
					assertStateMatch("kong_route.test", "service_id", "kong_service.test", "id"),
					assertStringValue("kong_route.test", "path_handling", "v1"),
					assertInt32Value("kong_route.test", "https_redirect_status_code", 301),
					assertBoolValue("kong_route.test", "request_buffering", false),
					assertBoolValue("kong_route.test", "response_buffering", false),
					assertStringArrayValue("kong_route.test", "tags", []string{"test", "app:fake"}),
					statecheck.ExpectKnownValue("kong_route.test", tfjsonpath.New("header"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("header-test"),
							"values": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("1"),
								knownvalue.StringExact("2"),
							}),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("header-test2"),
							"values": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("1"),
								knownvalue.StringExact("2"),
							}),
						}),
					})),
				},
			},
			{
				Config: providerConfig + `
resource "kong_service" "test" {
  name = "test2"
  host = "example.com"
}

resource "kong_route" "test" {
  name  = "test2"
  paths = ["/foobar"]
  service_id = kong_service.test.id
}
`,
				ConfigStateChecks: expectedState,
			},
			{
				Config: providerConfig + `
resource "kong_service" "test" {
  name = "test2"
  host = "example.com"
}

resource "kong_route" "test" {
  name  = "test2"
  paths = ["/foobar"]
  service_id = kong_service.test.id
  hosts      = []
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertStringArrayValue("kong_route.test", "hosts", []string{}),
				},
			},
		},
	})
}
