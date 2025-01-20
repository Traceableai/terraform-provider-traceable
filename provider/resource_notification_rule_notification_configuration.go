package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotificationRuleNotificationConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRuleNotificationConfigurationCreate,
		Read:   resourceNotificationRuleNotificationConfigurationRead,
		Update: resourceNotificationRuleNotificationConfigurationUpdate,
		Delete: resourceNotificationRuleNotificationConfigurationDelete,

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
			"notification_config_types": {
				Type:        schema.TypeSet,
				Description: "For which label config change we need notification (NOTIFICATION_CHANNEL/NOTIFICATION_RULE)",
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
				Default:     "NOTIFICATION_CONFIG_CHANGE_EVENT",
			},
		},
	}
}

func resourceNotificationRuleNotificationConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	channel_id := d.Get("channel_id").(string)
	event_types := d.Get("event_types").(*schema.Set).List()
	notification_config_types := d.Get("notification_config_types").(*schema.Set).List()
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
					notificationConfigChangeCondition: {
						notificationConfigChangeTypes: %s
						notificationTypes:" %s
					}
				}
				channelId: "%s"
				%s
			}
		) {
			ruleId
		}
	}`, category, name, event_types, notification_config_types, channel_id, frequencyString)
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
func resourceNotificationRuleNotificationConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	readQuery := `{
	notificationRules {
		results {
		ruleId
		ruleName
		channelId
		integrationTarget {
			type
			integrationId
		}
		category
		eventConditions {
			notificationConfigChangeCondition {
			notificationConfigChangeTypes
			notificationTypes
			}
		}
		rateLimitIntervalDuration
		}
	}
	}`
	var response map[string]interface{}
	responseStr, err := ExecuteQuery(readQuery, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", readQuery)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error:%s", err)
	}
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "notificationRules", id, "ruleId", "ruleName")
	if len(ruleDetails) == 0 {
		d.SetId("")
		return nil
	}
	d.Set("name", ruleDetails["ruleName"])
	d.Set("category", ruleDetails["category"])
	d.Set("channel_id", ruleDetails["channelId"])
	eventConditions := ruleDetails["eventConditions"]
	log.Printf("logss %s", eventConditions)
	notificationConfigChangeCondition := eventConditions.(map[string]interface{})["notificationConfigChangeCondition"]
	if notificationConfigChangeCondition != nil {
		notificationConfigChangeTypes := notificationConfigChangeCondition.(map[string]interface{})["notificationConfigChangeTypes"].([]interface{})
		notificationTypes := notificationConfigChangeCondition.(map[string]interface{})["notificationTypes"].([]interface{})
		d.Set("event_types", schema.NewSet(schema.HashString, notificationConfigChangeTypes))
		d.Set("notification_config_types", schema.NewSet(schema.HashString, notificationTypes))
	}

	if val, ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency", val)
	}
	return nil
}

func resourceNotificationRuleNotificationConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId := d.Id()
	name := d.Get("name").(string)
	channel_id := d.Get("channel_id").(string)
	event_types := d.Get("event_types").(*schema.Set).List()
	notification_config_types := d.Get("notification_config_types").(*schema.Set).List()
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
					notificationConfigChangeCondition: {
						notificationConfigChangeTypes: %s
						notificationTypes:" %s
					}
				}
				channelId: "%s"
				%s
			}
		) {
			ruleId
		}
	}`, category, ruleId, name, event_types, notification_config_types, channel_id, frequencyString)
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

func resourceNotificationRuleNotificationConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
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
