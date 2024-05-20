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
