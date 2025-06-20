---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "traceable_enumeration Resource - terraform-provider-traceable"
subcategory: ""
description: |-
  Traceable Enumeration Resource
---

# traceable_enumeration (Resource)

Traceable Enumeration Resource

## Example Usage
```
resource "traceable_enumeration" "test"{
    name = "enumeration-test-rule"
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
        duration="PT1M"
        event_severity = "MEDIUM"    
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
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (Attributes) (see [below for nested schema](#nestedatt--action))
- `enabled` (Boolean) Enable the Enumeration Rule
- `name` (String) Name of the Enumeration Rule.

### Optional

- `description` (String) Description of the Enumeration Rule
- `environments` (Set of String) Environments the rule is applicable to
- `sources` (Attributes) (see [below for nested schema](#nestedatt--sources))
- `threshold_configs` (Attributes Set) Threshold configs for the rule (see [below for nested schema](#nestedatt--threshold_configs))

### Read-Only

- `id` (String) Identifier of the Enumeration Rule

<a id="nestedatt--action"></a>
### Nested Schema for `action`

Required:

- `action_type` (String) ALERT,BLOCK,MARK FOR TESTING

Optional:

- `duration` (String) Duration for the action (PT60S/PT1M)
- `event_severity` (String) LOW,MEDIUM,HIGH,CRITICAL
- `header_injections` (Attributes Set) Header fields to be injected (see [below for nested schema](#nestedatt--action--header_injections))

<a id="nestedatt--action--header_injections"></a>
### Nested Schema for `action.header_injections`

Optional:

- `key` (String) The header field name to inject (e.g., 'X-Custom-Header')
- `value` (String) The value to set for the header field



<a id="nestedatt--sources"></a>
### Nested Schema for `sources`

Optional:

- `data_set` (Attributes Set) Request/response attributes as source (see [below for nested schema](#nestedatt--sources--data_set))
- `data_type` (Attributes Set) Request/response attributes as source (see [below for nested schema](#nestedatt--sources--data_type))
- `email_domain` (Attributes) Email domain as source, It will be a list of email domain regexes (see [below for nested schema](#nestedatt--sources--email_domain))
- `endpoint_labels` (Set of String) Filter endpoints by labels you want to apply this rule
- `endpoints` (Set of String) List of endpoint ids
- `ip_abuse_velocity` (String) Ip abuse velocity as source (LOW/MEDIUM/HIGH)
- `ip_address` (Attributes) Ip address as source (LIST_OF_IP's/ALL_EXTERNAL) (see [below for nested schema](#nestedatt--sources--ip_address))
- `ip_asn` (Attributes) (see [below for nested schema](#nestedatt--sources--ip_asn))
- `ip_connection_type` (Attributes) Ip connection type as source, It will be a list of ip connection type (see [below for nested schema](#nestedatt--sources--ip_connection_type))
- `ip_location_type` (Attributes) Ip location type as source (see [below for nested schema](#nestedatt--sources--ip_location_type))
- `ip_organisation` (Attributes) Ip organisation as source, It will be a list of ip organisation (see [below for nested schema](#nestedatt--sources--ip_organisation))
- `ip_reputation` (String) Ip reputation source (LOW/MEDIUM/HIGH/CRITICAL)
- `regions` (Attributes) Regions as source, It will be a list region ids (AX,DZ) (see [below for nested schema](#nestedatt--sources--regions))
- `request_response` (Attributes Set) Request/response attributes as source (see [below for nested schema](#nestedatt--sources--request_response))
- `scanner` (Attributes) Scanner as source, It will be a list of scanner type (see [below for nested schema](#nestedatt--sources--scanner))
- `user_agents` (Attributes) User agents as source, It will be a list of user agents (see [below for nested schema](#nestedatt--sources--user_agents))
- `user_id` (Attributes) User id as source (see [below for nested schema](#nestedatt--sources--user_id))

<a id="nestedatt--sources--data_set"></a>
### Nested Schema for `sources.data_set`

Optional:

- `data_location` (String) Specifies which metadata type to include (`REQUEST`, `RESPONSE`). If not defined, applies on both
- `data_sets_ids` (Set of String) Which operator to use


<a id="nestedatt--sources--data_type"></a>
### Nested Schema for `sources.data_type`

Optional:

- `data_location` (String) Which Metadatype to include
- `data_types_ids` (Set of String) Which operator to use


<a id="nestedatt--sources--email_domain"></a>
### Nested Schema for `sources.email_domain`

Required:

- `email_domain_regexes` (Set of String) It will be a list of email domain regexes
- `exclude` (Boolean) Set it to true to exclude given email domains regexes


<a id="nestedatt--sources--ip_address"></a>
### Nested Schema for `sources.ip_address`

Optional:

- `exclude` (Boolean) Set it to true to exclude given ip addresses
- `ip_address_list` (Set of String) List of ip addresses
- `ip_address_type` (String) Accepts ALL_EXTERNAL,ALL_INTERNAL


<a id="nestedatt--sources--ip_asn"></a>
### Nested Schema for `sources.ip_asn`

Required:

- `exclude` (Boolean) Set it to true to exclude given IP ASN
- `ip_asn_regexes` (Set of String) It will be a list of IP ASNs


<a id="nestedatt--sources--ip_connection_type"></a>
### Nested Schema for `sources.ip_connection_type`

Required:

- `exclude` (Boolean) Set it to true to exclude given IP connection
- `ip_connection_type_list` (Set of String) It will be a list of IP connection types(RESIDENTIAL,MOBILE,CORPORATE,DATA_CENTER,EDUCATION)


<a id="nestedatt--sources--ip_location_type"></a>
### Nested Schema for `sources.ip_location_type`

Required:

- `exclude` (Boolean) Set it to true to exclude given ip location types
- `ip_location_types` (Set of String) Ip location type as source ([BOT,ANONYMOUS_VPN,HOSTING_PROVIDER,TOR_EXIT_NODE, PUBLIC_PROXY,SCANNER])


<a id="nestedatt--sources--ip_organisation"></a>
### Nested Schema for `sources.ip_organisation`

Required:

- `exclude` (Boolean) Set it to true to exclude given ip organisation
- `ip_organisation_regexes` (Set of String) It will be a list of ip organisations


<a id="nestedatt--sources--regions"></a>
### Nested Schema for `sources.regions`

Required:

- `exclude` (Boolean) Set it to true to exclude given regions
- `regions_ids` (Set of String) It will be a list of regions ids in countryIsoCode


<a id="nestedatt--sources--request_response"></a>
### Nested Schema for `sources.request_response`

Required:

- `metadata_type` (String) Which Metadatype to include

Optional:

- `key_operator` (String) Which operator to use
- `key_value` (String) Value to match
- `value` (String) Value to match
- `value_operator` (String) Which operator to use


<a id="nestedatt--sources--scanner"></a>
### Nested Schema for `sources.scanner`

Required:

- `exclude` (Boolean) Set it to true to exclude given scaner types
- `scanner_types_list` (Set of String) It will be a list of scanner types(Traceable AST,Qualys,Rapid7 InsightAppSec,Invicti,Tenable)


<a id="nestedatt--sources--user_agents"></a>
### Nested Schema for `sources.user_agents`

Required:

- `exclude` (Boolean) Set it to true to exclude given user agents
- `user_agents_list` (Set of String) It will be a list of user agents


<a id="nestedatt--sources--user_id"></a>
### Nested Schema for `sources.user_id`

Required:

- `exclude` (Boolean) Set it to true to exclude given user id

Optional:

- `user_id_regexes` (Set of String) It will be a list of user id regexes
- `user_ids` (Set of String) List of user ids



<a id="nestedatt--threshold_configs"></a>
### Nested Schema for `threshold_configs`

Required:

- `threshold_config_type` (String) Threshold config type(ROLLING_WINDOW,VALUE_BASED,DYNAMIC)

Optional:

- `api_aggregate_type` (String) API aggregate type(PER_ENDPOINT,ACROSS_ENDPOINT)
- `duration` (String) Duration for the rule (PT60S/PT1M)
- `dynamic_duration` (String) Dynamic duration (PT60S/PT1M)
- `dynamic_mean_calculation_duration` (String) Dynamic mean calculation duration (PT60S/PT1M)
- `dynamic_percentage_exceding_mean_allowed` (Number) Dynamic percentage exceeding mean allowed
- `rolling_window_count_allowed` (Number) Rolling window count allowed
- `rolling_window_duration` (String) Rolling window duration (PT60S/PT1M)
- `sensitive_params_evaluation_type` (String) Sensitive params evaluation type (ALL/SELECTED_DATA_TYPES)
- `unique_values_allowed` (Number) Unique values allowed
- `user_aggregate_type` (String) User aggregate type(PER_USER,ACROSS_USER)
- `value_type` (String) Value type (REQUEST_BODY/SENSITIVE_PARAMS/PATH_PARAMS)
