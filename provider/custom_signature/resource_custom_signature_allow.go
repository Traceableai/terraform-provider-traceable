package custom_signature

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceCustomSignatureAllowRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourceCustomSignatureAllowCreate,
		Read:   ResourceCustomSignatureAllowRead,
		Update: ResourceCustomSignatureAllowUpdate,
		Delete: ResourceCustomSignatureAllowDelete,

		Schema: map[string]*schema.Schema{
			"rule_type": {
				Type:        schema.TypeString,
				Description: "Type of custom signature rule",
				Optional:    true,
				Default:     "ALLOW",
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the custom signature allow rule",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the custom signature allow rule",
				Optional:    true,
			},
			"environments": {
				Type:        schema.TypeSet,
				Description: "Environment of the custom signature allow rule",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"req_res_conditions": {
				Type:        schema.TypeList,
				Description: "Request/Response conditions for the rule",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"match_category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"match_operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"match_value": {
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
			"allow_expiry_duration": {
				Type:        schema.TypeString,
				Description: "Time to allow the rule (Time is seconds)",
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
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func ResourceCustomSignatureAllowCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	rule_type := d.Get("rule_type").(string)
	description := d.Get("description").(string)
	environments := d.Get("environments").(*schema.Set).List()
	// 	disabled := d.Get("disabled").(bool)
	req_res_conditions := d.Get("req_res_conditions").([]interface{})
	custom_sec_rule := d.Get("custom_sec_rule").(string)
	custom_sec_rule = strings.TrimSpace(EscapeString(custom_sec_rule))
	allow_expiry_duration := d.Get("allow_expiry_duration").(string)

	envQuery := ReturnEnvScopedQuery(environments)
	finalReqResConditionsQuery := ReturnReqResConditionsQuery(req_res_conditions)
	customSecRuleQuery := ReturnCustomSecRuleQuery(custom_sec_rule)

	if finalReqResConditionsQuery == "" && custom_sec_rule == "" {
		return fmt.Errorf("please provide one of req_res_conditions or custom_sec_rule")
	}

	exipiryDurationString := ReturnExipiryDuration(allow_expiry_duration)

	query := fmt.Sprintf(ALLOW_CREATE_QUERY, name, description, rule_type, finalReqResConditionsQuery, customSecRuleQuery, envQuery, exipiryDurationString)

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
	id := response["data"].(map[string]interface{})["createCustomSignatureRule"].(map[string]interface{})["id"].(string)

	d.SetId(id)
	return nil
}

func ResourceCustomSignatureAllowRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("rule_type", "ALLOW")
	d.Set("disabled", ruleDetails["disabled"].(bool))
	reqResConditions := []map[string]interface{}{}
	customSecRuleFlag := true
	if ruleDefinition, ok := ruleDetails["ruleDefinition"].(map[string]interface{}); ok {
		if clauseGroup, ok := ruleDefinition["clauseGroup"].(map[string]interface{}); ok {
			if clauses, ok := clauseGroup["clauses"].([]interface{}); ok {
				for _, clause := range clauses {
					if clauseMap, ok := clause.(map[string]interface{}); ok {
						if clauseType, exists := clauseMap["clauseType"].(string); exists && clauseType == "MATCH_EXPRESSION" {
							if matchExpression, ok := clauseMap["matchExpression"].(map[string]interface{}); ok {
								reqResCondition := map[string]interface{}{
									"match_key":      matchExpression["matchKey"],
									"match_category": matchExpression["matchCategory"],
									"match_operator": matchExpression["matchOperator"],
									"match_value":    matchExpression["matchValue"],
								}
								reqResConditions = append(reqResConditions, reqResCondition)
							}
						} else if clauseType == "CUSTOM_SEC_RULE" {
							customSecRule := (clauseMap["customSecRule"].(map[string]interface{})["inputSecRuleString"].(string))
							d.Set("custom_sec_rule", EscapeString(customSecRule))
							customSecRuleFlag = false
						}
					}
				}
				d.Set("req_res_conditions", reqResConditions)
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
	if blockingExpirationDuration, ok := ruleDetails["blockingExpirationDuration"].(string); ok {
		blockingExpirationDuration, _ := common.ConvertDurationToSeconds(blockingExpirationDuration)
		d.Set("allow_expiry_duration", blockingExpirationDuration)
	}

	return nil
}

func ResourceCustomSignatureAllowUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	rule_type := d.Get("rule_type").(string)
	id := d.Id()
	description := d.Get("description").(string)
	environments := d.Get("environments").(*schema.Set).List()
	req_res_conditions := d.Get("req_res_conditions").([]interface{})
	custom_sec_rule := d.Get("custom_sec_rule").(string)
	if !strings.Contains(custom_sec_rule, `\n`) {
		custom_sec_rule = strings.TrimSpace(EscapeString(custom_sec_rule))
	}
	allow_expiry_duration := d.Get("allow_expiry_duration").(string)

	disabled := d.Get("disabled").(bool)
	
	envQuery := ReturnEnvScopedQuery(environments)
	finalReqResConditionsQuery := ReturnReqResConditionsQuery(req_res_conditions)

	if finalReqResConditionsQuery == "" && custom_sec_rule == "" {
		return fmt.Errorf("please provide on of finalReqResConditionsQuery or custom_sec_rule")
	}

	customSecRuleQuery := ReturnCustomSecRuleQuery(custom_sec_rule)
	exipiryDurationString := ReturnExipiryDuration(allow_expiry_duration)

	query := fmt.Sprintf(ALLOW_UPDATE_QUERY, name, description, rule_type, finalReqResConditionsQuery, customSecRuleQuery, envQuery, exipiryDurationString, id, disabled)

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
	updatedId := response["data"].(map[string]interface{})["updateCustomSignatureRule"].(map[string]interface{})["id"].(string)

	d.SetId(updatedId)
	return nil
}

func ResourceCustomSignatureAllowDelete(d *schema.ResourceData, meta interface{}) error {
	err := DeleteCustomSignatureRule(d,meta)
	if err!=nil {
		return fmt.Errorf("error %s",err)
	}
	return nil
}