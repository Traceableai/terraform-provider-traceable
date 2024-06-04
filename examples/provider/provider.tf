provider "traceable" {
  platform_url="https://api-dev.traceable.ai/graphql"
  api_token=jsondecode(data.aws_secretsmanager_secret_version.api_token.secret_string)["api_token"]
}