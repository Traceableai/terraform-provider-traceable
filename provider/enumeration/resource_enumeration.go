package enumeration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/rate_limiting"
)

func ResourceEnumerationRule() *schema.Resource {
	return &schema.Resource{
		Create:        resourceEnumerationCreate,
		Read:          resourceEnumerationRead,
		Update:        resourceEnumerationUpdate,
		Delete:        resourceEnumerationDelete,
		CustomizeDiff: validateSchema,
		Schema: map[string]*schema.Schema{
			"rule_type": {
				Type:        schema.TypeString,
				Description: "ALERT or BLOCK",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the enumeration rule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the enumeration rule",
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
			"expiry_duration": {
				Type:        schema.TypeString,
				Description: "Block for a given period",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := common.ConvertDurationToSeconds(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be a valid duration in seconds or ISO 8601 format: %s", key, err))
					}
					return
				},
				StateFunc: func(val interface{}) string {
					v := val.(string)
					converted, _ := common.ConvertDurationToSeconds(v)
					return converted
				},
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
			"request_response_single_valued_conditions": {
				Type:        schema.TypeList,
				Description: "Request payload single valued conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request_location": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Host/Http Method/User Agent/Request Body",
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"request_response_multi_valued_conditions": {
				Type:        schema.TypeList,
				Description: "Request payload multi valued conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request_location": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Query Param/Request Body Param/Request Cookie",
						},
						"key_patterns": {
							Type:        schema.TypeList,
							Description: "key operator and value",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "key operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "value for key",
										Required:    true,
									},
								},
							},
						},
						"value_patterns": {
							Type:        schema.TypeList,
							Description: "value operator and value",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "value operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"data_types_conditions": {
				Type:        schema.TypeList,
				Description: "Datatypes/Datasets conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_type_ids": {
							Type:        schema.TypeList,
							Description: "Datatypes you want to include",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"data_sets_ids": {
							Type:        schema.TypeList,
							Description: "Datasets you want to include",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"data_location": {
							Type:        schema.TypeString,
							Description: "Where to look",
							Required:    true,
						},
					},
				},
			},
			"threshold_configs": {
				Type:        schema.TypeList,
				Description: "Threshold configs types for the rule",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_aggregate_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ACROSS_ENDPOINTS/PER_ENDPOINT",
						},
						"user_aggregate_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "PER_USER/ACROSS_USERS",
						},
						"unique_values_allowed": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Count of calls",
						},
						"duration": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Duration for the total call count in 1min(PT60S)",
						},
						"threshold_config_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SENSITIVE_PARAMS/PATH_PARAMS/REQUEST_BODY",
						},
						"sensitive_param_evaluation_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SELECTED_DATA_TYPES/ALL",
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
	thresholdConfigs := d.Get("threshold_configs").([]interface{})
	dataTypesConditions := d.Get("data_types_conditions").([]interface{})
	attributeBasedConditions := d.Get("attribute_based_conditions").([]interface{})
	ipAddress := d.Get("ip_address").([]interface{})
	userId := d.Get("user_id").([]interface{})
	ruleType := d.Get("rule_type")

	expiryDuration := d.Get("expiry_duration").(string)
	if expiryDuration != "" && ruleType != "BLOCK"{
		return fmt.Errorf("expiry_duration not expected here")
	}


	isDataTypesConditionsEmpty := true
	for _, data := range dataTypesConditions {
		isDataTypesConditionsEmpty = false
		dataTypeIds := data.(map[string]interface{})["data_type_ids"].([]interface{})
		dataSetIds := data.(map[string]interface{})["data_sets_ids"].([]interface{})
		if len(dataSetIds) == 0 && len(dataTypeIds) == 0 {
			return fmt.Errorf("atmost one is required")
		}
		if len(dataSetIds) > 0 && len(dataTypeIds) > 0 {
			return fmt.Errorf("atmost one is required")
		}
	}
	if len(userId) > 0 {
		flag1 := false
		flag2 := false
		if userIdRegexes, ok := userId[0].(map[string]interface{})["user_id_regexes"].([]interface{}); ok {
			fmt.Printf("this is len useridregex %d", len(userIdRegexes))
			if len(userIdRegexes) > 0 {
				flag1 = true
			}
		}
		if userIds, ok := userId[0].(map[string]interface{})["user_ids"].([]interface{}); ok {
			fmt.Printf("this is len userid %d", len(userIds))
			if len(userIds) > 0 {
				flag2 = true
			}
		}

		if flag1 && flag2 {
			return fmt.Errorf("required one of user_id_regexes or user_ids")
		}
	}
	if len(ipAddress) > 0 {
		flag1 := false
		flag2 := false
		if IpAddressList, ok := ipAddress[0].(map[string]interface{})["ip_address_list"].([]interface{}); ok {
			if len(IpAddressList) > 0 {
				flag1 = true
			}
		}
		if ipAddressConditionType, ok := ipAddress[0].(map[string]interface{})["ip_address_type"].(string); ok {
			if ipAddressConditionType != "" {
				flag2 = true
			}
		}
		if flag1 && flag2 {
			return fmt.Errorf("required only one from ip_address_list or ip_address_type")
		}
	}

	if len(labelScope) > 0 && len(endpointScope) > 0 {
		return fmt.Errorf("require on of `label_id_scope` or `endpoint_id_scope`")
	}
	for _, thresholdConfig := range thresholdConfigs {
		thresholdConfigData := thresholdConfig.(map[string]interface{})
		thresholdConfigType := thresholdConfigData["threshold_config_type"]
		sensitiveParamEvaluationType := thresholdConfigData["sensitive_param_evaluation_type"].(string)
		if (thresholdConfigType == "REQUEST_BODY" || thresholdConfigType == "PATH_PARAMS") && sensitiveParamEvaluationType != "" {
			return fmt.Errorf("not valid here sensitive_param_evaluation_type")
		} else if thresholdConfigType == "SENSITIVE_PARAMS" {
			if sensitiveParamEvaluationType == "" {
				return fmt.Errorf("required sensitive_param_evaluation_type for SENSITIVE_PARAMS threshold_config_type")
			} else if sensitiveParamEvaluationType == "SELECTED_DATA_TYPES" && isDataTypesConditionsEmpty {
				return fmt.Errorf("no data types selected")
			}
		}
	}
	for _, attBasedCondition := range attributeBasedConditions {
		valueConditionOperator := attBasedCondition.(map[string]interface{})["value_condition_operator"]
		valueConditionValue := attBasedCondition.(map[string]interface{})["value_condition_value"]
		if (valueConditionOperator != "" && valueConditionValue == "") || (valueConditionValue != "" && valueConditionOperator == "") {
			return fmt.Errorf("required both values value_condition_value and value_condition_operator")
		}
	}
	return nil
}

func resourceEnumerationCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	ruleType := d.Get("rule_type").(string)
	description := d.Get("description").(string)
	environments := d.Get("environments").([]interface{})
	enabled := d.Get("enabled").(bool)
	thresholdConfigs := d.Get("threshold_configs").([]interface{})
	expiryDuration := d.Get("expiry_duration").(string)
	alertSeverity := d.Get("alert_severity").(string)
	ipReputation := d.Get("ip_reputation").(string)
	ipAbuseVelocity := d.Get("ip_abuse_velocity").(string)
	labelIdScope := d.Get("label_id_scope").([]interface{})
	endpointIdScope := d.Get("endpoint_id_scope").([]interface{})
	requestResponseSingleValuedConditions := d.Get("request_response_single_valued_conditions").([]interface{})
	requestResponseMultiValuedConditions := d.Get("request_response_multi_valued_conditions").([]interface{})
	dataTypesConditions := d.Get("data_types_conditions").([]interface{})
	attributeBasedConditions := d.Get("attribute_based_conditions").([]interface{})
	ipLocationType := d.Get("ip_location_type").([]interface{})
	ipAddress := d.Get("ip_address").([]interface{})
	emailDomain := d.Get("email_domain").([]interface{})
	regions := d.Get("regions").([]interface{})
	userAgents := d.Get("user_agents").([]interface{})
	ipOrganisation := d.Get("ip_organisation").([]interface{})
	ipAsn := d.Get("ip_asn").([]interface{})
	ipConnectionType := d.Get("ip_connection_type").([]interface{})
	requestScannerType := d.Get("request_scanner_type").([]interface{})
	userId := d.Get("user_id").([]interface{})

	finalThresholdConfigQuery, err := ReturnFinalThresholdConfigQueryEnumeration(thresholdConfigs)
	if err != nil {
		return fmt.Errorf("err %s", err)
	}

	finalConditionsQuery, err := rate_limiting.ReturnConditionsStringRateLimit(
		ipReputation,
		ipAbuseVelocity,
		labelIdScope,
		endpointIdScope,
		requestResponseSingleValuedConditions,
		requestResponseMultiValuedConditions,
		attributeBasedConditions,
		ipLocationType,
		ipAddress,
		emailDomain,
		userAgents,
		regions,
		ipOrganisation,
		ipAsn,
		ipConnectionType,
		requestScannerType,
		userId,
		dataTypesConditions,
	)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	if finalConditionsQuery == "" {
		return fmt.Errorf("required at least one scope condition")
	}
	finalEnvironmentQuery := ""
	if len(environments) > 0 {
		finalEnvironmentQuery = fmt.Sprintf(rate_limiting.ENVIRONMENT_SCOPE_QUERY, common.ReturnQuotedStringList(environments))
	}
	actionsBlockQuery := fmt.Sprintf(`{ eventSeverity: %s }`, alertSeverity)
	if expiryDuration != "" {
		actionsBlockQuery = fmt.Sprintf(`{ eventSeverity: %s, duration: "%s" }`, alertSeverity, expiryDuration)
	}
	createEnumerationQuery := fmt.Sprintf(rate_limiting.RATE_LIMITING_CREATE_QUERY, ENUMERATION_QUERY_KEY, finalConditionsQuery, enabled, name, ruleType, strings.ToLower(ruleType), actionsBlockQuery, finalThresholdConfigQuery, finalEnvironmentQuery, description)
	responseStr, err := common.CallExecuteQuery(createEnumerationQuery, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", createEnumerationQuery)
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createRateLimitingRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	d.SetId(id)
	return nil
}

func resourceEnumerationRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	var response map[string]interface{}
	readQuery := fmt.Sprintf(rate_limiting.FETCH_RATE_LIMIT_RULES, ENUMERATION_QUERY_KEY)
	responseStr, err := common.CallExecuteQuery(readQuery, meta)
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql query %s", readQuery)
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
	d.Set("name", ruleDetails["name"].(string))
	d.Set("enabled", ruleDetails["enabled"].(bool))
	d.Set("description", ruleDetails["description"].(string))
	if thresholdActionConfigs, ok := ruleDetails["thresholdActionConfigs"].([]interface{}); ok {
		firstThresholdActionConfigs := thresholdActionConfigs[0].(map[string]interface{})
		thresholdActions := firstThresholdActionConfigs["actions"].([]interface{})
		firstThresholdActions := thresholdActions[0].(map[string]interface{})
		actionType := firstThresholdActions["actionType"].(string)
		d.Set("rule_type",actionType)
		if ruleTypeConfig, ok := firstThresholdActions[strings.ToLower(actionType)].(map[string]interface{}); ok {
			if duration,ok := ruleTypeConfig["duration"].(string); ok{
				d.Set("expiry_duration", duration)
			}else{
				d.Set("expiry_duration","")
			}
			if alertSev, ok := ruleTypeConfig["eventSeverity"].(string); ok {
				if alertSev != "" {
					d.Set("alert_severity", alertSev)
				}
			}
		}
		thresholdConfigs := firstThresholdActionConfigs["thresholdConfigs"].([]interface{})
		finalThresholdConfigs := []map[string]interface{}{}
		for _, thresholdConfig := range thresholdConfigs {
			var uniqueValuesAllowed float64
			duration := ""
			thresholdConfigType := ""
			sensitiveParamsEvaluationTypeVal := ""
			var thresholdConfigDataMap map[string]interface{}
			thresholdConfigData := thresholdConfig.(map[string]interface{})
			if valueBasedThresholdConfig, ok := thresholdConfigData["valueBasedThresholdConfig"].(map[string]interface{}); ok {
				uniqueValuesAllowed = valueBasedThresholdConfig["uniqueValuesAllowed"].(float64)
				duration = valueBasedThresholdConfig["duration"].(string)
				thresholdConfigType = valueBasedThresholdConfig["valueType"].(string)
				if sensitiveParamsEvaluationType, ok := thresholdConfigData["sensitiveParamsEvaluationType"].(string); ok {
					sensitiveParamsEvaluationTypeVal = sensitiveParamsEvaluationType
				}
			}
			if sensitiveParamsEvaluationTypeVal == "" {
				thresholdConfigDataMap = map[string]interface{}{
					"api_aggregate_type":    thresholdConfigData["apiAggregateType"].(string),
					"user_aggregate_type":   thresholdConfigData["userAggregateType"].(string),
					"threshold_config_type": thresholdConfigType,
					"unique_values_allowed": uniqueValuesAllowed,
					"duration":              duration,
				}
			} else {
				thresholdConfigDataMap = map[string]interface{}{
					"api_aggregate_type":              thresholdConfigData["apiAggregateType"].(string),
					"user_aggregate_type":             thresholdConfigData["userAggregateType"].(string),
					"threshold_config_type":           thresholdConfigType,
					"unique_values_allowed":           uniqueValuesAllowed,
					"duration":                        duration,
					"sensitive_param_evaluation_type": sensitiveParamsEvaluationTypeVal,
				}
			}
			finalThresholdConfigs = append(finalThresholdConfigs, thresholdConfigDataMap)
		}
		d.Set("threshold_configs", finalThresholdConfigs)
	}
	conditionsArray := ruleDetails["conditions"].([]interface{})
	finalReqResSingleValueConditionState := []map[string]interface{}{}
	finalReqResMultiValueConditionState := []map[string]interface{}{}
	finalAttributeBasedConditionState := []map[string]interface{}{}
	finalDataTypeConditionState := []map[string]interface{}{}

	labelIdScopeFlag, endPointIdScopeFlag, ipReputationScopeFlag, ipLocationTypeScopeFlag, ipAbuseVelFlag, ipAddressFlag, emailDomainFlag, userAgentFlag, regionFlag, ipOrgFlag, ipAsnFlag, ipConnTypeFlag, reqScannerFlag, userIdFlag := true, true, true, true, true, true, true, true, true, true, true, true, true, true
	for _, condition := range conditionsArray {
		leafCondition := condition.(map[string]interface{})["leafCondition"].(map[string]interface{})
		conditionType := leafCondition["conditionType"].(string)
		finalConditionState := []map[string]interface{}{}

		switch conditionType {
		case "IP_REPUTATION":
			minIpReputationSeverity := leafCondition["ipReputationCondition"].(map[string]interface{})["minIpReputationSeverity"].(string)
			d.Set("ip_reputation", minIpReputationSeverity)
			ipReputationScopeFlag = false

		case "IP_ABUSE_VELOCITY":
			minIpAbuseVelocity := leafCondition["ipAbuseVelocityCondition"].(map[string]interface{})["minIpAbuseVelocity"].(string)
			d.Set("ip_abuse_velocity", minIpAbuseVelocity)
			ipAbuseVelFlag = false

		case "IP_LOCATION_TYPE":
			ipLocationTypeCondition := leafCondition["ipLocationTypeCondition"].(map[string]interface{})
			ipLocationTypes := ipLocationTypeCondition["ipLocationTypes"].([]interface{})
			excludeIpLocationType := ipLocationTypeCondition["exclude"].(bool)
			ipLocationType := map[string]interface{}{
				"ip_location_types": ipLocationTypes,
				"exclude":           excludeIpLocationType,
			}
			finalConditionState = append(finalConditionState, ipLocationType)
			d.Set("ip_location_type", finalConditionState)
			ipLocationTypeScopeFlag = false

		case "IP_ADDRESS":
			ipAddressCondition := leafCondition["ipAddressCondition"].(map[string]interface{})
			excludeIpAddress := ipAddressCondition["exclude"].(bool)
			var ipAddressObj map[string]interface{}
			if ipAddressConditionType, ok := ipAddressCondition["ipAddressConditionType"].(string); ok {
				ipAddressObj = map[string]interface{}{
					"ip_address_type": ipAddressConditionType,
					"exclude":         excludeIpAddress,
					"ip_address_list": []interface{}{},
				}
			} else {
				rawInputIpData := ipAddressCondition["rawInputIpData"].([]interface{})
				ipAddressObj = map[string]interface{}{
					"ip_address_list": rawInputIpData,
					"exclude":         excludeIpAddress,
					"ip_address_type": "",
				}
			}
			finalConditionState = append(finalConditionState, ipAddressObj)
			d.Set("ip_address", finalConditionState)
			ipAddressFlag = false

		case "EMAIL_DOMAIN":
			emailDomainCondition := leafCondition["emailDomainCondition"].(map[string]interface{})
			emailRegexes := emailDomainCondition["emailRegexes"].([]interface{})
			excludeEmailRegex := emailDomainCondition["exclude"].(bool)
			emailDomainObj := map[string]interface{}{
				"email_domain_regexes": emailRegexes,
				"exclude":              excludeEmailRegex,
			}
			finalConditionState = append(finalConditionState, emailDomainObj)
			d.Set("email_domain", finalConditionState)
			emailDomainFlag = false

		case "USER_ID":
			userIdCondition := leafCondition["userIdCondition"].(map[string]interface{})
			userIdRegexes := userIdCondition["userIdRegexes"].([]interface{})
			userIds := userIdCondition["userIds"].([]interface{})
			excludeUserIdRegexes := userIdCondition["exclude"].(bool)
			var userIdObj map[string]interface{}
			if len(userIdRegexes) > 0 {
				userIdObj = map[string]interface{}{
					"user_id_regexes": userIdRegexes,
					"exclude":         excludeUserIdRegexes,
				}
			} else if len(userIds) > 0 {
				userIdObj = map[string]interface{}{
					"user_ids": userIds,
					"exclude":  excludeUserIdRegexes,
				}
			}
			finalConditionState = append(finalConditionState, userIdObj)
			d.Set("user_id", finalConditionState)
			userIdFlag = false

		case "USER_AGENT":
			userAgentCondition := leafCondition["userAgentCondition"].(map[string]interface{})
			userAgentRegexes := userAgentCondition["userAgentRegexes"].([]interface{})
			excludeUserAgents := userAgentCondition["exclude"].(bool)
			userAgentObj := map[string]interface{}{
				"user_agents_list": userAgentRegexes,
				"exclude":          excludeUserAgents,
			}
			finalConditionState = append(finalConditionState, userAgentObj)
			d.Set("user_agents", finalConditionState)
			userAgentFlag = false

		case "IP_ORGANISATION":
			ipOrganisationCondition := leafCondition["ipOrganisationCondition"].(map[string]interface{})
			ipOrganisationRegexes := ipOrganisationCondition["ipOrganisationRegexes"].([]interface{})
			excludeIpOrg := ipOrganisationCondition["exclude"].(bool)
			ipOrgObj := map[string]interface{}{
				"ip_organisation_regexes": ipOrganisationRegexes,
				"exclude":                 excludeIpOrg,
			}
			finalConditionState = append(finalConditionState, ipOrgObj)
			d.Set("ip_organisation", finalConditionState)
			ipOrgFlag = false

		case "IP_ASN":
			ipAsnCondition := leafCondition["ipAsnCondition"].(map[string]interface{})
			ipAsnRegexes := ipAsnCondition["ipAsnRegexes"].([]interface{})
			excludeIpAsn := ipAsnCondition["exclude"].(bool)
			ipAsnObj := map[string]interface{}{
				"ip_asn_regexes": ipAsnRegexes,
				"exclude":        excludeIpAsn,
			}
			finalConditionState = append(finalConditionState, ipAsnObj)
			d.Set("ip_asn", finalConditionState)
			ipAsnFlag = false

		case "IP_CONNECTION_TYPE":
			ipConnectionTypeCondition := leafCondition["ipConnectionTypeCondition"].(map[string]interface{})
			ipConnectionTypes := ipConnectionTypeCondition["ipConnectionTypes"].([]interface{})
			excludeIpConnection := ipConnectionTypeCondition["exclude"].(bool)
			ipConnectionObj := map[string]interface{}{
				"ip_connection_type_list": ipConnectionTypes,
				"exclude":                 excludeIpConnection,
			}
			finalConditionState = append(finalConditionState, ipConnectionObj)
			d.Set("ip_connection_type", finalConditionState)
			ipConnTypeFlag = false

		case "DATATYPE":
			dataTypeConditionMap := leafCondition["datatypeCondition"].(map[string]interface{})
			dataSetsIds := dataTypeConditionMap["datasetIds"].([]interface{})
			datatypeIds := dataTypeConditionMap["datatypeIds"].([]interface{})
			location := "REQUEST_RESPONSE"
			if dataLocation, ok := dataTypeConditionMap["dataLocation"].(string); ok {
				location = dataLocation
			}
			dataTypeConditionsObj := map[string]interface{}{
				"data_type_ids": datatypeIds,
				"data_sets_ids": dataSetsIds,
				"data_location": location,
			}
			finalDataTypeConditionState = append(finalDataTypeConditionState, dataTypeConditionsObj)

		case "REQUEST_SCANNER_TYPE":
			requestScannerTypeCondition := leafCondition["requestScannerTypeCondition"].(map[string]interface{})
			scannerTypes := requestScannerTypeCondition["scannerTypes"].([]interface{})
			excludeScanner := requestScannerTypeCondition["exclude"].(bool)
			reqScannerObj := map[string]interface{}{
				"scanner_types_list": scannerTypes,
				"exclude":            excludeScanner,
			}
			finalConditionState = append(finalConditionState, reqScannerObj)
			d.Set("request_scanner_type", finalConditionState)
			reqScannerFlag = false

		case "REGION":
			regionCondition := leafCondition["regionCondition"].(map[string]interface{})
			regionIdentifiers := regionCondition["regionIdentifiers"].([]interface{})
			excludeRegion := regionCondition["exclude"].(bool)
			var regionCodes []interface{}
			for _, region := range regionIdentifiers {
				regionCodes = append(regionCodes, region.(map[string]interface{})["countryIsoCode"])
			}
			regionObj := map[string]interface{}{
				"regions_ids": regionCodes,
				"exclude":     excludeRegion,
			}
			finalConditionState = append(finalConditionState, regionObj)
			d.Set("regions", finalConditionState)
			regionFlag = false

		case "KEY_VALUE":
			keyValueCondition := leafCondition["keyValueCondition"].(map[string]interface{})
			metadataType := keyValueCondition["metadataType"].(string)
			if metadataType == "TAG" {
				keyCondition := keyValueCondition["keyCondition"].(map[string]interface{})
				keyConditionOperator := keyCondition["operator"].(string)
				keyConditionValue := keyCondition["value"].(string)
				if valueCondition, ok := keyValueCondition["valueCondition"].(map[string]interface{}); ok {
					if valueConditionOperator, ok := valueCondition["operator"].(string); ok {
						valueConditionValue := valueCondition["value"].(string)
						keyValueObj := map[string]interface{}{
							"key_condition_operator":   keyConditionOperator,
							"key_condition_value":      keyConditionValue,
							"value_condition_operator": valueConditionOperator,
							"value_condition_value":    valueConditionValue,
						}
						finalAttributeBasedConditionState = append(finalAttributeBasedConditionState, keyValueObj)
					}
				} else {
					keyValueObj := map[string]interface{}{
						"key_condition_operator": keyConditionOperator,
						"key_condition_value":    keyConditionValue,
					}
					finalAttributeBasedConditionState = append(finalAttributeBasedConditionState, keyValueObj)
				}
			} else {
				valuePatternObjSlice := []map[string]interface{}{}
				keyPatternObjSlice := []map[string]interface{}{}
				if keyCondition,ok := keyValueCondition["keyCondition"].(map[string]interface{});ok{
					keyPatternObj := map[string]interface{}{
						"operator" : keyCondition["operator"].(string),
						"value" : keyCondition["value"].(string),
					}
					keyPatternObjSlice = append(keyPatternObjSlice, keyPatternObj)
					if valueCondition,ok := keyValueCondition["valueCondition"].(map[string]interface{});ok{
						valuePatternObj := map[string]interface{}{
							"operator" : valueCondition["operator"].(string),
							"value" : valueCondition["value"].(string),
						}
						valuePatternObjSlice = append(valuePatternObjSlice, valuePatternObj)
					}
					reqPayloadMultiValuedObj := map[string]interface{}{
						"request_location": metadataType,
						"key_patterns" : keyPatternObjSlice,
						"value_patterns" : valuePatternObjSlice,
					}
					finalReqResMultiValueConditionState = append(finalReqResMultiValueConditionState, reqPayloadMultiValuedObj)
				}else{
					valueCondition := keyValueCondition["valueCondition"].(map[string]interface{})
					operator := valueCondition["operator"].(string)
					value := valueCondition["value"].(string)
					reqPayloadSingleValuedObj := map[string]interface{}{
						"request_location": metadataType,
						"operator": operator,
						"value":value,
					}
					finalReqResSingleValueConditionState = append(finalReqResSingleValueConditionState, reqPayloadSingleValuedObj)
				}
			}

		case "SCOPE":
			scopeCondition := leafCondition["scopeCondition"].(map[string]interface{})
			scopeType := scopeCondition["scopeType"].(string)
			if scopeType == "LABEL" {
				labelScope := scopeCondition["labelScope"].(map[string]interface{})
				labelIds := labelScope["labelIds"].([]interface{})
				d.Set("label_id_scope", labelIds)
				labelIdScopeFlag = false
			} else if scopeType == "ENTITY" {
				entityScope := scopeCondition["entityScope"].(map[string]interface{})
				entityIds := entityScope["entityIds"].([]interface{})
				d.Set("endpoint_id_scope", entityIds)
				endPointIdScopeFlag = false
			}
		}
	}

	if ipAddressFlag {
		d.Set("ip_address", []interface{}{})
	}
	if labelIdScopeFlag {
		d.Set("label_id_scope", []interface{}{})
	}
	if endPointIdScopeFlag {
		d.Set("endpoint_id_scope", []interface{}{})
	}
	if ipReputationScopeFlag {
		d.Set("ip_reputation", "")
	}
	if ipLocationTypeScopeFlag {
		d.Set("ip_location_type", []interface{}{})
	}
	if ipAbuseVelFlag {
		d.Set("ip_abuse_velocity", "")
	}
	if emailDomainFlag {
		d.Set("email_domain", []interface{}{})
	}
	if userAgentFlag {
		d.Set("user_agents", []interface{}{})
	}
	if regionFlag {
		d.Set("regions", []interface{}{})
	}
	if ipOrgFlag {
		d.Set("ip_organisation", []interface{}{})
	}
	if userIdFlag {
		d.Set("user_id", []interface{}{})
	}
	if reqScannerFlag {
		d.Set("request_scanner_type", []interface{}{})
	}
	if ipConnTypeFlag {
		d.Set("ip_connection_type", []interface{}{})
	}
	if ipAsnFlag {
		d.Set("ip_asn", []interface{}{})
	}

	var envList []interface{}
	if ruleConfigScope, ok := ruleDetails["ruleConfigScope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleConfigScope["environmentScope"].(map[string]interface{}); ok {
			envList = environmentScope["environmentIds"].([]interface{})
		}
	}
	d.Set("environments", envList)
	d.Set("request_response_single_valued_conditions", finalReqResMultiValueConditionState)
	d.Set("request_response_multi_valued_conditions", finalReqResSingleValueConditionState)
	d.Set("attribute_based_conditions", finalAttributeBasedConditionState)
	d.Set("data_types_conditions", finalDataTypeConditionState)

	return nil
}

func resourceEnumerationUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	ruleType := d.Get("rule_type").(string)
	description := d.Get("description").(string)
	environments := d.Get("environments").([]interface{})
	enabled := d.Get("enabled").(bool)
	thresholdConfigs := d.Get("threshold_configs").([]interface{})
	dataTypesConditions := d.Get("data_types_conditions").([]interface{})
	expiryDuration := d.Get("expiry_duration").(string)
	alertSeverity := d.Get("alert_severity").(string)
	ipReputation := d.Get("ip_reputation").(string)
	ipAbuseVelocity := d.Get("ip_abuse_velocity").(string)
	labelIdScope := d.Get("label_id_scope").([]interface{})
	endpointIdScope := d.Get("endpoint_id_scope").([]interface{})
	requestResponseSingleValuedConditions := d.Get("request_response_single_valued_conditions").([]interface{})
	requestResponseMultiValuedConditions := d.Get("request_response_multi_valued_conditions").([]interface{})
	attributeBasedConditions := d.Get("attribute_based_conditions").([]interface{})
	ipLocationType := d.Get("ip_location_type").([]interface{})
	ipAddress := d.Get("ip_address").([]interface{})
	emailDomain := d.Get("email_domain").([]interface{})
	regions := d.Get("regions").([]interface{})
	userAgents := d.Get("user_agents").([]interface{})
	ipOrganisation := d.Get("ip_organisation").([]interface{})
	ipAsn := d.Get("ip_asn").([]interface{})
	ipConnectionType := d.Get("ip_connection_type").([]interface{})
	requestScannerType := d.Get("request_scanner_type").([]interface{})
	userId := d.Get("user_id").([]interface{})
	id := d.Id()

	finalThresholdConfigQuery, err := ReturnFinalThresholdConfigQueryEnumeration(thresholdConfigs)
	if err != nil {
		return fmt.Errorf("err %s", err)
	}

	finalConditionsQuery, err := rate_limiting.ReturnConditionsStringRateLimit(
		ipReputation,
		ipAbuseVelocity,
		labelIdScope,
		endpointIdScope,
		requestResponseSingleValuedConditions,
		requestResponseMultiValuedConditions,
		attributeBasedConditions,
		ipLocationType,
		ipAddress,
		emailDomain,
		userAgents,
		regions,
		ipOrganisation,
		ipAsn,
		ipConnectionType,
		requestScannerType,
		userId,
		dataTypesConditions,
	)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	if finalConditionsQuery == "" {
		return fmt.Errorf("required at least one scope condition")
	}
	finalEnvironmentQuery := ""
	if len(environments) > 0 {
		finalEnvironmentQuery = fmt.Sprintf(rate_limiting.ENVIRONMENT_SCOPE_QUERY, common.ReturnQuotedStringList(environments))
	}
	actionsBlockQuery := fmt.Sprintf(`{ eventSeverity: %s }`, alertSeverity)
	if expiryDuration != "" {
		actionsBlockQuery = fmt.Sprintf(`{ eventSeverity: %s, duration: "%s" }`, alertSeverity, expiryDuration)
	}
	updateRateLimitQuery := fmt.Sprintf(rate_limiting.RATE_LIMITING_UPDATE_QUERY, id, ENUMERATION_QUERY_KEY, finalConditionsQuery, enabled, name, ruleType, strings.ToLower(ruleType), actionsBlockQuery, finalThresholdConfigQuery, finalEnvironmentQuery, description)
	responseStr, err := common.CallExecuteQuery(updateRateLimitQuery, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", updateRateLimitQuery)
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateRateLimitingRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	d.SetId(updatedId)

	return nil
}

func resourceEnumerationDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(rate_limiting.DELETE_RATE_LIMIT_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	log.Println(query)
	d.SetId("")
	return nil
}
