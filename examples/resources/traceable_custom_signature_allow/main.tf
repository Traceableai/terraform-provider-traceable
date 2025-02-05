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
  description = "Name of Environemts on which to be enabled"
  type        = list(string)
}

variable "allow_expiry_duration" {
  description = "Description of the dataset"
  type        = string
}

# variable "custom_sec_rule" {
#   description = "Custom security rule"
#   type        = string
# }


variable "req_res_conditions" {
  description = "Request response conditions"
  type        = list(object({
    match_key      = string
    match_category = string
    match_operator = string
    match_value    = string
  }))
}



provider "traceable" {
  platform_url = var.platform_url
  api_token    = var.traceable_api_key
}

resource "traceable_custom_signature_allow" "cs_allow" {
  name                  = var.name
  description           = var.description
  environments          = var.environments
  allow_expiry_duration = var.allow_expiry_duration
  disabled = true

  dynamic "req_res_conditions" {
    for_each = var.req_res_conditions
    content {
      match_key      = req_res_conditions.value.match_key
      match_category = req_res_conditions.value.match_category
      match_operator = req_res_conditions.value.match_operator
      match_value    = req_res_conditions.value.match_value
    }
  }

}



output "custom_id" {
  value = traceable_custom_signature_allow.cs_allow.id
}


output "custom_name" {
  value = traceable_custom_signature_allow.cs_allow.name
}

output "custom_description" {
  value = traceable_custom_signature_allow.cs_allow.description
}

output "custom_environments" {
  value = traceable_custom_signature_allow.cs_allow.environments
}

output "custom_allow_expiry_duration" {
  value = traceable_custom_signature_allow.cs_allow.allow_expiry_duration
}

# output "custom_sec_rule" {
#   value = traceable_custom_signature_allow.cs_allow.custom_sec_rule
# }




