provider "kong" {}

data "kong_consumer" "test" {
  name = "test"
}

data "kong_consumer" "test_id" {
  id = "50c86b96-b973-4c8a-933f-7f48f6f49896"
}
