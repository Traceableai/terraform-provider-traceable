package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotificationRulePostureEvents() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRulePostureEventsCreate,
		Read:   resourceNotificationRulePostureEventsRead,
		Update: resourceNotificationRulePostureEventsUpdate,
		Delete: resourceNotificationRulePostureEventsDelete,

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
			"posture_events": {
				Type:        schema.TypeSet,
				Description: "Type of posture event want notification",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"risk_deltas": {
				Type:        schema.TypeSet,
				Description: "Applicable for risk score posture events (INCREASE,DECREASE)",
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
				Default:     "POSTURE_EVENT",
			},
		},
	}
}

func resourceNotificationRulePostureEventsCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	posture_events := d.Get("posture_events").(*schema.Set).List()
	risk_deltas := d.Get("risk_deltas").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

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
					postureEventCondition: {
						postureEvents: %s
						riskDeltas: %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, category, name, posture_events, risk_deltas, channel_id, frequencyString, envString)
	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Println("Error:", err)
	}
	id := response["data"].(map[string]interface{})["createNotificationRule"].(map[string]interface{})["ruleId"].(string)
	d.SetId(id)
	return nil
}
func resourceNotificationRulePostureEventsRead(d *schema.ResourceData, meta interface{}) error {
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
				postureEventCondition {
				postureEvents
				}
			}
			rateLimitIntervalDuration
			}
		}
	}
	`
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
	postureEventCondition := eventConditions.(map[string]interface{})["postureEventCondition"]
	if postureEventCondition != nil {
		postureEvents := postureEventCondition.(map[string]interface{})["postureEvents"].([]interface{})
		d.Set("posture_events", schema.NewSet(schema.HashString, postureEvents))
		if riskDeltas, ok := postureEventCondition.(map[string]interface{})["riskDeltas"].([]interface{}); ok {
			d.Set("risk_deltas", schema.NewSet(schema.HashString, riskDeltas))
		}
	}
	if val, ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency", val)
	}
	return nil
}

func resourceNotificationRulePostureEventsUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId := d.Id()
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	posture_events := d.Get("posture_events").(*schema.Set).List()
	risk_deltas := d.Get("risk_deltas").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

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
				ruleId: "%s"
				category: %s
				ruleName: "%s"
				eventConditions: {
					postureEventCondition: {
						postureEvents: %s
						riskDeltas: %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, ruleId, category, name, posture_events, risk_deltas, channel_id, frequencyString, envString)
	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		fmt.Println("Error:", err)
	}
	id := response["data"].(map[string]interface{})["updateNotificationRule"].(map[string]interface{})["ruleId"].(string)
	d.SetId(id)
	return nil
}

func resourceNotificationRulePostureEventsDelete(d *schema.ResourceData, meta interface{}) error {
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
