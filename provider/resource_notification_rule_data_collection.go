package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotificationRuleDataCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRuleDataCollectionCreate,
		Read:   resourceNotificationRuleDataCollectionRead,
		Update: resourceNotificationRuleDataCollectionUpdate,
		Delete: resourceNotificationRuleDataCollectionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the notification rule",
				Required:    true,
			},
			"environments": {
				Type:        schema.TypeSet,
				Description: "Environments where rule will be applicable",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"channel_id": {
				Type:        schema.TypeString,
				Description: "Reporting channel for this notification rule",
				Required:    true,
			},
			"agent_activity_type": {
				Type:        schema.TypeString,
				Description: "Agent activity type for which you want notification",
				Required:    true,
			},
			"agent_status_changes": {
				Type:        schema.TypeSet,
				Description: "Agent status change for which you want notification (ACTIVE/IDLE/OFFLINE)",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notification_frequency": {
				Type:        schema.TypeString,
				Description: "No more than one notification every configured notification_frequency (should be in this format PT1H for 1 hr)",
				Optional:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "Type of notification rule",
				Optional:    true,
				Default:     "DATA_COLLECTION_CHANGE_EVENT",
			},
		},
	}
}

func resourceNotificationRuleDataCollectionCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	agent_activity_type := d.Get("agent_activity_type").(string)
	agent_status_changes := d.Get("agent_status_changes").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	agentStatusChangesString := ""
	if agent_activity_type == "STATUS_CHANGE" {
		if len(agent_status_changes) == 0 {
			return fmt.Errorf("agent_status_changes is required here")
		}
		agentStatusChangesString = fmt.Sprintf(`agentStatusChanges: %s`, agent_status_changes)
	}
	if agent_activity_type != "STATUS_CHANGE" {
		if len(agent_status_changes) > 0 {
			return fmt.Errorf("agent_status_changes is not expected here")
		}
	}
	frequencyString := ""
	if notification_frequency != "" {
		frequencyString = fmt.Sprintf(`rateLimitIntervalDuration: "%s"`, notification_frequency)
	}
	envArrayString := "["
	for _, v := range environments {
		envArrayString += fmt.Sprintf(`"%s"`, v.(string))
		envArrayString += ","
	}
	envArrayString = envArrayString[:len(envArrayString)-1]
	envArrayString += "]"
	envString := fmt.Sprintf(`environmentScope: { environments: %s }`, envArrayString)

	if len(environments) == 0 || (len(environments) == 1 && environments[0] == "") {
		envString = ""
	}
	query := fmt.Sprintf(`mutation {
		createNotificationRule(
			input: {
				category: %s
				ruleName: "%s"
				eventConditions: {
					dataCollectionChangeCondition: {
						agentType: ""
						agentActivityType: %s
						%s
                	}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, category, name, agent_activity_type, agentStatusChangesString, channel_id, frequencyString, envString)
	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	id := response["data"].(map[string]interface{})["createNotificationRule"].(map[string]interface{})["ruleId"].(string)
	d.SetId(id)
	return nil
}
func resourceNotificationRuleDataCollectionRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	readQuery := `{
	notificationRules {
		results {
		ruleId
		ruleName
		environmentScope {
			environments
		}
		channelId
		integrationTarget {
			type
			integrationId
		}
		category
		eventConditions {
			dataCollectionChangeCondition {
			agentType
			agentActivityType
			agentStatusChanges 
			}
		}
		rateLimitIntervalDuration 
		}
	}
	}`
	var response map[string]interface{}
	responseStr, err := ExecuteQuery(readQuery, meta)
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql query %s", readQuery)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "notificationRules", id, "ruleId", "ruleName")
	if len(ruleDetails) == 0 {
		d.SetId("")
		return nil
	}
	d.Set("name", ruleDetails["ruleName"])
	d.Set("category", ruleDetails["category"])
	d.Set("channel_id", ruleDetails["channelId"])
	envs := ruleDetails["environmentScope"].(map[string]interface{})["environments"]
	d.Set("environments", schema.NewSet(schema.HashString, envs.([]interface{})))
	eventConditions := ruleDetails["eventConditions"]
	log.Printf("logss %s", eventConditions)
	dataCollectionChangeCondition := eventConditions.(map[string]interface{})["dataCollectionChangeCondition"]
	if dataCollectionChangeCondition != nil {
		agentActivityType := dataCollectionChangeCondition.(map[string]interface{})["agentActivityType"].(string)
		d.Set("agent_activity_type", agentActivityType)
		if agentActivityType == "STATUS_CHANGE" {
			agentStatusChanges := dataCollectionChangeCondition.(map[string]interface{})["agentStatusChanges"].([]interface{})
			if len(agentStatusChanges) == 0 {
				d.Set("agent_status_changes", schema.NewSet(schema.HashString, []interface{}{}))
			} else {
				d.Set("agent_status_changes", schema.NewSet(schema.HashString, agentStatusChanges))
			}
		}
	}

	if val, ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency", val)
	}
	return nil
}

func resourceNotificationRuleDataCollectionUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId := d.Id()
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	agent_activity_type := d.Get("agent_activity_type").(string)
	agent_status_changes := d.Get("agent_status_changes").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	agentStatusChangesString := ""
	if agent_activity_type == "STATUS_CHANGE" {
		if len(agent_status_changes) == 0 {
			return fmt.Errorf("agent_status_changes is required here")
		}
		agentStatusChangesString = fmt.Sprintf(`agentStatusChanges: %s`, agent_status_changes)
	}
	if agent_activity_type != "STATUS_CHANGE" {
		if len(agent_status_changes) > 0 {
			return fmt.Errorf("agent_status_changes is not expected here")
		}
	}
	frequencyString := ""
	if notification_frequency != "" {
		frequencyString = fmt.Sprintf(`rateLimitIntervalDuration: "%s"`, notification_frequency)
	}
	envArrayString := "["
	for _, v := range environments {
		envArrayString += fmt.Sprintf(`"%s"`, v.(string))
		envArrayString += ","
	}
	envArrayString = envArrayString[:len(envArrayString)-1]
	envArrayString += "]"
	envString := fmt.Sprintf(`environmentScope: { environments: %s }`, envArrayString)

	if len(environments) == 0 || (len(environments) == 1 && environments[0] == "") {
		envString = ""
	}
	query := fmt.Sprintf(`mutation {
		updateNotificationRule(
			input: {
				category: %s
				ruleId: "%s"
				ruleName: "%s"
				eventConditions: {
					dataCollectionChangeCondition: {
						agentType: ""
						agentActivityType: %s
						%s
                	}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, category, ruleId, name, agent_activity_type, agentStatusChangesString, channel_id, frequencyString, envString)
	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	id := response["data"].(map[string]interface{})["updateNotificationRule"].(map[string]interface{})["ruleId"].(string)
	d.SetId(id)
	return nil
}

func resourceNotificationRuleDataCollectionDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(`mutation {
		deleteNotificationRule(input: {ruleId: "%s"}) {
		  success
		}
	  }`, id)
	_, err := ExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	log.Println(query)
	d.SetId("")
	return nil
}
