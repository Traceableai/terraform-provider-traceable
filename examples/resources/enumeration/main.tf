data "traceable_endpoint_labels" "sample"{
    labels = ["CDE","Login APIs","test-anish"]
}

data "traceable_datasets" "sample"{
    data_sets = ["dlp-automation-dataset","exclude-dlp-automation-dataset"]
}

data "traceable_data_types" "sample"{
    data_types = ["exclude-dlp-automation-datatype","dlp_automation_datatype_23"]
}


data "traceable_endpoints" "sample"{
    endpoints = ["GET /precedence_check_test","GET /test-call-from-dev-random-user"]
}


resource "traceable_enumeration" "sample"{
    name = "enumerationrule"
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
       action_type = "ALERT"
        duration = "PT1M"
        event_severity = "LOW"    
    }

    sources = {
        endpoints = data.traceable_endpoints.sample3.endpoint_ids 

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
        endpoint_labels = data.traceable_endpoint_labels.sample1.label_ids
       
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
            data_sets_ids=data.traceable_datasets.sample2.data_set_ids
            data_location = "REQUEST"

          },
            {
            data_sets_ids=data.traceable_datasets.sample2.data_set_ids
            data_location = "RESPONSE"
          },
           {
            data_sets_ids=data.traceable_datasets.sample2.data_set_ids
           }
        ] 

        data_type = [
          {
            data_types_ids=data.traceable_data_types.sample.data_type_ids
            data_location = "REQUEST"

          },
           {
            data_types_ids=data.traceable_data_types.sample.data_type_ids
            data_location = "RESPONSE"
          },
          {
            data_types_ids=data.traceable_data_types.sample.data_type_ids
          }
        ]
    }
    depends_on = ["data.traceable_endpoint_labels.sample","data.traceable_endpoints.sample","data.traceable_datasets.sample","data.traceable_data_types.sample"]
}
