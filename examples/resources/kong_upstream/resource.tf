terraform {
  required_providers {
    kong = {
      source = "granular-oss/kong"
    }
  }
}

provider "kong" {}

resource "kong_upstream" "minimal" {
  name = "example"
}

resource "kong_upstream" "complete" {
  name                 = "test2"
  slots                = 10000
  hash_on              = "none"
  hash_fallback        = "none"
  hash_on_header       = "header"
  host_header          = "host_header"
  tags                 = ["tag"]
  hash_fallback_header = "fallback"
  hash_on_cookie       = "cookie"
  hash_on_cookie_path  = "/cookie"
  healthchecks {
    active {
      type                     = "http"
      concurrency              = 10
      http_path                = "/"
      https_verify_certificate = true
      timeout                  = 60
      healthy {
        http_statuses = [200]
        interval      = 10
        successes     = 2
      }
      unhealthy {
        http_failures = 2
        interval      = 10
        timeouts      = 2
        tcp_failures  = 2
        http_statuses = [429, 500, 503]
      }
    }
    passive {
      type = "http"
      healthy {
        http_statuses = [200]
        successes     = 2
      }
      unhealthy {
        http_failures = 2
        timeouts      = 2
        tcp_failures  = 2
        http_statuses = [429, 500, 503]
      }
    }
  }
}
