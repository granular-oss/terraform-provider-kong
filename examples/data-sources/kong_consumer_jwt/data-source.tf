provider "kong" {}

data "kong_consumer_jwt" "ids" {
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
  id          = "da8b0cd6-f1d3-4731-ba20-873254d9d474"
}

data "kong_consumer_jwt" "username_id" {
  consumer_username = "consumer"
  id                = "da8b0cd6-f1d3-4731-ba20-873254d9d474"
}

data "kong_consumer_jwt" "id_key" {
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
  key         = "some_key"
}

data "kong_consumer_jwt" "username_key" {
  consumer_username = "consumer"
  key               = "some_key"
}
