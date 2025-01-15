provider "kong" {}

data "kong_upstream" "example" {
  name = "upstream-name"
}

data "kong_upstream" "example_id" {
  id = "50c86b96-b973-4c8a-933f-7f48f6f49896"
}
