provider "kong" {}

resource "kong_consumer" "consumer" {
  username = "foobar"
}

resource "kong_consumer_key_auth" "example" {
  consumer_id = kong_consumer.consumer.id
  key         = "SECRET_KEY"
  tags        = ["test"]
}
