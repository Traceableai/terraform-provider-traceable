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
        action_type = "BLOCK"
        event_severity = "%s"
        duration="%s"
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
        action_type = "BLOCK"
        duration="%s"
        event_severity = "%s"    
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
    }
}`

const CUSTOM_SIGNATURE_RESOURCE = `resource "traceable_custom_signature" "test" {
  name         = "%s"         
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
    custom_sec_rule = <<EOF
    SecRule REQUEST_HEADERS:key-sec "@rx val-sec" \
    "id:92100120,\
    phase:2,\
    block,\
    msg:'Test sec Rule',\
    tag:'attack-protocol',\
    tag:'traceable/labels/OWASP_2021:A4,CWE:444,OWASP_API_2019:API8',\
    tag:'traceable/severity/HIGH',\
    tag:'traceable/type/safe,block',\
    severity:'CRITICAL',\
    setvar:'tx.anomaly_score_pl1=+tx.critical_anomaly_score'"
    EOF
  }

  action = {
    action_type    = "DETECTION_AND_BLOCKING" 
    duration       = "%s"  
    event_severity = "%s"   
  }
}`

const IP_TYPE_RESOURCE=`resource "traceable_malicious_ip_type" "test"{
    name = "%s"
    description = "terraform unit test"
    enabled = true
    event_severity = "%s"
    duration = "%s"
    action = "BLOCK"
    environments = ["1260"]
    ip_type = ["ANONYMOUS_VPN","BOT"]
}`

const EMAIL_DOMAIN=`resource "traceable_malicious_email_domain" "test" {
    name = "%s"
    description = "example rule"
    enabled = true
    event_severity = "%s"
    duration = "%s"
    action = "BLOCK"
    environments = ["1260"]
    email_domains_list = ["traceable.ai", "harness.io"]
    apply_rule_to_data_leaked_email = true
    min_email_fraud_score_level = "CRITICAL"
    email_regexes_list = ["traceable.ai", "harness.io"]
    apply_rule_to_disposable_email_domains = true
}`

const DATA_SET=`resource "traceable_data_set" "test" {
    name = "%s"
    description = "%s"
    icon_type = "Protect"
}`

const DLP_USER_BASED=`resource "traceable_data_loss_prevention_user_based" "test" {
  name         = "%s"
  description  = "Example DLP user-based rule"
  enabled      = true
  environments = ["dev", "prod"]

  action = {
    action_type    = "BLOCK"  
    duration     = "%s"
    event_severity = "%s"    
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
        data_location = "RESPONSE"

        }
    ] 

    data_type = [
        {
        data_types_ids=["test-data-type"]
        data_location = "RESPONSE"
        }
    ]
  }
}`

const DLP_REQ_BASED=`resource "traceable_data_loss_prevention_request_based" "test" {
  name         = "%s"
  description  = "Example DLP request-based rule"
  enabled      = true
  environments = ["dev", "prod"]

  action = {
    action_type    = "BLOCK"   
    duration     = "%s"   
    event_severity = "%s"     
  }

  sources = {
    service_scope = {
      service_ids = ["172d77d2-ae7a-3355-9a09-2db0803f3523"]
    }

    url_regex_scope = {
      url_regexes = ["/private/.*"]
    }

    regions = {
      region_ids = ["US", "CA"]
    }

    ip_location_type = {
      ip_location_types = ["BOT", "ANONYMOUS_VPN"]
    }

    ip_address = {
      ip_address_list = ["192.0.2.10","1.1.1.1"]
    }

    request_payload = [
      {
        metadata_type   = "REQUEST_HEADER" 
        key_operator    = "EQUALS"
        key_value       = "Content-Type"
        value_operator  = "CONTAINS"
        value           = "json"
      },
      {
        metadata_type  = "HTTP_METHOD"
        value_operator = "EQUALS"
        value          = "POST"
      }
    ]

    dateset_datatype_filter = {
      dateset_datatype_id = {
        data_sets_ids  = ["4c41bbe3-92a3-42c2-aa78-3d69cef1ef49"]
      }


      data_type_matching = {
        metadata_type = "REQUEST_BODY"
      }
    }
  }
}`