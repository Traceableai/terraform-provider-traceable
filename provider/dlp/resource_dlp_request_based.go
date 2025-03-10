package dlp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/rate_limiting"
	"log"
)

func ResourceDlpRequestBasedRule() *schema.Resource {
	return &schema.Resource{
		Create:        ResourceDlpRequestBasedCreate,
		Read:          ResourceDlpRequestBasedRead,
		Update:        ResourceDlpRequestBasedUpdate,
		Delete:        ResourceDlpRequestBasedDelete,
		CustomizeDiff: validateSchemaRequestBased,
		Schema: map[string]*schema.Schema{
			"rule_type": {
				Type:        schema.TypeString,
				Description: "ALERT or BLOCK or ALLOW",
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
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable/disable the rule",
				Required:    true,
			},
			"expiry_duration": {
				Type:        schema.TypeString,
				Description: "Block/Allow for a given period of time",
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
			"environments": {
				Type:        schema.TypeList,
				Description: "List of environments ids",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"request_payload_single_valued_conditions": {
				Type:        schema.TypeList,
				Description: "Request payload single valued conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request_location": {
							Type:        schema.TypeString,
							Required:    true,
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
			"request_payload_multi_valued_conditions": {
				Type:        schema.TypeList,
				Description: "Request payload multi valued conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"request_location": {
							Type:        schema.TypeString,
							Required:    true,
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
										Type:     schema.TypeString,
										Required: true,
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
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_types": {
							Type:        schema.TypeList,
							Description: "Data types to include",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_type_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"data_sets_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},

						"custom_location_data_type_key_value_pair_matching": {
							Type:        schema.TypeBool,
							Description: "Enable to look for data type key in a custom location",
							Optional:    true,
							Default:     false,
						},
						"custom_location_attribute": {
							Type:        schema.TypeString,
							Description: "Custom location attribute key",
							Optional:    true,
						},
						"custom_location_attribute_key_operator": {
							Type:        schema.TypeString,
							Description: "Match regex/Match exactly",
							Optional:    true,
						},
						"custom_location_attribute_value": {
							Type:        schema.TypeString,
							Description: "Custom location attribute key",
							Optional:    true,
						},
					},
				},
			},
			"ip_address": {
				Type:        schema.TypeList,
				Description: "Ip address as source (LIST_OF_IP's)",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"regions": {
				Type:        schema.TypeList,
				Description: "Regions as source, It will be a list region ids (AX,DZ)",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ip_location_type": {
				Type:        schema.TypeList,
				Description: "Ip location type as source ([BOT, TOR_EXIT_NODE, PUBLIC_PROXY])",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"target_scope": {
				Type:        schema.TypeList,
				Description: "Service scope and url regex",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_ids": {
							Type:        schema.TypeList,
							Description: "Service ids",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"url_regex": {
							Type:        schema.TypeList,
							Description: "specify url regex",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}
func validateSchemaRequestBased(ctx context.Context, rData *schema.ResourceDiff, meta interface{}) error {
	dataTypesConditions := rData.Get("data_types_conditions").([]interface{})
	ruleType := rData.Get("rule_type").(string)
	alertSeverity := rData.Get("alert_severity").(string)
	expiryDuration := rData.Get("expiry_duration").(string)
	if expiryDuration != "" && ruleType == "ALERT" {
		return fmt.Errorf("expiry_duration not required here")
	}
	if ruleType == "ALLOW" && alertSeverity != "" {
		return fmt.Errorf("alert_severity not required here")
	}
	if ruleType != "ALLOW" && alertSeverity == "" {
		return fmt.Errorf("required feild missing alert_severity")
	}

	for _, data := range dataTypesConditions {
		customLocationDataTypeKeyValuePairMatching := data.(map[string]interface{})["custom_location_data_type_key_value_pair_matching"].(bool)
		customLocationAttribute := data.(map[string]interface{})["custom_location_attribute"].(string)
		customLocationAttributeKeyOperator := data.(map[string]interface{})["custom_location_attribute_key_operator"].(string)
		customLocationAttributeKeyValue := data.(map[string]interface{})["custom_location_attribute_value"].(string)

		if customLocationDataTypeKeyValuePairMatching {
			if customLocationAttribute == "REQUEST_BODY" {
				if customLocationAttributeKeyOperator != "" || customLocationAttributeKeyValue != "" {
					return fmt.Errorf("attributes not required in this context")
				}
			} else if customLocationAttribute == "" || customLocationAttributeKeyOperator == "" || customLocationAttributeKeyValue == "" {
				return fmt.Errorf("required attributes are missing")
			}
		} else {
			if customLocationAttribute != "" || customLocationAttributeKeyOperator != "" || customLocationAttributeKeyValue != "" {
				return fmt.Errorf("attributes not expected here")
			}
		}
	}
	return nil
}

func ResourceDlpRequestBasedCreate(rData *schema.ResourceData, meta interface{}) error {
	name := rData.Get("name").(string)
	ruleType := rData.Get("rule_type").(string)
	description := rData.Get("description").(string)
	environments := rData.Get("environments").([]interface{})
	enabled := rData.Get("enabled").(bool)
	expiryDuration := rData.Get("expiry_duration").(string)
	alertSeverity := rData.Get("alert_severity").(string)
	dataTypesConditions := rData.Get("data_types_conditions").([]interface{})
	requestPayloadSingleValuedConditions := rData.Get("request_payload_single_valued_conditions").([]interface{})
	requestPayloadMultiValuedConditions := rData.Get("request_payload_multi_valued_conditions").([]interface{})
	ipLocationType := rData.Get("ip_location_type").([]interface{})
	ipAddress := rData.Get("ip_address").([]interface{})
	targetScope := rData.Get("target_scope").([]interface{})
	regions := rData.Get("regions").([]interface{})

	finalConditionsQuery := GetConditionsStringDlp(
		targetScope,
		ipAddress,
		regions,
		ipLocationType,
		dataTypesConditions,
		requestPayloadSingleValuedConditions,
		requestPayloadMultiValuedConditions,
	)
	if finalConditionsQuery == "" {
		return fmt.Errorf("required at least one scope condition")
	}
	finalEnvironmentQuery := ""
	if len(environments) > 0 {
		finalEnvironmentQuery = fmt.Sprintf(rate_limiting.ENVIRONMENT_SCOPE_QUERY, common.ReturnQuotedStringList(environments))
	}
	transactionActionConfigs := GetTransactionActionConfigs(ruleType, alertSeverity, expiryDuration)
	createDlpReqBasedQuery := fmt.Sprintf(DLP_REQUEST_BASED_QUERY_CREATE, finalConditionsQuery, enabled, name, description, transactionActionConfigs, finalEnvironmentQuery)
	log.Printf("This is the graphql query %s", createDlpReqBasedQuery)
	responseStr, err := common.CallExecuteQuery(createDlpReqBasedQuery, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createRateLimitingRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	rData.SetId(id)
	return nil
}

func ResourceDlpRequestBasedRead(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	var response map[string]interface{}
	readQuery := fmt.Sprintf(rate_limiting.FETCH_RATE_LIMIT_RULES, DLP_KEY)
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
	if transactionActionConfigs, ok := ruleDetails["transactionActionConfigs"].(map[string]interface{}); ok {
		action := transactionActionConfigs["action"].(map[string]interface{})
		actionType := action["actionType"].(string)
		rData.Set("rule_type", actionType)
		switch actionType {
		case "ALLOW":
			rData.Set("alert_severity", "")
			allow := action["allow"].(map[string]interface{})
			if duration, ok := allow["duration"].(string); ok {
				rData.Set("expiry_duration", duration)
			} else {
				rData.Set("expiry_duration", "")
			}
		case "ALERT":
			rData.Set("expiry_duration", "")
			alert := action["alert"].(map[string]interface{})
			severity := alert["eventSeverity"].(string)
			rData.Set("alert_severity", severity)
		case "BLOCK":
			block := action["block"].(map[string]interface{})
			severity := block["eventSeverity"].(string)
			rData.Set("alert_severity", severity)
			if duration, ok := block["duration"].(string); ok {
				rData.Set("expiry_duration", duration)
			} else {
				rData.Set("expiry_duration", "")
			}
		}
	}
	conditionsArray := ruleDetails["conditions"].([]interface{})
	finalReqPayloadSingleValueConditionState := []map[string]interface{}{}
	finalReqPayloadMultiValueConditionState := []map[string]interface{}{}
	urlRegexes := []interface{}{}
	serviceIds := []interface{}{}

	ipLocationTypeScopeFlag, ipAddressFlag, regionFlag := true, true, true
	for _, condition := range conditionsArray {
		leafCondition := condition.(map[string]interface{})["leafCondition"].(map[string]interface{})
		conditionType := leafCondition["conditionType"].(string)

		switch conditionType {

		case "IP_LOCATION_TYPE":
			ipLocationTypeCondition := leafCondition["ipLocationTypeCondition"].(map[string]interface{})
			ipLocationTypes := ipLocationTypeCondition["ipLocationTypes"].([]interface{})
			rData.Set("ip_location_type", ipLocationTypes)
			ipLocationTypeScopeFlag = false

		case "IP_ADDRESS":
			ipAddressCondition := leafCondition["ipAddressCondition"].(map[string]interface{})
			rawInputIpData := ipAddressCondition["rawInputIpData"].([]interface{})
			rData.Set("ip_address", rawInputIpData)
			ipAddressFlag = false

		case "DATATYPE":
			dataTypeConditionMap := leafCondition["datatypeCondition"].(map[string]interface{})
			dataSetsIds := dataTypeConditionMap["datasetIds"].([]interface{})
			datatypeIds := dataTypeConditionMap["datatypeIds"].([]interface{})
			dataTypesObj := map[string]interface{}{
				"data_type_ids": datatypeIds,
				"data_sets_ids": dataSetsIds,
			}
			dataTypesObjSlice := []map[string]interface{}{}
			dataTypesObjSlice = append(dataTypesObjSlice, dataTypesObj)
			finalDataTypeConditionsObj := []map[string]interface{}{}
			var dataTypeConditionsObj map[string]interface{}
			if datatypeMatching, ok := dataTypeConditionMap["datatypeMatching"].(map[string]interface{}); ok {
				rData.Set("custom_location_data_type_key_value_pair_matching", true)
				regexBasedMatching := datatypeMatching["regexBasedMatching"].(map[string]interface{})
				customMatchingLocation := regexBasedMatching["customMatchingLocation"].(map[string]interface{})
				metadataType := customMatchingLocation["metadataType"].(string)
				if metadataType == "REQUEST_BODY" {
					dataTypeConditionsObj = map[string]interface{}{
						"data_types": dataTypesObjSlice,
						"custom_location_data_type_key_value_pair_matching": true,
						"custom_location_attribute":                         metadataType,
						"custom_location_attribute_key_operator":            "",
						"custom_location_attribute_value":                   "",
					}
				} else {
					keyCondition := customMatchingLocation["keyCondition"].(map[string]interface{})
					operator := keyCondition["operator"].(string)
					value := keyCondition["value"].(string)
					dataTypeConditionsObj = map[string]interface{}{
						"data_types": dataTypesObjSlice,
						"custom_location_data_type_key_value_pair_matching": true,
						"custom_location_attribute":                         metadataType,
						"custom_location_attribute_key_operator":            operator,
						"custom_location_attribute_value":                   value,
					}
				}
			} else {
				dataTypeConditionsObj = map[string]interface{}{
					"data_types": dataTypesObjSlice,
					"custom_location_data_type_key_value_pair_matching": false,
					"custom_location_attribute":                         "",
					"custom_location_attribute_key_operator":            "",
					"custom_location_attribute_value":                   "",
				}
			}
			finalDataTypeConditionsObj = append(finalDataTypeConditionsObj, dataTypeConditionsObj)
			rData.Set("data_types_conditions", finalDataTypeConditionsObj)

		case "REGION":
			regionCondition := leafCondition["regionCondition"].(map[string]interface{})
			regionIdentifiers := regionCondition["regionIdentifiers"].([]interface{})
			rData.Set("regions", regionIdentifiers)
			regionFlag = false

		case "KEY_VALUE":
			keyValueCondition := leafCondition["keyValueCondition"].(map[string]interface{})
			metadataType := keyValueCondition["metadataType"].(string)
			valuePatternObjSlice := []map[string]interface{}{}
			keyPatternObjSlice := []map[string]interface{}{}
			if keyCondition, ok := keyValueCondition["keyCondition"].(map[string]interface{}); ok {
				keyPatternObj := map[string]interface{}{
					"operator": keyCondition["operator"].(string),
					"value":    keyCondition["value"].(string),
				}
				keyPatternObjSlice = append(keyPatternObjSlice, keyPatternObj)
				if valueCondition, ok := keyValueCondition["valueCondition"].(map[string]interface{}); ok {
					valuePatternObj := map[string]interface{}{
						"operator": valueCondition["operator"].(string),
						"value":    valueCondition["value"].(string),
					}
					valuePatternObjSlice = append(valuePatternObjSlice, valuePatternObj)
				}
				reqPayloadMultiValuedObj := map[string]interface{}{
					"request_location": metadataType,
					"key_patterns":     keyPatternObjSlice,
					"value_patterns":   valuePatternObjSlice,
				}
				finalReqPayloadMultiValueConditionState = append(finalReqPayloadMultiValueConditionState, reqPayloadMultiValuedObj)
			} else {
				valueCondition := keyValueCondition["valueCondition"].(map[string]interface{})
				operator := valueCondition["operator"].(string)
				value := valueCondition["value"].(string)
				reqPayloadSingleValuedObj := map[string]interface{}{
					"request_location": metadataType,
					"operator":         operator,
					"value":            value,
				}
				finalReqPayloadSingleValueConditionState = append(finalReqPayloadSingleValueConditionState, reqPayloadSingleValuedObj)
			}

		case "SCOPE":
			scopeCondition := leafCondition["scopeCondition"].(map[string]interface{})
			scopeType := scopeCondition["scopeType"].(string)
			if scopeType == "URL" {
				urlScope := scopeCondition["urlScope"].(map[string]interface{})
				urlRegexes = urlScope["urlRegexes"].([]interface{})

			} else if scopeType == "ENTITY" {
				entityScope := scopeCondition["entityScope"].(map[string]interface{})
				serviceIds = entityScope["entityIds"].([]interface{})
			}
		}
	}
	targetScopeObj := map[string]interface{}{
		"service_ids": serviceIds,
		"url_regex":   urlRegexes,
	}
	targetScopeSlice := []map[string]interface{}{}
	targetScopeSlice = append(targetScopeSlice, targetScopeObj)
	rData.Set("target_scope", targetScopeSlice)
	if ipAddressFlag {
		rData.Set("ip_address", []interface{}{})
	}
	if ipLocationTypeScopeFlag {
		rData.Set("ip_location_type", []interface{}{})
	}
	if regionFlag {
		rData.Set("regions", []interface{}{})
	}

	var envList []interface{}
	if ruleConfigScope, ok := ruleDetails["ruleConfigScope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleConfigScope["environmentScope"].(map[string]interface{}); ok {
			envList = environmentScope["environmentIds"].([]interface{})
		}
	}
	rData.Set("environments", envList)
	rData.Set("request_payload_multi_valued_conditions", finalReqPayloadMultiValueConditionState)
	rData.Set("request_payload_single_valued_conditions", finalReqPayloadSingleValueConditionState)

	return nil
}

func ResourceDlpRequestBasedUpdate(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	name := rData.Get("name").(string)
	ruleType := rData.Get("rule_type").(string)
	description := rData.Get("description").(string)
	environments := rData.Get("environments").([]interface{})
	enabled := rData.Get("enabled").(bool)
	expiryDuration := rData.Get("expiry_duration").(string)
	alertSeverity := rData.Get("alert_severity").(string)
	dataTypesConditions := rData.Get("data_types_conditions").([]interface{})
	requestPayloadSingleValuedConditions := rData.Get("request_payload_single_valued_conditions").([]interface{})
	requestPayloadMultiValuedConditions := rData.Get("request_payload_multi_valued_conditions").([]interface{})
	ipLocationType := rData.Get("ip_location_type").([]interface{})
	ipAddress := rData.Get("ip_address").([]interface{})
	targetScope := rData.Get("target_scope").([]interface{})
	regions := rData.Get("regions").([]interface{})

	finalConditionsQuery := GetConditionsStringDlp(
		targetScope,
		ipAddress,
		regions,
		ipLocationType,
		dataTypesConditions,
		requestPayloadSingleValuedConditions,
		requestPayloadMultiValuedConditions,
	)
	if finalConditionsQuery == "" {
		return fmt.Errorf("required at least one scope condition")
	}
	finalEnvironmentQuery := ""
	if len(environments) > 0 {
		finalEnvironmentQuery = fmt.Sprintf(rate_limiting.ENVIRONMENT_SCOPE_QUERY, common.ReturnQuotedStringList(environments))
	}
	transactionActionConfigs := GetTransactionActionConfigs(ruleType, alertSeverity, expiryDuration)
	createDlpReqBasedQuery := fmt.Sprintf(DLP_REQUEST_BASED_QUERY_UPDATE, id, finalConditionsQuery, enabled, name, description, transactionActionConfigs, finalEnvironmentQuery)
	log.Printf("This is the graphql query %s", createDlpReqBasedQuery)
	responseStr, err := common.CallExecuteQuery(createDlpReqBasedQuery, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateRateLimitingRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	rData.SetId(updatedId)
	return nil
}

func ResourceDlpRequestBasedDelete(rData *schema.ResourceData, meta interface{}) error {
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
