package malicious_sources

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceRegionRuleAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegionRuleAlertCreate,
		Read:   resourceRegionRuleAlertRead,
		Update: resourceRegionRuleAlertUpdate,
		Delete: resourceRegionRuleAlertDelete,

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
				Default: "ALERT",
			},
			"event_severity": {
				Type:        schema.TypeString,
				Description: "Generated event severity among LOW,MEDIUM,HIGH,CRITICAL",
				Required:    true,
			},
			"environment": {
				Type:        schema.TypeSet,
				Description: "environment where it will be applied",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"regions": {
				Type:        schema.TypeSet,
				Description: "Region identifiers for rule",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRegionRuleAlertCreate(d *schema.ResourceData, meta interface{}) error {
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	regions := d.Get("regions").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()

	envQuery := custom_signature.ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(CREATE_REGION_RULE_ALERT,name,strings.Join(common.InterfaceToStringSlice(regions), ","),rule_action,description,event_severity,envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	id,err := common.GetIdFromResponse(responseStr,"createRegionRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	d.SetId(id)
	return nil
}

func resourceRegionRuleAlertRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(REGION_READ_QUERY, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "regionRules", id)
	if len(ruleData) == 0 {
		d.SetId("")
		return nil
	}
	ruleDetails := ruleData["ruleDetails"].(map[string]interface{})
	d.Set("name", ruleDetails["name"].(string))
	d.Set("description", ruleDetails["description"].(string))
	d.Set("rule_action", ruleDetails["ruleAction"].(string))
	event_severity, ok := ruleDetails["eventSeverity"].(string)
	if ok {
		d.Set("event_severity", event_severity)
	} else {
		d.Set("event_severity", nil)
	}
	d.Set("rule_action", ruleDetails["ruleAction"].(string))
	
	rawIpRangeData := ruleDetails["rawIpRangeData"].([]interface{})
	d.Set("regions", rawIpRangeData)

	if ruleScope, ok := ruleData["ruleScope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
			if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
				d.Set("environment", environmentIds)
			}else{
				d.Set("environment",[]interface{}{})
			}
		}
	}
	return nil
}

func resourceRegionRuleAlertUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	regions := d.Get("regions").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()

	envQuery := custom_signature.ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(UPDATE_REGION_RULE_ALERT,id,name,strings.Join(common.InterfaceToStringSlice(regions), ","),rule_action,description,event_severity,envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	updatedId,err := common.GetIdFromResponse(responseStr,"updateRegionRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	d.SetId(updatedId)
	return nil
}

func resourceRegionRuleAlertDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteRegionRule(d,meta)
	return nil
}
