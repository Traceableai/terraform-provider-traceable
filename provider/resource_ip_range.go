package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIpRangeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpRangeRuleCreate,
		Read:   resourceIpRangeRuleRead,
		Update: resourceIpRangeRuleUpdate,
		Delete: resourceIpRangeRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "name of the policy",
				Required:    true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "description of the policy",
				Optional:    true,
			},
			"rule_action": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Need to provide the action to be performed",
				Required:    true,
			},
			"event_severity": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Generated event severity among LOW,MEDIUM,HIGH,CRITICAL",
				Required:    true,
			},
			"expiration": &schema.Schema{
				Type:        schema.TypeString,
				Description: "expiration for Allow and Block actions",
				Optional:    true,
			},
			"environment": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "environment where it will be applied",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"raw_ip_range_data": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "IPV4/V6 range to be Alerted/blcoked/Allowed ",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceIpRangeRuleCreate(d *schema.ResourceData, meta interface{}) error {
	expiry := d.Get("expiration").(string)
	eventSeverity := d.Get("event_severity").(string)
	Name := d.Get("name").(string)
	RawIpRangeData := toStringSlice(d.Get("raw_ip_range_data").(*schema.Set).List())
	RuleAction := d.Get("rule_action").(string)
	Description := d.Get("description").(string)

	envList := toStringSlice(d.Get("environment").(*schema.Set).List())

	var allEnv bool
	allEnv = false
	for _, env := range envList {
		if env == "" {
			allEnv = true
			break
		}
	}

	queryInput := map[string]string{
		"name":               Name,
		"eventSeverity":      eventSeverity,
		"description":        Description,
		"rawIpRangeData":     listToString(RawIpRangeData),
		"environmentScope":   listToString(envList),
		"expirationDuration": expiry,
		"ruleAction":         RuleAction,
	}

	query := "mutation { createIpRangeRule( create:{ ruleDetails:{"
	env := ""
	for key, value := range queryInput {
		if value != "" && !allEnv {
			if key == "environmentScope" {
				env = fmt.Sprintf(`ruleScope: {%s:{environmentIds:[%s]}}`, key, value)
			} else {
				if key == "rawIpRangeData" {
					query += fmt.Sprintf(`%s:[%s],`, key, value)
				} else if key == "ruleAction" || key == "eventSeverity" {
					query += fmt.Sprintf(`%s:%s,`, key, value)
				} else {
					query += fmt.Sprintf(`%s:"%s",`, key, value)
				}
			}
		}
		if value == "" && key == "description" {
			query += fmt.Sprintf(`%s:"%s",`, key, value)
		}
	}
	query = query[:len(query)-1]
	query += "}"
	if env != "" {
		query += "," + env
	}
	query += "}) {id}}"

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	log.Println("This is the graphql query %s", query)
	log.Println("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Println("Error:", err)
	}
	id := response["data"].(map[string]interface{})["createIpRangeRule"].(map[string]interface{})["id"].(string)
	d.SetId(id)
	return nil
}

func resourceIpRangeRuleRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Println("Id from read ", id)
	query := fmt.Sprintf("{ ipRangeRules { results { id internal disabled ruleDetails { name description rawIpRangeData ruleAction eventSeverity expirationDetails { expirationDuration expirationTimestamp } } ruleScope { environmentScope { environmentIds } } } } }	")

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ipRangeRules := response["data"].(map[string]interface{})["ipRangeRules"].(map[string]interface{})
	results := ipRangeRules["results"].([]interface{})
	for _, rule := range results {
		ruleData := rule.(map[string]interface{})
		rule_id := ruleData["id"].(string)

		if rule_id == id {
			ruleDetails := ruleData["ruleDetails"].(map[string]interface{})
			d.Set("name", ruleDetails["name"].(string))
			d.Set("description", ruleDetails["description"].(string))
			d.Set("rule_action", ruleDetails["ruleAction"].(string))
			d.Set("event_severity", ruleDetails["eventSeverity"].(string))
			d.Set("rule_action", ruleDetails["ruleAction"].(string))
			expiration, ok := ruleDetails["expiration"].(string)
			if ok {
				d.Set("expiration", expiration)
			} else {
				d.Set("expiration", nil)
			}

			rawIpRangeData := ruleDetails["rawIpRangeData"].([]interface{})
			d.Set("raw_ip_range_data", schema.NewSet(schema.HashString, convertToStringSlice(rawIpRangeData)))

			if ruleScope, ok := ruleData["ruleScope"].(map[string]interface{}); ok {
				if environmentScope, ok := ruleScope["environmentScope"].(map[string]interface{}); ok {
					if environmentIds, ok := environmentScope["environmentIds"].([]interface{}); ok {
						d.Set("environment", schema.NewSet(schema.HashString, convertToStringSlice(environmentIds)))
					}
				}
			}

			break
		}
	}
	return nil
}

func resourceIpRangeRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	expiry := d.Get("expiration").(string)
	eventSeverity := d.Get("event_severity").(string)
	Name := d.Get("name").(string)
	RawIpRangeData := toStringSlice(d.Get("raw_ip_range_data").(*schema.Set).List())
	RuleAction := d.Get("rule_action").(string)
	Description := d.Get("description").(string)

	query := "mutation { updateIpRangeRule( update:{ ruleDetails:{"
	envList := toStringSlice(d.Get("environment").(*schema.Set).List())

	queryInput := map[string]string{
		"name":               Name,
		"eventSeverity":      eventSeverity,
		"description":        Description,
		"rawIpRangeData":     listToString(RawIpRangeData),
		"environmentScope":   listToString(envList),
		"expirationDuration": expiry,
		"ruleAction":         RuleAction,
	}

	env := ""
	for key, value := range queryInput {
		if value != "" {
			if key == "environmentScope" {
				env = fmt.Sprintf(`ruleScope: {%s:{environmentIds:[%s]}}`, key, value)
			} else {
				if key == "rawIpRangeData" {
					query += fmt.Sprintf(`%s:[%s],`, key, value)
				} else if key == "ruleAction" || key == "eventSeverity" {
					query += fmt.Sprintf(`%s:%s,`, key, value)
				} else {
					query += fmt.Sprintf(`%s:"%s",`, key, value)
				}
			}
		}
		if value == "" && key == "description" {
			query += fmt.Sprintf(`%s:"%s",`, key, value)
		}
	}
	query = query[:len(query)-1]
	query += "}"
	if env != "" {
		query += "," + env
	}
	query += "," + fmt.Sprintf(`id:"%s"`, id)
	query += "}) {id}}"

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	log.Println("This is the graphql query %s", query)
	log.Println("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Println("Error:", err)
	}
	updatedId := response["data"].(map[string]interface{})["updateIpRangeRule"].(map[string]interface{})["id"].(string)
	d.SetId(updatedId)
	return nil
}

func resourceIpRangeRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(" mutation {deleteIpRangeRule(delete: {id: \"%s\"}) {success}}", id)
	_, err := executeQuery(query, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
