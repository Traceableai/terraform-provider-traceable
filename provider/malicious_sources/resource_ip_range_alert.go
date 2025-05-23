package malicious_sources

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"log"
)

func ResourceIpRangeRuleAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpRangeRuleAlertCreate,
		Read:   resourceIpRangeRuleAlertRead,
		Update: resourceIpRangeRuleAlertUpdate,
		Delete: resourceIpRangeRuleAlertDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the policy",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the policy",
				Optional:    true,
			},
			"event_severity": {
				Type:        schema.TypeString,
				Description: "Generated event severity among LOW,MEDIUM,HIGH,CRITICAL",
				Required:    true,
			},
			"rule_action": {
				Type:        schema.TypeString,
				Description: "Need to provide the action to be performed",
				Optional:    true,
				Default:     "RULE_ACTION_ALERT",
			},
			"environment": {
				Type:        schema.TypeSet,
				Description: "environment where it will be applied",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"raw_ip_range_data": {
				Type:        schema.TypeSet,
				Description: "IPV4/V6 range for the rule",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
		},
	}
}

func resourceIpRangeRuleAlertCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	raw_ip_range_data := d.Get("raw_ip_range_data").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	event_severity := d.Get("event_severity").(string)
	environment := d.Get("environment").(*schema.Set).List()
	injectRequestHeaders := d.Get("inject_request_headers").([]interface{})
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	finalAgentEffectQuery := custom_signature.ReturnfinalAgentEffectQuery(injectRequestHeaders)

	query := fmt.Sprintf(CREATE_IP_RANGE_ALERT, name, common.InterfaceToStringSlice(raw_ip_range_data), rule_action, description, event_severity, finalAgentEffectQuery, envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createIpRangeRule")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId(id)
	return nil
}

func resourceIpRangeRuleAlertRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(IP_RANGE_READ_QUERY, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "ipRangeRules", id)
	if len(ruleData) == 0 {
		d.SetId("")
		return nil
	}
	ruleDetails := ruleData["ruleDetails"].(map[string]interface{})
	d.Set("name", ruleDetails["name"].(string))
	d.Set("description", ruleDetails["description"].(string))
	d.Set("rule_action", ruleDetails["ruleAction"].(string))
	event_severity, ok := ruleDetails["eventSeverity"]
	if ok {
		d.Set("event_severity", event_severity)
	} else {
		d.Set("event_severity", "")
	}

	d.Set("rule_action", ruleDetails["ruleAction"].(string))

	rawIpRangeData := ruleDetails["rawIpRangeData"].([]interface{})
	d.Set("raw_ip_range_data", rawIpRangeData)

	envFlag := true
	if ruleScope, ok := ruleDetails["ruleScope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
			if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
				d.Set("environment", environmentIds)
				envFlag = false
			}
		}
	}
	if envFlag {
		d.Set("environment", []interface{}{})
	}
	injectedHeader := SetInjectedHeaders(ruleDetails)
	d.Set("inject_request_headers", injectedHeader)
	return nil
}

func resourceIpRangeRuleAlertUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	raw_ip_range_data := d.Get("raw_ip_range_data").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	event_severity := d.Get("event_severity").(string)
	environment := d.Get("environment").(*schema.Set).List()
	injectRequestHeaders := d.Get("inject_request_headers").([]interface{})
	finalAgentEffectQuery := custom_signature.ReturnfinalAgentEffectQuery(injectRequestHeaders)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(UPDATE_IP_RANGE_ALERT, id, name, common.InterfaceToStringSlice(raw_ip_range_data), rule_action, description, event_severity, finalAgentEffectQuery, envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateIpRangeRule")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId(updatedId)
	return nil
}

func resourceIpRangeRuleAlertDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteIPRangeRule(d, meta)
	return nil
}
