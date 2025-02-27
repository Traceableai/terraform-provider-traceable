package dlp

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

func ResourceDlpUserBasedRule() *schema.Resource {
	return &schema.Resource{
		Create:        ResourceDlpUserBasedCreate,
		Read:          ResourceDlpUserBasedRead,
		Update:        ResourceDlpUserBasedUpdate,
		Delete:        ResourceDlpUserBasedDelete,
		CustomizeDiff: validateSchema,
		Schema: map[string]*schema.Schema{
			"rule_type": {
				Type:        schema.TypeString,
				Description: "ALERT or BLOCK",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the dlp rule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the dlp rule",
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
				Description: "Block for a given period of time",
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
				Required:    true,
				MinItems: 1,
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
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rolling_window_threshold_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_aggregate_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "User aggregation type: PER_USER/ACROSS_USERS",
									},
									"api_aggregate_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "API aggregation type: PER_USER/ACROSS_USERS",
									},
									"count_allowed": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Maximum number of calls allowed",
									},
									"duration": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Duration for the total call count (e.g., PT60S for 1 minute)",
									},
								},
							},
						},
						"dynamic_threshold_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_aggregate_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "User aggregation type: PER_USER/ACROSS_USERS",
									},
									"api_aggregate_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "API aggregation type: PER_USER/ACROSS_USERS",
									},
									"percentage_exceeding_mean_allowed": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Percentage exceeding the calculated mean allowed",
									},
									"mean_calculation_duration": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Duration for mean calculation (e.g., PT60S for 1 minute)",
									},
									"duration": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Overall duration for evaluation (e.g., PT60S for 1 minute)",
									},
								},
							},
						},
						"value_based_threshold_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_aggregate_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "User aggregation type: PER_USER/ACROSS_USERS",
									},
									"api_aggregate_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "API aggregation type: PER_USER/ACROSS_USERS",
									},
									"unique_values_allowed": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Number of unique values allowed",
									},
									"duration": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Duration for the total evaluation period (e.g., PT60S for 1 minute)",
									},
									"sensitive_params_evaluation_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of sensitive parameter evaluation",
									},
								},
							},
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
func validateSchema(ctx context.Context, rData *schema.ResourceDiff, meta interface{}) error {
	labelScope := rData.Get("label_id_scope").([]interface{})
	endpointScope := rData.Get("endpoint_id_scope").([]interface{})
	dataTypesConditions := rData.Get("data_types_conditions").([]interface{})
	attributeBasedConditions := rData.Get("attribute_based_conditions").([]interface{})
	ipAddress := rData.Get("ip_address").([]interface{})
	userId := rData.Get("user_id").([]interface{})
	ruleType := rData.Get("rule_type")
	
	expiryDuration := rData.Get("expiry_duration").(string)
	if expiryDuration != "" && ruleType != "BLOCK"{
		return fmt.Errorf("expiry_duration not expected here")
	}

	for _, data := range dataTypesConditions {
		dataTypeIds := data.(map[string]interface{})["data_type_ids"]
		dataSetIds := data.(map[string]interface{})["data_sets_ids"]
		if dataSetIds == "" && dataTypeIds == "" {
			return fmt.Errorf("atmost one is required")
		}
		if dataSetIds != "" && dataTypeIds != "" {
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
	
	for _, attBasedCondition := range attributeBasedConditions {
		valueConditionOperator := attBasedCondition.(map[string]interface{})["value_condition_operator"]
		valueConditionValue := attBasedCondition.(map[string]interface{})["value_condition_value"]
		if (valueConditionOperator != "" && valueConditionValue == "") || (valueConditionValue != "" && valueConditionOperator == "") {
			return fmt.Errorf("required both values value_condition_value and value_condition_operator")
		}
	}
	return nil
}

func ResourceDlpUserBasedCreate(rData *schema.ResourceData, meta interface{}) error {
	name := rData.Get("name").(string)
	ruleType := rData.Get("rule_type").(string)
	description := rData.Get("description").(string)
	environments := rData.Get("environments").([]interface{})
	enabled := rData.Get("enabled").(bool)
	thresholdConfigs := rData.Get("threshold_configs").([]interface{})
	expiryDuration := rData.Get("expiry_duration").(string)
	alertSeverity := rData.Get("alert_severity").(string)
	ipReputation := rData.Get("ip_reputation").(string)
	ipAbuseVelocity := rData.Get("ip_abuse_velocity").(string)
	labelIdScope := rData.Get("label_id_scope").([]interface{})
	endpointIdScope := rData.Get("endpoint_id_scope").([]interface{})
	requestResponseSingleValuedConditions := rData.Get("request_response_single_valued_conditions").([]interface{})
	requestResponseMultiValuedConditions := rData.Get("request_response_multi_valued_conditions").([]interface{})
	dataTypesConditions := rData.Get("data_types_conditions").([]interface{})
	attributeBasedConditions := rData.Get("attribute_based_conditions").([]interface{})
	ipLocationType := rData.Get("ip_location_type").([]interface{})
	ipAddress := rData.Get("ip_address").([]interface{})
	emailDomain := rData.Get("email_domain").([]interface{})
	regions := rData.Get("regions").([]interface{})
	userAgents := rData.Get("user_agents").([]interface{})
	ipOrganisation := rData.Get("ip_organisation").([]interface{})
	ipAsn := rData.Get("ip_asn").([]interface{})
	ipConnectionType := rData.Get("ip_connection_type").([]interface{})
	requestScannerType := rData.Get("request_scanner_type").([]interface{})
	userId := rData.Get("user_id").([]interface{})

	finalThresholdConfigQuery, err := ReturnFinalThresholdConfigQueryDlp(thresholdConfigs)
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
	createEnumerationQuery := fmt.Sprintf(rate_limiting.RATE_LIMITING_CREATE_QUERY,DLP_KEY, finalConditionsQuery, enabled, name, ruleType, strings.ToLower(ruleType), actionsBlockQuery, finalThresholdConfigQuery, finalEnvironmentQuery, description)
	log.Printf("This is the graphql query %s", createEnumerationQuery)
	responseStr, err := common.CallExecuteQuery(createEnumerationQuery, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	id,err := common.GetIdFromResponse(responseStr,"")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	rData.SetId(id)
	return nil
}

func ResourceDlpUserBasedRead(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	var response map[string]interface{}
	readQuery := fmt.Sprintf(rate_limiting.FETCH_RATE_LIMIT_RULES,DLP_KEY)
	log.Printf("This is the graphql query %s", readQuery)
	responseStr, err := common.CallExecuteQuery(readQuery, meta)
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	ruleDetails := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "rateLimitingRules", id)
	if len(ruleDetails) == 0 {
		rData.SetId("")
		return nil
	}
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	log.Printf("fetching from read %s", ruleDetails)
	rData.Set("name", ruleDetails["name"].(string))
	rData.Set("enabled", ruleDetails["enabled"].(bool))
	rData.Set("description", ruleDetails["description"].(string))
	if thresholdActionConfigs, ok := ruleDetails["thresholdActionConfigs"].([]interface{}); ok {
		firstThresholdActionConfigs := thresholdActionConfigs[0].(map[string]interface{})
		thresholdActions := firstThresholdActionConfigs["actions"].([]interface{})
		firstThresholdActions := thresholdActions[0].(map[string]interface{})
		actionType := firstThresholdActions["actionType"].(string)
		rData.Set("rule_type",actionType)
		if ruleTypeConfig, ok := firstThresholdActions[strings.ToLower(actionType)].(map[string]interface{}); ok {
			if duration,ok := ruleTypeConfig["duration"].(string); ok{
				rData.Set("expiry_duration", duration)
			}else{
				rData.Set("expiry_duration","")
			}
			if alertSev, ok := ruleTypeConfig["eventSeverity"].(string); ok {
				if alertSev != "" {
					rData.Set("alert_severity", alertSev)
				}
			}
		}
		thresholdConfigs := firstThresholdActionConfigs["thresholdConfigs"].([]interface{})
		finalThresholdConfigs := []map[string]interface{}{}
		rollingWindowThresholdConfigData := []map[string]interface{}{}
		dynamicThresholdConfigData := []map[string]interface{}{}
		valueBasedThresholdConfigData := []map[string]interface{}{}
		for _, thresholdConfig := range thresholdConfigs {
			thresholdConfigData := thresholdConfig.(map[string]interface{})
			thresholdConfigType := thresholdConfigData["thresholdConfigData"].(string)
			switch thresholdConfigType{
				case "ROLLING_WINDOW":
					apiAggregateType := thresholdConfigData["apiAggregateType"].(string)
					userAggregateType := thresholdConfigData["userAggregateType"].(string)
					rollingWindowThresholdConfig := thresholdConfigData["rollingWindowThresholdConfig"].(map[string]interface{})
					countAllowed := rollingWindowThresholdConfig["countAllowed"].(float64)
					duration := rollingWindowThresholdConfig["duration"].(string)
					rollingWinObj := map[string]interface{}{
						"user_aggregate_type":userAggregateType,
						"api_aggregate_type":apiAggregateType,
						"count_allowed":countAllowed,
						"duration":duration,
					}
					rollingWindowThresholdConfigData = append(rollingWindowThresholdConfigData,rollingWinObj)
				case "DYNAMIC":
					apiAggregateType := thresholdConfigData["apiAggregateType"].(string)
					userAggregateType := thresholdConfigData["userAggregateType"].(string)
					dynamicThresholdConfig := thresholdConfigData["dynamicThresholdConfig"].(map[string]interface{})
					percentageExceedingMeanAllowed := dynamicThresholdConfig["percentageExceedingMeanAllowed"].(float64)
					meanCalculationDuration := dynamicThresholdConfig["meanCalculationDuration"].(string)
					duration := dynamicThresholdConfig["duration"].(string)
					dynamicThresholdConfigObj := map[string]interface{}{
						"user_aggregate_type":userAggregateType,
						"api_aggregate_type":apiAggregateType,
						"percentage_exceeding_mean_allowed":percentageExceedingMeanAllowed,
						"mean_calculation_duration":meanCalculationDuration,
						"duration":duration,
					}
					dynamicThresholdConfigData = append(dynamicThresholdConfigData,dynamicThresholdConfigObj)
				case "VALUE_BASED":
					apiAggregateType := thresholdConfigData["apiAggregateType"].(string)
					userAggregateType := thresholdConfigData["userAggregateType"].(string)
					valueBasedThresholdConfig := thresholdConfigData["valueBasedThresholdConfig"].(map[string]interface{})
					uniqueValuesAllowed := valueBasedThresholdConfig["uniqueValuesAllowed"].(float64)
					sensitiveParamsEvaluationType := valueBasedThresholdConfig["sensitiveParamsEvaluationType"].(string)
					duration := valueBasedThresholdConfig["duration"].(string)
					valueBasedThresholdConfigObj := map[string]interface{}{
						"user_aggregate_type":userAggregateType,
						"api_aggregate_type":apiAggregateType,
						"unique_values_allowed":uniqueValuesAllowed,
						"duration":duration,
						"sensitive_params_evaluation_type":sensitiveParamsEvaluationType,
					}
					valueBasedThresholdConfigData = append(valueBasedThresholdConfigData,valueBasedThresholdConfigObj)
			}
		}
		finalThresholdConfigsObj := map[string]interface{}{
			"rolling_window_threshold_config":rollingWindowThresholdConfigData,
			"dynamic_threshold_config":dynamicThresholdConfigData,
			"value_based_threshold_config": valueBasedThresholdConfigData,
		}
		finalThresholdConfigs = append(finalThresholdConfigs, finalThresholdConfigsObj)
		rData.Set("threshold_configs",finalThresholdConfigs)
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
			rData.Set("ip_reputation", minIpReputationSeverity)
			ipReputationScopeFlag = false

		case "IP_ABUSE_VELOCITY":
			minIpAbuseVelocity := leafCondition["ipAbuseVelocityCondition"].(map[string]interface{})["minIpAbuseVelocity"].(string)
			rData.Set("ip_abuse_velocity", minIpAbuseVelocity)
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
			rData.Set("ip_location_type", finalConditionState)
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
			rData.Set("ip_address", finalConditionState)
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
			rData.Set("email_domain", finalConditionState)
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
			rData.Set("user_id", finalConditionState)
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
			rData.Set("user_agents", finalConditionState)
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
			rData.Set("ip_organisation", finalConditionState)
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
			rData.Set("ip_asn", finalConditionState)
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
			rData.Set("ip_connection_type", finalConditionState)
			ipConnTypeFlag = false

		case "DATATYPE":
			dataTypeConditionMap := leafCondition["datatypeCondition"].(map[string]interface{})
			dataSetsIds := dataTypeConditionMap["datasetIds"].([]interface{})
			datatypeIds := dataTypeConditionMap["datatypeIds"].([]interface{})
			location := "REQUEST_RESPONSE"
			if dataLocation,ok := dataTypeConditionMap["dataLocation"].(string); ok {
				location=dataLocation
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
			rData.Set("request_scanner_type", finalConditionState)
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
			rData.Set("regions", finalConditionState)
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
				rData.Set("label_id_scope", labelIds)
				labelIdScopeFlag = false
			} else if scopeType == "ENTITY" {
				entityScope := scopeCondition["entityScope"].(map[string]interface{})
				entityIds := entityScope["entityIds"].([]interface{})
				rData.Set("endpoint_id_scope", entityIds)
				endPointIdScopeFlag = false
			}
		}
	}

	if ipAddressFlag {
		rData.Set("ip_address", []interface{}{})
	}
	if labelIdScopeFlag {
		rData.Set("label_id_scope", []interface{}{})
	}
	if endPointIdScopeFlag {
		rData.Set("endpoint_id_scope", []interface{}{})
	}
	if ipReputationScopeFlag {
		rData.Set("ip_reputation", "")
	}
	if ipLocationTypeScopeFlag {
		rData.Set("ip_location_type", []interface{}{})
	}
	if ipAbuseVelFlag {
		rData.Set("ip_abuse_velocity", "")
	}
	if emailDomainFlag {
		rData.Set("email_domain", []interface{}{})
	}
	if userAgentFlag {
		rData.Set("user_agents", []interface{}{})
	}
	if regionFlag {
		rData.Set("regions", []interface{}{})
	}
	if ipOrgFlag {
		rData.Set("ip_organisation", []interface{}{})
	}
	if userIdFlag {
		rData.Set("user_id", []interface{}{})
	}
	if reqScannerFlag {
		rData.Set("request_scanner_type", []interface{}{})
	}
	if ipConnTypeFlag {
		rData.Set("ip_connection_type", []interface{}{})
	}
	if ipAsnFlag {
		rData.Set("ip_asn", []interface{}{})
	}

	var envList []interface{}
	if ruleConfigScope, ok := ruleDetails["ruleConfigScope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleConfigScope["environmentScope"].(map[string]interface{}); ok {
			envList = environmentScope["environmentIds"].([]interface{})
		}
	}
	rData.Set("environments", envList)
	rData.Set("request_response_single_valued_conditions", finalReqResMultiValueConditionState)
	rData.Set("request_response_multi_valued_conditions", finalReqResSingleValueConditionState)
	rData.Set("attribute_based_conditions", finalAttributeBasedConditionState)
	rData.Set("data_types_conditions", finalDataTypeConditionState)

	return nil
}

