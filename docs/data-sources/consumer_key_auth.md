---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kong_consumer_key_auth Data Source - kong"
subcategory: ""
description: |-
  
---

# kong_consumer_key_auth (Data Source)



## Example Usage

```terraform
provider "kong" {}

data "kong_consumer_key_auth" "consumer_id_id" {
  id          = "068141a0-93c3-4286-8329-97a95b9fe22d"
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
}

data "kong_consumer_key_auth" "consumer_user_id" {
  id                = "068141a0-93c3-4286-8329-97a95b9fe22d"
  consumer_username = "example"
}

data "kong_consumer_basic_auth" "consumer_id_key" {
  key         = "user"
  consumer_id = "068141a0-93c3-4286-8329-97a95b9fe22d"
}

data "kong_consumer_basic_auth" "consumer_user_key" {
  key               = "user"
  consumer_username = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `consumer_id` (String)
- `consumer_username` (String)
- `key` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `tags` (List of String)
