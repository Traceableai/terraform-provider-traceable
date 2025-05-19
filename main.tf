terraform {
  required_providers {
    traceable = {
      source  = "traceableai/traceable"
      version = "0.0.1"
    }
  }
}

//object empty ka case check karna hai 
provider "traceable" {
  platform_url = "https://api-staging.traceable.ai"
  api_token    = ""
}

# data "traceable_endpoint_labels" "sample1"{
#     labels = ["CDE","Login APIs","test-anish"]
# }
# output "traceable_labels_sample" {
#   value = data.traceable_endpoint_labels.sample1.label_ids
# }
# data "traceable_datasets" "sample2"{
#     data_sets = ["dlp-automation-dataset","exclude-dlp-automation-dataset"]
# }
# output "traceable_datasets_sample2" {
#   value = data.traceable_datasets.sample2.data_set_ids
# }
# data "traceable_data_types" "sample"{
#     data_types = ["exclude-dlp-automation-datatype","dlp_automation_datatype_23"]
# }
# output "traceable_data_types_sample" {
#   value = data.traceable_data_types.sample.data_type_ids
# }

data "traceable_endpoints" "sample3"{
    endpoints = ["POST /test1","POST /test2","GET /sriniapigee","GET /tQ1f"]
}







output "traceable_endpoint_sample" {
  value = data.traceable_endpoints.sample3.endpoint_ids
}


# data "traceable_services" "sample4"{
#     services = ["java-for-automation-tests-20240719_101628","nginx-for-automation-tests"]
# # }
# output "traceable_services_sample4" {
#   value = data.traceable_services.sample4.service_ids
# }




# resource "traceable_rate_limiting" "sample"{
#     name = "gen"
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
#             user_ids = ["123"]
#             exclude = true  
#         }
#         endpoint_labels = data.traceable_endpoint_labels.sample1.label_ids
#         # endpoints = data.traceable_endpoints.sample3.endpoint_ids 
       
#         ip_reputation = "LOW"
#         ip_location_type = {
#             ip_location_types = ["BOT"]
#             exclude = true  
#         }
#         ip_abuse_velocity = "LOW"
#         # ip_address= {
#         #     ip_address_list = ["192.168.1.1"]
#         #     exclude = true  
#         # }
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
#             },
#             {
#               metadata_type ="TAG"
#               key_operator = "EQUALS"
#               key_value = "key"
#               value_operator = "EQUALS"
#               value = "value"
#             }

#         ]
#     }
#     depends_on = ["data.traceable_endpoint_labels.sample1","data.traceable_endpoints.sample3"]
# }

# resource "traceable_enumeration" "sample"{
#     name = "gen"
#     description = "revamp"
#     enabled = true
#     environments = ["env2","env1"]
#     threshold_configs=[
#       {
#         api_aggregate_type = "PER_ENDPOINT"
#         value_type ="SENSITIVE_PARAMS"
#         unique_values_allowed = 10
#         duration ="PT1M"
#         threshold_config_type = "VALUE_BASED"
#         user_aggregate_type = "PER_USER"
#         sensitive_params_evaluation_type = "ALL"
#     },
#     {
#         api_aggregate_type = "PER_ENDPOINT"
#         value_type ="REQUEST_BODY"
#         unique_values_allowed = 10
#         duration ="PT1M"
#         threshold_config_type = "VALUE_BASED"
#         user_aggregate_type = "PER_USER"
#     },
#     {
#         api_aggregate_type = "PER_ENDPOINT"
#         value_type ="PATH_PARAMS"
#         unique_values_allowed = 10
#         duration ="PT1M"
#         threshold_config_type = "VALUE_BASED"
#         user_aggregate_type = "PER_USER"
#     },
#     #  {
#     #     api_aggregate_type = "PER_ENDPOINT"
#     #     user_aggregate_type = "PER_USER"
#     #     rolling_window_count_allowed = 10
#     #     rolling_window_duration = "PT1M"
#     #     threshold_config_type = "ROLLING_WINDOW"
#     # },
#     # {
       
