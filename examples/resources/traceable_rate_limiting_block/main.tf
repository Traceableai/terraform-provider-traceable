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

variable "description" {
  type        = string
  default     = null

}
variable "alert_severity" {
  type        = string
}
variable "enabled" {
  type        = bool
}
variable "block_expiry_duration" {
  type        = string
  default     = null
}
variable "label_id_scope" {
  type        = list(string)
  default     = null
}
variable "endpoint_id_scope" {
  type        = list(string)
  default     = null
}
variable "environments" {
  type        = list(string)
  default     = null
}


variable "request_response_single_valued_conditions" {
  type        = list(object({
    request_location      = string
    operator = string
    value = string
  }))
  default = []
}

variable "request_response_multi_valued_conditions" {
  type        = list(object({
    request_location      = string
    key_patterns = list(object({
          operator      = string
          value      = string
    }))
     value_patterns = list(object({
          operator      = string
          value      = string
    }))
  }))
  default = []
}

variable "threshold_configs" {
  type        = list(object({
    api_aggregate_type      = string
    user_aggregate_type      = string
    rolling_window_count_allowed      = string
    rolling_window_duration        = string
    threshold_config_type   = string
    dynamic_mean_calculation_duration  = string 

  }))
  default = []
}
variable "attribute_based_conditions" {
  type        = list(object({
    key_condition_operator      = string
    key_condition_value      = string
    value_condition_operator     = string
    value_condition_value        = string
  }))
  default = []
}

// 12 sources 

variable "ip_reputation" {
  type        = string
  default     = null
}
variable "ip_location_type" {
  type        = list(object({
   ip_location_types      = list(string)
   exclude    = bool
  }))
  default = []
}
variable "ip_abuse_velocity" {
  type        = string
  default     = null
}
variable "ip_address"{
  type        = list(object({
   ip_address_list      = list(string)
   exclude    = bool
  }))
  default = []
}
variable "email_domain" {
  type        = list(object({
  email_domain_regexes      = list(string)
   exclude    = bool
  }))
  default = []
}
variable "user_agents" {
  type        = list(object({
   user_agents_list      = list(string)
   exclude    = bool
  }))
  default = []
}
variable "regions" {
  type        = list(object({
   regions_ids      = list(string)
   exclude    = bool
  }))
  default = []
}
variable "ip_organisation" {
  type        = list(object({
   ip_organisation_regexes      = list(string)
   exclude    = bool
  }))
  default = []
}
variable "ip_asn" {
  type        = list(object({
   ip_asn_regexes     = list(string)
   exclude    = bool
  }))
  default = []
}
variable "ip_connection_type" {
  type        = list(object({
  ip_connection_type_list     = list(string)
   exclude    = bool
  }))
  default = []
}
variable "request_scanner_type" {
  type        = list(object({
   scanner_types_list     = list(string)
   exclude    = bool
  }))
  default = []
}
variable "user_id" {
  type        = list(object({
  user_id_regexes     = list(string)
   exclude    = bool
   user_ids  = list(string)
     }))
  default = []
}

provider "traceable" {
  platform_url = var.platform_url
  api_token    = var.traceable_api_key
   provider_version="terraform/v1.0.1"
}

resource "traceable_rate_limiting_block" "sample_rule" {
    name=var.name
    environments=var.environments
    enabled=var.enabled
    alert_severity=var.alert_severity
dynamic "threshold_configs" {
  for_each = var.threshold_configs
  content {
    api_aggregate_type      = threshold_configs.value.api_aggregate_type
    user_aggregate_type     = threshold_configs.value.user_aggregate_type
    rolling_window_duration = threshold_configs.value.rolling_window_duration
    rolling_window_count_allowed = threshold_configs.value.rolling_window_count_allowed
    dynamic_mean_calculation_duration = threshold_configs.value.dynamic_mean_calculation_duration
    threshold_config_type = threshold_configs.value.threshold_config_type
  }
}

    
     dynamic "ip_address" {
    for_each = var.ip_address
    content {
      ip_address_list = ip_address.value.ip_address_list
      exclude         = ip_address.value.exclude
    }
    
  }
    dynamic "ip_location_type" {
    for_each = var.ip_location_type
    content {
      ip_location_types = ip_location_type.value.ip_location_types
      exclude         = ip_location_type.value.exclude
    }
    
  }
    dynamic "regions" {
    for_each = var.regions
    content {
      regions_ids   = regions.value.regions_ids 
      exclude         = regions.value.exclude
    }
    
  }
   dynamic "ip_organisation" {
    for_each = var.ip_organisation
    content {
      ip_organisation_regexes =  ip_organisation.value.ip_organisation_regexes
      exclude         = ip_organisation.value.exclude
    }
    
  }
   dynamic "ip_asn" {
    for_each = var.ip_asn
    content {
      ip_asn_regexes =   ip_asn.value.ip_asn_regexes
      exclude         = ip_asn.value.exclude
    }
    
  }
     dynamic "ip_connection_type" {
    for_each = var.ip_connection_type
    content {
      ip_connection_type_list =   ip_connection_type.value.ip_connection_type_list
      exclude         = ip_connection_type.value.exclude
    }
    
  }
   dynamic "request_scanner_type" {
    for_each = var.request_scanner_type
    content {
      scanner_types_list =  request_scanner_type.value.scanner_types_list
      exclude         =  request_scanner_type.value.exclude
    }
    
  }
  dynamic "user_id" {
    for_each = var.user_id
    content {
      user_id_regexes =  user_id.value.user_id_regexes 
      user_ids =  user_id.value.user_ids 
      exclude         = user_id.value.exclude
    }
    
  }
   dynamic "user_agents" {
    for_each = var.user_agents
    content {
      user_agents_list =    user_agents.value.user_agents_list
      exclude         = user_agents.value.exclude
    }
    
  }

   dynamic "request_response_single_valued_conditions" {
    for_each = var.request_response_single_valued_conditions
    content {
    request_location  = request_response_single_valued_conditions.value.request_location
      operator        =request_response_single_valued_conditions.value.operator
      value =request_response_single_valued_conditions.value.value
    }
    
  }

  dynamic "request_response_multi_valued_conditions" {
  for_each = var.request_response_multi_valued_conditions
  content {
    request_location = request_response_multi_valued_conditions.value.request_location

    dynamic "key_patterns" {
      for_each = request_response_multi_valued_conditions.value.key_patterns
      content {
        operator = key_patterns.value.operator
        value    = key_patterns.value.value
      }
    }

    dynamic "value_patterns" {
      for_each = request_response_multi_valued_conditions.value.value_patterns
      content {
        operator = value_patterns.value.operator
        value    = value_patterns.value.value
      }
    }
  }
}
dynamic "attribute_based_conditions" {
  for_each = var.attribute_based_conditions
  content {
    key_condition_operator   = attribute_based_conditions.value.key_condition_operator
    key_condition_value      = attribute_based_conditions.value.key_condition_value
    value_condition_operator = attribute_based_conditions.value.value_condition_operator
    value_condition_value    = attribute_based_conditions.value.value_condition_value
  }
}
dynamic "email_domain" {
  for_each = var.email_domain
  content {
    email_domain_regexes =    email_domain.value.email_domain_regexes 
    exclude = email_domain.value.exclude
  
  }
}

  ip_reputation = var.ip_reputation
  ip_abuse_velocity = var.ip_abuse_velocity
  endpoint_id_scope = var.endpoint_id_scope
  label_id_scope = var.label_id_scope

   
}

output "traceable_rate_limiting_block"{
  value = traceable_rate_limiting_block.sample_rule
}






















