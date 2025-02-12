---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_notification_rule_logged_threat_activity Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  
---

# traceable_notification_rule_logged_threat_activity (Resource)



## Example Usage

```terraform
resource "traceable_notification_rule_logged_threat_activity" "rule1" {
  name                    = "example_notification_rule"
  environments            = []
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  threat_types            = ["SQLInjection","bola"]
  severities              = ["HIGH", "MEDIUM","LOW","CRITICAL"]
  impact                  = ["LOW", "HIGH"]
  confidence              = ["HIGH", "MEDIUM"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `channel_id` (String) Reporting channel for this notification rule
- `confidence` (Set of String) Confidence of threat events you want to notify (LOW,MEDIUM,HIGH)
- `environments` (Set of String) Environments where rule will be applicable
- `impact` (Set of String) Impact of threat events you want to notify (LOW,MEDIUM,HIGH)
- `name` (String) Name of the notification rule
- `severities` (Set of String) Severites of threat events you want to notify (LOW,MEDIUM,HIGH,CRITICAL)
- `threat_types` (Set of String) Threat types for which you want notification

### Optional

- `category` (String) Type of notification rule
- `notification_frequency` (String) No more than one notification every configured notification_frequency (should be in this format PT1H for 1 hr)

### Read-Only

- `id` (String) The ID of this resource.
