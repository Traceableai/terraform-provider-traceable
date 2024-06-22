terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
      version = "0.0.1"
    }
    # aws={
    #   source = "hashicorp/aws"
    #   version = "5.35.0"
    # }
  }
  #   backend "s3" {
  #   bucket         = "traceable-provider-store"
  #   key            = "traceable-provider-store"
  #   region         = "us-west-2"
  # }
}

# data "aws_secretsmanager_secret_version" "api_token" {
#   secret_id = "your secret manager arn where api token is stored"
# }

# output "api_token" {
#   value=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
#   sensitive = true
# }

provider "traceable" {
  platform_url="https://app.traceable.ai/graphql"
  # api_token=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
  api_token="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlFqYzBNell4UkRKRVJUSkVOMFZGTlVRMk4wVXlOVFZCTlVVME1rVTBSVUl6T0VZNVF6VTFPQSJ9.eyJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlc192MiI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvY3VzdG9tZXJfaWQiOiI2NGJmM2Q0Zi01ZGVjLTRkZWEtYWEzNC03MzJhYWNmZjBkNjEiLCJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlcyI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvanRpIjoiMzg1MDI1YWMtOWViNi00ZWI3LThmZjgtYzFjZGFmNzQ4NDJhIiwiaHR0cHM6Ly90cmFjZWFibGUuYWkvcmljaF9yb2xlcyI6W3siZW52cyI6W10sImlkIjoidHJhY2VhYmxlIn1dLCJuaWNrbmFtZSI6InN0YWdpbmdfYXV0b21hdGlvbnRlc3QiLCJuYW1lIjoic3RhZ2luZ19hdXRvbWF0aW9udGVzdEB0cmFjZWFibGV0ZXN0LmNvbSIsInBpY3R1cmUiOiJodHRwczovL3MuZ3JhdmF0YXIuY29tL2F2YXRhci81YmU4ZDBkMjA1NTc3MWZjNDg5NzdmYmMwZjQyMDlkOT9zPTQ4MCZyPXBnJmQ9aHR0cHMlM0ElMkYlMkZjZG4uYXV0aDAuY29tJTJGYXZhdGFycyUyRnN0LnBuZyIsInVwZGF0ZWRfYXQiOiIyMDI0LTA2LTIyVDAxOjQ5OjE3LjA1M1oiLCJlbWFpbCI6InN0YWdpbmdfYXV0b21hdGlvbnRlc3RAdHJhY2VhYmxldGVzdC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6Ly9hdXRoLnRyYWNlYWJsZS5haS8iLCJhdWQiOiJ1czVrZGJueGNlM05oZUxiekxDeHVacVlJUVlnUWdtOCIsImlhdCI6MTcxOTAzMzE5NywiZXhwIjoxNzE5MDY5MTk3LCJzdWIiOiJhdXRoMHw2MTkzNjRhMjBhOTk3ZjAwNjk2ZDM4ODciLCJzaWQiOiJPSzhabzNVME1KVmlEZmV2X20wc1NyVXpEcnRQSUNORSJ9.lWqDgOXgwN6pkX9DQNPW8vsq-iawdJLKwkyhAaW8v9tzuLBvMzwD6NnL1d7omxiyxBCU9trXr8c_2vJn7djQyB5zooYxBsz4P9k1EZwP2PqeWxoXGOFSwlyJ07RhZIp888J_KTpCJLd4t7onxsv212UEDNBkgYPEFgsHTQqU2bfc52zg9CI0uvkn9lOdtfEkstTwkTci6fgiF9IlbX_PZtNJyFlr5-Zyqg5UXenV2gJ2g07XTi8jQgQ98pBG2aGJYcwDMjW3nHZrBrUBC0qX_c4syyqSCMkXYcbZQC3w2K_cQ4O8ioYvrIKXqkYR1nNvhE2IUgWo-ZKVxZ3QUV2u4A"
}

# resource "traceable_user_attribution_rule_basic_auth" "test1" {
#   name = "traceable_user_attribution_rule_basic_auth_1"
#   scope_type = "SYSTEM_WIDE"
# }

# resource "traceable_user_attribution_rule_req_header" "test2" {
#   name = "traceable_user_attribution_rule_req_header"
#   scope_type = "CUSTOM"
#   url_regex = "abcd"
#   auth_type = "test"
#   user_id_location = "test"
#   user_role_location="yes"
#   role_location_regex_capture_group="test"
# }

# resource "traceable_user_attribution_rule_jwt_authentication" "test3" {
#   name = "jwtauth"
#   scope_type = "CUSTOM"
#   url_regex="sfdsf"
#   jwt_location = "COOKIE"
#   jwt_key = "abcd"
#   user_id_claim = "aditya"
# }

# resource "traceable_user_attribution_rule_response_body" "test4" {
#   name = "resbody"
#   url_regex="sfdsf"
#   user_id_location_json_path="test"
#   auth_type="sadsak"
#   user_role_location_json_path="hjasa"
# }

# resource "traceable_user_attribution_rule_custom_json" "test5" {
#   name = "traceable_user_attribution_rule_custom_json"
#   scope_type="CUSTOM"
#   url_regex="sfdsf"
#   auth_type_json=jsonencode(file("authType.json"))
#   user_id_json=jsonencode(file("authType.json"))
# }

