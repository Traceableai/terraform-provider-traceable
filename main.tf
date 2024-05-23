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
  platform_url="https://api-dev.traceable.ai/graphql"
  api_token="ZjRmOGJmNTktYmY4Ni00M2RiLThlNTItYjA5Zjc4ZDgxYWIx"
}

resource "traceable_label_application_rule" "meg-test1" {
  name        = "newtest-meg-tf"
  description = "optional desc"
  enabled     = true

  condition_list {
    key      = "test1"
    operator = "OPERATOR_EQUALS"
    value    = "sdfghj"
  }

  condition_list {
    key      = "test2"
    operator = "OPERATOR_MATCHES_REGEX"
    value    = "sdfg"
  }

  condition_list {
    key      = "test3"
    operator = "OPERATOR_EQUALS"
  }

  condition_list {
    key      = "test4"
    operator = "OPERATOR_MATCHES_IPS"
    values   = ["1.0.0.0", "9.0.8.0"]
  }

  condition_list {
    key      = "test5"
    operator = "OPERATOR_NOT_MATCHES_IPS"
    values   = ["6.9.3.5", "5.5.6.5"]
  }

  action {
    type             = "DYNAMIC_LABEL_KEY"
    entity_types     = ["API", "SERVICE", "BACKEND"]
    operation        = "OPERATION_MERGE"
    dynamic_label_key = "hellotest"
  }

  
}
# resource "traceable_user_attribution_rule_basic_auth" "test1" {
#   name = "aditya-2"
#   scope_type = "SYSTEM_WIDE"
# }

# resource "traceable_ip_range_rule" "my_ip_range" {
#     name     = "first_rule_2"
#     rule_action     = "RULE_ACTION_ALERT"
#     event_severity     = "LOW"
#     raw_ip_range_data = [
#         "1.1.1.1",
#         "3.3.3.3"
#     ]
#     environment=[] #all env
#     description="rule created from custom provider"
# }

# resource "traceable_user_attribution_rule_basic_auth" "test1" {
#   name = "traceable_user_attribution_rule_basic_auth"
#   scope_type = "SYSTEM_WIDE"
#   url_regex = "abcd"
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

# # resource "traceable_user_attribution_rule_custom_token" "test6" {
# #   name = "traceable_user_attribution_rule_custom_token"
# #   scope_type="SYSTEM_WIDE"
# #   auth_type="test"
# #   location="REQUEST_COOKIE"
# #   token_name="test"
# }

# resource "traceable_api_naming_rule" "example_naming_rule" {
#   name             = "test-rule-naming"
#   disabled         = false
#   regexes          = ["hello", "test", "123"]
#   values           = ["hello", "test", "number"]
#   service_names    = [""]
#   environment_names = [""]
# }

# resource "traceable_api_exclusion_rule" "example_exclusion_rule" {
#   name =  "test-rule-exclusion"
#   disabled= true
#   regexes=  "hello/test/6785"
#   service_names=  [""]
#   environment_names=  [""]
# # }

# data "traceable_syslog_integration" "syslog" {
#   name="prer-test"
# }
# data "traceable_endpoint_id" "endpooint" {
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