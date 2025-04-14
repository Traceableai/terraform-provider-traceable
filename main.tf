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
  platform_url ="https://app-dev.traceable.ai"
  api_token    =""


  
}




# resource "traceable_rate_limiting" "sample"{
#     name = "shreyanshrevamp12"
#     description = "revamp"
#     enabled = true
#     environments = ["env2","env1"]
#     threshold_configs=[
#       {
#         api_aggregate_type = "PER_ENDPOINT"
#         user_aggregate_type = "PER_USER"
#         rolling_window_count_allowed = 10
#         rolling_window_duration = "PT1M"
#         threshold_config_type = "ROLLING_WINDOW"
#     },
#     {
       
#        dynamic_duration = "PT1M"
#        dynamic_percentage_exceding_mean_allowed = 10
#        dynamic_mean_calculation_duration ="PT24H"
#         threshold_config_type = "DYNAMIC"
#     }

#     ]
#     action = {
#        action_type = "ALERT"
#         duration = "PT1M"
#         event_severity = "LOW"    
#         # header_injections = [
#         #     {
#         #         key = "X-Custom-Header"
#         #         value = "CustomValue"
#         #     },
#         #     {
#         #         key = "X-Custom-Header2"
#         #         value = "CustomValue2"
#         #     }
#         # ] 
#     }
#     sources = {
#         ip_asn = {
#             ip_asn_regexes = ["vmv"]
#             exclude = true  
#         }
#         ip_connection_type = {
#             ip_connection_type_list = ["RESIDENTIAL","MOBILE"]
#             exclude = true  
#         }
#         user_id = {
#             exclude = true  
#         }
#         # endpoint_labels = ["/abc123"]
#         # endpoints = ["/abc123"]
       
        
    

#         ip_reputation = "LOW"
#         ip_location_type = {
#             ip_location_types = ["BOT"]
#             exclude = true  
#         }
#         ip_abuse_velocity = "LOW"
#         ip_address= {
#             ip_address_list = ["192.168.1.1"]
#             exclude = true  
#         }
#         email_domain = {
#             email_domain_regexes = ["abc*.*gmail.com"]
#             exclude = true  
#         }

#         user_agents = {
#             user_agents_list = ["chrome"]
#             exclude = true  
#         }
#         regions = {
#             regions_ids = ["AF","AW"]
#             # regions_ids =["AF","AQ"]
#             exclude = true  
#         }
#         ip_organisation = {
#             ip_organisation_regexes = ["192.168.1.1"]
#             exclude = true  
#         }
#         request_response = [  
#             {
#                 metadata_type = "REQUEST_HEADER"
#                 key_operator = "EQUALS"
#                 key_value = "key"
#                 value_operator = "EQUALS"
#                 value = "value"
#             },
#             {
#                 metadata_type = "URL"
#                 value_operator = "EQUALS"
#                 value = "https://example.com/abc"
#             }
#         ]
#     }
# }

    
       

    
resource "traceable_malicious_ip_range" "sample"{
    name = "shreyanshrevamp12"
    description = "revamp"
    enabled = true
    event_severity = "LOW"
    duration = "PT1M"
    action = "BLOCK_ALL_EXCEPT"
    ip_range = ["192.168.1.1","192.168.1.2"]
    environments = ["env1","env2"]
}