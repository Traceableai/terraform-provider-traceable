package notification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceNotificationRuleExcludeRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRuleExcludeRuleCreate,
		Read:   resourceNotificationRuleExcludeRuleRead,
		Update: resourceNotificationRuleExcludeRuleUpdate,
		Delete: resourceNotificationRuleExcludeRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the notification rule",
				Required:    true,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Description: "Reporting channel for this notification rule",
				Required:    true,
			},
			"event_types": {
				Type:        schema.TypeSet,
				Description: "For which operation we need notification (CREATE/UPDATE/DELETE)",
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
				Default:     "EXCLUDE_SPAN_RULE_CONFIG_CHANGE_EVENT",
			},
		},
	}
}

func resourceNotificationRuleExcludeRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	channel_id := d.Get("channel_id").(string)
	event_types := d.Get("event_types").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	frequencyString := ""
	if notification_frequency != "" {
		frequencyString = fmt.Sprintf(`rateLimitIntervalDuration: "%s"`, notification_frequency)
	}
	query := fmt.Sprintf(`mutation {
		createNotificationRule(
			input: {
				category: %s
				ruleName: "%s"
				eventConditions: {
					excludeSpanRuleConfigChangeEventCondition: {
						excludeSpanRuleConfigChangeTypes: %s
					}
				}
				channelId: "%s"
				%s
			}
		) {
			ruleId
		}
	}`, category, name, event_types, channel_id, frequencyString)
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
func resourceNotificationRuleExcludeRuleRead(d *schema.ResourceData, meta interface{}) error {
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
	eventConditions := ruleDetails["eventConditions"]
	log.Printf("logss %s", eventConditions)
	excludeSpanRuleConfigChangeEventCondition := eventConditions.(map[string]interface{})["excludeSpanRuleConfigChangeEventCondition"]
	if excludeSpanRuleConfigChangeEventCondition != nil {
		excludeSpanRuleConfigChangeTypes := excludeSpanRuleConfigChangeEventCondition.(map[string]interface{})["excludeSpanRuleConfigChangeTypes"].([]interface{})
		if len(excludeSpanRuleConfigChangeTypes) == 0 {
			d.Set("event_types", schema.NewSet(schema.HashString, []interface{}{}))
		} else {
			d.Set("event_types", schema.NewSet(schema.HashString, excludeSpanRuleConfigChangeTypes))
		}
	}

	if val, ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency", val)
	}
	return nil
}

func resourceNotificationRuleExcludeRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId := d.Id()
	name := d.Get("name").(string)
	channel_id := d.Get("channel_id").(string)
	event_types := d.Get("event_types").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	frequencyString := ""
	if notification_frequency != "" {
		frequencyString = fmt.Sprintf(`rateLimitIntervalDuration: "%s"`, notification_frequency)
	}
	query := fmt.Sprintf(`mutation {
		updateNotificationRule(
			input: {
				category: %s
				ruleId: "%s"
				ruleName: "%s"
				eventConditions: {
					excludeSpanRuleConfigChangeEventCondition: {
						excludeSpanRuleConfigChangeTypes: %s
					}
				}
				channelId: "%s"
				%s
			}
		) {
			ruleId
		}
	}`, category, ruleId, name, event_types, channel_id, frequencyString)
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

func resourceNotificationRuleExcludeRuleDelete(d *schema.ResourceData, meta interface{}) error {
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
