package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSessionIdentificationRequestRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSessionIdentificationRequestRuleCreate,
		Read:   resourceSessionIdentificationRequestRuleRead,
		Update: resourceSessionIdentificationRequestRuleUpdate,
		Delete: resourceSessionIdentificationRequestRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Session Identification Rule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the Session Identification Rule",
				Optional:    true,
			},
			"environment_names": {
				Type:        schema.TypeList,
				Description: "List of environment names",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_names": {
				Type:        schema.TypeList,
				Description: "List of service names",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"url_match_regexes": {
				Type:        schema.TypeList,
				Description: "List of URL match regexes",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"token_extraction_condition_list": {
				Type:        schema.TypeList,
				Description: "Conditions to satisfy for extracting Session Token",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition_request_header": {
							Type:        schema.TypeSet,
							Description: "Attribute type request header",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Description: "request header key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "request header value",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressConditionListDiff,
						},
						"condition_request_cookie": {
							Type:        schema.TypeSet,
							Description: "Attribute type request cookie",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Description: "request cookie key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "request cookie value",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressConditionListDiff,
						},
						"condition_request_query_param": {
							Type:        schema.TypeSet,
							Description: "Attribute type request query param",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Description: "request query param key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "request query param value",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressConditionListDiff,
						},
					},
				},
			},
			"session_token_details": {
				Type:        schema.TypeList,
				Description: "Details of the session token of type request",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token_request_header": {
							Type:        schema.TypeSet,
							Description: "request header for token",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token_key": {
										Type:        schema.TypeString,
										Description: "Test header key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressSessionTokenDetailsDiff,
						},
						"token_request_cookie": {
							Type:        schema.TypeSet,
							Description: "request cookie for token",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token_key": {
										Type:        schema.TypeString,
										Description: "Test cookie key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressSessionTokenDetailsDiff,
						},
						"token_request_query_param": {
							Type:        schema.TypeSet,
							Description: "request query param for token",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"token_key": {
										Type:        schema.TypeString,
										Description: "Test query param key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressSessionTokenDetailsDiff,
						},
						"token_request_body": {
							Type:             schema.TypeBool,
							Description:      "request body for token",
							Optional:         true,
							DiffSuppressFunc: suppressSessionTokenDetailsDiff,
						},
					},
				},
			},
			"obfuscation": {
				Type:        schema.TypeBool,
				Description: "If the obfuscation strategy of HASH to be used",
				Required:    true,
			},
			"expiration_type": {
				Type:        schema.TypeString,
				Description: "expiration is jwt based or not applicable",
				Optional:    true,
			},
			"token_value_transformation_list": {
				Type:        schema.TypeList,
				Description: "Conditions to satisfy for extracting Session Token",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"json_path": {
							Type:             schema.TypeString,
							Description:      "the json path group for value transformation",
							Optional:         true,
							DiffSuppressFunc: suppressValueTransformationListDiff,
						},
						"regex_capture_group": {
							Type:             schema.TypeString,
							Description:      "the regex capture group for value transformation",
							Optional:         true,
							DiffSuppressFunc: suppressValueTransformationListDiff,
						},
						"jwt_payload_claim": {
							Type:             schema.TypeString,
							Description:      "the jwt payload claim for value transformation",
							Optional:         true,
							DiffSuppressFunc: suppressValueTransformationListDiff,
						},
						"base64": {
							Type:             schema.TypeBool,
							Description:      "whether we use the base64 value transformation",
							Optional:         true,
							DiffSuppressFunc: suppressValueTransformationListDiff,
						},
					},
				},
			},
		},
	}
}

func suppressConditionListDiff(k, old, new string, d *schema.ResourceData) bool {
	return suppressListDiff(old, new)
}

func suppressSessionTokenDetailsDiff(k, old, new string, d *schema.ResourceData) bool {
	return suppressListDiff(old, new)
}

func suppressValueTransformationListDiff(k, old, new string, d *schema.ResourceData) bool {
	return suppressListDiff(old, new)
}

