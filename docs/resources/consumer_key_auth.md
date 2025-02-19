---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kong_consumer_key_auth Resource - kong"
subcategory: ""
description: |-
  
---

# kong_consumer_key_auth (Resource)



## Example Usage

```terraform
provider "kong" {}

resource "kong_consumer" "consumer" {
  username = "foobar"
}

resource "kong_consumer_key_auth" "example" {
  consumer_id = kong_consumer.consumer.id
  key         = "SECRET_KEY"
  tags        = ["test"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `consumer_id` (String)
- `key` (String)

### Optional

- `tags` (List of String)

### Read-Only

- `id` (String) The ID of this resource.
- `kong_id` (String)

## Import

Import is supported using the following syntax:

```shell
# Consumer ACLs can be import by consumer id/username:id/username
terraform import kong_consumer.example <consumer_id>:<key_auth_id>
terraform import kong_consumer.example <consumer_username>:<key_auth_id>
terraform import kong_consumer.example <consumer_id>:<key>
terraform import kong_consumer.example <consumer_username>:<key>
```
