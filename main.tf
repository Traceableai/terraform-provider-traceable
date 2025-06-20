terraform {
  required_providers {
    traceable = {
      source = "terraform.local/local/traceable"
      version = "0.0.1"
    }
  }
}

variable "API_TOKEN" {
  type = string
}

provider "traceable" {
  platform_url="https://api-dev.traceable.ai"
  api_token = var.API_TOKEN
}

resource "traceable_api_naming" "test" {
  name = "adityatf27"
  disabled = false
  service_names=["nginx"]
  environment_names=[]
  values=["someval"]
  regexes=["nginx-traceshop"]
}