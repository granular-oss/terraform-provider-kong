---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "kong_certificate Data Source - kong"
subcategory: ""
description: |-
  
---

# kong_certificate (Data Source)



## Example Usage

```terraform
provider "kong" {}

data "kong_certificate" "example" {
  id = "50c86b96-b973-4c8a-933f-7f48f6f49896"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `cert` (String, Sensitive)
- `id` (String) The ID of this resource.
- `key` (String, Sensitive)
- `snis` (List of String)
- `tags` (List of String)
