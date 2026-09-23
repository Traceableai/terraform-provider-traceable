terraform {
  required_providers {
    traceable = {
      source  = "Traceableai/traceable"
      version = "1.0.12"
    }
  }
}

variable "API_TOKEN" {
  type = string
}

provider "traceable" {
  platform_url = "https://api.traceable.ai"
  api_token    = var.API_TOKEN
}
