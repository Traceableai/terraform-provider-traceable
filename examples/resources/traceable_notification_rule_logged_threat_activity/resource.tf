resource "traceable_notification_rule_logged_threat_activity" "rule1" {
  name                    = "example_notification_rule"
  environments            = []
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  threat_types            = ["SQLInjection","bola"]
  severities              = ["HIGH", "MEDIUM","LOW","CRITICAL"]
  impact                  = ["LOW", "HIGH"]
  confidence              = ["HIGH", "MEDIUM"]
}