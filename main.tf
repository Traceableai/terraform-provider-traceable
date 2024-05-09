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

resource "example_label_creation_rule" "example_label_create_rule" {
  key="meg-test-tf-2"
  description="testing rule for tf"
  color="#E295E9"
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
