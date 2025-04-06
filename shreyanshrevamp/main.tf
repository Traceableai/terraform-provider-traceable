terraform {
  required_providers {
    traceable = {
      source  = "traceableai/traceable"
      version = "0.0.1"
    }
    # aws = {
    #   source  = "hashicorp/aws"
    #   version = "5.35.0"
    # }
  }
}

//object empty ka case check karna hai 
provider "traceable" {
  platform_url ="https://app-dev.traceable.ai/graphql"
  api_token    ="Bearer "
  
}

# resource "traceable_data_set" "sampledataset"{
#              description = "hello I am good"
#             icon_type = "Password"
#             name= "shreyansh12"
# }


resource "traceable_rate_limiting" "sample"{
    name = "shreyanshrevamp"
    description = "revamp"
    enabled = true
    environments = ["env1","env2"]
    threshold_configs=[
      {
        api_aggregate_type = "PER_ENDPOINT"
        user_aggregate_type = "PER_USER"
        rolling_window_count_allowed = 10
        rolling_window_duration = "PT60s"
        threshold_config_type = "ROLLING_WINDOW"
    },
    {
       
       dynamic_duration = "PT60S"
       dynamic_percentage_exceding_mean_allowed = 10
       dynamic_mean_calculation_duration ="PT60S"
        threshold_config_type = "DYNAMIC"
    }

    ]
    action = {
        action_type = "ALERT"
        # duration = "PT60s"
        event_severity = "LOW"    
        # header_injections = [
        #     {
        #         key = "X-Custom-Header"
        #         value = "CustomValue"
        #     },
        #     {
        #         key = "X-Custom-Header2"
        #         value = "CustomValue2"
        #     }
        # ] 
    }
    sources = {
        # ip_asn = {
        #     ip_asn_regexes = ["vmv"]
        #     exclude = true  
        # }
        ip_connection_type = {
            ip_connection_type_list = ["RESIDENTIAL","MOBILE"]
            exclude = true  
        }
        user_id = {
            user_id_regexes = ["192.168.1.1"]
            exclude = true  
        }
        endpoint_labels = ["/abc123"]
        endpoints = ["/abc123"]
        attribute = [{
            key_condition_operator = "EQUALS"
            key_condition_value = "key"
            value_condition_operator = "EQUALS"
            value_condition_value = "value" 
        }]
        
    

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
            # regions_ids = ["AF","AW"]
            regions_ids =["AF","AQ"]
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
            }
        ]
    }
}


       

    
