---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_notification_rule_threat_actor_status Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  
---

# traceable_notification_rule_threat_actor_status (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `actor_states` (Set of String) Actor states for which you want notification
- `channel_id` (String) Reporting channel for this notification rule
- `environments` (Set of String) Environments where rule will be applicable
- `name` (String) Name of the notification rule

### Optional

- `category` (String) Type of notification rule
- `notification_frequency` (String) No more than one notification every configured notification_frequency (should be in this format PT1H for 1 hr)

### Read-Only

- `id` (String) The ID of this resource.
