package rate_limiting

import (
	"encoding/json"
	"fmt"
	"log"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceRateLimitingRuleBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceRateLimitingRuleBlockCreate,
		Read:   resourceRateLimitingRuleBlockRead,
		Update: resourceRateLimitingRuleBlockUpdate,
		Delete: resourceRateLimitingRuleBlockDelete,
		CustomizeDiff: validateSchema,
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
func validateSchema(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	labelScope := d.Get("label_id_scope").([]interface{})
	endpointScope := d.Get("endpoint_id_scope").([]interface{})
	threshold_configs := d.Get("threshold_configs").([]interface{})
	attribute_based_conditions := d.Get("attribute_based_conditions").([]interface{})
	ip_address := d.Get("ip_address").([]interface{})
	user_id := d.Get("user_id").([]interface{})

	if len(user_id)>0{
		flag1 := false
		flag2 :=false
		if userIdRegexes,ok := user_id[0].(map[string]interface{})["user_id_regexes"].([]interface{}) ; ok {
			fmt.Printf("this is len useridregex %d",len(userIdRegexes))
			if len(userIdRegexes)>0{
				flag1=true
			}
		}
		if userIds,ok := user_id[0].(map[string]interface{})["user_ids"].([]interface{}); ok {
			fmt.Printf("this is len userid %d",len(userIds))
			if len(userIds)>0{
				flag2=true
			}
		}
		
		if flag1 && flag2 {
			return fmt.Errorf("required one of user_id_regexes or user_ids")
		}
	}
	if len(ip_address)>0{
		flag1:=false
		flag2:=false
		if IpAddressList,ok := ip_address[0].(map[string]interface{})["ip_address_list"].([]interface{}); ok {
			if len(IpAddressList)>0{
				flag1=true
			}
		}
		if ipAddressConditionType,ok := ip_address[0].(map[string]interface{})["ip_address_type"].(string); ok {
			if ipAddressConditionType!=""{
				flag2=true
			}
		}
		if flag1 && flag2 {
			return fmt.Errorf("required only one from ip_address_list or ip_address_type")
		}
	}

	if len(labelScope) > 0 && len(endpointScope) > 0 {
		return fmt.Errorf("require on of `label_id_scope` or `endpoint_id_scope`")
	}
	for _, thresholdConfig := range threshold_configs {
		thresholdConfigData := thresholdConfig.(map[string]interface{})
		thresholdConfigType := thresholdConfigData["threshold_config_type"]
		dynamicMeanCalculationDuration := thresholdConfigData["dynamic_mean_calculation_duration"].(string)
		if thresholdConfigType == "ROLLING_WINDOW" &&  dynamicMeanCalculationDuration != ""{
			return fmt.Errorf("not valid here dynamicMeanCalculationDuration")
		} else if thresholdConfigType == "DYNAMIC" {
			if dynamicMeanCalculationDuration == "" {
				return fmt.Errorf("required dynamic_mean_calculation_duration for dynamic threshold_config_type")
			}
		}
	}
	for _,attBasedCondition := range attribute_based_conditions {
		valueConditionOperator := attBasedCondition.(map[string]interface{})["value_condition_operator"]
		valueConditionValue := attBasedCondition.(map[string]interface{})["value_condition_value"]
		if (valueConditionOperator!="" && valueConditionValue=="") || (valueConditionValue!="" && valueConditionOperator==""){
			return fmt.Errorf("required both values value_condition_value and value_condition_operator")
		}
	}
	return nil
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

	finalThresholdConfigQuery, err := returnFinalThresholdConfigQuery(threshold_configs)
	if err != nil {
		return fmt.Errorf("err %s", err)
	}

	finalConditionsQuery, err := returnConditionsStringRateLimit(
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
	if err != nil {
		return fmt.Errorf("error %s", err)
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
	responseStr, err := common.CallExecuteQuery(createRateLimitQuery, meta)
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
	id := d.Id()
	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(FETCH_RATE_LIMIT_RULES, meta)
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql query %s", FETCH_RATE_LIMIT_RULES)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	ruleDetails := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "rateLimitingRules", id)
	if len(ruleDetails) == 0 {
		d.SetId("")
		return nil
	}
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
    log.Printf("fetching from read %s", ruleDetails)
	d.Set("name",ruleDetails["name"].(string))
	d.Set("enabled",ruleDetails["enabled"].(bool))
	d.Set("description",ruleDetails["description"].(string))
	if thresholdActionConfigs,ok := ruleDetails["thresholdActionConfigs"].([]interface{}); ok {
		firstThresholdActionConfigs := thresholdActionConfigs[0].(map[string]interface{})
		thresholdActions := firstThresholdActionConfigs["actions"].([]interface{})
		firstThresholdActions := thresholdActions[0].(map[string]interface{})
		d.Set("rule_type",firstThresholdActions["actionType"])
		if blockingConfig,ok := firstThresholdActions["block"].(map[string]interface{}); ok {
			d.Set("block_expiry_duration",blockingConfig["duration"])
			
			if blockingSeverity,ok := blockingConfig["eventSeverity"].(string); ok {
				if blockingSeverity!=""{
					d.Set("alert_severity",blockingSeverity)
				}
			}
		}
		thresholdConfigs := firstThresholdActionConfigs["thresholdConfigs"].([]interface{})
		finalThresholdConfigs := []map[string]interface{}{}
		for _,thresholdConfig := range thresholdConfigs {
			var count_allowed float64
			duration := ""
			dynamic_mean_calculation_duration:=""
			thresholdConfigData := thresholdConfig.(map[string]interface{})
			if rollingWindowThresholdConfig,ok := thresholdConfigData["rollingWindowThresholdConfig"].(map[string]interface{}); ok {
				if rollingWindowCountAllowed,ok := rollingWindowThresholdConfig["countAllowed"].(float64); ok{
					count_allowed=rollingWindowCountAllowed
				}
				if rollingWindowDuration,ok := rollingWindowThresholdConfig["duration"].(string); ok{
					duration=rollingWindowDuration
				}

			}
			if dynamicThresholdConfig,ok := thresholdConfigData["dynamicThresholdConfig"].(map[string]interface{}); ok {
				if dynamicCountAllowed,ok := dynamicThresholdConfig["percentageExceedingMeanAllowed"].(float64); ok{
					count_allowed=dynamicCountAllowed
				}
				if dynamicDuration,ok := dynamicThresholdConfig["duration"].(string); ok{
					duration=dynamicDuration
				}
				if dynamicMeanCalculationDuration,ok := dynamicThresholdConfig["meanCalculationDuration"].(string); ok{
					dynamic_mean_calculation_duration=dynamicMeanCalculationDuration
				}
			}
			thresholdConfigDataMap := map[string]interface{}{
				"api_aggregate_type" : thresholdConfigData["apiAggregateType"].(string),
				"rolling_window_count_allowed" : count_allowed,
				"rolling_window_duration" : duration,
				"threshold_config_type" : thresholdConfigData["thresholdConfigType"].(string),
				"dynamic_mean_calculation_duration" : dynamic_mean_calculation_duration,
			}
			finalThresholdConfigs=append(finalThresholdConfigs,thresholdConfigDataMap)
		}
		d.Set("threshold_configs",finalThresholdConfigs)
	}
	conditionsArray := ruleDetails["conditions"].([]interface{})
	finalReqResConditionsState := []map[string]interface{}{}
	finalAttributeBasedConditionState := []map[string]interface{}{}

	labelIdScopeFlag, endPointIdScopeFlag, ipReputationScopeFlag, ipLocationTypeScopeFlag, ipAbuseVelFlag,ipAddressFlag, emailDomainFlag, userAgentFlag, regionFlag, ipOrgFlag, ipAsnFlag, ipConnTypeFlag,reqScannerFlag, userIdFlag := true, true, true, true, true, true, true, true, true, true, true, true, true, true
	for _,condition := range conditionsArray {
		leafCondition := condition.(map[string]interface{})["leafCondition"].(map[string]interface{})
		conditionType := leafCondition["conditionType"].(string)
		finalConditionState := []map[string]interface{}{}
		if conditionType == "IP_REPUTATION"{
			minIpReputationSeverity := leafCondition["ipReputationCondition"].(map[string]interface{})["minIpReputationSeverity"].(string)
			d.Set("ip_reputation",minIpReputationSeverity)
			ipReputationScopeFlag=false
		}else if conditionType == "IP_ABUSE_VELOCITY" {
			minIpAbuseVelocity := leafCondition["ipAbuseVelocityCondition"].(map[string]interface{})["minIpAbuseVelocity"].(string)
			d.Set("ip_abuse_velocity",minIpAbuseVelocity)
			ipAbuseVelFlag=false
		}else if conditionType == "IP_LOCATION_TYPE" {
			ipLocationTypeCondition := leafCondition["ipLocationTypeCondition"].(map[string]interface{})
			ipLocationTypes := ipLocationTypeCondition["ipLocationTypes"].([]interface{})
			excludeIpLocationType := ipLocationTypeCondition["exclude"].(bool)
			ip_location_type := map[string]interface{}{
				"ip_location_types":ipLocationTypes,
				"exclude":excludeIpLocationType,	
			}
			finalConditionState=append(finalConditionState,ip_location_type)
			d.Set("ip_location_type",finalConditionState)
			ipLocationTypeScopeFlag=false
		}else if conditionType == "IP_ADDRESS" {
			ipAddressCondition := leafCondition["ipAddressCondition"].(map[string]interface{})
			excludeIpAddress := ipAddressCondition["exclude"].(bool)
			var ipAddressObj map[string]interface{}
			if ipAddressConditionType,ok := ipAddressCondition["ipAddressConditionType"].(string) ; ok{
				ipAddressObj = map[string]interface{}{
					"ip_address_type": ipAddressConditionType,
					"exclude": excludeIpAddress,
					"ip_address_list": []interface{}{},
				}
			}else {
				rawInputIpData := ipAddressCondition["rawInputIpData"].([]interface{})
				ipAddressObj = map[string]interface{}{
					"ip_address_list": rawInputIpData,
					"exclude": excludeIpAddress,
					"ip_address_type": "",
				}
			}
			finalConditionState=append(finalConditionState,ipAddressObj)
			d.Set("ip_address",finalConditionState)
			ipAddressFlag=false
		}else if conditionType == "EMAIL_DOMAIN" {
			emailDomainCondition :=	leafCondition["emailDomainCondition"].(map[string]interface{})
			emailRegexes := emailDomainCondition["emailRegexes"].([]interface{})
			excludeEmailRegex := emailDomainCondition["exclude"].(bool)
			emailDomainObj := map[string]interface{}{
				"email_domain_regexes": emailRegexes,
				"exclude": excludeEmailRegex,
			}
			finalConditionState=append(finalConditionState, emailDomainObj)
			d.Set("email_domain",finalConditionState)
			emailDomainFlag=false
		}else if conditionType == "USER_ID"{
			userIdCondition := leafCondition["userIdCondition"].(map[string]interface{})
			userIdRegexes := userIdCondition["userIdRegexes"].([]interface{})
			userIds := userIdCondition["userIds"].([]interface{})
			excludeUserIdRegexes := userIdCondition["exclude"].(bool)
			userIdObj := map[string]interface{}{}
			if len(userIdRegexes)>0{
				userIdObj = map[string]interface{}{
					"user_id_regexes": userIdRegexes,
					"exclude": excludeUserIdRegexes,
				}
			} else if len(userIds)>0 {
				userIdObj = map[string]interface{}{
					"user_ids": userIds,
					"exclude": excludeUserIdRegexes,
				}
			}
			finalConditionState=append(finalConditionState, userIdObj)
			d.Set("user_id",finalConditionState)
			userIdFlag=false
		}else if conditionType == "USER_AGENT" {
			userAgentCondition :=	leafCondition["userAgentCondition"].(map[string]interface{})
			userAgentRegexes := userAgentCondition["userAgentRegexes"].([]interface{})
			excludeUserAgents := userAgentCondition["exclude"].(bool)
			userAgentObj := map[string]interface{}{
				"user_agents_list": userAgentRegexes,
				"exclude": excludeUserAgents,
			}
			finalConditionState=append(finalConditionState, userAgentObj)
			d.Set("user_agents",finalConditionState)
			userAgentFlag=false
		}else if conditionType == "IP_ORGANISATION" {
			ipOrganisationCondition :=	leafCondition["ipOrganisationCondition"].(map[string]interface{})
			ipOrganisationRegexes := ipOrganisationCondition["ipOrganisationRegexes"].([]interface{})
			excludeIpOrg := ipOrganisationCondition["exclude"].(bool)
			ipOrgObj := map[string]interface{}{
				"ip_organisation_regexes": ipOrganisationRegexes,
				"exclude": excludeIpOrg,
			}
			finalConditionState=append(finalConditionState, ipOrgObj)
			d.Set("ip_organisation",finalConditionState)
			ipOrgFlag=false
		}else if conditionType == "IP_ASN" {
			ipAsnCondition :=	leafCondition["ipAsnCondition"].(map[string]interface{})
			ipAsnRegexes := ipAsnCondition["ipAsnRegexes"].([]interface{})
			excludeIpAsn := ipAsnCondition["exclude"].(bool)
			ipAsnObj := map[string]interface{}{
				"ip_asn_regexes": ipAsnRegexes,
				"exclude": excludeIpAsn,
			}
			finalConditionState=append(finalConditionState, ipAsnObj)
			d.Set("ip_asn",finalConditionState)
			ipAsnFlag=false
		}else if conditionType == "IP_CONNECTION_TYPE" {
			ipConnectionTypeCondition := leafCondition["ipConnectionTypeCondition"].(map[string]interface{})
			ipConnectionTypes := ipConnectionTypeCondition["ipConnectionTypes"].([]interface{})
			excludeIpConnection := ipConnectionTypeCondition["exclude"].(bool)
			ipConnectionObj := map[string]interface{}{
				"ip_connection_type_list": ipConnectionTypes,
				"exclude": excludeIpConnection,
			}
			finalConditionState=append(finalConditionState, ipConnectionObj)
			d.Set("ip_connection_type",finalConditionState)
			ipConnTypeFlag=false
		}else if conditionType == "REQUEST_SCANNER_TYPE" {
			requestScannerTypeCondition := leafCondition["requestScannerTypeCondition"].(map[string]interface{})
			scannerTypes := requestScannerTypeCondition["scannerTypes"].([]interface{})
			excludeScanner := requestScannerTypeCondition["exclude"].(bool)
			reqScannerObj := map[string]interface{}{
				"scanner_types_list": scannerTypes,
				"exclude": excludeScanner,
			}
			finalConditionState=append(finalConditionState, reqScannerObj)
			d.Set("request_scanner_type",finalConditionState)
			reqScannerFlag=false
		}else if conditionType == "REGION"{
			regionCondition := leafCondition["regionCondition"].(map[string]interface{})
			regionIdentifiers := regionCondition["regionIdentifiers"].([]interface{})
			excludeRegion := regionCondition["exclude"].(bool)
			var regionCodes []interface{}
			for _,region := range regionIdentifiers {
				regionCodes=append(regionCodes,region.(map[string]interface{})["countryIsoCode"])
			}
			regionObj := map[string]interface{}{
				"regions_ids": regionCodes,
				"exclude": excludeRegion,
			}
			finalConditionState = append(finalConditionState, regionObj)
			d.Set("regions",finalConditionState)
			regionFlag=false
		}else if conditionType == "KEY_VALUE" {
			keyValueCondition := leafCondition["keyValueCondition"].(map[string]interface{})
			metadataType := keyValueCondition["metadataType"].(string)
			if metadataType == "TAG" {
				keyCondition := keyValueCondition["keyCondition"].(map[string]interface{})
				keyConditionOperator := keyCondition["operator"].(string)
				keyConditionValue := keyCondition["value"].(string)
				if valueCondition,ok := keyValueCondition["valueCondition"].(map[string]interface{}); ok {
					if valueConditionOperator,ok := valueCondition["operator"].(string); ok {
						valueConditionValue := valueCondition["value"].(string)
						keyValueObj := map[string]interface{}{	
							"key_condition_operator": keyConditionOperator,
							"key_condition_value": keyConditionValue,
							"value_condition_operator": valueConditionOperator,
							"value_condition_value": valueConditionValue,
						}
						finalAttributeBasedConditionState=append(finalAttributeBasedConditionState, keyValueObj)
					}
				}else{
					keyValueObj := map[string]interface{}{	
						"key_condition_operator": keyConditionOperator,
						"key_condition_value": keyConditionValue,
					}
					finalAttributeBasedConditionState=append(finalAttributeBasedConditionState, keyValueObj)
				}
			}else {
				valueCondition := keyValueCondition["valueCondition"].(map[string]interface{})
				valueConditionValue := valueCondition["value"].(string)
				valueConditionKey := valueCondition["operator"].(string)
				reqResObj := map[string]interface{}{	
					"metadata_type": metadataType,
					"req_res_operator": valueConditionKey,
					"req_res_value": valueConditionValue,
				}
				finalReqResConditionsState=append(finalReqResConditionsState,reqResObj)
				
			}
		}else if conditionType == "SCOPE" {
			scopeCondition := leafCondition["scopeCondition"].(map[string]interface{})
			scopeType := scopeCondition["scopeType"].(string)
			if scopeType == "LABEL" {
				labelScope := scopeCondition["labelScope"].(map[string]interface{})
				labelIds := labelScope["labelIds"].([]interface{})
				d.Set("label_id_scope",labelIds)
				labelIdScopeFlag=false
			}else if scopeType == "ENTITY" {
				entityScope := scopeCondition["entityScope"].(map[string]interface{})
				entityIds := entityScope["entityIds"].([]interface{})
				d.Set("endpoint_id_scope",entityIds)
				endPointIdScopeFlag=false
			}
		}
	}
	if ipAddressFlag{
		d.Set("ip_address",[]interface{}{})
	}
	if labelIdScopeFlag{
		d.Set("label_id_scope",[]interface{}{})
	}
	if endPointIdScopeFlag {
		d.Set("endpoint_id_scope",[]interface{}{})
	}
	if ipReputationScopeFlag {
		d.Set("ip_reputation","")
	}
	if ipLocationTypeScopeFlag {
		d.Set("ip_location_type",[]interface{}{})
	}
	if ipAbuseVelFlag {
		d.Set("ip_abuse_velocity","")
	}
	if  emailDomainFlag {
		d.Set("email_domain",[]interface{}{})
	}
	if  userAgentFlag {
		d.Set("user_agents",[]interface{}{})
	}
	if  regionFlag {
		d.Set("regions",[]interface{}{})
	}
	if  ipOrgFlag {
		d.Set("ip_organisation",[]interface{}{})
	}
	if  userIdFlag {
		d.Set("user_id",[]interface{}{})
	}
	if  reqScannerFlag {
		d.Set("request_scanner_type",[]interface{}{})
	}
	if  ipConnTypeFlag {
		d.Set("ip_connection_type",[]interface{}{})
	}
	if  ipAsnFlag {
		d.Set("ip_asn",[]interface{}{})
	}

	var envList []interface{}
	if ruleConfigScope,ok := ruleDetails["ruleConfigScope"].(map[string]interface{}); ok{
		if environmentScope,ok := ruleConfigScope["environmentScope"].(map[string]interface{}); ok {
			envList = environmentScope["environmentIds"].([]interface{})
		}
	}
	d.Set("environments",envList)
	d.Set("req_res_conditions",finalReqResConditionsState)
	d.Set("attribute_based_conditions",finalAttributeBasedConditionState)
	
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

	finalThresholdConfigQuery, err := returnFinalThresholdConfigQuery(threshold_configs)
	if err != nil {
		return fmt.Errorf("err %s", err)
	}

	finalConditionsQuery, err := returnConditionsStringRateLimit(
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
	if err != nil {
		return fmt.Errorf("error %s", err)
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
	updateRateLimitQuery := fmt.Sprintf(RATE_LIMITING_UPDATE_QUERY, id, finalConditionsQuery, enabled, name, rule_type, actionsBlockQuery, finalThresholdConfigQuery, finalEnvironmentQuery, description)
	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(updateRateLimitQuery, meta)
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
