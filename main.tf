terraform {
  required_providers {
    traceable = {
      source = "registry.terraform.io/traceableai/traceable"
      version = "0.0.1"
    }
  }
}

variable "API_TOKEN" {
  type = string
}

provider "traceable" {
  platform_url="https://api.traceable.ai"
  api_token = var.API_TOKEN
}

resource "traceable_api_naming" "test" {
  name = "adityatf27-mvk"
  disabled = false
  service_names=["nginx"]
  environment_names=[]
  values=["someval"]
  regexes=["nginx-traceshop"]
}

resource "traceable_malicious_ip_range" "block_sample1"{
    name = "traceable-source-mvk1"
    description = "traceable-des1"
    enabled = true
    event_severity = "LOW"
    duration = "PT1M"
    action = "BLOCK"
    ip_range = ["192.168.22.1","192.168.33.2"]
    environments = ["env1","env2"]
}