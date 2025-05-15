package acctest

const RATE_LIMIT_CREATE = `resource "traceable_rate_limiting" "test"{
    name = "%s"
    description = "tf rate limit t1"
    enabled = true
    environments = ["fintech-app"]
    threshold_configs=[
      {
        api_aggregate_type = "PER_ENDPOINT"
        user_aggregate_type = "PER_USER"
        rolling_window_count_allowed = 10
        rolling_window_duration = "PT1M"
        threshold_config_type = "ROLLING_WINDOW"
    },
    {
       
       dynamic_duration = "PT1M"
       dynamic_percentage_exceding_mean_allowed = 10
       dynamic_mean_calculation_duration ="PT1M"
        threshold_config_type = "DYNAMIC"
    }

    ]
    action = {
        action_type = "%s"
        event_severity = "LOW"
    }
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
            user_id_regexes = ["192.168.1.1"]
            exclude = true  
        }
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
            regions_ids =["NU","NF"]
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
                metadata_type = "TAG"
                key_operator = "EQUALS"
                key_value = "key"
                value_operator = "NOT_EQUAL"
                value = "value"
            },
            {
                metadata_type = "TAG"
                key_operator = "EQUALS"
                key_value = "key"
            },
            {
                metadata_type = "URL"
                value_operator = "EQUALS"
                value = "https://example.com/abc"
            }
        ]
    }
}`

const ENUMERATION_RESOURCE=`resource "traceable_enumeration" "test"{
    name = "%s"
    description = "rule-description"
    enabled = true
    environments = ["env2","env1"]
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
    action = {
        action_type = "%s"
        event_severity = "LOW"    
    }

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
        endpoint_labels = ["External","Critical"]
       
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
            data_sets_ids=["0aa32d26-90b9-dbd3-d2a6-6ec3bd1717c3"]
            data_location = "REQUEST"

          },
            {
            data_sets_ids=["0aa32d26-90b9-dbd3-d2a6-6ec3bd1717c3"]
            data_location = "RESPONSE"
          },
           {
            data_sets_ids=["0aa32d26-90b9-dbd3-d2a6-6ec3bd1717c3"]
           }
        ] 

        data_type = [
          {
            data_types_ids=["bba8b92b-6e1f-627e-ed1f-e38b87585159"]
            data_location = "REQUEST"

          },
           {
            data_types_ids=["bba8b92b-6e1f-627e-ed1f-e38b87585159"]
            data_location = "RESPONSE"
          },
          {
            data_types_ids=["bba8b92b-6e1f-627e-ed1f-e38b87585159"]
          }
        ]
    }
}`