func suppressListDiff(old, new string) bool {
	oldList := strings.FieldsFunc(old, func(r rune) bool { return r == ',' || r == '[' || r == ']' || r == '{' || r == '}' || r == '"' })
	newList := strings.FieldsFunc(new, func(r rune) bool { return r == ',' || r == '[' || r == ']' || r == '{' || r == '}' || r == '"' })
	if len(oldList) != len(newList) {
		return false
	}
	oldMap := make(map[string]bool)
	for _, v := range oldList {
		oldMap[v] = true
	}
	for _, v := range newList {
		if !oldMap[v] {
			return false
		}
	}
	return true
}

func resourceSessionIdentificationRequestRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	descriptionStr := ""
	if v, ok := d.GetOk("description"); ok {
		descriptionStr = fmt.Sprintf(`description: "%s",`, v.(string))
	}

	scopeStr := ""
	envNames := d.Get("environment_names").([]interface{})
	serviceNames := d.Get("service_names").([]interface{})
	urlMatchRegexes := d.Get("url_match_regexes").([]interface{})

	envNamesStr := buildStringArray(interfaceSliceToStringSlice(envNames))
	serviceNamesStr := buildStringArray(interfaceSliceToStringSlice(serviceNames))
	urlMatchRegexesStr := buildStringArray(interfaceSliceToStringSlice(urlMatchRegexes))

	if len(envNames) > 0 {
		scopeStr = fmt.Sprintf(`scope: {
			environmentNames: %s`, envNamesStr)
		if len(serviceNames) > 0 {
			scopeStr += fmt.Sprintf(`, serviceNames: %s`, serviceNamesStr)
		}
		if len(urlMatchRegexes) > 0 {
			scopeStr += fmt.Sprintf(`, urlMatchRegexes: %s`, urlMatchRegexesStr)
		}
		scopeStr += "},"
	} else if len(serviceNames) > 0 && (len(serviceNames) != 1 || serviceNames[0] != "") {
		scopeStr = fmt.Sprintf(`scope: {
			serviceNames: %s
		},`, serviceNamesStr)
	}

	conditionListStr := ""
	if v, ok := d.GetOk("token_extraction_condition_list"); ok {
		conditions := v.([]interface{})
		if len(conditions) > 0 {
			conditionListStr = `predicate: { predicateType: LOGICAL, logicalPredicate: { operator: AND, children: [`

			for _, condition := range conditions {
				conditionMap := condition.(map[string]interface{})

				if v, ok := conditionMap["condition_request_header"]; ok {
					conditionListStr += buildConditionList("HEADER", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_cookie"]; ok {
					conditionListStr += buildConditionList("COOKIE", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_query_param"]; ok {
					conditionListStr += buildConditionList("QUERY_PARAMETER", v.(*schema.Set).List())
				}
			}

			conditionListStr += `]}},`
		}
	}

	// the trailing comma is removed if no conditions are present
	if conditionListStr == `predicate: { predicateType: LOGICAL, logicalPredicate: { operator: AND, children: []}}, ` {
		conditionListStr = ""
	}

	expirationTypeStr := ""
	if v, ok := d.GetOk("expiration_type"); ok {
		expirationTypeStr = fmt.Sprintf(`expirationType: %s,`, v.(string))
	}

	obfuscationStr := ""
	if v, ok := d.GetOk("obfuscation"); ok {
		if v.(bool) {
			obfuscationStr = `obfuscationStrategy: HASH,`
		}
	}

	requestAttributeKeyLocationStr := ""
	if v, ok := d.GetOk("session_token_details"); ok {
		sessionTokenDetails := v.([]interface{})[0].(map[string]interface{})

		if tokenRequestHeader, ok := sessionTokenDetails["token_request_header"]; ok {
			if len(tokenRequestHeader.(*schema.Set).List()) > 0 {
				requestAttributeKeyLocationStr = "HEADER"
			}
		} else if tokenRequestCookie, ok := sessionTokenDetails["token_request_cookie"]; ok {
			if len(tokenRequestCookie.(*schema.Set).List()) > 0 {
				requestAttributeKeyLocationStr = "COOKIE"
			}
		} else if tokenRequestQueryParam, ok := sessionTokenDetails["token_request_query_param"]; ok {
			if len(tokenRequestQueryParam.(*schema.Set).List()) > 0 {
				requestAttributeKeyLocationStr = "QUERY_PARAMETER"
			}
		} else if tokenRequestBody, ok := sessionTokenDetails["token_request_body"].(bool); ok && tokenRequestBody {
			requestAttributeKeyLocationStr = "BODY"
		}
	}

	tokenMatchConditionStr := ""
	if requestAttributeKeyLocationStr != "BODY" && requestAttributeKeyLocationStr != "" {
		if v, ok := d.GetOk("session_token_details"); ok {
			sessionTokenDetails := v.([]interface{})[0].(map[string]interface{})
			if tokenRequestHeader, ok := sessionTokenDetails["token_request_header"].([]interface{}); ok && len(tokenRequestHeader) > 0 {
				tokenKey := tokenRequestHeader[0].(map[string]interface{})["token_key"].(string)
				operator := tokenRequestHeader[0].(map[string]interface{})["operator"].(string)
				tokenMatchConditionStr = fmt.Sprintf(`matchCondition: { matchOperator: %s, stringValue: "%s" },`, operator, tokenKey)
			} else if tokenRequestCookie, ok := sessionTokenDetails["token_request_cookie"].([]interface{}); ok && len(tokenRequestCookie) > 0 {
				tokenKey := tokenRequestCookie[0].(map[string]interface{})["token_key"].(string)
				operator := tokenRequestCookie[0].(map[string]interface{})["operator"].(string)
				tokenMatchConditionStr = fmt.Sprintf(`matchCondition: { matchOperator: %s, stringValue: "%s" },`, operator, tokenKey)
			} else if tokenRequestQueryParam, ok := sessionTokenDetails["token_request_query_param"].([]interface{}); ok && len(tokenRequestQueryParam) > 0 {
				tokenKey := tokenRequestQueryParam[0].(map[string]interface{})["token_key"].(string)
				operator := tokenRequestQueryParam[0].(map[string]interface{})["operator"].(string)
				tokenMatchConditionStr = fmt.Sprintf(`matchCondition: { matchOperator: %s, stringValue: "%s" },`, operator, tokenKey)
			}
		}
	}

	valueProjectionsStr := ""
	if v, ok := d.GetOk("token_value_transformation_list"); ok {
		valueTransformations := v.([]interface{})
		var valueProjections []string
		for _, transformation := range valueTransformations {
			transformationMap := transformation.(map[string]interface{})
			if jsonPath, ok := transformationMap["json_path"].(string); ok && jsonPath != "" {
				valueProjections = append(valueProjections, fmt.Sprintf(`{
				valueProjectionType: JSON_PATH,
				jsonPathProjection: { path: "%s" }
			}`, jsonPath))
			}
			if regexCaptureGroup, ok := transformationMap["regex_capture_group"].(string); ok && regexCaptureGroup != "" {
				valueProjections = append(valueProjections, fmt.Sprintf(`{
				valueProjectionType: REGEX_CAPTURE_GROUP,
				regexCaptureGroupProjection: { regexCaptureGroup: "%s" }
			}`, regexCaptureGroup))
			}
			if jwtPayloadClaim, ok := transformationMap["jwt_payload_claim"].(string); ok && jwtPayloadClaim != "" {
				valueProjections = append(valueProjections, fmt.Sprintf(`{
				valueProjectionType: JWT_PAYLOAD_CLAIM,
				jwtPayloadClaimProjection: { claim: "%s" }
			}`, jwtPayloadClaim))
			}
			if base64, ok := transformationMap["base64"].(bool); ok && base64 {
				valueProjections = append(valueProjections, `{ valueProjectionType: BASE64 }`)
			}
		}
		if len(valueProjections) > 0 {
			valueProjectionsStr = fmt.Sprintf("valueProjections: [%s]", strings.Join(valueProjections, ", "))
		}
	}

	query := fmt.Sprintf(`mutation {
		createSessionIdentificationRuleV2(
		  create: {
			name: "%s"
			%s
			%s
			sessionTokenRules: [
			  {
				%s
				tokenType: REQUEST
				requestSessionTokenDetails: {
				  requestAttributeKeyLocation: %s
				  %s
				}
				sessionTokenValueRule: {
				  %s
				  projectionRoot: {
					projectionType: ATTRIBUTE
					attributeProjection: {
					  %s
					  %s
					}
				  }
				}
			  }
			]
		  }
		) {
		  id
		}
	  }
	`, name, descriptionStr, scopeStr, conditionListStr, requestAttributeKeyLocationStr, expirationTypeStr, obfuscationStr, tokenMatchConditionStr, valueProjectionsStr)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)

	}

	log.Printf("GraphQL response: %s", responseStr)

	if response["data"] != nil && response["data"].(map[string]interface{})["createSessionIdentificationRuleV2"] != nil {
		id := response["data"].(map[string]interface{})["createSessionIdentificationRuleV2"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf("could not create Session Identification request rule, no ID returned")
	}

	return nil

}

