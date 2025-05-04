resource "traceable_custom_signature" "example" {
  name         = "example_signature"         
  disabled     = false                       
  description  = "Sample custom signature"   
  environments = ["dlp-test-aditya"]            

  payload_criteria = {

    request_response = [
      {
        key_value_tag        = "REQUEST_HEADER"  
        match_category       = "EQUALS"          
        key_match_operator   = "EQUALS"          
        match_key            = "X-Custom-Key"    
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
    action_type    = "ALERT" 
    duration       = "PT1M"  
    event_severity = "LOW"   
  }
}