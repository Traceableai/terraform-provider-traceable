terraform {
  required_providers {
    example = {
      source  = "terraform.local/local/example"
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

provider "example" {
  platform_url="https://api-dev.traceable.ai/graphql"
  api_token=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
}

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

resource "traceable_user_attribution_rule_custom_token" "test6" {
  name = "traceable_user_attribution_rule_custom_token"
  scope_type="SYSTEM_WIDE"
  auth_type="test"
  location="REQUEST_COOKIE"
  token_name="test"
}

# resource "traceable_api_naming_rule" "example_naming_rule" {
#   name             = "test-rule-naming"
#   disabled         = false
#   regexes          = ["hello", "test", "123"]
#   values           = ["hello", "test", "number"]
#   service_names    = ["example-svc"]
#   environment_names = ["example-env"]
# }

# resource "traceable_api_exclusion_rule" "example_exclusion_rule" {
#   name =  "test-rule-exclusion"
#   disabled= true
#   regexes=  "hello/test/6785"
#   service_names=  ["example-svc"]
#   environment_names=  ["example-env"]
# }