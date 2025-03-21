---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_ip_type_rule_alert Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  
---

# traceable_ip_type_rule_alert (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment` (List of String) environment where it will be applied
- `event_severity` (String) Generated event severity among LOW,MEDIUM,HIGH,CRITICAL
- `ip_types` (List of String) Ip types to include for the rule among ANONYMOUS VPN,HOSTING PROVIDER,PUBLIC PROXY,TOR EXIT NODE,BOT
- `name` (String) name of the policy

### Optional

- `description` (String) description of the policy
- `inject_request_headers` (Block List) Inject Data in Request header? (see [below for nested schema](#nestedblock--inject_request_headers))
- `rule_action` (String) Need to provide the action to be performed

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--inject_request_headers"></a>
### Nested Schema for `inject_request_headers`

Required:

- `header_key` (String)
- `header_value` (String)
