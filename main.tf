terraform {
  required_providers {
    example = {
      source  = "terraform.local/local/example"
      version = "0.0.1"
    }
  }
}

provider "example" {
  platform_url="https://api-dev.traceable.ai/graphql"
  api_token="NWFkNjc5NWItNDA3Ny00MGQwLWE1MjMtZTc4ODYzMWRlMDZm"
}

resource "example_ip_range_rule" "my_ip_range" {
    name     = "first_rule_2"
    rule_action     = "RULE_ACTION_ALLOW"
    event_severity     = "LOW"
    raw_ip_range_data = [
        "1.1.1.1",
        "3.3.3.3"
    ]
    environment = ["meghana-mock-env"]
    # expiration = "PT600S"
    description="rule created from custom provider"
}
