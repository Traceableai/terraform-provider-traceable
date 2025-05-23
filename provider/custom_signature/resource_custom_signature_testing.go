package custom_signature

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceCustomSignatureTestingRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourceCustomSignatureTestingCreate,
		Read:   ResourceCustomSignatureTestingRead,
		Update: ResourceCustomSignatureTestingUpdate,
		Delete: ResourceCustomSignatureTestingDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the custom signature allow rule",
				Required:    true,
			},
			"rule_type": {
				Type:        schema.TypeString,
				Description: "Type of custom signature rule",
				Optional:    true,
				Default:     "TESTING_DETECTION",
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the custom signature allow rule",
				Optional:    true,
			},
			"environments": {
				Type:        schema.TypeList,
				Description: "Environment of the custom signature allow rule",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"request_payload_single_valued_conditions": {
				Type:        schema.TypeList,
				Description: "Request payload single valued conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_category": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Accepts these two values REQUEST/RESPONSE",
						},
						"match_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Possible values HTTP_METHOD/PARAMETER_VALUE/PARAMETER_NAME/HEADER_VALUE",
						},
						"match_operator": {
							Type:        schema.TypeString,
							Description: "These oprators are applied on match_key they varied based on the match_key",
							Required:    true,
						},
						"match_value": {
							Type:        schema.TypeString,
							Description: "Value on which the operator will be applied",
							Required:    true,
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
						"match_category": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Accepts these two values REQUEST/RESPONSE",
						},
						"key_value_tag": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Accept these values COOKIE/PARAMETER/HEADER",
						},
						"key_match_operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "These oprators are applied on match_key they varied based on the match_key",
						},
						"match_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "String value for the key",
						},
						"value_match_operator": {
							Type:        schema.TypeString,
							Description: "Operator for the value possible values EQUALS/CONTAINS/NOT_EQUAL",
							Required:    true,
						},
						"match_value": {
							Type:        schema.TypeString,
							Description: "Value on which the operator will be applied",
							Required:    true,
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
							Type:        schema.TypeString,
							Description: "key operator",
							Required:    true,
						},
						"key_condition_value": {
							Type:        schema.TypeString,
							Description: "key value",
							Required:    true,
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
			"inject_request_headers": {
				Type:        schema.TypeList,
				Description: "Inject Data in Request header?",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"header_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"custom_sec_rule": {
				Type:        schema.TypeString,
				Description: "Custom sec rule",
				Optional:    true,
				StateFunc: func(val interface{}) string {
					return strings.TrimSpace(EscapeString(val.(string)))
				},
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Required:    true,
			},
		},
	}
}

func ResourceCustomSignatureTestingCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	environments := d.Get("environments").([]interface{})
	rule_type := d.Get("rule_type").(string)
	disabled := d.Get("disabled").(bool)
	requestPayloadSingleValuedConditions := d.Get("request_payload_single_valued_conditions").([]interface{})
	requestPayloadMultiValuedConditions := d.Get("request_payload_multi_valued_conditions").([]interface{})
	attribute_based_conditions := d.Get("attribute_based_conditions").([]interface{})
	custom_sec_rule := d.Get("custom_sec_rule").(string)
	inject_request_headers := d.Get("inject_request_headers").([]interface{})
	custom_sec_rule = strings.TrimSpace(EscapeString(custom_sec_rule))

	envQuery := ReturnEnvScopedQuery(environments)
	finalReqResConditionsQuery := ReturnReqResConditionsQuery(requestPayloadSingleValuedConditions, requestPayloadMultiValuedConditions)
	finalAttributeBasedConditionsQuery, _ := ReturnAttributeBasedConditionsQuery(attribute_based_conditions)

	if finalReqResConditionsQuery == "" && custom_sec_rule == "" && finalAttributeBasedConditionsQuery == "" {
		return fmt.Errorf("please provide on of finalReqResConditionsQuery or custom_sec_rule")
	}

	customSecRuleQuery := ReturnCustomSecRuleQuery(custom_sec_rule)

	finalAgentEffectQuery := ReturnfinalAgentEffectQuery(inject_request_headers)

	query := fmt.Sprintf(TEST_CREATE_QUERY, name, description, disabled, rule_type, finalAgentEffectQuery, finalReqResConditionsQuery, customSecRuleQuery, finalAttributeBasedConditionsQuery, envQuery)

	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	id, err := common.GetIdFromResponse(responseStr, "createCustomSignatureRule")
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	d.SetId(id)
	return nil
}

func ResourceCustomSignatureTestingRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("Id from read %s", id)

	responseStr, err := common.CallExecuteQuery(READ_QUERY, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}
	ruleDetails := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "customSignatureRules", id)
	if len(ruleDetails) == 0 {
		d.SetId("")
		return nil
	}
	d.Set("name", ruleDetails["name"].(string))
	d.Set("disabled", ruleDetails["disabled"].(bool))
	d.Set("rule_type", ruleDetails["ruleEffect"].(map[string]interface{})["eventType"].(string))

	singleValuedReqResConditions := []map[string]interface{}{}
	multiValuedReqResConditions := []map[string]interface{}{}
	injectedHeaders := []map[string]interface{}{}
	attributeBasedConditions := []map[string]interface{}{}
	if ruleEffect, ok := ruleDetails["ruleEffect"].(map[string]interface{}); ok {
		if effects, ok := ruleEffect["effects"].([]interface{}); ok {
			for _, effect := range effects {
				effectMap := effect.(map[string]interface{})
				if effectMap["ruleEffectType"].(string) == "AGENT_EFFECT" {
					agentEffect := effectMap["agentEffect"].(map[string]interface{})
					agentModifications := agentEffect["agentModifications"].([]interface{})
					for _, agentModification := range agentModifications {
						agentModificationMap := agentModification.(map[string]interface{})
						injectedHeader := map[string]interface{}{
							"header_key":   agentModificationMap["headerInjection"].(map[string]interface{})["key"].(string),
							"header_value": agentModificationMap["headerInjection"].(map[string]interface{})["value"].(string),
						}
						injectedHeaders = append(injectedHeaders, injectedHeader)
					}
					d.Set("inject_request_headers", injectedHeaders)
				}
			}
		}
	}
	customSecRuleFlag := true
	if ruleDefinition, ok := ruleDetails["ruleDefinition"].(map[string]interface{}); ok {
		if clauseGroup, ok := ruleDefinition["clauseGroup"].(map[string]interface{}); ok {
			if clauses, ok := clauseGroup["clauses"].([]interface{}); ok {
				for _, clause := range clauses {
					if clauseMap, ok := clause.(map[string]interface{}); ok {
						if clauseType, exists := clauseMap["clauseType"].(string); exists && (clauseType == "MATCH_EXPRESSION" || clauseType == "KEY_VALUE_EXPRESSION") {
							if clauseType == "MATCH_EXPRESSION" {
								if matchExpression, ok := clauseMap["matchExpression"].(map[string]interface{}); ok {
									reqResCondition := map[string]interface{}{
										"match_key":      matchExpression["matchKey"],
										"match_category": matchExpression["matchCategory"],
										"match_operator": matchExpression["matchOperator"],
										"match_value":    matchExpression["matchValue"],
									}
									singleValuedReqResConditions = append(singleValuedReqResConditions, reqResCondition)
								}
							} else if clauseType == "KEY_VALUE_EXPRESSION" {
								if keyValueExpression, ok := clauseMap["keyValueExpression"].(map[string]interface{}); ok {
									reqResCondition := map[string]interface{}{
										"match_key":            keyValueExpression["matchKey"],
										"match_value":          keyValueExpression["matchValue"],
										"key_value_tag":        keyValueExpression["keyValueTag"],
										"value_match_operator": keyValueExpression["valueMatchOperator"],
										"match_category":       keyValueExpression["matchCategory"],
										"key_match_operator":   keyValueExpression["keyMatchOperator"],
									}
									multiValuedReqResConditions = append(multiValuedReqResConditions, reqResCondition)
								}
							}
						} else if clauseType, exists := clauseMap["clauseType"].(string); exists && clauseType == "CUSTOM_SEC_RULE" {
							d.Set("custom_sec_rule", strings.TrimSpace(EscapeString(clauseMap["customSecRule"].(map[string]interface{})["inputSecRuleString"].(string))))
							customSecRuleFlag = false
						} else if clauseType, exists := clauseMap["clauseType"].(string); exists && clauseType == "ATTRIBUTE_KEY_VALUE_EXPRESSION" {
							if attributeKeyValueExpression, ok := clauseMap["attributeKeyValueExpression"].(map[string]interface{}); ok {
								if valueKey, ok := attributeKeyValueExpression["valueCondition"].(map[string]interface{}); ok {
									attributeBasedCondition := map[string]interface{}{
										"key_condition_operator":   attributeKeyValueExpression["keyCondition"].(map[string]interface{})["operator"],
										"key_condition_value":      attributeKeyValueExpression["keyCondition"].(map[string]interface{})["value"],
										"value_condition_value":    valueKey["value"],
										"value_condition_operator": valueKey["operator"],
									}
									attributeBasedConditions = append(attributeBasedConditions, attributeBasedCondition)
								} else {
									attributeBasedCondition := map[string]interface{}{
										"key_condition_operator": attributeKeyValueExpression["keyCondition"].(map[string]interface{})["operator"],
										"key_condition_value":    attributeKeyValueExpression["keyCondition"].(map[string]interface{})["value"],
									}
									attributeBasedConditions = append(attributeBasedConditions, attributeBasedCondition)
								}

							}
						}
					}
				}
				d.Set("attribute_based_conditions", attributeBasedConditions)
				d.Set("request_payload_single_valued_conditions", singleValuedReqResConditions)
				d.Set("request_payload_multi_valued_conditions", multiValuedReqResConditions)
			}
		}
	}
	if customSecRuleFlag {
		d.Set("custom_sec_rule", "")
	}
	environments := []string{}

	// Extract environment IDs from ruleScope.environmentScope
	if ruleScope, ok := ruleDetails["ruleScope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
			if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
				for _, envID := range environmentIds {
					if envStr, ok := envID.(string); ok {
						environments = append(environments, envStr)
					}
				}
			}
		}
	}

	d.Set("environments", environments)

	return nil
}

func ResourceCustomSignatureTestingUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	rule_type := d.Get("rule_type").(string)
	disabled := d.Get("disabled").(bool)
	environments := d.Get("environments").([]interface{})
	requestPayloadSingleValuedConditions := d.Get("request_payload_single_valued_conditions").([]interface{})
	requestPayloadMultiValuedConditions := d.Get("request_payload_multi_valued_conditions").([]interface{})
	attribute_based_conditions := d.Get("attribute_based_conditions").([]interface{})
	custom_sec_rule := d.Get("custom_sec_rule").(string)
	inject_request_headers := d.Get("inject_request_headers").([]interface{})
	if !strings.Contains(custom_sec_rule, `\n`) {
		custom_sec_rule = strings.TrimSpace(EscapeString(custom_sec_rule))
	}

	envQuery := ReturnEnvScopedQuery(environments)
	finalReqResConditionsQuery := ReturnReqResConditionsQuery(requestPayloadSingleValuedConditions, requestPayloadMultiValuedConditions)
	finalAttributeBasedConditionsQuery, _ := ReturnAttributeBasedConditionsQuery(attribute_based_conditions)
	if finalReqResConditionsQuery == "" && custom_sec_rule == "" && finalAttributeBasedConditionsQuery == "" {
		return fmt.Errorf("please provide on of finalReqResConditionsQuery or custom_sec_rule")
	}

	customSecRuleQuery := ReturnCustomSecRuleQuery(custom_sec_rule)

	finalAgentEffectQuery := ReturnfinalAgentEffectQuery(inject_request_headers)
	query := fmt.Sprintf(TEST_UPDATE_QUERY, id, disabled, name, description, rule_type, finalAgentEffectQuery, finalReqResConditionsQuery, customSecRuleQuery, finalAttributeBasedConditionsQuery, envQuery)

	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	updatedId, err := common.GetIdFromResponse(responseStr, "updateCustomSignatureRule")
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	d.SetId(updatedId)
	return nil
}

func ResourceCustomSignatureTestingDelete(d *schema.ResourceData, meta interface{}) error {
	err := DeleteCustomSignatureRule(d, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	return nil
}
