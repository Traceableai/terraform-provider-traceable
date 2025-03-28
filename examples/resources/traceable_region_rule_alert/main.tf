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
  default = "ALERT"
}

variable "event_severity" {
  type = string
}

variable "environment" {
  type = set(string)
}

variable "regions" {
  type = set(string)
}

provider "traceable" {
  platform_url     = var.platform_url
  api_token        = var.traceable_api_key
  provider_version = "terraform/v1.0.1"
}

resource "traceable_region_rule_alert" "sample_rule" {
  name           = var.name
  description    = var.description
  rule_action    = var.rule_action
  event_severity = var.event_severity
  environment    = var.environment
  regions        = var.regions
}

output "traceable_region_rule_alert" {
  value = traceable_region_rule_alert.sample_rule
}
