resource "traceable_data_loss_prevention_user_based" "example" {
  name         = "example_dlp_user_rule_2"
  description  = "Example DLP user-based rule"
  enabled      = true
  environments = ["dev", "prod"]

  action = {
    action_type    = "BLOCK"   # ALERT | BLOCK | ALLOW
    event_severity = "LOW"     # LOW | MEDIUM | HIGH | CRITICAL (only BLOCK / ALERT)
    # duration     = "PT60S"   # Allowed only with ALLOW / BLOCK
  }
     threshold_configs=[
      {
        api_aggregate_type = "PER_ENDPOINT"
        value_type ="SENSITIVE_PARAMS"
        unique_values_allowed = 10
        duration ="PT1M"
        threshold_config_type = "VALUE_BASED"
        user_aggregate_type = "PER_USER"
        sensitive_params_evaluation_type = "ALL"
    },
    {
        api_aggregate_type = "PER_ENDPOINT"
        value_type ="REQUEST_BODY"
        unique_values_allowed = 10
        duration ="PT1M"
        threshold_config_type = "VALUE_BASED"
        user_aggregate_type = "PER_USER"
    },
    {
        api_aggregate_type = "PER_ENDPOINT"
        value_type ="PATH_PARAMS"
        unique_values_allowed = 10
        duration ="PT1M"
        threshold_config_type = "VALUE_BASED"
        user_aggregate_type = "PER_USER"
    }
    ]
  sources = {
        ip_asn = {
            ip_asn_regexes = ["vmv"]
            exclude = true  
        }
        ip_connection_type = {
            ip_connection_type_list = ["RESIDENTIAL","MOBILE"]
            exclude = true  
        }
        user_id = {
            user_ids = ["123"]
            exclude = true  
        }
        endpoint_labels = ["External"]
       
        ip_reputation = "LOW"
        ip_location_type = {
            ip_location_types = ["BOT"]
            exclude = true  
        }
        ip_abuse_velocity = "LOW"
        ip_address= {
            ip_address_list = ["192.168.1.1"]
            exclude = true  
        }
        email_domain = {
            email_domain_regexes = ["abc*.*gmail.com"]
            exclude = true  
        }

        user_agents = {
            user_agents_list = ["chrome"]
            exclude = true  
        }
        regions = {
            regions_ids = ["AF","AW"]
          
            exclude = true  
        }
        ip_organisation = {
            ip_organisation_regexes = ["192.168.1.1"]
            exclude = true  
        }
        request_response = [  
            {
                metadata_type = "REQUEST_HEADER"
                key_operator = "EQUALS"
                key_value = "key"
                value_operator = "EQUALS"
                value = "value"
            },
            {
                metadata_type = "URL"
                value_operator = "EQUALS"
                value = "https://example.com/abc"
            },
            {
              metadata_type ="TAG"
              key_operator = "EQUALS"
              key_value = "key"
              value_operator = "EQUALS"
              value = "value"
            }

        ]

     data_set = [
        {
        data_sets_ids=["test-data-set"]
        data_location = "REQUEST"

        }
    ] 

    data_type = [
        {
        data_types_ids=["test-data-type"]
        data_location = "REQUEST"
        }
    ]
  }
}