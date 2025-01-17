terraform {
  required_providers {
    kong = {
      source = "granular-oss/kong"
    }
  }
}

provider "kong" {}

resource "kong_upstream" "upstream" {
  name = "example"
}

resource "kong_target" "example_min" {
  target      = "example.com:8000"
  upstream_id = kong_upstream.upstream.id
}

resource "kong_target" "example_complete" {
  target      = "example.com:8000"
  upstream_id = kong_upstream.upstream.id
  weight      = 70
  tags        = ["tag"]
}
