provider "kong" {}

resource "kong_consumer" "consumer" {
  username = "foobar"
}

resource "kong_consumer_jwt_auth" "rsa_example" {
  consumer_id    = kong_consumer.consumer.id
  rsa_public_key = "fake_public_key"
  algorithm      = "rs256"
}

resource "kong_consumer_jwt_auth" "hs256_example" {
  consumer_id = kong_consumer.consumer.id
  secret      = "fake_secret"
  algorithm   = "hs256"
}
