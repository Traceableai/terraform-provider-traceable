package notification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceNotificationRuleBlockedThreatActivity() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRuleBlockedThreatActivityCreate,
		Read:   resourceNotificationRuleBlockedThreatActivityRead,
		Update: resourceNotificationRuleBlockedThreatActivityUpdate,
		Delete: resourceNotificationRuleBlockedThreatActivityDelete,

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
			"threat_types": {
				Type:        schema.TypeSet,
				Description: "Threat types for which you want notification",
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
				Default:     "BLOCKED_EVENT",
			},
		},
	}
}

func resourceNotificationRuleBlockedThreatActivityCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	threat_types := d.Get("threat_types").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	threatTypesString := "["
	for _, v := range threat_types {
		str := ""
		if IsCustomThreatEvent(v.(string)) {
			str = fmt.Sprintf(`{
				blockedThreatActivityConditionType: CUSTOM
				customDetectionCondition: { customDetectionType: %s }
			}`, v)
		} else if ok, val := IsPreDefinedThreatEvent(v.(string)); ok {
			str = fmt.Sprintf(`{
				blockedThreatActivityConditionType: PRE_DEFINED
				preDefinedBlockingCondition: { anomalyRuleId: "%s" }
			}`, val)
		}
		if str == "" {
			return fmt.Errorf("threat type %s not expected", v)
		}
		threatTypesString += str
	}
	threatTypesString += "]"

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
					blockedEventCondition: {
						blockedThreatActivityConditions: %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, category, name, threatTypesString, channel_id, frequencyString, envString)
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
func resourceNotificationRuleBlockedThreatActivityRead(d *schema.ResourceData, meta interface{}) error {
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
	blockedSecurityEventCondition := eventConditions.(map[string]interface{})["blockedEventCondition"]
	if blockedSecurityEventCondition == nil {
		d.Set("threat_types", schema.NewSet(schema.HashString, []interface{}{""}))
	}

	if val, ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency", val)
	}
	var threat_types []interface{}
	blockedThreatActivityConditions := blockedSecurityEventCondition.(map[string]interface{})["blockedThreatActivityConditions"].([]interface{})
	if blockedThreatActivityConditions != nil {
		for _, val := range blockedThreatActivityConditions {
			isCustom := val.(map[string]interface{})["blockedThreatActivityConditionType"]
			if isCustom == "CUSTOM" {
				customBlockingType := val.(map[string]interface{})["customBlockingCondition"].(map[string]interface{})["customBlockingType"]
				threat_types = append(threat_types, customBlockingType.(string))
			} else {
				preDefinedBlockingCondition := val.(map[string]interface{})["preDefinedBlockingCondition"].(map[string]interface{})["anomalyRuleId"]
				threat_types = append(threat_types, FindThreatByCrsId(preDefinedBlockingCondition.(string)))
			}
		}
		d.Set("threat_types", schema.NewSet(schema.HashString, threat_types))
	}
	return nil
}

func resourceNotificationRuleBlockedThreatActivityUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId := d.Id()
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	threat_types := d.Get("threat_types").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	threatTypesString := "["
	for _, v := range threat_types {
		str := ""
		if IsCustomThreatEvent(v.(string)) {
			str = fmt.Sprintf(`{
				blockedThreatActivityConditionType: CUSTOM
				customDetectionCondition: { customDetectionType: %s }
			}`, v)
		} else if ok, val := IsPreDefinedThreatEvent(v.(string)); ok {
			str = fmt.Sprintf(`{
				blockedThreatActivityConditionType: PRE_DEFINED
				preDefinedBlockingCondition: { anomalyRuleId: "%s" }
			}`, val)
		}
		if str == "" {
			return fmt.Errorf("threat type %s not expected", v)
		}
		threatTypesString += str
	}
	threatTypesString += "]"

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
					blockedEventCondition: {
						blockedThreatActivityConditions: %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`, ruleId, category, name, threatTypesString, channel_id, frequencyString, envString)
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

func resourceNotificationRuleBlockedThreatActivityDelete(d *schema.ResourceData, meta interface{}) error {
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
