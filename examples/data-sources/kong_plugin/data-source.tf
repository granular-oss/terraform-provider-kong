provider "kong" {}

data "kong_plugin" "service_id" {
  service_id = "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"
  name       = "cors"
}

data "kong_plugin" "service_name" {
  service_name = "example"
  name         = "cors"
}

data "kong_plugin" "route_id" {
  route_id = "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"
  name     = "cors"
}

data "kong_plugin" "route_name" {
  route_name = "example"
  name       = "cors"
}

data "kong_plugin" "consumer_id" {
  consumer_id = "204ee4c2-6dc0-44fb-93e1-4fbeb24489c7"
  name        = "aws-lambda"
}

data "kong_plugin" "consumer_name" {
  consumer_name = "example"
  name          = "aws-lambda"
}
