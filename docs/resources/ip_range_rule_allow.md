---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_ip_range_rule_allow Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  
---

# traceable_ip_range_rule_allow (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment` (Set of String) environment where it will be applied
- `name` (String) name of the policy
- `raw_ip_range_data` (Set of String) IPV4/V6 range for the rule

### Optional

- `description` (String) description of the policy
- `expiration` (String) expiration for allow actions
- `rule_action` (String) Need to provide the action to be performed

### Read-Only

- `id` (String) The ID of this resource.
