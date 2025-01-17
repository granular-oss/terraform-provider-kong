provider "kong" {}

resource "kong_consumer" "consumer" {
  username = "foobar"
}

resource "kong_consumer_oauth2" "example" {
  consumer_id   = kong_consumer.consumer.id
  name          = "example_oauth2"
  client_id     = ""
  client_secret = "super_secret"
  hash_secret   = true
  redirect_uris = ["http://example.com"]
  tags          = ["tag"]
}
