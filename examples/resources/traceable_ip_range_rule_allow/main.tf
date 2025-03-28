terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
      version = "0.0.1"
    }
  }
}

variable "platform_url" {
  type        = string
  description = "Traceable Platform URL"
}

variable "traceable_api_key" {
  type        = string
  description = "Traceable API Key"
  sensitive   = true
}

variable "name" {
  type = string
}

variable "description" {
  type    = string
  default = null
}

variable "rule_action" {
  type    = string
  default = "RULE_ACTION_ALLOW"
}

variable "expiration" {
  type    = string
  default = null
}

variable "environment" {
  type = set(string)
}

variable "raw_ip_range_data" {
  type = set(string)
}

variable "inject_request_headers" {
  type = list(object({
    header_key   = string
    header_value = string
  }))
  default = []
}

provider "traceable" {
  platform_url     = var.platform_url
  api_token        = var.traceable_api_key
  provider_version = "terraform/v1.0.1"
}

resource "traceable_ip_range_rule_allow" "sample_rule" {
  name              = var.name
  description       = var.description
  rule_action       = var.rule_action
  expiration        = var.expiration
  environment       = var.environment
  raw_ip_range_data = var.raw_ip_range_data

  # dynamic "inject_request_headers" {
  #   for_each = var.inject_request_headers
  #   content {
  #     header_key   = inject_request_headers.value.header_key
  #     header_value = inject_request_headers.value.header_value
  #   }
  # }
}

output "traceable_ip_range_rule_allow" {
  value = traceable_ip_range_rule_allow.sample_rule
}
