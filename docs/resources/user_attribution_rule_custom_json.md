---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_user_attribution_rule_custom_json Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  
---

# traceable_user_attribution_rule_custom_json (Resource)



## Example Usage

```terraform
resource "traceable_user_attribution_rule_custom_json" "test5" {
  name = "traceable_user_attribution_rule_custom_json"
  scope_type="CUSTOM"
  url_regex="sfdsf"
  auth_type_json=jsonencode(file("authType.json"))
  user_id_json=jsonencode(file("authType.json"))
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auth_type_json` (String) auth type json
- `name` (String) name of the user attribution rule
- `scope_type` (String) system wide, environment, url regex

### Optional

- `disabled` (Boolean) Flag to enable or disable the rule
- `environment` (String) environement of rule
- `url_regex` (String) url regex
- `user_id_json` (String) user id json
- `user_role_json` (String) user role json

### Read-Only

- `id` (String) The ID of this resource.
