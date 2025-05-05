resource "traceable_custom_signature" "example" {
  name         = "example_signature"         
  disabled     = false                       
  description  = "Sample custom signature"   
  environments = ["1260"]            

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
    action_type    = "NORMAL_DETECTION" 
    duration       = "PT1M"  
    event_severity = "LOW"   
  }
}