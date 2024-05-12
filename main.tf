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
  platform_url="cluster_url"
  api_token=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
}


resource "example_api_naming_rule" "example_naming_rule" {
  name             = "meg-test-rule-all-env"
  disabled         = true
  regexes          = ["hello", "test", "123467"]
  values           = ["hello", "test", "testing"]
  service_names    = [""]
  environment_names = [""]
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

