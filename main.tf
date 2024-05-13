terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
      version = "0.0.1"
    }
  }
}

provider "traceable" {
  platform_url="https://api-dev.traceable.ai/graphql"
  api_token="OTk3ZDA2NWMtNWEyYS00NjQyLWFkNGItOTYzNWExZGU5YjEy"
}

resource "traceable_user_attribution_rule_basic_auth" "test1" {
  name = "aditya-2"
  scope_type = "SYSTEM_WIDE"
}

resource "traceable_user_attribution_rule_req_header" "test2" {
  name = "aditya-3"
  scope_type = "SYSTEM_WIDE"
  auth_type = "test"
  user_id_location = "test"
  user_role_location="yes"
  role_location_regex_capture_group="test"
}