terraform {
  required_providers {
    traceable = {
      source  = "terraform.local/local/traceable"
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
  name         = "example_signature_adi_test_2"         
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
    custom_sec_rule = <<EOR
    SecRule REQUEST_HEADERS:key-sec "@rx val-sec" \
    "id:92100120,\
    phase:2,\
    block,\
    msg:'Test sec Rule',\
    logdata:'Matched Data: %%{TX.0} found within %%{MATCHED_VAR_NAME}: %%{MATCHED_VAR}',\
    tag:'attack-protocol',\
    tag:'traceable/labels/OWASP_2021:A4,CWE:444,OWASP_API_2019:API8',\
    tag:'traceable/severity/HIGH',\
    tag:'traceable/type/safe,block',\
    severity:'CRITICAL',\
    setvar:'tx.anomaly_score_pl1=+%%{tx.critical_anomaly_score}'"
    EOR
  }

  action = {
    action_type    = "NORMAL_DETECTION" 
    event_severity = "LOW"   
  }
}