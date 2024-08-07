terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
      version = "0.0.1"
    }
    aws={
      source = "hashicorp/aws"
      version = "5.35.0"
    }
  }
    backend "s3" {
    bucket         = "traceable-provider-store"
    key            = "traceable-provider-store"
    region         = "us-west-2"
  }
}

data "aws_secretsmanager_secret_version" "api_token" {
  secret_id = "your secret manager arn where api token is stored"
}

output "api_token" {
  value=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
  sensitive = true
}

provider "traceable" {
  platform_url="https://api-dev.traceable.ai/graphql"
  api_token=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
}

resource "traceable_user_attribution_rule_basic_auth" "test1" {
  name = "traceable_user_attribution_rule_basic_auth"
  scope_type = "SYSTEM_WIDE"
  url_regex = "abcd"
}

resource "traceable_user_attribution_rule_req_header" "test2" {
  name = "traceable_user_attribution_rule_req_header"
  scope_type = "CUSTOM"
  url_regex = "abcd"
  auth_type = "test"
  user_id_location = "test"
  user_role_location="yes"
  role_location_regex_capture_group="test"
}

resource "traceable_user_attribution_rule_jwt_authentication" "test3" {
  name = "jwtauth"
  scope_type = "CUSTOM"
  url_regex="sfdsf"
  jwt_location = "COOKIE"
  jwt_key = "abcd"
  user_id_claim = "testuser"
}

resource "traceable_user_attribution_rule_response_body" "test4" {
  name = "resbody"
  url_regex="sfdsf"
  user_id_location_json_path="test"
  auth_type="sadsak"
  user_role_location_json_path="hjasa"
}

resource "traceable_user_attribution_rule_custom_json" "test5" {
  name = "traceable_user_attribution_rule_custom_json"
  scope_type="CUSTOM"
  url_regex="sfdsf"
  auth_type_json=jsonencode(file("authType.json"))
  user_id_json=jsonencode(file("authType.json"))
}

resource "traceable_user_attribution_rule_custom_token" "test6" {
  name = "traceable_user_attribution_rule_custom_token"
  scope_type="SYSTEM_WIDE"
  auth_type="test"
  location="REQUEST_COOKIE"
  token_name="test"
}

data "traceable_syslog_integration" "syslog" {
  name="test"
}

data "traceable_splunk_integration" "splunk" {
  name="test"
}

data "traceable_endpoint_id" "endpoint" {
  name="POST /Unauthenticated_Modification_of_external_APIs"
  service_name="test-service"
  enviroment_name="test-env"
}

data "traceable_service_id" "endpoint" {
  service_name="test-service"
  enviroment_name="test-env"
}

output "traceable_service_id" {
  value = data.traceable_service_id.endpoint.service_id
}

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

data "traceable_notification_channels" "mychannel"{
  name = "example_channel1"
}

resource "traceable_notification_rule_logged_threat_activity" "rule1" {
  name                    = "example_notification_rule"
  environments            = []
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  threat_types            = ["SQLInjection","bola"]
  severities              = ["HIGH", "MEDIUM","LOW","CRITICAL"]
  impact                  = ["LOW", "HIGH"]
  confidence              = ["HIGH", "MEDIUM"]
}

resource "traceable_label_creation_rule" "example_label_create_rule" {
  key="test-rule-create-label-test1"
  description="test rule to create a label"
  color="#E295E9"
}

resource "traceable_agent_token" "example" {
  name = "tf-provider-token-testing-resource-latest"
}

data "traceable_agent_token" "example" {
  name = traceable_agent_token.example.name
  depends_on = [traceable_agent_token.example]
}

output "agent_token" {
  value = traceable_agent_token.example.token
}

output "agent_token_creation_timestamp" {
  value = data.traceable_agent_token.example.creation_timestamp
}



resource "traceable_ip_range_rule" "my_ip_range" {
    name     = "first_rule"
    rule_action     = "RULE_ACTION_ALERT"
    event_severity     = "LOW"
    raw_ip_range_data = [
        "1.1.1.1",
        "3.3.3.3"
    ]
    environment=[] #all env
    description="rule created from custom provider"
}
resource "traceable_notification_rule_blocked_threat_activity" "rule1" {
  name                    = "example_notification_rule3"
  environments            = []
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  threat_types            = []
  notification_frequency  = "PT1H"
}

resource "traceable_notification_rule_threat_actor_status" "rule1" {
  name                    = "terraform_threat_actor_status"
  environments            = ["fintech-1"]
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  actor_states            = ["NORMAL"]
}
resource "traceable_notification_rule_actor_severity_change" "rule1" {
  name                    = "terraform_threat_actor_severity2"
  environments            = ["fintech-1"]
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  actor_severities            = []
  actor_ip_reputation_levels            = ["HIGH"]

}