func buildStringArray(input []string) string {
	if len(input) == 0 {
		return "[]"
	}
	output := "["
	for _, v := range input {
		output += fmt.Sprintf(`"%s",`, v)
	}
	output = output[:len(output)-1] // Remove trailing comma
	output += "]"
	return output
}

func interfaceSliceToStringSlice(input []interface{}) []string {
	var output []string
	for _, v := range input {
		output = append(output, v.(string))
	}
	return output
}

func buildConditionList(attributeType string, conditions []interface{}) string {
	conditionStr := ""
	for _, condition := range conditions {
		conditionMap := condition.(map[string]interface{})
		key := conditionMap["key"].(string)
		operator := conditionMap["operator"].(string)
		value := conditionMap["value"].(string)
		conditionStr += fmt.Sprintf(`{
			predicateType: ATTRIBUTE,
			attributePredicate: {
				attributeProjection: {
					matchCondition: {
						matchOperator: %s,
						stringValue: "%s"
					}
				},
				matchCondition: {
					matchOperator: %s,
					stringValue: "%s"
				},
				attributeKeyLocationType: REQUEST,
				requestAttributeKeyLocation: %s
			}
		},`, operator, key, operator, value, attributeType)
	}
	// Remove trailing comma
	if len(conditionStr) > 0 {
		conditionStr = conditionStr[:len(conditionStr)-1]
	}
	return conditionStr
}

