provider "kong" {}

data "kong_consumer_key_auth" "consumer_id_id" {
  id          = "068141a0-93c3-4286-8329-97a95b9fe22d"
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
}

data "kong_consumer_key_auth" "consumer_user_id" {
  id                = "068141a0-93c3-4286-8329-97a95b9fe22d"
  consumer_username = "example"
}

data "kong_consumer_basic_auth" "consumer_id_key" {
  key         = "user"
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
}

data "kong_consumer_basic_auth" "consumer_user_key" {
  key               = "user"
  consumer_username = "example"
}
