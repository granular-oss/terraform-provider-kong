provider "kong" {}

data "kong_upstream" "example_id_upstream_name" {
  id            = "921687b2-ac57-4fa5-a99f-ae08fbf30eff"
  upstream_name = "test-upstream"
}

data "kong_upstream" "example_id_upstream_id" {
  id          = "921687b2-ac57-4fa5-a99f-ae08fbf30eff"
  upstream_id = "921687b2-ac57-4fa5-a99f-ae08fbf30eff"
}

data "kong_upstream" "example_target_upstream_name" {
  target        = "example.com:8000"
  upstream_name = "test-upstream"
}

data "kong_upstream" "example_target_upstream_id" {
  target      = "example.com:8000"
  upstream_id = "921687b2-ac57-4fa5-a99f-ae08fbf30eff"
}