#     #    dynamic_duration = "PT1M"
#     #    dynamic_percentage_exceding_mean_allowed = 10
#     #    dynamic_mean_calculation_duration ="PT24H"
#     #     threshold_config_type = "DYNAMIC"
#     # }
   

#     ]
#     action = {
#        action_type = "ALERT"
#         event_severity = "LOW"    
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
#             user_ids = ["123"]
#             exclude = true  
#         }
#         # endpoint_labels = data.traceable_endpoint_labels.sample1.label_ids
#         # endpoints = data.traceable_endpoints.sample3.endpoint_ids 
       
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
#             },
#             {
#               metadata_type ="TAG"
#               key_operator = "EQUALS"
#               key_value = "key"
#               value_operator = "EQUALS"
#               value = "value"
#             }

#         ]
#         # data_set = [
#         #   {
#         #     data_sets_ids=data.traceable_datasets.sample2.data_set_ids
#         #     data_location = "REQUEST"

#         #   },
#         #     {
#         #     data_sets_ids=data.traceable_datasets.sample2.data_set_ids
#         #     data_location = "RESPONSE"
#         #   },
#         #    {
#         #     data_sets_ids=data.traceable_datasets.sample2.data_set_ids
#         #    }
#         # ] 

#         # data_type = [
#         #   {
#         #     data_types_ids=data.traceable_data_types.sample.data_type_ids
#         #     data_location = "REQUEST"

#         #   },
#         #    {
#         #     data_types_ids=data.traceable_data_types.sample.data_type_ids
#         #     data_location = "RESPONSE"
#         #   },
#         #   {
#         #     data_types_ids=data.traceable_data_types.sample.data_type_ids
#         #   }
#         # ]
#     }
#     # depends_on = ["data.traceable_endpoint_labels.sample","data.traceable_endpoints.sample","data.traceable_datasets.sample","data.traceable_data_types.sample"]
# }



# resource "traceable_data_loss_prevention_user_based" "sample"{
#     name = "gen"
#     description = "revamp"
#     enabled = true
#     environments = ["env2","env1"]
#     threshold_configs=[
#       {
#         api_aggregate_type = "PER_ENDPOINT"
#         value_type ="SENSITIVE_PARAMS"
#         unique_values_allowed = 10
#         duration ="PT1M"
#         threshold_config_type = "VALUE_BASED"
#         user_aggregate_type = "PER_USER"
#         sensitive_params_evaluation_type = "ALL"
#     },
#     {
#         api_aggregate_type = "PER_ENDPOINT"
#         value_type ="REQUEST_BODY"
#         unique_values_allowed = 10
#         duration ="PT1M"
#         threshold_config_type = "VALUE_BASED"
#         user_aggregate_type = "PER_USER"
#     },
#     {
#         api_aggregate_type = "PER_ENDPOINT"
#         value_type ="PATH_PARAMS"
#         unique_values_allowed = 10
#         duration ="PT1M"
#         threshold_config_type = "VALUE_BASED"
#         user_aggregate_type = "PER_USER"
#     },
   

#     ]
#     action = {
#        action_type = "ALERT"
#         duration = "PT1M"
#         event_severity = "LOW"    
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
#             user_ids = ["123"]
#             exclude = true  
#         }
#         endpoint_labels = data.traceable_endpoint_labels.sample1.label_ids
#         # endpoints = data.traceable_endpoints.sample3.endpoint_ids 
       
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
#             },
#             {
#               metadata_type ="TAG"
#               key_operator = "EQUALS"
#               key_value = "key"
#               value_operator = "EQUALS"
#               value = "value"
#             }

#         ]
#         data_set = [
#           {
#             data_sets_ids=data.traceable_datasets.sample2.data_set_ids
#             data_location = "REQUEST"

#           },
#             {
#             data_sets_ids=data.traceable_datasets.sample2.data_set_ids
#             data_location = "RESPONSE"
#           },
#            {
#             data_sets_ids=data.traceable_datasets.sample2.data_set_ids
#            }
#         ] 

