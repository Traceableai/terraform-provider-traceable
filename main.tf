terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
      version = "0.0.1"
    }
  }
}


provider "traceable" {
  platform_url="platform graphql url"
  api_token="platform api token"
}

resource "traceable_user_attribution_rule_basic_auth" "test1" {
  name = "aditya-2"
  scope_type = "SYSTEM_WIDE"
}