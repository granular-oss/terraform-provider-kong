provider "kong" {}

resource "kong_consumer" "consumer" {
  username = "foobar"
}

resource "kong_consumer_basic_auth" "example" {
  consumer_id = kong_consumer.consumer.id
  username    = "player1"
  password    = "ready"
  tags        = ["test"]
}
