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

variable "rule_type" {
  type= string

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
variable "expiry_duration" {
  type        = string
  default     = null
}
variable "environments" {
  type        = list(string)
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
variable "data_location"{
  type        = list(object({
    data_type_ids = string
  }))
  default = []
}


variable "request_payload_single_valued_conditions" {
  type        = list(object({
    request_location      = string
    operator = string
    value = string
  }))
  default = []
}

variable "request_payload_multi_valued_conditions" {
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


variable "rolling_window_threshold_config" {
  type        = list(object({
    api_aggregate_type      = string
    user_aggregate_type      = string
    count_allowed = number
    duration = string
  
  }))
  default = []
}
variable "dynamic_threshold_config" {
  type        = list(object({
    api_aggregate_type      = string
    user_aggregate_type      = string
    percentage_exceeding_mean_allowed=number
    mean_calculation_duration=string
    duration=string
  }))
  default = []
}
variable "value_based_threshold_config" {
  type        = list(object({
      api_aggregate_type      = string
    user_aggregate_type      = string
    unique_values_allowed =number
    duration=string
    sensitive_params_evaluation_type=string
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
     }))
  default = []
}
variable data_set_name{
  type = string
}


provider "traceable" {
  platform_url = var.platform_url
  api_token    = var.traceable_api_key
   provider_version="terraform/v1.0.1"
}


data "traceable_data_set_id" "datasetId" {
  name = var.data_set_name
}


resource "traceable_dlp_user_based" "sample_rule"{
   rule_type = var.rule_type
    name=var.name
    description=var.description
    alert_severity=var.alert_severity
    enabled=var.enabled
    expiry_duration=var.expiry_duration
    endpoint_id_scope = var.endpoint_id_scope
    label_id_scope=var.label_id_scope
    environments=var.environments
    dynamic "request_response_single_valued_conditions" {
    for_each = var.request_payload_single_valued_conditions
    content {
    request_location  = request_response_single_valued_conditions.value.request_location
      operator        =request_response_single_valued_conditions.value.operator
      value =request_response_single_valued_conditions.value.value
    } 
  }
    dynamic "request_response_multi_valued_conditions" {
  for_each = var.request_payload_multi_valued_conditions
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
dynamic "data_types_conditions" {
  for_each = var.data_location
  content {
  data_location=   data_types_conditions.value.data_type_ids
   data_type_ids = [data.traceable_data_set_id.datasetId.id]
  }
}
 threshold_configs {

  dynamic "rolling_window_threshold_config"{
   for_each = var.rolling_window_threshold_config
   content{
    user_aggregate_type = rolling_window_threshold_config.value.user_aggregate_type
    api_aggregate_type = rolling_window_threshold_config.value.api_aggregate_type
    count_allowed = rolling_window_threshold_config.value.count_allowed
    duration = rolling_window_threshold_config.value.duration  
   }
  }

  dynamic "dynamic_threshold_config" {
  for_each = var.dynamic_threshold_config
  content {
    api_aggregate_type      = dynamic_threshold_config.value.api_aggregate_type
    user_aggregate_type     = dynamic_threshold_config.value.user_aggregate_type
    percentage_exceeding_mean_allowed = dynamic_threshold_config.value.percentage_exceeding_mean_allowed
    mean_calculation_duration= dynamic_threshold_config.value.mean_calculation_duration
    duration = dynamic_threshold_config.value.duration
  }
 }

    dynamic "value_based_threshold_config" {
  for_each = var.value_based_threshold_config
  content {
       api_aggregate_type   = value_based_threshold_config.value.api_aggregate_type
       user_aggregate_type     = value_based_threshold_config.value.user_aggregate_type
       unique_values_allowed = value_based_threshold_config.value.unique_values_allowed
       duration = value_based_threshold_config.value.duration
       sensitive_params_evaluation_type =value_based_threshold_config.value.sensitive_params_evaluation_type
  }
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
  dynamic "email_domain" {
  for_each = var.email_domain
  content {
    email_domain_regexes =    email_domain.value.email_domain_regexes 
    exclude = email_domain.value.exclude
  
  }
 
  }

  ip_reputation = var.ip_reputation
  ip_abuse_velocity = var.ip_abuse_velocity
  depends_on = [ data.traceable_data_set_id.datasetId ]


}



  

output "traceable_dlp_user_based"{
  value = traceable_dlp_user_based.sample_rule
}






















