---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_user_attribution_rule_jwt_authentication Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  
---

# traceable_user_attribution_rule_jwt_authentication (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `jwt_key` (String) header name for jwt in header or cookie
- `jwt_location` (String) header or cookie
- `name` (String) name of the user attribution rule
- `scope_type` (String) system wide, environment, url regex
- `user_id_claim` (String) user id claim

### Optional

- `auth_type` (String) auth type of the user attribution rule
- `environment` (String) environment
- `token_capture_group` (String) token capture group
- `url_regex` (String) url regex
- `user_id_location_json_path` (String) user id location json path
- `user_role_claim` (String) user role claim
- `user_role_location_json_path` (String) user role location json path

### Read-Only

- `id` (String) The ID of this resource.