func resourceSessionIdentificationRequestRuleRead(d *schema.ResourceData, meta interface{}) error {
	readQuery := `{sessionIdentificationRulesV2{count results{id scope{environmentNames serviceNames urlMatchRegexes}description name sessionTokenRules{predicate{attributePredicate{attributeKeyLocationType attributeProjection{matchCondition{matchOperator stringValue}valueProjections{jsonPathProjection{path}jwtPayloadClaimProjection{claim}regexCaptureGroupProjection{regexCaptureGroup}valueProjectionType}}matchCondition{matchOperator stringValue}requestAttributeKeyLocation responseAttributeKeyLocation}customProjection{customJson}logicalPredicate{children{attributePredicate{attributeKeyLocationType attributeProjection{matchCondition{matchOperator stringValue}valueProjections{jsonPathProjection{path}jwtPayloadClaimProjection{claim}regexCaptureGroupProjection{regexCaptureGroup}valueProjectionType}}matchCondition{matchOperator stringValue}requestAttributeKeyLocation responseAttributeKeyLocation}predicateType}operator}predicateType}requestSessionTokenDetails{requestAttributeKeyLocation expirationType}responseSessionTokenDetails{attributeExpiration{expirationFormat projectionRoot{attributeProjection{matchCondition{matchOperator stringValue}valueProjections{jsonPathProjection{path}jwtPayloadClaimProjection{claim}regexCaptureGroupProjection{regexCaptureGroup}valueProjectionType}}customProjection{customJson}projectionType}responseAttributeKeyLocation}expirationType responseAttributeKeyLocation}sessionTokenValueRule{obfuscationStrategy projectionRoot{attributeProjection{matchCondition{matchOperator stringValue}valueProjections{jsonPathProjection{path}jwtPayloadClaimProjection{claim}regexCaptureGroupProjection{regexCaptureGroup}valueProjectionType}}customProjection{customJson}projectionType}}tokenType}status{disabled}}total}}`

	var response map[string]interface{}
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	log.Printf("This is the GraphQL query: %s", readQuery)
	log.Printf("This is the GraphQL response: %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	id := d.Id()
	ruleDetails := getRuleDetailsFromRulesListUsingIdName(response, "sessionIdentificationRulesV2", id, "id", "name")
	if len(ruleDetails) == 0 {
		d.SetId("")
		return nil
	}
	log.Printf("Session Identification Rule: %s", ruleDetails)

	d.Set("name", ruleDetails["name"])
	d.Set("description", ruleDetails["description"])

	scope := ruleDetails["scope"].(map[string]interface{})
	if environmentNames, ok := scope["environmentNames"]; ok {
		d.Set("environment_names", environmentNames)
	}
	if serviceNames, ok := scope["serviceNames"]; ok {
		d.Set("service_names", serviceNames)
	}
	if urlMatchRegexes, ok := scope["urlMatchRegexes"]; ok {
		d.Set("url_match_regexes", urlMatchRegexes)
	}

	sessionTokenRules := ruleDetails["sessionTokenRules"].([]interface{})
	if len(sessionTokenRules) > 0 {
		sessionTokenRule := sessionTokenRules[0].(map[string]interface{})
		if predicate, ok := sessionTokenRule["predicate"].(map[string]interface{}); ok {
			if logicalPredicate, ok := predicate["logicalPredicate"].(map[string]interface{}); ok {
				if children, ok := logicalPredicate["children"].([]interface{}); ok {
					for _, child := range children {
						childMap := child.(map[string]interface{})
						if attributePredicate, ok := childMap["attributePredicate"].(map[string]interface{}); ok {
							if attributeProjection, ok := attributePredicate["attributeProjection"].(map[string]interface{}); ok {
								if matchCondition, ok := attributeProjection["matchCondition"].(map[string]interface{}); ok {
									d.Set("condition_request_header", map[string]interface{}{
										"key":      matchCondition["stringValue"].(string),
										"operator": matchCondition["matchOperator"].(string),
									})
								}
							}
						}
					}
				}
			}
		}

		if requestSessionTokenDetails, ok := sessionTokenRule["requestSessionTokenDetails"].(map[string]interface{}); ok {
			requestAttributeKeyLocation := requestSessionTokenDetails["requestAttributeKeyLocation"].(string)
			d.Set("session_token_details", map[string]interface{}{
				"token_request_header": map[string]interface{}{
					"requestAttributeKeyLocation": requestAttributeKeyLocation,
				},
			})
		}

		if sessionTokenValueRule, ok := sessionTokenRule["sessionTokenValueRule"].(map[string]interface{}); ok {
			obfuscationStrategy := sessionTokenValueRule["obfuscationStrategy"].(string)
			if obfuscationStrategy == "HASH" {
				d.Set("obfuscation", true)
			} else {
				d.Set("obfuscation", false)
			}

			if projectionRoot, ok := sessionTokenValueRule["projectionRoot"].(map[string]interface{}); ok {
				if attributeProjection, ok := projectionRoot["attributeProjection"].(map[string]interface{}); ok {
					if valueProjections, ok := attributeProjection["valueProjections"].([]interface{}); ok {
						var transformations []interface{}
						for _, valueProjection := range valueProjections {
							valueProjectionMap := valueProjection.(map[string]interface{})
							transformation := make(map[string]interface{})
							switch valueProjectionMap["valueProjectionType"].(string) {
							case "JSON_PATH":
								transformation["json_path"] = valueProjectionMap["jsonPathProjection"].(map[string]interface{})["path"].(string)
							case "REGEX_CAPTURE_GROUP":
								transformation["regex_capture_group"] = valueProjectionMap["regexCaptureGroupProjection"].(map[string]interface{})["regexCaptureGroup"].(string)
							case "JWT_PAYLOAD_CLAIM":
								transformation["jwt_payload_claim"] = valueProjectionMap["jwtPayloadClaimProjection"].(map[string]interface{})["claim"].(string)
							case "BASE64":
								transformation["base64"] = true
							}
							transformations = append(transformations, transformation)
						}
						d.Set("token_value_transformation_list", transformations)
					}
				}
			}
		}
	}

	return nil
}

func resourceSessionIdentificationRequestRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)

	descriptionStr := ""
	if v, ok := d.GetOk("description"); ok {
		descriptionStr = fmt.Sprintf(`description: "%s",`, v.(string))
	}

	scopeStr := ""
	envNames := d.Get("environment_names").([]interface{})
	serviceNames := d.Get("service_names").([]interface{})
	urlMatchRegexes := d.Get("url_match_regexes").([]interface{})

	envNamesStr := buildStringArray(interfaceSliceToStringSlice(envNames))
	serviceNamesStr := buildStringArray(interfaceSliceToStringSlice(serviceNames))
	urlMatchRegexesStr := buildStringArray(interfaceSliceToStringSlice(urlMatchRegexes))

	if len(envNames) > 0 && (len(envNames) != 1 || envNames[0] != "") {
		scopeStr = fmt.Sprintf(`scope: {
			environmentNames: %s`, envNamesStr)
		if len(serviceNames) > 0 {
			scopeStr += fmt.Sprintf(`, serviceNames: %s`, serviceNamesStr)
		}
		if len(urlMatchRegexes) > 0 {
			scopeStr += fmt.Sprintf(`, urlMatchRegexes: %s`, urlMatchRegexesStr)
		}
		scopeStr += "},"
	} else if len(serviceNames) > 0 && (len(serviceNames) != 1 || serviceNames[0] != "") {
		scopeStr = fmt.Sprintf(`scope: {
			serviceNames: %s
		},`, serviceNamesStr)
	}

	conditionListStr := ""
	if v, ok := d.GetOk("token_extraction_condition_list"); ok {
		conditions := v.([]interface{})
		if len(conditions) > 0 {
			conditionListStr = `predicate: { predicateType: LOGICAL, logicalPredicate: { operator: AND, children: [`

			for _, condition := range conditions {
				conditionMap := condition.(map[string]interface{})

				if v, ok := conditionMap["condition_request_header"]; ok {
					conditionListStr += buildConditionList("HEADER", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_cookie"]; ok {
					conditionListStr += buildConditionList("COOKIE", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_query_param"]; ok {
					conditionListStr += buildConditionList("QUERY_PARAMETER", v.(*schema.Set).List())
				}
			}

			conditionListStr += `]}},`
		}
	}

	if conditionListStr == `predicate: { predicateType: LOGICAL, logicalPredicate: { operator: AND, children: []}},` {
		conditionListStr = ""
	}

	expirationTypeStr := ""
	if v, ok := d.GetOk("expiration_type"); ok {
		expirationTypeStr = fmt.Sprintf(`expirationType: %s,`, v.(string))
	}

	obfuscationStr := ""
	if v, ok := d.GetOk("obfuscation"); ok {
		if v.(bool) {
			obfuscationStr = `obfuscationStrategy: HASH,`
		}
	}

	requestAttributeKeyLocationStr := ""
	if v, ok := d.GetOk("session_token_details"); ok {
		sessionTokenDetails := v.([]interface{})[0].(map[string]interface{})

		if tokenRequestHeader, ok := sessionTokenDetails["token_request_header"]; ok {
			if len(tokenRequestHeader.(*schema.Set).List()) > 0 {
				requestAttributeKeyLocationStr = "HEADER"
			}
		} else if tokenRequestCookie, ok := sessionTokenDetails["token_request_cookie"]; ok {
			if len(tokenRequestCookie.(*schema.Set).List()) > 0 {
				requestAttributeKeyLocationStr = "COOKIE"
			}
		} else if tokenRequestQueryParam, ok := sessionTokenDetails["token_request_query_param"]; ok {
			if len(tokenRequestQueryParam.(*schema.Set).List()) > 0 {
				requestAttributeKeyLocationStr = "QUERY_PARAMETER"
			}
		} else if tokenRequestBody, ok := sessionTokenDetails["token_request_body"].(bool); ok && tokenRequestBody {
			requestAttributeKeyLocationStr = "BODY"
		}
	}

	tokenMatchConditionStr := ""
	if requestAttributeKeyLocationStr != "BODY" && requestAttributeKeyLocationStr != "" {
		if v, ok := d.GetOk("session_token_details"); ok {
			sessionTokenDetails := v.([]interface{})[0].(map[string]interface{})

			if tokenRequestHeader, ok := sessionTokenDetails["token_request_header"].([]interface{}); ok && len(tokenRequestHeader) > 0 {
				tokenKey := tokenRequestHeader[0].(map[string]interface{})["token_key"].(string)
				operator := tokenRequestHeader[0].(map[string]interface{})["operator"].(string)
				tokenMatchConditionStr = fmt.Sprintf(`matchCondition: { matchOperator: %s, stringValue: "%s" },`, operator, tokenKey)
			} else if tokenRequestCookie, ok := sessionTokenDetails["token_request_cookie"].([]interface{}); ok && len(tokenRequestCookie) > 0 {
				tokenKey := tokenRequestCookie[0].(map[string]interface{})["token_key"].(string)
				operator := tokenRequestCookie[0].(map[string]interface{})["operator"].(string)
				tokenMatchConditionStr = fmt.Sprintf(`matchCondition: { matchOperator: %s, stringValue: "%s" },`, operator, tokenKey)
			} else if tokenRequestQueryParam, ok := sessionTokenDetails["token_request_query_param"].([]interface{}); ok && len(tokenRequestQueryParam) > 0 {
				tokenKey := tokenRequestQueryParam[0].(map[string]interface{})["token_key"].(string)
				operator := tokenRequestQueryParam[0].(map[string]interface{})["operator"].(string)
				tokenMatchConditionStr = fmt.Sprintf(`matchCondition: { matchOperator: %s, stringValue: "%s" },`, operator, tokenKey)
			}
		}
	}

	valueProjectionsStr := ""
	if v, ok := d.GetOk("token_value_transformation_list"); ok {
		valueTransformations := v.([]interface{})
		var valueProjections []string

		for _, transformation := range valueTransformations {
			transformationMap := transformation.(map[string]interface{})

			if jsonPath, ok := transformationMap["json_path"].(string); ok && jsonPath != "" {
				valueProjections = append(valueProjections, fmt.Sprintf(`{
				valueProjectionType: JSON_PATH,
				jsonPathProjection: { path: "%s" }
			}`, jsonPath))
			}

			if regexCaptureGroup, ok := transformationMap["regex_capture_group"].(string); ok && regexCaptureGroup != "" {
				valueProjections = append(valueProjections, fmt.Sprintf(`{
				valueProjectionType: REGEX_CAPTURE_GROUP,
				regexCaptureGroupProjection: { regexCaptureGroup: "%s" }
			}`, regexCaptureGroup))
			}

			if jwtPayloadClaim, ok := transformationMap["jwt_payload_claim"].(string); ok && jwtPayloadClaim != "" {
				valueProjections = append(valueProjections, fmt.Sprintf(`{
				valueProjectionType: JWT_PAYLOAD_CLAIM,
				jwtPayloadClaimProjection: { claim: "%s" }
			}`, jwtPayloadClaim))
			}

			if base64, ok := transformationMap["base64"].(bool); ok && base64 {
				valueProjections = append(valueProjections, `{ valueProjectionType: BASE64 }`)
			}
		}

		if len(valueProjections) > 0 {
			valueProjectionsStr = fmt.Sprintf("valueProjections: [%s]", strings.Join(valueProjections, ", "))
		}
	}

	query := fmt.Sprintf(`mutation {
		updateSessionIdentificationRuleV2(
		  update: {
			id: "%s"
			name: "%s"
			%s
			%s
			sessionTokenRules: [
			  {
				%s
				tokenType: REQUEST
				requestSessionTokenDetails: {
				  requestAttributeKeyLocation: %s
				  %s
				}
				sessionTokenValueRule: {
				  %s
				  projectionRoot: {
					projectionType: ATTRIBUTE
					attributeProjection: {
					  %s
					  %s
					}
				  }
				}
			  }
			]
		  }
		) {
		  id
		}
	  }
	`, id, name, descriptionStr, scopeStr, conditionListStr, requestAttributeKeyLocationStr, expirationTypeStr, obfuscationStr, tokenMatchConditionStr, valueProjectionsStr)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	log.Printf("GraphQL response: %s", responseStr)

	if response["data"] != nil && response["data"].(map[string]interface{})["updateSessionIdentificationRuleV2"] != nil {
		d.SetId(id)
	} else {
		return fmt.Errorf("could not update Session Identification request rule, no ID returned")
	}

	return nil
}

func resourceSessionIdentificationRequestRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	query := fmt.Sprintf(`mutation {
		deleteSessionIdentificationRuleV2(
		  delete: { id: "%s" }
		) {
		  success
		}
	  }
	`, id)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	log.Printf("GraphQL response: %s", responseStr)

	success, ok := response["data"].(map[string]interface{})["deleteSessionIdentificationRuleV2"].(map[string]interface{})["success"].(bool)
	if !ok || !success {
		return fmt.Errorf("failed to delete Session Identification request rule")
	}

	d.SetId("")
	return nil
}
