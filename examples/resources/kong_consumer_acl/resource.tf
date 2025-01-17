terraform {
  required_providers {
    kong = {
      source = "granular-oss/kong"
    }
  }
}

provider "kong" {}

resource "kong_consumer" "consumer" {
  username = "foobar"
}

resource "kong_consumer_acl" "example" {
  consumer_id = kong_consumer.consumer.id
  group       = "test"
  tags        = ["test"]
}
