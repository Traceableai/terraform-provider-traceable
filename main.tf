terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
      version = "0.0.1"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "5.35.0"
    }
  }
  backend "s3" {
    bucket = "traceable-provider-store"
    key    = "traceable-provider-store"
    region = "us-west-2"
  }
}

data "aws_secretsmanager_secret_version" "api_token" {
  secret_id = "your secret manager arn where api token is stored"
}

output "api_token" {
  value     = jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
  sensitive = true
}

variable "API_TOKEN" {
  default = ""
}

provider "traceable" {
  platform_url = "platform url"
  api_token    = jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
  provider_version="terraform/v1.0.1"
}

resource "traceable_user_attribution_rule_basic_auth" "test1" {
  name       = "traceable_user_attribution_rule_basic_auth_1"
  scope_type = "SYSTEM_WIDE"
}

resource "traceable_user_attribution_rule_req_header" "test2" {
  name                              = "traceable_user_attribution_rule_req_header"
  scope_type                        = "CUSTOM"
  url_regex                         = "abcd"
  auth_type                         = "test"
  user_id_location                  = "test"
  user_role_location                = "yes"
  role_location_regex_capture_group = "test"
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
  name                         = "resbody"
  url_regex                    = "sfdsf"
  user_id_location_json_path   = "test"
  auth_type                    = "sadsak"
  user_role_location_json_path = "hjasa"
}

resource "traceable_user_attribution_rule_custom_json" "test5" {
  name           = "traceable_user_attribution_rule_custom_json"
  scope_type     = "CUSTOM"
  url_regex      = "sfdsf"
  auth_type_json = jsonencode(file("authType.json"))
  user_id_json   = jsonencode(file("authType.json"))
}

resource "traceable_user_attribution_rule_custom_token" "test6" {
  name       = "traceable_user_attribution_rule_custom_token"
  scope_type = "SYSTEM_WIDE"
  auth_type  = "test"
  location   = "REQUEST_COOKIE"
  token_name = "test"
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
  splunk_id     = data.traceable_splunk_integration.splunk.splunk_id
  # syslog_id=""
  custom_webhook {
    webhook_url = "https://example.com/webhook"
    custom_webhook_headers {
      key       = "Authorization"
      value     = "Bearer token123"
      is_secret = false
    }
    custom_webhook_headers {
      key       = "Authorization1"
      value     = "Bearer token1232"
      is_secret = true
    }
    custom_webhook_headers {
      key       = "tets"
      value     = "Bearer"
      is_secret = false
    }
  }

  s3_webhook {
    bucket_name = "your-s3-bucket"
    region      = "us-west-2"
    bucket_arn  = "arn:aws:s3:::your-s3-bucket"
  }
}

data "traceable_notification_channels" "mychannel" {
  name = "helloworld"
}

resource "traceable_notification_rule_logged_threat_activity" "rule1" {
  name         = "example_notification_rule"
  environments = []
  channel_id   = data.traceable_notification_channels.mychannel.channel_id
  threat_types = ["SQLInjection", "bola"]
  severities   = ["HIGH", "MEDIUM", "LOW", "CRITICAL"]
  impact       = ["LOW", "HIGH"]
  confidence   = ["HIGH", "MEDIUM"]
}

resource "traceable_label_creation_rule" "example_label_create_rule" {
  key         = "test-rule-create-label-test1"
  description = "test rule to create a label"
  color       = "#E295E9"
}

resource "traceable_notification_rule_protection_configuration_change" "protection_config" {
  name                         = "aditya"
  environments                 = ["3095423142-ip-blocking"]
  channel_id                   = data.traceable_notification_channels.mychannel.channel_id
  security_configuration_types = []
  notification_frequency       = "PT1H"
}

resource "traceable_notification_rule_team_activity" "team_activity" {
  name                   = "team-activity-1"
  channel_id             = data.traceable_notification_channels.mychannel.channel_id
  user_change_types      = []
  notification_frequency = "PT1H"
}

output "agent_token" {
  value = traceable_agent_token.example.token
}