# resource "traceable_user_attribution_rule_custom_token" "test6" {
#   name = "traceable_user_attribution_rule_custom_token"
#   scope_type="SYSTEM_WIDE"
#   auth_type="test"
#   location="REQUEST_COOKIE"
#   token_name="test"
# }

# data "traceable_syslog_integration" "syslog" {
#   name="prer-test"
# }

# data "traceable_splunk_integration" "splunk" {
#   name="aditya"
# }

# data "traceable_endpoint_id" "endpoint" {
#   name="POST /Unauthenticated_Modification_of_external_APIs"
#   service_name="nginx-automation-test"
#   enviroment_name="fintech-1"
# }

# data "traceable_service_id" "endpoint" {
#   service_name="nginx-automation-test"
#   enviroment_name="fintech-1"
# }

# output "traceable_service_id" {
#   value = data.traceable_service_id.endpoint.service_id
# }

# resource "traceable_notification_channel" "testchannel" {
#   channel_name = "example_channel1"

#   email = [
#     "example4@example.com",
#     "example2@example.com"
#   ]

#   slack_webhook = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
#   splunk_id=data.traceable_splunk_integration.splunk.splunk_id
#   # syslog_id=""
#   custom_webhook  {
#     webhook_url = "https://example.com/webhook"
#     custom_webhook_headers  {
#       key       = "Authorization"
#       value     = "Bearer token123"
#       is_secret = false
#     }
#     custom_webhook_headers  {
#       key       = "Authorization1"
#       value     = "Bearer token1232"
#       is_secret = true
#     }
#     custom_webhook_headers  {
#       key       = "tets"
#       value     = "Bearer"
#       is_secret = false
#     }
#   }

#   s3_webhook  {
#     bucket_name = "your-s3-bucket"
#     region      = "us-west-2"
#     bucket_arn  = "arn:aws:s3:::your-s3-bucket"
#   }
# }

data "traceable_notification_channels" "mychannel"{
  name = "example_channel1"
}

# resource "traceable_notification_rule_logged_threat_activity" "rule1" {
#   name                    = "example_notification_rule"
#   environments            = []
#   channel_id              = data.traceable_notification_channels.mychannel.channel_id
#   threat_types            = ["SQLInjection","bola"]
#   severities              = ["HIGH", "MEDIUM","LOW","CRITICAL"]
#   impact                  = ["LOW", "HIGH"]
#   confidence              = ["HIGH", "MEDIUM"]
# }

# resource "traceable_label_creation_rule" "example_label_create_rule" {
#   key="test-rule-create-label-test1"
#   description="test rule to create a label"
#   color="#E295E9"
# }

# resource "traceable_agent_token" "example" {
#   name = "tf-provider-token-testing"
# }

# output "agent_token" {
#   value = traceable_agent_token.example.token
#   sensitive = true
# }

# output "agent_token_creation_timestamp" {
#   value = traceable_agent_token.example.creation_timestamp
# }

# data "traceable_agent_token" "example" {
#   name = "tf-provider-token-testing"
# }

# output "agent_token" {
#   value = data.traceable_agent_token.example.token
#   sensitive = true
# }

# output "agent_token_creation_timestamp" {
#   value = data.traceable_agent_token.example.creation_timestamp

# }

resource "traceable_notification_rule_protection_configuration_change" "protection_config" {
  name="aditya"
  environments=["3095423142-ip-blocking"]
  channel_id=data.traceable_notification_channels.mychannel.channel_id
  security_configuration_types=[]
  notification_frequency="PT1H"
}
resource "traceable_notification_rule_team_activity" "team_activity" {
  name="team-activity-1"
  channel_id=data.traceable_notification_channels.mychannel.channel_id
  user_change_types=[]
  notification_frequency="PT1H"
}

# resource "traceable_notification_rule_api_naming" "api_naming" {
#     name = "traceable_notification_rule_api_naming"
#     channel_id = data.traceable_notification_channels.mychannel.channel_id
#     event_types = []
#     # notification_frequency = "PT1H"
# }

# resource "traceable_notification_rule_api_documentation" "api_documentation" {
#     name = "traceable_notification_rule_api_documentation"
#     channel_id = data.traceable_notification_channels.mychannel.channel_id
#     event_types = []
#     # notification_frequency = "PT1H"
# }

# resource "traceable_notification_rule_data_collection" "data_collection" {
#     name = "traceable_notification_rule_data_collection"
#     environments = ["3095423142-ip-blocking"]
#     channel_id = data.traceable_notification_channels.mychannel.channel_id
#     agent_activity_type = "NO_DATA_IN_ENVIRONMENT"
#     notification_frequency = "PT1H"
# }

# resource "traceable_notification_rule_risk_scoring" "risk_scoring" {
#     name = "traceable_notification_rule_risk_scoring"
#     environments = ["3095423142-ip-blocking"]
#     channel_id = data.traceable_notification_channels.mychannel.channel_id
#     event_types = ["UPDATE"]
#     # notification_frequency = "PT1H"
# }

# resource "traceable_notification_rule_exclude_rule" "exclude_rule" {
#     name = "traceable_notification_rule_exclude_rule"
#     channel_id = data.traceable_notification_channels.mychannel.channel_id
#     event_types = ["CREATE"]
#     # notification_frequency = "PT1H"
# }

output "agent_token" {
  value = data.traceable_agent_token.example.token
  sensitive = true
}

output "agent_token_creation_timestamp" {
  value = data.traceable_agent_token.example.creation_timestamp

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
