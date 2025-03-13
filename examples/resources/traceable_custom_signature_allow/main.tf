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
  type        = string
}

variable "rule_type"{
  type = string
  default ="ALLOW"
}
variable "description" {
  type        = string
  default = null
}

variable "environments" {
  type        = list(string)
  default     = null
}

variable "request_payload_single_valued_conditions" {
  type        = list(object({
    match_category      = string
    match_key = string
    match_operator = string
    match_value = string
  }))
  default = []
}

variable "request_payload_multi_valued_conditions" {
  type        = list(object({
    match_category = string
    key_value_tag = string
    key_match_operator = string
    match_key = string
    value_match_operator = string
    match_value = string
  }))
  default = []
}

variable "custom_sec_rule"{
  type        = string
   default     = null

}

variable "disabled"{
  type = bool
  default     = true
}

variable "allow_expiry_duration"{
  type = string 
  default = null
}

provider "traceable" {
  platform_url = var.platform_url
  api_token    = var.traceable_api_key
   provider_version="terraform/v1.0.1"
}

resource "traceable_custom_signature_allow" "sample_rule" {
  name                  = var.name
  description           = var.description
  rule_type              = var.rule_type
  environments          = var.environments

 dynamic "request_payload_single_valued_conditions" {
  for_each = var.request_payload_single_valued_conditions
  content {
    match_category = request_payload_single_valued_conditions.value.match_category
   match_key   = request_payload_single_valued_conditions.value.match_key
   match_operator = request_payload_single_valued_conditions.value.match_operator
     match_value =  request_payload_single_valued_conditions.value.match_value
  }
}
dynamic "request_payload_multi_valued_conditions" {
  for_each = var.request_payload_multi_valued_conditions
  content {
    match_category   = request_payload_multi_valued_conditions.value.match_category
    key_value_tag      = request_payload_multi_valued_conditions.value.key_value_tag
    key_match_operator = request_payload_multi_valued_conditions.value.key_match_operator
    match_key    = request_payload_multi_valued_conditions.value.match_key
    value_match_operator    = request_payload_multi_valued_conditions.value.value_match_operator
    match_value    = request_payload_multi_valued_conditions.value.match_value
  }
}

custom_sec_rule=var.custom_sec_rule
disabled= var.disabled
allow_expiry_duration = var.allow_expiry_duration
}

output "traceable_custom_signature_allow" {
  value = traceable_custom_signature_allow.sample_rule
}


