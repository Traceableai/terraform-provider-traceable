---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_notification_rule_team_activity Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  
---

# traceable_notification_rule_team_activity (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `channel_id` (String) Reporting channel for this notification rule
- `name` (String) Name of the notification rule
- `user_change_types` (Set of String) User change types for which you want notification

### Optional

- `category` (String) Type of notification rule
- `notification_frequency` (String) No more than one notification every configured notification_frequency (should be in this format PT1H for 1 hr)

### Read-Only

- `id` (String) The ID of this resource.