output "agent_token_creation_timestamp" {
  value = traceable_agent_token.example.creation_timestamp
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

resource "traceable_notification_rule_protection_configuration_change" "protection_config" {
  name                         = "aditya"
  environments                 = ["3095423142-ip-blocking"]
  channel_id                   = data.traceable_notification_channels.mychannel.channel_id
  security_configuration_types = []
  notification_frequency       = "PT1H"
}
resource "traceable_notification_rule_team_activity" "team_activity" {
  name                   = "team-activity-1"
  channel_id             = data.traceable_notification_channels.mychannel.channel_id
  user_change_types      = []
  notification_frequency = "PT1H"
}

resource "traceable_notification_rule_api_naming" "api_naming" {
  name        = "traceable_notification_rule_api_naming"
  channel_id  = data.traceable_notification_channels.mychannel.channel_id
  event_types = []
  # notification_frequency = "PT1H"
}

resource "traceable_notification_rule_api_documentation" "api_documentation" {
  name        = "traceable_notification_rule_api_documentation"
  channel_id  = data.traceable_notification_channels.mychannel.channel_id
  event_types = []
  # notification_frequency = "PT1H"
}

resource "traceable_notification_rule_data_collection" "data_collection" {
  name                   = "traceable_notification_rule_data_collection"
  environments           = ["3095423142-ip-blocking"]
  channel_id             = data.traceable_notification_channels.mychannel.channel_id
  agent_activity_type    = "NO_DATA_IN_ENVIRONMENT"
  notification_frequency = "PT1H"
}

resource "traceable_notification_rule_risk_scoring" "risk_scoring" {
  name         = "traceable_notification_rule_risk_scoring"
  environments = ["3095423142-ip-blocking"]
  channel_id   = data.traceable_notification_channels.mychannel.channel_id
  event_types  = ["UPDATE"]
  # notification_frequency = "PT1H"
}

resource "traceable_notification_rule_exclude_rule" "exclude_rule" {
  name                   = "traceable_notification_rule_exclude_rule"
  channel_id             = data.traceable_notification_channels.mychannel.channel_id
  event_types            = ["CREATE"]
  notification_frequency = "PT1H"
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
  name                   = "example_notification_rule3"
  environments           = []
  channel_id             = data.traceable_notification_channels.mychannel.channel_id
  threat_types           = []
  notification_frequency = "PT1H"
}

resource "traceable_notification_rule_threat_actor_status" "rule1" {
  name         = "terraform_threat_actor_status"
  environments = ["fintech-1"]
  channel_id   = data.traceable_notification_channels.mychannel.channel_id
  actor_states = ["NORMAL"]
}

resource "traceable_notification_rule_actor_severity_change" "rule1" {
  name                    = "terraform_threat_actor_severity"
  environments            = ["fintech-1"]
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  actor_severities            = []
  actor_ip_reputation_levels  = ["HIGH"]
}
resource "traceable_notification_rule_posture_events" "rule1" {
  name                    = "terraform_notification_posture_events"
  environments            = ["fintech"]
  channel_id              = data.traceable_notification_channels.mychannel.channel_id
  posture_events            = ["RISK_SCORE_CHANGE"]
  risk_deltas  = ["INCREASE"]
}


resource traceable_custom_signature_allow "csruletf"{
    name="testtf2"
    description="test1"
    disabled = true
    environments=["fintech-1","demo-test"]
    custom_sec_rule=<<EOT
     SecRule REQUEST_HEADERS:key-sec "@rx val-sec" \
     "id:92100120,\
     phase:2,\
     block,\
     msg:'Test sec Rule',\
     logdata:'Matched Data: %%{TX.0} found within %%{MATCHED_VAR_NAME}: %%{MATCHED_VAR}',\
     tag:'attack-protocol',\
     tag:'traceable/labels/OWASP_2021:A4,CWE:444,OWASP_API_2019:API8',\
     tag:'traceable/severity/HIGH',\
     tag:'traceable/type/safe,block',\
     severity:'CRITICAL',\
     setvar:'tx.anomaly_score_pl1=+%%{tx.critical_anomaly_score}'"
     EOT
    req_res_conditions{
        match_key="HEADER_NAME"
        match_category="REQUEST"
        match_operator="EQUALS"
        match_value="req_header"
    }
}


data "traceable_endpoint_id" "endpoint" {
  name="POST /EOaml"
  service_name="service_UxDyUNPq"
  enviroment_name="jatinenv"
}

resource "traceable_rate_limiting_block" "sample_rule" {
    name="please one last time 1"
    environments=["utkarsh_21"]
    enabled=true
    alert_severity="HIGH"
    threshold_configs {
        api_aggregate_type="ACROSS_ENDPOINTS"
        rolling_window_count_allowed=10
        rolling_window_duration="PT60S"
        threshold_config_type="ROLLING_WINDOW"
    }
    threshold_configs {
        api_aggregate_type="PER_ENDPOINT"
        rolling_window_count_allowed=100
        rolling_window_duration="PT120S"
        threshold_config_type="ROLLING_WINDOW"
    }
    threshold_configs {
        api_aggregate_type="ACROSS_ENDPOINTS"
        rolling_window_count_allowed=100000
        rolling_window_duration="PT300S"
        threshold_config_type="DYNAMIC"
        dynamic_mean_calculation_duration="PT86400S"  
    }
    ip_address {
      ip_address_list = ["1.1.1.1","1.1.1.1/32"]
      exclude = true
    }
    regions{
      regions_ids = [ "AX", "DZ" ]
      exclude = true
    }
    user_id {
        user_id_regexes = ["sample","test"]
        exclude=false
    }
    req_res_conditions {
          metadata_type    = "HTTP_METHOD"
          req_res_operator = "EQUALS"
          req_res_value    = "zfcd"
        }
    label_id_scope = ["External"]
}

resource "traceable_detection_policies" "enablelfi"{
    config_name="typeAnomaly"
    disabled=true
}
resource "traceable_detection_policies" "enablelfi"{
    config_name="typeAnomaly"
    disabled=true
}

resource "traceable_custom_signature_allow" "cs_allow" {
  name = "tf-cs-allow-test-1"
  description = "test rule created from tf"
  environments = ["crapi-test"]
  disabled = true
  allow_expiry_duration="PT30M"
  custom_sec_rule=<<EOT
     SecRule REQUEST_HEADERS:key-sec "@rx val-sec" \
     "id:92100120,\
     phase:2,\
     block,\
     msg:'Test sec Rule',\
     logdata:'Matched Data: %%{TX.0} found within %%{MATCHED_VAR_NAME}: %%{MATCHED_VAR}',\
     tag:'attack-protocol',\
     tag:'traceable/labels/OWASP_2021:A4,CWE:444,OWASP_API_2019:API8',\
     tag:'traceable/severity/HIGH',\
     tag:'traceable/type/safe,block',\
     severity:'CRITICAL',\
     setvar:'tx.anomaly_score_pl1=+%%{tx.critical_anomaly_score}'"
  EOT
  req_res_conditions{
        match_key="HEADER_NAME"
        match_category="REQUEST"
        match_operator="EQUALS"
        match_value="req_header"
    }
  req_res_conditions{
        match_key="HEADER_NAME"
        match_category="REQUEST"
        match_operator="EQUALS"
        match_value="req_header_test"
    }
}

resource "traceable_custom_signature_alert" "cs_alert" {
  name = "tf-cs-alert-test-1"
  description = "test rule created from tf"
  environments = []
  disabled = true
  custom_sec_rule=<<EOT
     SecRule REQUEST_HEADERS:key-sec "@rx val-sec" \
     "id:92100120,\
     phase:2,\
     block,\
     msg:'Test sec Rule',\
     logdata:'Matched Data: %%{TX.0} found within %%{MATCHED_VAR_NAME}: %%{MATCHED_VAR}',\
     tag:'attack-protocol',\
     tag:'traceable/labels/OWASP_2021:A4,CWE:444,OWASP_API_2019:API8',\
     tag:'traceable/severity/HIGH',\
     tag:'traceable/type/safe,block',\
     severity:'CRITICAL',\
     setvar:'tx.anomaly_score_pl1=+%%{tx.critical_anomaly_score}'"
  EOT
  req_res_conditions{
        match_key="HEADER_NAME"
        match_category="REQUEST"
        match_operator="EQUALS"
        match_value="req_header"
    }
  req_res_conditions{
        match_key="HEADER_NAME"
        match_category="REQUEST"
        match_operator="EQUALS"
        match_value="req_header_test"
    }
    alert_severity = "HIGH"
}

resource "traceable_custom_signature_block" "cs_block" {
  name = "tf-cs-block-test-1"
  description = "test rule created from tf"
  environments = []
  disabled = true
  custom_sec_rule=<<EOT
     SecRule REQUEST_HEADERS:key-sec "@rx val-sec" \
     "id:92100120,\
     phase:2,\
     block,\
     msg:'Test sec Rule',\
     logdata:'Matched Data: %%{TX.0} found within %%{MATCHED_VAR_NAME}: %%{MATCHED_VAR}',\
     tag:'attack-protocol',\
     tag:'traceable/labels/OWASP_2021:A4,CWE:444,OWASP_API_2019:API8',\
     tag:'traceable/severity/HIGH',\
     tag:'traceable/type/safe,block',\
     severity:'CRITICAL',\
     setvar:'tx.anomaly_score_pl1=+%%{tx.critical_anomaly_score}'"
  EOT
  req_res_conditions{
        match_key="HEADER_NAME"
        match_category="REQUEST"
        match_operator="EQUALS"
        match_value="req_header"
    }
  req_res_conditions{
        match_key="HEADER_NAME"
        match_category="REQUEST"
        match_operator="EQUALS"
        match_value="req_header_test"
    }
    alert_severity = "HIGH"
}

resource "traceable_label_management_label_rule" "example_rule" {
  name        = "Example Label Application Rule"
  description = "An example rule for applying labels based on conditions"
  enabled     = true

  condition_list {
    key      = "request.header.content-type"
    operator = "OPERATOR_EQUALS"
    value    = "application/json"
  }

  condition_list {
    key      = "request.query.param"
    operator = "OPERATOR_EQUALS"
    values   = ["value1", "value2"]
  }

  condition_list {
    key      = "response.status"
    operator = "OPERATOR_EQUALS"
  }

  action {
    type         = "DYNAMIC_LABEL_KEY"
    entity_types = ["request", "response"]
    operation    = "OPERATION_MERGE"
    dynamic_label_key = "dynamic_label_key_example"
  }
}