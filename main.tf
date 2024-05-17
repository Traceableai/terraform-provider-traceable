terraform {
  required_providers {
    example = {
      source  = "terraform.local/local/example"
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
#   secret_id = "you secret manager arn"
# }

# output "api_token" {
#   value=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
#   sensitive = true
# }

provider "example" {
  platform_url=""
  api_token=""
}

resource "example_label_application_rule" "example_label_application" {
  name        = "testrule-1"
  description = "This is a test rule for applying labels based on conditions."
  enabled     = true

  condition_list {
    key_condition {
      operator = "OPERATOR_EQUALS"
      value    = "dfghj"
    }
    value_condition {
      value_condition_type = "STRING_CONDITION"
      string_condition {
        operator                   = "OPERATOR_EQUALS"
        string_condition_value_type = "VALUE"
        value                     = "dfgh"
      }
    }
  }

  condition_list {
    key_condition {
      operator = "OPERATOR_EQUALS"
      value    = "fgh"
    }
    value_condition {
      value_condition_type = "STRING_CONDITION"
      string_condition {
        operator                   = "OPERATOR_MATCHES_REGEX"
        string_condition_value_type = "VALUE"
        value                     = "sdfgh"
      }
    }
  }

  condition_list {
    key_condition {
      operator = "OPERATOR_EQUALS"
      value    = "gsdu"
    }
    value_condition {
      value_condition_type = "UNARY_CONDITION"
      unary_condition {
        operator = "OPERATOR_EXISTS"
      }
    }
  }

  condition_list {
    key_condition {
      operator = "OPERATOR_EQUALS"
      value    = "wrrtxxg"
    }
    value_condition {
      value_condition_type = "STRING_CONDITION"
      string_condition {
        operator                   = "OPERATOR_MATCHES_IPS"
        string_condition_value_type = "VALUES"
        values                     = ["4.9.07.9"]
      }
    }
  }

  condition_list {
    key_condition {
      operator = "OPERATOR_EQUALS"
      value    = "dtcg"
    }
    value_condition {
      value_condition_type = "STRING_CONDITION"
      string_condition {
        operator                   = "OPERATOR_NOT_MATCHES_IPS"
        string_condition_value_type = "VALUES"
        values                     = ["9.0.9.8"]
      }
    }
  }

  action {
    type            = "DYNAMIC_LABEL_KEY"
    entity_types    = ["API", "SERVICE", "BACKEND"]
    operation       = "OPERATION_MERGE"
    dynamic_label_key = "testtf"
  }
}


# resource "example_ip_range_rule" "my_ip_range" {
#     name     = "first_rule_2"
#     rule_action     = "RULE_ACTION_ALERT"
#     event_severity     = "LOW"
#     raw_ip_range_data = [
#         "1.1.1.1",
#         "3.3.3.3"
#     ]
#     environment=[]
#     expiration = "PT600S"
#     description="rule created from custom provider"
# }
# resource "example_rate_limit_rule" "my_rate_limit" {
#     name     = "first_rule_2"
#     rule_action     = "RULE_ACTION_ALLOW"
#     event_severity     = "LOW"
#     raw_ip_range_data = [
#         "1.1.1.1",
#         "3.3.3.3"
#     ]
#     conditions=[
#       {
#         name     = "server1"
#         size     = "t2.micro"
#         location = "us-east-1"
#       },
#       {
#         type     = "server1"
#         size     = "t2.micro"
#         location = "us-east-1"
#       }
#     ]
#     # expiration = "PT600S"
#     description="rule created from custom provider"
# }
