terraform {
  required_providers {
    kong = {
      source = "granular-oss/kong"
    }
  }
}

provider "kong" {}

resource "kong_service" "minimal" {
  name = "test2"
  host = "fake.example.com"
}

resource "kong_service" "complete" {
  name                  = "test2"
  host                  = "fake.example.com"
  protocol              = "https"
  read_timeout          = 30000
  write_timeout         = 30000
  connect_timeout       = 30000
  retries               = 10
  port                  = 443
  ca_certificate_ids    = ["8bb0f886-209f-488a-ace5-1cdd535fbb40"]
  client_certificate_id = "d0e79808-9136-4e8f-af96-e1740548c03c"
  path                  = "/example"
  tags                  = ["test1", "test:2"]
  tls_depth             = 2
  tls_verify            = false
}
