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
  description = "Name of the dataset"
  type        = string
}
variable "description" {
  description = "Description of the dataset"
  type        = string
}
variable "environments" {
  description = "Name of Environments on which to be enabled"
  type        = list(string)
}
variable "expiration" {
  description = "Description of the dataset"
  type        = string
}
variable "rule_action" {
  description = "Set to default for (RULE_ACTION_ALLOW)"
  type        = string
  default     = "RULE_ACTION_ALLOW"
}
variable "event_severity" {
  description = "Description of the dataset"
  type        = string
  sensitive   = true
}
variable "raw_ip_range_data" {
  description = "IPV4/V6 range for the rule"
  type        = set(string)
  sensitive   = true
}
provider "traceable" {
  platform_url     = var.platform_url
  api_token        = var.traceable_api_key
  provider_version = "terraform/v1.0.1"
}

resource "traceable_ip_range_rule_block" "sample_rule" {
  name               = var.name
  description        = var.description
  rule_action        = var.rule_action
  event_severity     = var.event_severity
  expiration         = var.expiration
  environment        = var.environments
  raw_ip_range_data  = var.raw_ip_range_data
}

output "traceable_ip_range_rule_block" {
  value = traceable_ip_range_rule_block.sample_rule
  sensitive = true
}