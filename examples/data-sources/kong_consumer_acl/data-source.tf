provider "kong" {}

data "kong_consumer" "test" {
  name = "test"
}

data "kong_consumer_acl" "test_id" {
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
  group       = "test"
}

data "kong_consumer_acl" "test_username" {
  consumer_username = "test"
  group             = "test"
}
