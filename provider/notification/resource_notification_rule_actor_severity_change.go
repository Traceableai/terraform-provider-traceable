package notification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceNotificationRuleActorSeverityChange() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRuleActorSeverityChangeCreate,
		Read:   resourceNotificationRuleActorSeverityChangeRead,
		Update: resourceNotificationRuleActorSeverityChangeUpdate,
		Delete: resourceNotificationRuleActorSeverityChangeDelete,

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
			"actor_severities": {
				Type:        schema.TypeSet,
				Description: "Threat types for which you want notification",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"actor_ip_reputation_levels": {
				Type:        schema.TypeSet,
				Description: "Severites of threat events you want to notify (LOW,MEDIUM,HIGH,CRITICAL)",
				Required:    true,
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
				Default:     "ACTOR_SEVERITY_STATE_CHANGE_EVENT",
			},
		},
	}
}

func resourceNotificationRuleActorSeverityChangeCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	actor_severities := d.Get("actor_severities").(*schema.Set).List()
	actor_ip_reputation_levels := d.Get("actor_ip_reputation_levels").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	actorSeveritiesString := "["
	for _, v := range actor_severities {
		actorSeveritiesString += v.(string)
		actorSeveritiesString += ","
	}
	if len(actorSeveritiesString) > 1 {
		actorSeveritiesString = actorSeveritiesString[:len(actorSeveritiesString)-1]
	}

	actorSeveritiesString += "]"
	if len(actor_severities) == 4 {
		actorSeveritiesString = ""
	}

	actorIpRepString := "["
	for _, v := range actor_ip_reputation_levels {
		actorIpRepString += v.(string)
		actorIpRepString += ","
	}
	if len(actorIpRepString) > 1 {
		actorIpRepString = actorIpRepString[:len(actorIpRepString)-1]
	}

	actorIpRepString += "]"
	if len(actor_ip_reputation_levels) == 4 {
		actorIpRepString = ""
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
					actorSeverityStateChangeEventCondition: {
						actorSeverities: %s
						actorIpReputationLevels: %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, category, name, actorSeveritiesString, actorIpRepString, channel_id, frequencyString, envString)
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
	id := response["data"].(map[string]interface{})["createNotificationRule"].(map[string]interface{})["ruleId"].(string)
	d.SetId(id)
	return nil
}
func resourceNotificationRuleActorSeverityChangeRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	var response map[string]interface{}
	responseStr, err := common.CallExecuteQuery(NOTIFICATION_RULE_READ, meta)
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		_ = fmt.Errorf("Error:%s", err)
	}
	ruleDetails := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "notificationRules", id, "ruleId", "ruleName")
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
	actorSeverityStateChangeEventCondition := eventConditions.(map[string]interface{})["actorSeverityStateChangeEventCondition"]
	if actorSeverityStateChangeEventCondition != nil {
		actorSeverities := actorSeverityStateChangeEventCondition.(map[string]interface{})["actorSeverities"].([]interface{})
		d.Set("actor_severities", schema.NewSet(schema.HashString, actorSeverities))
		actorIpReputationLevels := actorSeverityStateChangeEventCondition.(map[string]interface{})["actorIpReputationLevels"].([]interface{})
		d.Set("actor_ip_reputation_levels", schema.NewSet(schema.HashString, actorIpReputationLevels))
	}
	if val, ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency", val)
	}
	return nil
}

func resourceNotificationRuleActorSeverityChangeUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId := d.Id()
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	actor_severities := d.Get("actor_severities").(*schema.Set).List()
	actor_ip_reputation_levels := d.Get("actor_ip_reputation_levels").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	actorSeveritiesString := "["
	for _, v := range actor_severities {
		actorSeveritiesString += v.(string)
		actorSeveritiesString += ","
	}
	if len(actorSeveritiesString) > 1 {
		actorSeveritiesString = actorSeveritiesString[:len(actorSeveritiesString)-1]
	}
	actorSeveritiesString += "]"
	if len(actor_severities) == 4 {
		actorSeveritiesString = ""
	}

	actorIpRepString := "["
	for _, v := range actor_ip_reputation_levels {
		actorIpRepString += v.(string)
		actorIpRepString += ","
	}
	if len(actorIpRepString) > 1 {
		actorIpRepString = actorIpRepString[:len(actorIpRepString)-1]
	}
	actorIpRepString += "]"
	if len(actor_ip_reputation_levels) == 4 {
		actorIpRepString = ""
	}

	frequencyString := ""
	if notification_frequency != "" {
		frequencyString = fmt.Sprintf(`rateLimitIntervalDuration: "%s"`, notification_frequency)
	}
	envString := fmt.Sprintf(`environmentScope: { environments: [%s] }`, common.InterfaceToStringSlice(environments))

	if len(environments) == 0 || (len(environments) == 1 && environments[0] == "") {
		envString = ""
	}
	query := fmt.Sprintf(`mutation {
		updateNotificationRule(
			input: {
				ruleId: "%s"
				category: %s
				ruleName: "%s"
				eventConditions: {
					actorSeverityStateChangeEventCondition: {
						actorSeverities: %s
						actorIpReputationLevels: %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, ruleId, category, name, actorSeveritiesString, actorIpRepString, channel_id, frequencyString, envString)
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
	id := response["data"].(map[string]interface{})["updateNotificationRule"].(map[string]interface{})["ruleId"].(string)
	d.SetId(id)
	return nil
}

func resourceNotificationRuleActorSeverityChangeDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_NOTIFICATION_RULE, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return err
	}
	log.Println(query)
	d.SetId("")
	return nil
}
