provider "kong" {}

resource "kong_consumer" "example" {
  username  = "foobar"
  custom_id = "fake"
  tags      = ["tag"]
}
