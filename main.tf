terraform {
  required_providers {
    traceable = {
      source  = "traceableai/traceable"
      version = "0.0.1"
    }
  }
}

variable "API_TOKEN" {
  type = string
}

provider "traceable" {
  platform_url = "https://api-dev.traceable.ai"
  api_token    = var.API_TOKEN
}

resource "traceable_custom_signature" "example" {
  name         = "example_signature_test_aditya_tf_1"         
  disabled     = false                       
  description  = "Sample custom signature"   

  payload_criteria = {
    request_response = [
      {
        match_category       = "REQUEST"         
        match_key            = "URL"    
        value_match_operator = "EQUALS"          
        match_value          = "secret"          
      }
    ]
    attributes = [
      {
        key_condition_operator   = "EQUALS"      
        key_condition_value      = "http.status" 
        value_condition_operator = "EQUALS"      
        value_condition_value    = "500"         
      }
    ]  
  }

  action = {
    action_type    = "TESTING_DETECTION" 
  }
}