provider "kong" {}

data "kong_consumer_oauth2" "consumer_id_id" {
  id          = "068141a0-93c3-4286-8329-97a95b9fe22d"
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
}

data "kong_consumer_oauth2" "consumer_user_id" {
  id                = "068141a0-93c3-4286-8329-97a95b9fe22d"
  consumer_username = "example"
}
