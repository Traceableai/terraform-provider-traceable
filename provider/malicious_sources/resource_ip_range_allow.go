package malicious_sources

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"log"
)

func ResourceIpRangeRuleAllow() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpRangeRuleAllowCreate,
		Read:   resourceIpRangeRuleAllowRead,
		Update: resourceIpRangeRuleAllowUpdate,
		Delete: resourceIpRangeRuleAllowDelete,

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
				Default:     "RULE_ACTION_ALLOW",
			},
			"expiration": {
				Type:        schema.TypeString,
				Description: "expiration for allow actions",
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
			"raw_ip_range_data": {
				Type:        schema.TypeSet,
				Description: "IPV4/V6 range for the rule",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceIpRangeRuleAllowCreate(d *schema.ResourceData, meta interface{}) error {
	expiration := d.Get("expiration").(string)
	name := d.Get("name").(string)
	raw_ip_range_data := d.Get("raw_ip_range_data").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()

	exipiryDurationString := ReturnExipiryDuration(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(CREATE_IP_RANGE_ALLOW, name, common.InterfaceToStringSlice(raw_ip_range_data), rule_action, description, exipiryDurationString, envQuery)
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

func resourceIpRangeRuleAllowRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("rule_action", ruleDetails["ruleAction"].(string))
	expiration, ok := ruleDetails["expiration"].(string)
	if ok {
		expiration, _ = common.ConvertDurationToSeconds(expiration)
		d.Set("expiration", expiration)
	} else {
		d.Set("expiration", "")
	}

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
	return nil
}

func resourceIpRangeRuleAllowUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	expiration := d.Get("expiration").(string)
	name := d.Get("name").(string)
	raw_ip_range_data := d.Get("raw_ip_range_data").(*schema.Set).List()
	rule_action := d.Get("rule_action").(string)
	description := d.Get("description").(string)
	environment := d.Get("environment").(*schema.Set).List()

	exipiryDurationString := ReturnExipiryDuration(expiration)
	envQuery := custom_signature.ReturnEnvScopedQuery(environment)

	query := fmt.Sprintf(UPDATE_IP_RANGE_ALLOW, id, name, common.InterfaceToStringSlice(raw_ip_range_data), rule_action, description, exipiryDurationString, envQuery)
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

func resourceIpRangeRuleAllowDelete(d *schema.ResourceData, meta interface{}) error {
	DeleteIPRangeRule(d, meta)
	return nil
}
