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

variable "environment" {
  type = string
}

variable "waap_config" {
  type = list(object({
    rule_id        = string
    rule_config = optional(list(object({
      disabled = bool
    })), [])  
    subrule_config = optional(list(object({
      sub_rule_id    = string
      sub_rule_action = string
    })), [])  
  }))
  default = []
}

provider "traceable" {
  platform_url     = var.platform_url
  api_token        = var.traceable_api_key
  provider_version = "terraform/v1.0.1"
}

resource "traceable_detection_policies" "sample_rule" {
  environment = var.environment

  dynamic "waap_config" {
    for_each = var.waap_config
    content {
      rule_id = waap_config.value.rule_id

      dynamic "rule_config" {
        for_each = waap_config.value.rule_config
        content {
          disabled = rule_config.value.disabled
        }
      }

      dynamic "subrule_config" {
        for_each = waap_config.value.subrule_config
        content {
          sub_rule_id    = subrule_config.value.sub_rule_id
          sub_rule_action = subrule_config.value.sub_rule_action
        }
      }
    }
  }
}

output "traceable_detection_policies" {
  value = {
    id          = traceable_detection_policies.sample_rule.id
    environment = traceable_detection_policies.sample_rule.environment
  }
}