package malicious_sources

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"log"
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
	environment := d.Get("environment").(*schema.Set).List()

	exipiryDurationString := ReturnMalicousSourcesExipiryDuration(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(CREATE_IP_TYPE_BLOCK, name, description, event_severity, rule_action, exipiryDurationString, common.InterfaceToEnumStringSlice(ip_types), envQuery)
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
		d.Set("rule_action", action["ruleActionType"].(string))
		expirationFlag := true
		if expirationDetails, ok := action["expirationDetails"].(map[string]interface{}); ok {
			d.Set("expiration", expirationDetails["expirationDuration"].(string))
			expirationFlag = false
		}
		if expirationFlag {
			d.Set("expiration", "")
		}
	}

	condition := ruleDetails["conditions"].([]interface{})[0].(map[string]interface{})
	ipLocationTypeCondition := condition["ipLocationTypeCondition"].(map[string]interface{})
	d.Set("ip_types", ipLocationTypeCondition["ipLocationTypes"].([]interface{}))

	envFlag := true
	if ruleScope, ok := ruleDetails["scope"].(map[string]interface{}); ok {
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
	environment := d.Get("environment").(*schema.Set).List()

	exipiryDurationString := ReturnMalicousSourcesExipiryDuration(expiration)
	envQuery := ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(UPDATE_IP_TYPE_BLOCK, id, name, description, event_severity, rule_action, exipiryDurationString, common.InterfaceToEnumStringSlice(ip_types), envQuery)
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
