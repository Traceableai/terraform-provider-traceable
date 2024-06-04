resource "traceable_notification_channel" "testchannel" {
  channel_name = "example_channel1"

  email = [
    "example4@example.com",
    "example2@example.com"
  ]

  slack_webhook = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  splunk_id=data.traceable_splunk_integration.splunk.splunk_id
  # syslog_id=""
  custom_webhook  {
    webhook_url = "https://example.com/webhook"
    custom_webhook_headers  {
      key       = "Authorization"
      value     = "Bearer token123"
      is_secret = false
    }
    custom_webhook_headers  {
      key       = "Authorization1"
      value     = "Bearer token1232"
      is_secret = true
    }
    custom_webhook_headers  {
      key       = "tets"
      value     = "Bearer"
      is_secret = false
    }
  }

  s3_webhook  {
    bucket_name = "your-s3-bucket"
    region      = "us-west-2"
    bucket_arn  = "arn:aws:s3:::your-s3-bucket"
  }
}