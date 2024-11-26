package rate_limiting

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"log"
)

func ResourceRateLimitingRuleBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceRateLimitingRuleBlockCreate,
		Read:   resourceRateLimitingRuleBlockRead,
		Update: resourceRateLimitingRuleBlockUpdate,
		Delete: resourceRateLimitingRuleBlockDelete,

		Schema: map[string]*schema.Schema{
			"rule_type": {
				Type:        schema.TypeString,
				Description: "ALERT or BLOCK",
				Optional:    true,
				Default:     "BLOCK",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the rate limiting block rule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the rate limiting rule",
				Optional:    true,
			},
			"alert_severity": {
				Type:        schema.TypeString,
				Description: "LOW/MEDIUM/HIGH/CRITICAL",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable/disable the rule",
				Required:    true,
			},
			"block_expiry_duration": {
				Type:        schema.TypeString,
				Description: "Block for a given period",
				Optional:    true,
			},
			"label_id_scope": {
				Type:        schema.TypeList,
				Description: "Filter endpoints by labels you want to apply this rule",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"endpoint_id_scope": {
				Type:        schema.TypeList,
				Description: "List of endpoint ids",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environments": {
				Type:        schema.TypeList,
				Description: "List of environments ids",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"req_res_conditions": {
				Type:        schema.TypeList,
				Description: "Request/Response conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"req_res_operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"req_res_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"threshold_configs": {
				Type:        schema.TypeList,
				Description: "Threshold configs for the rule",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_aggregate_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ACROSS_ENDPOINTS/PER_ENDPOINT",
						},
						"rolling_window_count_allowed": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Count of calls",
						},
						"rolling_window_duration": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Duration for the total call count in 1min(PT60S)",
						},
						"threshold_config_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DYNAMIC/ROLLING_WINDOW",
						},
						"dynamic_mean_calculation_duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Baseline is calculated over 1D(PT86400S)",
						},
					},
				},
			},
			"attribute_based_conditions": {
				Type:        schema.TypeList,
				Description: "Attribute based conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_condition_operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key_condition_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value_condition_operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_condition_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ip_reputation": {
				Type:        schema.TypeString,
				Description: "Ip reputation source (LOW/MEDIUM/HIGH/CRITICAL)",
				Optional:    true,
			},
			"ip_location_type": {
				Type:        schema.TypeList,
				Description: "Ip location type as source ([BOT, TOR_EXIT_NODE, PUBLIC_PROXY])",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_location_types": {
							Type:        schema.TypeList,
							Description: "It will be a list of ip location types",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given ip location",
							Required:    true,
						},
					},
				},
			},
			"ip_abuse_velocity": {
				Type:        schema.TypeString,
				Description: "Ip abuse velocity as source (LOW/MEDIUM/HIGH/CRITICAL)",
				Optional:    true,
			},
			"ip_address": {
				Type:        schema.TypeList,
				Description: "Ip address as source (LIST_OF_IP's/ALL_EXTERNAL)",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address_list": {
							Type:        schema.TypeList,
							Description: "List of ip addresses",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given ip addresses",
							Required:    true,
						},
						"ip_address_type": {
							Type:        schema.TypeString,
							Description: "Accepts one value ALL_EXTERNAL",
							Optional:    true,
						},
					},
				},
			},
			"email_domain": {
				Type:        schema.TypeList,
				Description: "Email domain as source, It will be a list of email domain regexes",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_domain_regexes": {
							Type:        schema.TypeList,
							Description: "It will be a list of email domain regexes",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given email domains regexes",
							Required:    true,
						},
					},
				},
			},
			"user_agents": {
				Type:        schema.TypeList,
				Description: "User agents as source, It will be a list of user agents",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_agents_list": {
							Type:        schema.TypeList,
							Description: "It will be a list of user agents",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given user agents",
							Required:    true,
						},
					},
				},
			},
			"regions": {
				Type:        schema.TypeList,
				Description: "Regions as source, It will be a list region ids (AX,DZ)",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"regions_ids": {
							Type:        schema.TypeList,
							Description: "It will be a list of regions ids in countryIsoCode",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given regions",
							Required:    true,
						},
					},
				},
			},
			"ip_organisation": {
				Type:        schema.TypeList,
				Description: "Ip organisation as source, It will be a list of ip organisation",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_organisation_regexes": {
							Type:        schema.TypeList,
							Description: "It will be a list of ip organisations",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given ip organisation",
							Required:    true,
						},
					},
				},
			},
			"ip_asn": {
				Type:        schema.TypeList,
				Description: "Ip ASN as source, It will be a list of IP ASN",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_asn_regexes": {
							Type:        schema.TypeList,
							Description: "It will be a list of IP ASN",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given IP ASN",
							Required:    true,
						},
					},
				},
			},
			"ip_connection_type": {
				Type:        schema.TypeList,
				Description: "Ip connection type as source, It will be a list of ip connection type",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_connection_type_list": {
							Type:        schema.TypeList,
							Description: "It will be a list of IP connection types",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given IP coonection",
							Required:    true,
						},
					},
				},
			},
			"request_scanner_type": {
				Type:        schema.TypeList,
				Description: "Scanner as source, It will be a list of scanner type",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scanner_types_list": {
							Type:        schema.TypeList,
							Description: "It will be a list of scanner types",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given scaner types",
							Required:    true,
						},
					},
				},
			},
			"user_id": {
				Type:        schema.TypeList,
				Description: "User id as source",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id_regexes": {
							Type:        schema.TypeList,
							Description: "It will be a list of user id regexes",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"exclude": {
							Type:        schema.TypeBool,
							Description: "Set it to true to exclude given user id",
							Required:    true,
						},
						"user_ids": {
							Type:        schema.TypeList,
							Description: "List of user ids",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceRateLimitingRuleBlockCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	rule_type := d.Get("rule_type").(string)
	description := d.Get("description").(string)
	environments := d.Get("environments").([]interface{})
	enabled := d.Get("enabled").(bool)
	threshold_configs := d.Get("threshold_configs").([]interface{})
	block_expiry_duration := d.Get("block_expiry_duration").(string)
	alert_severity := d.Get("alert_severity").(string)
	ip_reputation := d.Get("ip_reputation").(string)
	ip_abuse_velocity := d.Get("ip_abuse_velocity").(string)
	label_id_scope := d.Get("label_id_scope").([]interface{})
	endpoint_id_scope := d.Get("endpoint_id_scope").([]interface{})
	req_res_conditions := d.Get("req_res_conditions").([]interface{})
	attribute_based_conditions := d.Get("attribute_based_conditions").([]interface{})
	ip_location_type := d.Get("ip_location_type").([]interface{})
	ip_address := d.Get("ip_address").([]interface{})
	email_domain := d.Get("email_domain").([]interface{})
	user_agents := d.Get("user_agents").([]interface{})
	regions := d.Get("regions").([]interface{})
	ip_organisation := d.Get("ip_organisation").([]interface{})
	ip_asn := d.Get("ip_asn").([]interface{})
	ip_connection_type := d.Get("ip_connection_type").([]interface{})
	request_scanner_type := d.Get("request_scanner_type").([]interface{})
	user_id := d.Get("user_id").([]interface{})

	finalThresholdConfigQuery,err := returnFinalThresholdConfigQuery(threshold_configs)
    if err!=nil{
        return fmt.Errorf("err %s",err)
    }

	finalConditionsQuery,err := returnConditionsStringRateLimit(
		ip_reputation,
		ip_abuse_velocity,
		label_id_scope,
		endpoint_id_scope,
		req_res_conditions,
		attribute_based_conditions,
		ip_location_type,
		ip_address,
		email_domain,
		user_agents,
		regions,
		ip_organisation,
		ip_asn,
		ip_connection_type,
		request_scanner_type,
		user_id,
	)
    if err!=nil{
        return fmt.Errorf("error %s",err)
    }
	if finalConditionsQuery == "" {
		return fmt.Errorf("required at least one scope condition")
	}
	finalEnvironmentQuery := ""
	if len(environments) > 0 {
		finalEnvironmentQuery = fmt.Sprintf(ENVIRONMENT_SCOPE_QUERY, common.ReturnQuotedStringList(environments))
	}
	actionsBlockQuery := fmt.Sprintf(`{ eventSeverity: %s }`, alert_severity)
	if block_expiry_duration != "" {
		actionsBlockQuery = fmt.Sprintf(`{ eventSeverity: %s, duration: "%s" }`, alert_severity, block_expiry_duration)
	}
	createRateLimitQuery := fmt.Sprintf(RATE_LIMITING_CREATE_QUERY, finalConditionsQuery, enabled, name, rule_type, actionsBlockQuery, finalThresholdConfigQuery, finalEnvironmentQuery, description)
	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(createRateLimitQuery,meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", createRateLimitQuery)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	id := response["data"].(map[string]interface{})["createRateLimitingRule"].(map[string]interface{})["id"].(string)

	d.SetId(id)

	return nil
}

func resourceRateLimitingRuleBlockRead(d *schema.ResourceData, meta interface{}) error {
    
	return nil 
}

func resourceRateLimitingRuleBlockUpdate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)
	rule_type := d.Get("rule_type").(string)
	description := d.Get("description").(string)
	environments := d.Get("environments").([]interface{})
	enabled := d.Get("enabled").(bool)
	threshold_configs := d.Get("threshold_configs").([]interface{})
	block_expiry_duration := d.Get("block_expiry_duration").(string)
	alert_severity := d.Get("alert_severity").(string)
	ip_reputation := d.Get("ip_reputation").(string)
	ip_abuse_velocity := d.Get("ip_abuse_velocity").(string)
	label_id_scope := d.Get("label_id_scope").([]interface{})
	endpoint_id_scope := d.Get("endpoint_id_scope").([]interface{})
	req_res_conditions := d.Get("req_res_conditions").([]interface{})
	attribute_based_conditions := d.Get("attribute_based_conditions").([]interface{})
	ip_location_type := d.Get("ip_location_type").([]interface{})
	ip_address := d.Get("ip_address").([]interface{})
	email_domain := d.Get("email_domain").([]interface{})
	user_agents := d.Get("user_agents").([]interface{})
	regions := d.Get("regions").([]interface{})
	ip_organisation := d.Get("ip_organisation").([]interface{})
	ip_asn := d.Get("ip_asn").([]interface{})
	ip_connection_type := d.Get("ip_connection_type").([]interface{})
	request_scanner_type := d.Get("request_scanner_type").([]interface{})
	user_id := d.Get("user_id").([]interface{})
    id := d.Id()

	finalThresholdConfigQuery,err := returnFinalThresholdConfigQuery(threshold_configs)
    if err!=nil{
        return fmt.Errorf("err %s",err)
    }

	finalConditionsQuery,err := returnConditionsStringRateLimit(
		ip_reputation,
		ip_abuse_velocity,
		label_id_scope,
		endpoint_id_scope,
		req_res_conditions,
		attribute_based_conditions,
		ip_location_type,
		ip_address,
		email_domain,
		user_agents,
		regions,
		ip_organisation,
		ip_asn,
		ip_connection_type,
		request_scanner_type,
		user_id,
	)
    if err!=nil{
        return fmt.Errorf("error %s",err)
    }
	if finalConditionsQuery == "" {
		return fmt.Errorf("required at least one scope condition")
	}
	finalEnvironmentQuery := ""
	if len(environments) > 0 {
		finalEnvironmentQuery = fmt.Sprintf(ENVIRONMENT_SCOPE_QUERY, common.ReturnQuotedStringList(environments))
	}
	actionsBlockQuery := fmt.Sprintf(`{ eventSeverity: %s }`, alert_severity)
	if block_expiry_duration != "" {
		actionsBlockQuery = fmt.Sprintf(`{ eventSeverity: %s, duration: "%s" }`, alert_severity, block_expiry_duration)
	}
	updateRateLimitQuery := fmt.Sprintf(RATE_LIMITING_UPDATE_QUERY,id, finalConditionsQuery, enabled, name, rule_type, actionsBlockQuery, finalThresholdConfigQuery, finalEnvironmentQuery, description)
	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(updateRateLimitQuery,meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", updateRateLimitQuery)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	updatedId := response["data"].(map[string]interface{})["updateRateLimitingRule"].(map[string]interface{})["id"].(string)

	d.SetId(updatedId)

	return nil
}

func resourceRateLimitingRuleBlockDelete(d *schema.ResourceData, meta interface{}) error {
    id := d.Id()
	query := fmt.Sprintf(DELETE_RATE_LIMIT_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	log.Println(query)
	d.SetId("")
	return nil
}
