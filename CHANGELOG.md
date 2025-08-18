# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 8.1.0

- Update `kong_plugin` to be imported as `|` separated string to be consistent with the other object types
- Update `kong_target` to use the shared ID builder function
- Fix docs for imports so they show the correct ID delineator for imports

## 8.0.3

- Fix bug where provider fails if `kong_route.hosts` is set to empty array.

## 8.0.2

- Fix reading provider config so that it reads config not just environment variables.

## 8.0.1

- Update change that broke backwards compatibility in `kong_route`:
  - Header is a block object instead of attribute. See example for more details

## 8.0.0

- Update to use the `terraform-plugin-framework`
- Add `data-source` for all the corresponding `resource` types
- Update `resource` imports to allow for human readable properties where available in kong
