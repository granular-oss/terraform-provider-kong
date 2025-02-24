provider "kong" {}

resource "kong_service" "service" {
  name = "test2"
  host = "fake.example.com"
}

resource "kong_route" "minimal" {
  name       = "test2"
  paths      = ["/foobar"]
  service_id = kong_service.service.id
}

resource "kong_route" "complete" {
  name                       = "test2"
  paths                      = ["/foobaz"]
  service_id                 = kong_service.service.id
  protocols                  = ["http"]
  methods                    = ["GET"]
  hosts                      = ["test.example.com"]
  tags                       = ["test", "app:fake"]
  strip_path                 = false
  preserve_host              = true
  regex_priority             = 3
  path_handling              = "v1"
  https_redirect_status_code = 301
  request_buffering          = false
  response_buffering         = false
  header {
    name  = "header-test"
    value = ["1", "2"]
  }
}
