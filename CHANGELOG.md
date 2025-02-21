## 8.0.1

- Update change that broke backwards compatibility in `kong_route`:
  - Header is a block object instead of attribute. See example for more details

## 8.0.0

- Update to use the `terraform-plugin-framework`
- Add `data-source` for all the corresponding `resource` types
- Update `resource` imports to allow for human readable properties where available in kong
