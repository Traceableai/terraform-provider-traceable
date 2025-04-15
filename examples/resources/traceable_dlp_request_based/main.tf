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


variable "data_location"{
  type        = list(object({
    data_location = string
  }))
  default = []
}



variable "target_scope"{
   type        = list(object({
    service_ids = list(string)
    url_regex= list(string)
  }))
    default = []
} 


variable "regions" {
  type =list(string)
  default = []
}

variable "ip_location_type" {
  type=list(string)
  default = []
}
variable "ip_address"{
 type = list(string)
  default = []
}

variable "data_types_conditions"{
   type        = list(object({
   custom_location_data_type_key_value_pair_matching = bool
   custom_location_attribute = string
   custom_location_attribute_key_operator = string
   custom_location_attribute_value = string
  }))



}

variable "data_set_name"{
  type = string
}

variable "service_name"{
  type = string
}

variable "url_regexes"{
  type =list(string)
}


provider "traceable" {
  platform_url = var.platform_url
  api_token    = var.traceable_api_key
   provider_version="terraform/v1.0.1"
}


data "traceable_data_set_id" "datasetId" {
  name = var.data_set_name
}
data "traceable_service_id" "serviceId" {
  service_name = var.service_name
  enviroment_name = var.environments[0]
 
}

resource "traceable_dlp_request_based" "sample_rule" {
   rule_type = var.rule_type
    name=var.name
    environments=var.environments
    enabled=var.enabled
    alert_severity=var.alert_severity
   ip_address = var.ip_address
  regions = var.regions

   ip_location_type = var.ip_location_type 

    dynamic data_types_conditions {
      for_each = var.data_types_conditions
      content{
        data_types {
          data_type_ids=[data.traceable_data_set_id.datasetId.id]
        }
        custom_location_data_type_key_value_pair_matching = data_types_conditions.value.custom_location_data_type_key_value_pair_matching
        custom_location_attribute = data_types_conditions.value.custom_location_attribute
        custom_location_attribute_key_operator =data_types_conditions.value.custom_location_attribute_key_operator
        custom_location_attribute_value = data_types_conditions.value.custom_location_attribute_value
      }
    }
   target_scope {
      
        service_ids = [data.traceable_service_id.serviceId.service_id]
        url_regex = var.url_regexes
      }


  dynamic "request_payload_multi_valued_conditions" {
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
    dynamic "request_payload_single_valued_conditions" {
    for_each = var.request_payload_single_valued_conditions
    content {
    request_location  = request_response_single_valued_conditions.value.request_location
      operator        =request_response_single_valued_conditions.value.operator
      value =request_response_single_valued_conditions.value.value
    } 
  }
    depends_on = [data.traceable_data_set_id.datasetId,data.traceable_service_id.serviceId ]
}

output "traceable_dlp_request_based"{
  value = traceable_dlp_request_based.sample_rule
}






















