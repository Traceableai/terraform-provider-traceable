package malicious_sources

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"log"
)

func ResourceRegionRuleBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegionRuleBlockCreate,
		Read:   resourceRegionRuleBlockRead,
		Update: resourceRegionRuleBlockUpdate,
		Delete: resourceRegionRuleBlockDelete,

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
				Description: "Need to provide the action to be performed (BLOCK,BLOCK_ALL_EXCEPT)",
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

func resourceRegionRuleBlockCreate(d *schema.ResourceData, meta interface{}) error {
	expiration := d.Get("expiration").(string)
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	regions := d.Get("regions").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()
	exipiryDurationString := RegionRuleExpiryString(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	regionIds, ok := MapCountryNameToRegionId(regions, meta)
	if ok != nil {
		return fmt.Errorf("error: %s", ok)
	}
	query := fmt.Sprintf(CREATE_REGION_RULE_BLOCK, name, common.InterfaceToStringSlice(regionIds), rule_action, description, event_severity, exipiryDurationString, envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createRegionRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	d.SetId(id)
	return nil
}

func resourceRegionRuleBlockRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("name", ruleData["name"].(string))
	d.Set("description", ruleData["description"].(string))
	d.Set("rule_action", ruleData["type"].(string))
	event_severity, ok := ruleData["eventSeverity"].(string)
	if ok {
		d.Set("event_severity", event_severity)
	} else {
		d.Set("event_severity", nil)
	}
	expiration, ok := ruleData["expiration"]
	if ok {
		d.Set("expiration", expiration)
	} else {
		d.Set("expiration", "")
	}

	finalRegionsState := []interface{}{}
	regionsData := ruleData["regions"].([]interface{})
	for _, regions := range regionsData {
		regionsMap := regions.(map[string]interface{})
		countryName := regionsMap["name"].(string)
		finalRegionsState = append(finalRegionsState, countryName)
	}
	d.Set("regions", finalRegionsState)

	envFlag := false
	if ruleScope, ok := ruleData["ruleScope"].(map[string]interface{}); ok {
		if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
			if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
				d.Set("environment", environmentIds)
				envFlag = true
			}
		}
	}
	if !envFlag {
		d.Set("environment", []interface{}{})
	}
	return nil
}

func resourceRegionRuleBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	expiration := d.Get("expiration").(string)
	event_severity := d.Get("event_severity").(string)
	name := d.Get("name").(string)
	regions := d.Get("regions").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()

	exipiryDurationString := RegionRuleExpiryString(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)
	regionIds, ok := MapCountryNameToRegionId(regions, meta)
	if ok != nil {
		return fmt.Errorf("error: %s", ok)
	}
	query := fmt.Sprintf(UPDATE_REGION_RULE_BLOCK, id, name, common.InterfaceToStringSlice(regionIds), rule_action, description, event_severity, exipiryDurationString, envQuery)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateRegionRule")
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	d.SetId(updatedId)
	return nil
}

func resourceRegionRuleBlockDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteRegionRule(d, meta) //fix this
	return nil
}
