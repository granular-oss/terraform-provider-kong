---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kong_consumer_jwt Resource - kong"
subcategory: ""
description: |-
  
---

# kong_consumer_jwt (Resource)



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `consumer_id` (String)

### Optional

- `algorithm` (String)
- `key` (String)
- `rsa_public_key` (String)
- `secret` (String, Sensitive)
- `tags` (List of String)

### Read-Only

- `id` (String) The ID of this resource.
- `kong_id` (String)

## Import

Import is supported using the following syntax:

```shell
# Consumer ACLs can be import by consumer id/username:id/key
terraform import kong_consumer.example <consumer_id>:<jwt_id>
terraform import kong_consumer.example <consumer_username>:<jwt_id>
terraform import kong_consumer.example <consumer_id>:<key>
terraform import kong_consumer.example <consumer_username>:<key>
```
