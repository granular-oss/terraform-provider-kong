package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

func TestAccKongServiceDataSource(t *testing.T) {
	stateChecks := []statecheck.StateCheck{
		assertStringValue("data.kong_service.test", "id", "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"),
		assertStringValue("data.kong_service.test", "name", "test"),
		assertStringValue("data.kong_service.test", "host", "example.com"),
		assertStringValue("data.kong_service.test", "protocol", "http"),
		assertInt32Value("data.kong_service.test", "read_timeout", 60000),
		assertInt32Value("data.kong_service.test", "write_timeout", 60000),
		assertInt32Value("data.kong_service.test", "connect_timeout", 60000),
		assertInt32Value("data.kong_service.test", "retries", 5),
		assertInt32Value("data.kong_service.test", "port", 80),
		assertNullValue("data.kong_service.test", "ca_certificate_ids"),
		assertNullValue("data.kong_service.test", "client_certificate_id"),
		assertNullValue("data.kong_service.test", "path"),
		assertNullValue("data.kong_service.test", "tags"),
		assertNullValue("data.kong_service.test", "tls_depth"),
		assertNullValue("data.kong_service.test", "tls_verify"),
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            providerConfig + `data "kong_service" "test" { name = "test" }`,
				ConfigStateChecks: stateChecks,
			},
			{
				Config:            providerConfig + `data "kong_service" "test" { id = "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7" }`,
				ConfigStateChecks: stateChecks,
			},
		},
	})
}

func TestAccKongServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "kong_service" "test" {
  name = "test2"
  host = "example.com"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertNotNull("kong_service.test", "id"),
					assertStringValue("kong_service.test", "name", "test2"),
					assertStringValue("kong_service.test", "host", "example.com"),
					assertStringValue("kong_service.test", "protocol", "http"),
					assertInt32Value("kong_service.test", "read_timeout", 60000),
					assertInt32Value("kong_service.test", "write_timeout", 60000),
					assertInt32Value("kong_service.test", "connect_timeout", 60000),
					assertInt32Value("kong_service.test", "retries", 5),
					assertInt32Value("kong_service.test", "port", 80),
					assertNullValue("kong_service.test", "ca_certificate_ids"),
					assertNullValue("kong_service.test", "client_certificate_id"),
					assertNullValue("kong_service.test", "path"),
					assertNullValue("kong_service.test", "tags"),
					assertNullValue("kong_service.test", "tls_depth"),
					assertNullValue("kong_service.test", "tls_verify"),
				},
			},
			{
				ResourceName:      "kong_service.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
resource "kong_service" "test" {
  name = "test2"
  host = "fake.example.com"
  protocol = "https"
  read_timeout = 30000
  write_timeout = 30000
  connect_timeout = 30000
  retries = 10
  port = 443
  ca_certificate_ids = ["a0518569-74ea-4548-ba13-f852f7517191"]
  client_certificate_id = "56c71339-eda3-4c34-9aff-22ac90539119"
  path = "/example"
  tags = ["test1","test:2"]
  tls_depth = 2
  tls_verify = false
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertNotNull("kong_service.test", "id"),
					assertStringValue("kong_service.test", "name", "test2"),
					assertStringValue("kong_service.test", "host", "fake.example.com"),
					assertStringValue("kong_service.test", "protocol", "https"),
					assertInt32Value("kong_service.test", "read_timeout", 30000),
					assertInt32Value("kong_service.test", "write_timeout", 30000),
					assertInt32Value("kong_service.test", "connect_timeout", 30000),
					assertInt32Value("kong_service.test", "retries", 10),
					assertInt32Value("kong_service.test", "port", 443),
					assertStringArrayValue("kong_service.test", "ca_certificate_ids", []string{"a0518569-74ea-4548-ba13-f852f7517191"}),
					assertStringValue("kong_service.test", "client_certificate_id", "56c71339-eda3-4c34-9aff-22ac90539119"),
					assertStringValue("kong_service.test", "path", "/example"),
					assertStringArrayValue("kong_service.test", "tags", []string{"test1", "test:2"}),
					assertInt32Value("kong_service.test", "tls_depth", 2),
					assertBoolValue("kong_service.test", "tls_verify", false),
				},
			},
			// Verify that properties can be set back to null
			{
				Config: providerConfig + `
resource "kong_service" "test" {
  name = "test2"
  host = "example.com"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					assertNotNull("kong_service.test", "id"),
					assertStringValue("kong_service.test", "name", "test2"),
					assertStringValue("kong_service.test", "host", "example.com"),
					assertStringValue("kong_service.test", "protocol", "http"),
					assertInt32Value("kong_service.test", "read_timeout", 60000),
					assertInt32Value("kong_service.test", "write_timeout", 60000),
					assertInt32Value("kong_service.test", "connect_timeout", 60000),
					assertInt32Value("kong_service.test", "retries", 5),
					assertInt32Value("kong_service.test", "port", 80),
					assertNullValue("kong_service.test", "ca_certificate_ids"),
					assertNullValue("kong_service.test", "client_certificate_id"),
					assertNullValue("kong_service.test", "path"),
					assertNullValue("kong_service.test", "tags"),
					assertNullValue("kong_service.test", "tls_depth"),
					assertNullValue("kong_service.test", "tls_verify"),
				},
			},
		},
	})
}