#         data_type = [
#           {
#             data_types_ids=data.traceable_data_types.sample.data_type_ids
#             data_location = "REQUEST"

#           },
#            {
#             data_types_ids=data.traceable_data_types.sample.data_type_ids
#             data_location = "RESPONSE"
#           },
#           {
#             data_types_ids=data.traceable_data_types.sample.data_type_ids
#           }
#         ]
#     }
#     depends_on = ["data.traceable_endpoint_labels.sample1","data.traceable_endpoints.sample3","data.traceable_datasets.sample2","data.traceable_data_types.sample"]
# }




    
       

    
# resource "traceable_malicious_ip_range" "sample"{
#     name = "shreyanshrevamp12"
#     description = "revamp"
#     enabled = true
#     event_severity = "LOW"
#     duration = "PT1M"
#     action = "BLOCK"
#     ip_range = ["192.168.1.1","192.168.1.2"]
#     environments = ["env1","env2"]
# }
# resource "traceable_malicious_region" "sample"{
#     name = "tftest1234"
#     description = "revamp"
#     enabled = false
#     event_severity = "LOW"
#     duration = "PT1M"
#     action = "BLOCK"
#     environments = ["env1","env2"]

#     regions=[  "NU","AF",
#        ]
# }

# resource "traceable_malicious_ip_type" "sample"{
#     name = "tftest1234"
#     description = "revamp"
#     enabled = true
#     event_severity = "HIGH"
#     duration = "PT1M"
#     action = "ALERT"
#     environments = ["env1","env2"]
#     ip_type = ["ANONYMOUS_VPN","BOT",]
# }



# mutation {
#   createRegionRule(
#     input: {
#       name: "asC1"
#       regionIds: [
#         ""
#         "3fd065e2-d636-5c39-b7de-9634fb9b5cc9"
#         "bd0df17d-3921-554c-9d70-71e57fda09bc"
#       ]
#       type: BLOCK
#       description: "sC"
#       eventSeverity: LOW
#       effects: []
#     }
#   ) {
#     id
# #     __typename
# #   }
# # }


resource "traceable_rate_limiting" "test"{
    name = "test"
    description = "tf rate limit t1"
    enabled = true
    environments = ["fintech-app"]
    threshold_configs=[
      {
        api_aggregate_type = "PER_ENDPOINT"
        user_aggregate_type = "PER_USER"
        rolling_window_count_allowed = 10
        rolling_window_duration = "PT1H30M20S"
        threshold_config_type = "ROLLING_WINDOW"
    },
    {
       
       dynamic_duration = "PT60S"
       dynamic_percentage_exceding_mean_allowed = 10
       dynamic_mean_calculation_duration ="PT60S"
        threshold_config_type = "DYNAMIC"
        
    },
 


    ]
    action = {
        action_type = "ALERT"
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
            ip_address_type = "ALL_EXTERNAL"
            # ip_address_list = ["192.168.1.1"]
            # exclude = true  
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
                metadata_type = "URL"
                value_operator = "EQUALS"
                value = "https://example.com/abc"
            }
        ]
    }
}

# resource "traceable_custom_signature" "example" {
#   name         = "shreyansh_signature_1"         
#   disabled     = false                       
#   description  = "Sample custom signature"   
#   environments = ["1260"]            

#   payload_criteria = {
#     request_response = [
#       {
#         match_category       = "REQUEST"         
#         match_key            = "URL"    
#         value_match_operator = "EQUALS"          
#         match_value          = "secret"          
#       }
#     ]
#     attributes = [
#       {
#         key_condition_operator   = "EQUALS"      
#         key_condition_value      = "http.status" 
#         value_condition_operator = "EQUALS"      
#         value_condition_value    = "500"         
#       }
#     ]  
#   }

#   action = {
#     action_type    = "NORMAL_DETECTION" 
#     duration       = "PT1M"  
#     event_severity = "LOW"   
#   }
# }