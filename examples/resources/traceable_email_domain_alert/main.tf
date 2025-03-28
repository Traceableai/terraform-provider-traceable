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

variable "data_leaked_email" {
  type = bool
}

variable "disposable_email_domain" {
  type = bool
}

variable "email_domains" {
  type = set(string)
}

variable "email_regexes" {
  type = set(string)
}

variable "email_fraud_score" {
  type = string
}

provider "traceable" {
  platform_url     = var.platform_url
  api_token        = var.traceable_api_key
  provider_version = "terraform/v1.0.1"
}

resource "traceable_email_domain_alert" "sample_rule" {
  name                    = var.name
  description             = var.description
  rule_action             = var.rule_action
  event_severity          = var.event_severity
  environment             = var.environment
  data_leaked_email       = var.data_leaked_email
  disposable_email_domain = var.disposable_email_domain
  email_domains           = var.email_domains
  email_regexes           = var.email_regexes
  email_fraud_score       = var.email_fraud_score
}

output "traceable_email_domain_alert" {
  value = traceable_email_domain_alert.sample_rule
}
