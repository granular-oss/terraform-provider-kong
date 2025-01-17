resource "kong_certificate" "example" {
  certificate = file("<path to certificate>")
  private_key = file("<path to key>")
  snis        = ["sni1"]
  tags        = ["example_cert"]
}
