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

func ResourceIpTypeRuleBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpTypeRuleBlockCreate,
		Read:   resourceIpTypeRuleBlockRead,
		Update: resourceIpTypeRuleBlockUpdate,
		Delete: resourceIpTypeRuleBlockDelete,

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
				Default:     "BLOCK",
			},
			"event_severity": {
				Type:        schema.TypeString,
				Description: "Generated event severity among LOW,MEDIUM,HIGH,CRITICAL",
				Required:    true,
			},
			"expiration": {
				Type:        schema.TypeString,
				Description: "expiration for Block actions",
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
			"environment": {
				Type:        schema.TypeSet,
				Description: "environment where it will be applied",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_types": {
				Type:        schema.TypeSet,
				Description: "Ip types to include for the rule",
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

func resourceIpTypeRuleBlockCreate(d *schema.ResourceData, meta interface{}) error {
	expiration := d.Get("expiration").(string)
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	ip_types := d.Get("ip_types").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	inject_request_headers := d.Get("inject_request_headers").([]interface{})
	environment := d.Get("environment").(*schema.Set).List()

	exipiryDurationString := ReturnMalicousSourcesExipiryDuration(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	finalAgentEffectQuery := custom_signature.ReturnfinalAgentEffectQuery(inject_request_headers)

	query := fmt.Sprintf(CREATE_IP_TYPE_BLOCK, name, description, event_severity, rule_action, exipiryDurationString, finalAgentEffectQuery, strings.Join(common.InterfaceToEnumStringSlice(ip_types), ","), envQuery)
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

func resourceIpTypeRuleBlockRead(d *schema.ResourceData, meta interface{}) error {
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
		d.Set("rule_action", ruleDetails["ruleActionType"].(string))
		if expirationDetails, ok := action["expirationDetails"].(map[string]interface{}); ok {
			d.Set("expiration", expirationDetails["expirationDuration"].(string))
		} else {
			d.Set("expiration", "")
		}
		injectedHeaders := []map[string]interface{}{}
		if ruleEffects, ok := action["effects"].([]interface{}); ok {
			for _, ruleEffect := range ruleEffects {
				ruleEffectMap := ruleEffect.(map[string]interface{})
				if agentEffect, ok := ruleEffectMap["agentEffect"].(map[string]interface{}); ok {
					if agentModifications, ok := agentEffect["agentModifications"].([]interface{}); ok {
						for _, agentModification := range agentModifications {
							agentModificationMap := agentModification.(map[string]interface{})
							injectedHeader := map[string]interface{}{
								"header_key":   agentModificationMap["headerInjection"].(map[string]interface{})["key"].(string),
								"header_value": agentModificationMap["headerInjection"].(map[string]interface{})["value"].(string),
							}
							injectedHeaders = append(injectedHeaders, injectedHeader)
						}
					}
				}
			}
		}
		d.Set("inject_request_headers", injectedHeaders)
	}

	condition := ruleData["conditions"].([]interface{})[0].(map[string]interface{})
	ipLocationTypeCondition := condition["ipLocationTypeCondition"].(map[string]interface{})
	d.Set("ip_types", ipLocationTypeCondition["ipLocationTypes"].([]interface{}))

	if ruleScope, ok := ruleData["scope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
			if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
				d.Set("environment", environmentIds)
			} else {
				d.Set("environment", []interface{}{})
			}
		}
	}
	return nil
}

func resourceIpTypeRuleBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	expiration := d.Get("expiration").(string)
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	ip_types := d.Get("ip_types").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	inject_request_headers := d.Get("inject_request_headers").([]interface{})
	environment := d.Get("environment").(*schema.Set).List()

	exipiryDurationString := ReturnMalicousSourcesExipiryDuration(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	finalAgentEffectQuery := custom_signature.ReturnfinalAgentEffectQuery(inject_request_headers)

	query := fmt.Sprintf(UPDATE_IP_TYPE_BLOCK, id, name, description, event_severity, rule_action, exipiryDurationString, finalAgentEffectQuery, strings.Join(common.InterfaceToEnumStringSlice(ip_types), ","), envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateMaliciousSourcesRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	d.SetId(updatedId)
	return nil
}

func resourceIpTypeRuleBlockDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteMaliciousSourcesRule(d, meta)
	return nil
}