func ResourceDlpUserBasedUpdate(rData *schema.ResourceData, meta interface{}) error {
	name := rData.Get("name").(string)
	ruleType := rData.Get("rule_type").(string)
	description := rData.Get("description").(string)
	environments := rData.Get("environments").([]interface{})
	enabled := rData.Get("enabled").(bool)
	thresholdConfigs := rData.Get("threshold_configs").([]interface{})
	dataTypesConditions := rData.Get("data_types_conditions").([]interface{})
	expiryDuration := rData.Get("expiry_duration").(string)
	alertSeverity := rData.Get("alert_severity").(string)
	ipReputation := rData.Get("ip_reputation").(string)
	ipAbuseVelocity := rData.Get("ip_abuse_velocity").(string)
	labelIdScope := rData.Get("label_id_scope").([]interface{})
	endpointIdScope := rData.Get("endpoint_id_scope").([]interface{})
	requestResponseSingleValuedConditions := rData.Get("request_response_single_valued_conditions").([]interface{})
	requestResponseMultiValuedConditions := rData.Get("request_response_multi_valued_conditions").([]interface{})
	attributeBasedConditions := rData.Get("attribute_based_conditions").([]interface{})
	ipLocationType := rData.Get("ip_location_type").([]interface{})
	ipAddress := rData.Get("ip_address").([]interface{})
	emailDomain := rData.Get("email_domain").([]interface{})
	regions := rData.Get("regions").([]interface{})
	userAgents := rData.Get("user_agents").([]interface{})
	ipOrganisation := rData.Get("ip_organisation").([]interface{})
	ipAsn := rData.Get("ip_asn").([]interface{})
	ipConnectionType := rData.Get("ip_connection_type").([]interface{})
	requestScannerType := rData.Get("request_scanner_type").([]interface{})
	userId := rData.Get("user_id").([]interface{})
	id := rData.Id()

	finalThresholdConfigQuery, err := ReturnFinalThresholdConfigQueryDlp(thresholdConfigs)
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
	updateRateLimitQuery := fmt.Sprintf(rate_limiting.RATE_LIMITING_UPDATE_QUERY,DLP_KEY, id, finalConditionsQuery, enabled, name, ruleType, strings.ToLower(ruleType), actionsBlockQuery, finalThresholdConfigQuery, finalEnvironmentQuery, description)
	log.Printf("This is the graphql query %s", updateRateLimitQuery)
	responseStr, err := common.CallExecuteQuery(updateRateLimitQuery, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	updatedId,err := common.GetIdFromResponse(responseStr,"updateRateLimitingRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	rData.SetId(updatedId)

	return nil
}

func ResourceDlpUserBasedDelete(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	query := fmt.Sprintf(rate_limiting.DELETE_RATE_LIMIT_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	log.Println(query)
	rData.SetId("")
	return nil
}
