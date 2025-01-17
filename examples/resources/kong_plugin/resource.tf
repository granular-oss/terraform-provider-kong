provider "kong" {}

resource "kong_service" "service" {
  name = "example-service"
  host = "example.com"
}

resource "kong_plugin" "service_example" {
  service_id = kong_service.service.id
  name       = "cors"
  config_json = jsonencode(
    {
      origins = [".*\\.example2.com"]
      methods = ["GET"]
    }
  )
}

resource "kong_route" "route" {
  name       = "example-route"
  service_id = kong_service.service.id
}

resource "kong_plugin" "route_example" {
  route_id = kong_service.service.id
  name     = "aws-lambda"
  config_json = jsonencode(
    {
      function_name = "fake_lambda"
    }
  )
}

resource "kong_consumer" "plugin" {
  username = "consumer_plugin"
}
resource "kong_plugin" "test" {
  consumer_id = kong_consumer.plugin.id
  name        = "aws-lambda"
  config_json = jsonencode(
    {
      function_name = "test"
    }
  )
}
