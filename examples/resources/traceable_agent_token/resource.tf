resource "traceable_agent_token" "example" {
  name = "tf-provider-token-testing"
}

output "agent_token" {
  value = traceable_agent_token.example.token
  sensitive = true
}

output "agent_token_creation_timestamp" {
  value = traceable_agent_token.example.creation_timestamp
}