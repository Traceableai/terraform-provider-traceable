package malicious_sources

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
)

func ResourceIpTypeRuleAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpTypeRuleAlertCreate,
		Read:   resourceIpTypeRuleAlertRead,
		Update: resourceIpTypeRuleAlertUpdate,
		Delete: resourceIpTypeRuleAlertDelete,

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
			"rule_action": {
				Type:        schema.TypeString,
				Description: "Need to provide the action to be performed",
				Optional:    true,
				Default:     "ALERT",
			},
			"event_severity": {
				Type:        schema.TypeString,
				Description: "Generated event severity among LOW,MEDIUM,HIGH,CRITICAL",
				Required:    true,
			},
			"environment": {
				Type:        schema.TypeList,
				Description: "environment where it will be applied",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_types": {
				Type:        schema.TypeList,
				Description: "Ip types to include for the rule among ANONYMOUS VPN,HOSTING PROVIDER,PUBLIC PROXY,TOR EXIT NODE,BOT",
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

func resourceIpTypeRuleAlertCreate(d *schema.ResourceData, meta interface{}) error {
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	ip_types := d.Get("ip_types").([]interface{})
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").([]interface{})
	injectRequestHeaders := d.Get("inject_request_headers").([]interface{})
	finalAgentEffectQuery := custom_signature.ReturnfinalAgentEffectQuery(injectRequestHeaders)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(CREATE_IP_TYPE_ALERT, name, description, event_severity, rule_action, finalAgentEffectQuery, strings.Join(common.InterfaceToEnumStringSlice(ip_types), ","), envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createMaliciousSourcesRule")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId(id)
	return nil
}

func resourceIpTypeRuleAlertRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(MALICOUS_SOURCES_READ, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "maliciousSourcesRules", id)
	if len(ruleData) == 0 {
		d.SetId("")
		return nil
	}
	ruleDetails := ruleData["info"].(map[string]interface{})
	d.Set("name", ruleDetails["name"].(string))
	d.Set("description", ruleDetails["description"].(string))
	if action, ok := ruleDetails["action"].(map[string]interface{}); ok {
		d.Set("event_severity", action["eventSeverity"])
		d.Set("rule_action", action["ruleActionType"].(string))
	}

	condition := ruleDetails["conditions"].([]interface{})[0].(map[string]interface{})
	ipLocationTypeCondition := condition["ipLocationTypeCondition"].(map[string]interface{})
	d.Set("ip_types", ipLocationTypeCondition["ipLocationTypes"].([]interface{}))

	envFlag := true
	if ruleScope, ok := ruleData["scope"].(map[string]interface{}); ok {
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
	injectedHeader := SetInjectedHeaders(ruleData)
	d.Set("inject_request_headers", injectedHeader)
	return nil
}

func resourceIpTypeRuleAlertUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	ip_types := d.Get("ip_types").([]interface{})
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").([]interface{})
	injectRequestHeaders := d.Get("inject_request_headers").([]interface{})
	finalAgentEffectQuery := custom_signature.ReturnfinalAgentEffectQuery(injectRequestHeaders)
	envQuery := ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(UPDATE_IP_TYPE_ALERT, id, name, description, event_severity, rule_action, finalAgentEffectQuery, common.InterfaceToEnumStringSlice(ip_types), envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateMaliciousSourcesRule")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId(updatedId)
	return nil
}

func resourceIpTypeRuleAlertDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteMaliciousSourcesRule(d, meta)
	return nil
}
