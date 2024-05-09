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
  platform_url="https://app-dev.traceable.ai/graphql"
  api_token="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkxPbDdCcnhCVzUweUYxVERNYWRpZyJ9.eyJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlc192MiI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvY3VzdG9tZXJfaWQiOiJiOWVkYzU2Yi1mMmNlLTQ3M2MtYmNkNC1hYjczZWRmN2M2MTAiLCJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlcyI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvanRpIjoiMDE2MmVhNWItZDIyNi00MGE2LTg4NGQtMjk1NjE1N2UxNzMwIiwiaHR0cHM6Ly90cmFjZWFibGUuYWkvcmljaF9yb2xlcyI6W3siZW52cyI6W10sImlkIjoidHJhY2VhYmxlIn1dLCJuaWNrbmFtZSI6ImN1c3RvbWVyK2I5ZWRjNTZiLWYyY2UtNDczYy1iY2Q0LWFiNzNlZGY3YzYxMCIsIm5hbWUiOiJjdXN0b21lcitiOWVkYzU2Yi1mMmNlLTQ3M2MtYmNkNC1hYjczZWRmN2M2MTBAdHJhY2VhYmxlLmFpIiwicGljdHVyZSI6Imh0dHBzOi8vcy5ncmF2YXRhci5jb20vYXZhdGFyLzU5MTI5NjA2NWY5ZjQxZTM3ZDdiMWMxMGVjYjA0N2EzP3M9NDgwJnI9cGcmZD1odHRwcyUzQSUyRiUyRmNkbi5hdXRoMC5jb20lMkZhdmF0YXJzJTJGY3UucG5nIiwidXBkYXRlZF9hdCI6IjIwMjQtMDUtMDlUMDk6MDY6MTQuMzY2WiIsImVtYWlsIjoiY3VzdG9tZXIrYjllZGM1NmItZjJjZS00NzNjLWJjZDQtYWI3M2VkZjdjNjEwQHRyYWNlYWJsZS5haSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczovL3RyYWNlYWJsZS1kZXYudXMuYXV0aDAuY29tLyIsImF1ZCI6IkloeHg4QVBlc3dEaUlrb3JxSzIzYkt0OHVUa0pkaDA3IiwiaWF0IjoxNzE1MjUxNjIyLCJleHAiOjE3MTUyODc2MjIsInN1YiI6ImF1dGgwfDYzMjBiMGNkOWE3NmVjNjZhYmJhN2Q4YiIsInNpZCI6ImpzeW5heTdyQm5LTV9YZnRvWjdyLTMxSUs2azl1czhIIn0.um5FLNZ3Ke_ZOtwOfFKwl3McxPJJZr0XHkC2Rx8SlN01WUvL5PoKpxxTKpEnHisDCCUHO-W2TlrjNiq0gAn56LYhtk3V6BW4yrwywdKhsI1Rm5GwuWd9kslZATfWEjsqlev0GKKCK11v4g2LmoiUAV3Nd033UBHZt85a1fQPC_pYS9wOW3ognmf_NhuOPxf4Mj1OAj5cGianQBlRU94qBBPYuy-Tdtpyo-nwmaHJLCBBCF96zy_kL1pZJ7yZ3cBWuVX-RKl6EcDsP7P0W0uwFfiXQ9N5wpd0fdiDsMIiyRkBJyAuMiFoWTt02DD4rHWlmCk087TYXS38Xf2tmietjA"
}

resource "example_api_exclusion_rule" "example_exclusion_rule" {
  name ="meg-test"
  disabled=true
  regexes="test/hello1"
  service_names=[]
  environment_names=[]
}

# resource "example_api_naming_rule" "example_naming_rule" {
#   name             = "meg-test-rule-all-env"
#   disabled         = false
#   regexes          = ["hello", "test", "123"]
#   values           = ["hello", "test", "number"]
#   service_names    = [""]
#   environment_names = [""]
# }


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

