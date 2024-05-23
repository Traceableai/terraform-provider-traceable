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
  platform_url="https://app-dev.traceable.ai/graphql"
  # api_token=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
  api_token="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkxPbDdCcnhCVzUweUYxVERNYWRpZyJ9.eyJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlc192MiI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvY3VzdG9tZXJfaWQiOiI3Njc4NjA2Zi1jZTEwLTQyYjQtYWU1MC1mOTlkYzA5NzcxMWMiLCJodHRwczovL3RyYWNlYWJsZS5haS9yb2xlcyI6WyJ0cmFjZWFibGUiXSwiaHR0cHM6Ly90cmFjZWFibGUuYWkvanRpIjoiMWM2ZTU5NjMtZDY2ZC00NjBmLWE1MmMtM2JlNTEyODc5NjA4IiwiaHR0cHM6Ly90cmFjZWFibGUuYWkvcmljaF9yb2xlcyI6W3siZW52cyI6W10sImlkIjoidHJhY2VhYmxlIn1dLCJuaWNrbmFtZSI6InRyYWNlYWJsZV9pbnRlcm5hbF90ZXN0aW5nX3RyYWNlYWJsZSIsIm5hbWUiOiJ0cmFjZWFibGVfaW50ZXJuYWxfdGVzdGluZ190cmFjZWFibGVAdHJhY2VhYmxlLmFpIiwicGljdHVyZSI6Imh0dHBzOi8vcy5ncmF2YXRhci5jb20vYXZhdGFyLzIyZTA5NGQzNTFkOGM3NzNiZjU3NGQ4YTUwM2UwN2QwP3M9NDgwJnI9cGcmZD1odHRwcyUzQSUyRiUyRmNkbi5hdXRoMC5jb20lMkZhdmF0YXJzJTJGdHIucG5nIiwidXBkYXRlZF9hdCI6IjIwMjQtMDUtMjNUMTA6MDE6NTQuNzY5WiIsImVtYWlsIjoidHJhY2VhYmxlX2ludGVybmFsX3Rlc3RpbmdfdHJhY2VhYmxlQHRyYWNlYWJsZS5haSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczovL3RyYWNlYWJsZS1kZXYudXMuYXV0aDAuY29tLyIsImF1ZCI6IkloeHg4QVBlc3dEaUlrb3JxSzIzYkt0OHVUa0pkaDA3IiwiaWF0IjoxNzE2NDU4NTE2LCJleHAiOjE3MTY0OTQ1MTYsInN1YiI6ImF1dGgwfDYzODJmYmQxYjE0MzM5NDU4ZTlmZjVkOCIsInNpZCI6ImR5SXk4OHZaVDVkWV9qa09oVk9TZFhwRTVNVHZCVjgyIn0.DdcMbiWdhQiegYaC3NHO4NX0WhWwZGzAHIbEFBba2HHhdQcd1nWrPAWXyVvFRxE_EljLlxMqmLY1EjECvholaNgAJY0Q4hsJZtJ3JaqJV_a3A3y7gvzjj9v08o6pQH3od6ZqNZCZnbnavJwGvRTXXQW4lhgoz7RFw5K9HFmqsjWjrOOT1IG7S6S-B-n5sqrLzEOvntQ6f17qSaB1ZOSin5N1Zquz2zLk4a4gWNYXUz087kCOCsLwf-Ge96ZonCb9XTudqkr-8MKMMPxT6mZNsvzR202cNtGg2mX5QITlLMp9N4nAVEWv1gO3KP28sRbnYbYqe_hkkqvlJzd4KuSG4w"
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

# resource "traceable_user_attribution_rule_custom_token" "test6" {
#   name = "traceable_user_attribution_rule_custom_token"
#   scope_type="SYSTEM_WIDE"
#   auth_type="test"
#   location="REQUEST_COOKIE"
#   token_name="test"
# }

# data "traceable_syslog_integration" "syslog" {
#   name="prer-test"
# }
# data "traceable_splunk_integration" "splunk" {
#   name="aditya"
# }
# data "traceable_endpoint_id" "endpoint" {
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

resource "traceable_notification_channel" "testchannel" {
  channel_name = "example_channel1"

  email = [
    "example4@example.com",
    "example2@example.com"
  ]

  slack_webhook = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
  splunk_id="7d5c77ab-1039-4123-af32-e68f116b2b83"
  # syslog_id=""
  custom_webhook  {
    webhook_url = "https://example.com/webhook"
    custom_webhook_headers  {
      key       = "Authorization"
      value     = "Bearer token123"
      is_secret = true
    }
    custom_webhook_headers  {
      key       = "Authorization1"
      value     = "Bearer token1232"
      is_secret = false
    }
  }

  s3_webhook  {
    bucket_name = "your-s3-bucket"
    region      = "us-west-2"
    bucket_arn  = "arn:aws:s3:::your-s3-bucket"
  }
}