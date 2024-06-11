package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSessionIdentificationCustomRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSessionIdentificationCustomRuleCreate,
		Read:   resourceSessionIdentificationCustomRuleRead,
		Update: resourceSessionIdentificationCustomRuleUpdate,
		Delete: resourceSessionIdentificationCustomRuleDelete,

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
						"condition_response_header": {
							Type:        schema.TypeSet,
							Description: "Attribute type response header",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Description: "response header key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "response header value",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressConditionListDiff,
						},
						"condition_response_cookie": {
							Type:        schema.TypeSet,
							Description: "Attribute type response cookie",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Description: "response cookie key",
										Required:    true,
									},
									"operator": {
										Type:        schema.TypeString,
										Description: "match operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "response cookie value",
										Required:    true,
									},
								},
							},
							DiffSuppressFunc: suppressConditionListDiff,
						},
					},
				},
			},
			"custom_json": {
				Type:        schema.TypeString,
				Description: "Custom json session token",
				Required:    true,
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

func resourceSessionIdentificationCustomRuleCreate(d *schema.ResourceData, meta interface{}) error {
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

				if v, ok := conditionMap["condition_response_header"]; ok {
					conditionListStr += buildConditionListResponse("HEADER", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_response_cookie"]; ok {
					conditionListStr += buildConditionListResponse("COOKIE", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_header"]; ok {
					conditionListStr += buildConditionListRequest("HEADER", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_cookie"]; ok {
					conditionListStr += buildConditionListRequest("COOKIE", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_query_param"]; ok {
					conditionListStr += buildConditionListRequest("QUERY_PARAMETER", v.(*schema.Set).List())
				}
			}

			conditionListStr += `]}},`
		}
	}

	if conditionListStr == `predicate: { predicateType: LOGICAL, logicalPredicate: { operator: AND, children: []}}, ` {
		conditionListStr = ""
	}

	customJsonStr := ""
	if v, ok := d.GetOk("custom_json"); ok {
		customJsonStr = fmt.Sprintf(`customJson: "%s",`, v.(string))
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
				tokenType: CUSTOM
				sessionTokenValueRule: {
				  projectionRoot: {
					projectionType: CUSTOM
					customProjection: {
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
	`, name, descriptionStr, scopeStr, conditionListStr, customJsonStr)

	var response map[string]interface{}
	log.Printf(query)
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	log.Printf("GraphQL response: %s", responseStr)

	if response["data"] != nil && response["data"].(map[string]interface{})["createSessionIdentificationRuleV2"] != nil {
		id := response["data"].(map[string]interface{})["createSessionIdentificationRuleV2"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf("could not create Session Identification response rule, no ID returned")
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

func buildConditionListRequest(attributeType string, conditions []interface{}) string {
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
	if len(conditionStr) > 0 {
		conditionStr = conditionStr[:len(conditionStr)-1]
	}
	return conditionStr
}

func buildConditionListResponse(attributeType string, conditions []interface{}) string {
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
				attributeKeyLocationType: RESPONSE,
				responseAttributeKeyLocation: %s
			}
		},`, operator, key, operator, value, attributeType)
	}
	if len(conditionStr) > 0 {
		conditionStr = conditionStr[:len(conditionStr)-1]
	}
	return conditionStr
}

func resourceSessionIdentificationCustomRuleRead(d *schema.ResourceData, meta interface{}) error {
	readQuery := `{
		sessionIdentificationRulesV2 {
			count
			results {
				id
				scope {
					environmentNames
					serviceNames
					urlMatchRegexes
				}
				description
				name
				sessionTokenRules {
					predicate {
						customProjection {
							customJson
						}
						logicalPredicate {
							children {
								attributePredicate {
									attributeKeyLocationType
									attributeProjection {
										matchCondition {
											matchOperator
											stringValue
										}
									}
									matchCondition {
										matchOperator
										stringValue
									}
									requestAttributeKeyLocation
									responseAttributeKeyLocation
								}
								predicateType
							}
							operator
						}
						predicateType
					}
					tokenType
				}
				status {
					disabled
				}
			}
			total
		}
	}`

	var response map[string]interface{}
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}
	log.Printf("This is the GraphQL query: %s", readQuery)
	log.Printf("This is the GraphQL response: %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	id := d.Id()
	ruleDetails := getRuleDetailsFromRulesListUsingIdName(response, "sessionIdentificationRulesV2", id, "id", "name")
	log.Printf("Session Identification Rule: %s", ruleDetails)

	if ruleDetails == nil {
		return fmt.Errorf("Could not find rule with id %s", id)
	}

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
		if predicate, predicateOk := sessionTokenRule["predicate"].(map[string]interface{}); predicateOk {
			if customProjection, customProjectionOk := predicate["customProjection"].(map[string]interface{}); customProjectionOk {
				d.Set("custom_json", customProjection["customJson"].(string))
			}

			if logicalPredicate, logicalPredicateOk := predicate["logicalPredicate"].(map[string]interface{}); logicalPredicateOk {
				children := logicalPredicate["children"].([]interface{})
				var conditionList []interface{}
				for _, child := range children {
					childMap := child.(map[string]interface{})
					if attributePredicate, attributePredicateOk := childMap["attributePredicate"].(map[string]interface{}); attributePredicateOk {
						if attributeProjection, attributeProjectionOk := attributePredicate["attributeProjection"].(map[string]interface{}); attributeProjectionOk {
							matchCondition := attributeProjection["matchCondition"].(map[string]interface{})
							condition := map[string]interface{}{
								"key":      matchCondition["stringValue"].(string),
								"operator": matchCondition["matchOperator"].(string),
							}
							if attributePredicate["requestAttributeKeyLocation"] == "HEADER" {
								conditionList = append(conditionList, map[string]interface{}{"condition_request_header": []interface{}{condition}})
							} else if attributePredicate["requestAttributeKeyLocation"] == "COOKIE" {
								conditionList = append(conditionList, map[string]interface{}{"condition_request_cookie": []interface{}{condition}})
							} else if attributePredicate["requestAttributeKeyLocation"] == "QUERY_PARAMETER" {
								conditionList = append(conditionList, map[string]interface{}{"condition_request_query_param": []interface{}{condition}})
							} else if attributePredicate["responseAttributeKeyLocation"] == "HEADER" {
								conditionList = append(conditionList, map[string]interface{}{"condition_response_header": []interface{}{condition}})
							} else if attributePredicate["responseAttributeKeyLocation"] == "COOKIE" {
								conditionList = append(conditionList, map[string]interface{}{"condition_response_cookie": []interface{}{condition}})
							}
						}
					}
				}
				d.Set("token_extraction_condition_list", conditionList)
			}
		}
	}

	return nil
}

func resourceSessionIdentificationCustomRuleUpdate(d *schema.ResourceData, meta interface{}) error {
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

				if v, ok := conditionMap["condition_response_header"]; ok {
					conditionListStr += buildConditionListResponse("HEADER", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_response_cookie"]; ok {
					conditionListStr += buildConditionListResponse("COOKIE", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_header"]; ok {
					conditionListStr += buildConditionListRequest("HEADER", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_cookie"]; ok {
					conditionListStr += buildConditionListRequest("COOKIE", v.(*schema.Set).List())
				}
				if v, ok := conditionMap["condition_request_query_param"]; ok {
					conditionListStr += buildConditionListRequest("QUERY_PARAMETER", v.(*schema.Set).List())
				}
			}

			conditionListStr += `]}},`
		}
	}

	if conditionListStr == `predicate: { predicateType: LOGICAL, logicalPredicate: { operator: AND, children: []}}, ` {
		conditionListStr = ""
	}

	customJsonStr := ""
	if v, ok := d.GetOk("custom_json"); ok {
		customJsonStr = fmt.Sprintf(`customJson: "%s",`, v.(string))
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
				tokenType: CUSTOM
				sessionTokenValueRule: {
				  projectionRoot: {
					projectionType: CUSTOM
					customProjection: {
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
	`, id, name, descriptionStr, scopeStr, conditionListStr, customJsonStr)

	var response map[string]interface{}
	log.Printf(query)
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	log.Printf("GraphQL response: %s", responseStr)

	if response["data"] != nil && response["data"].(map[string]interface{})["updateSessionIdentificationRuleV2"] != nil {
		d.SetId(id)
	} else {
		return fmt.Errorf("could not update Session Identification response rule, no ID returned")
	}

	return nil
}

func resourceSessionIdentificationCustomRuleDelete(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("failed to delete Session Identification custom rule")
	}

	d.SetId("")
	return nil
}
