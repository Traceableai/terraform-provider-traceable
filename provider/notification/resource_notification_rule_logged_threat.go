package notification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceNotificationRuleLoggedThreatActivity() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRuleLoggedThreatActivityCreate,
		Read:   resourceNotificationRuleLoggedThreatActivityRead,
		Update: resourceNotificationRuleLoggedThreatActivityUpdate,
		Delete: resourceNotificationRuleLoggedThreatActivityDelete,

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
			"severities": {
				Type:        schema.TypeSet,
				Description: "Severites of threat events you want to notify (LOW,MEDIUM,HIGH,CRITICAL)",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"impact": {
				Type:        schema.TypeSet,
				Description: "Impact of threat events you want to notify (LOW,MEDIUM,HIGH)",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"confidence": {
				Type:        schema.TypeSet,
				Description: "Confidence of threat events you want to notify (LOW,MEDIUM,HIGH)",
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
				Default:     "DETECTED_SECURITY_EVENT",
			},
		},
	}
}

func resourceNotificationRuleLoggedThreatActivityCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	threat_types := d.Get("threat_types").(*schema.Set).List()
	severities := d.Get("severities").(*schema.Set).List()
	impact := d.Get("impact").(*schema.Set).List()
	confidence := d.Get("confidence").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	severitiesString := "severities: ["
	for _, v := range severities {
		severitiesString += v.(string)
		severitiesString += ","
	}
	severitiesString = severitiesString[:len(severitiesString)-1]
	severitiesString += "]"
	if len(severities) == 4 {
		severitiesString = ""
	}

	impactString := "impactLevels: ["
	for _, v := range impact {
		impactString += v.(string)
		impactString += ","
	}
	impactString = impactString[:len(impactString)-1]
	impactString += "]"
	if len(impact) == 3 {
		impactString = ""
	}

	confidenceString := "confidenceLevels: ["
	for _, v := range confidence {
		confidenceString += v.(string)
		confidenceString += ","
	}
	confidenceString = confidenceString[:len(confidenceString)-1]
	confidenceString += "]"
	if len(confidence) == 3 {
		confidenceString = ""
	}

	threatTypesString := "["
	for _, v := range threat_types {
		str := ""
		if IsCustomThreatEvent(v.(string)) {
			str = fmt.Sprintf(`{
				detectedThreatActivityConditionType: CUSTOM
				customDetectionCondition: { customDetectionType: %s }
			}`, v)
		} else if ok, val := IsPreDefinedThreatEvent(v.(string)); ok {
			str = fmt.Sprintf(`{
				detectedThreatActivityConditionType: PRE_DEFINED
				preDefinedDetectionCondition: { anomalyRuleId: "%s" }
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
					detectedSecurityEventCondition: {
						detectedThreatActivityConditions: %s
						%s
						%s
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
	}`, category, name, threatTypesString, severitiesString, impactString, confidenceString, channel_id, frequencyString, envString)
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
func resourceNotificationRuleLoggedThreatActivityRead(d *schema.ResourceData, meta interface{}) error {
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
	detectedSecurityEventCondition := eventConditions.(map[string]interface{})["detectedSecurityEventCondition"]
	if detectedSecurityEventCondition != nil {
		severities := detectedSecurityEventCondition.(map[string]interface{})["severities"].([]interface{})
		if len(severities) == 0 {
			d.Set("severities", schema.NewSet(schema.HashString, []interface{}{"LOW", "MEDIUM", "HIGH", "CRITICAL"}))
		} else {
			d.Set("severities", schema.NewSet(schema.HashString, severities))
		}

		impact := detectedSecurityEventCondition.(map[string]interface{})["impactLevels"].([]interface{})
		if len(impact) == 0 {
			d.Set("impact", schema.NewSet(schema.HashString, []interface{}{"LOW", "MEDIUM", "HIGH"}))
		} else {
			d.Set("impact", schema.NewSet(schema.HashString, impact))
		}

		confidence := detectedSecurityEventCondition.(map[string]interface{})["confidenceLevels"].([]interface{})
		if len(confidence) == 0 {
			d.Set("confidence", schema.NewSet(schema.HashString, []interface{}{"LOW", "MEDIUM", "HIGH"}))
		} else {
			d.Set("confidence", schema.NewSet(schema.HashString, confidence))
		}
		var threat_types []interface{}
		detectedThreatActivityConditions := detectedSecurityEventCondition.(map[string]interface{})["detectedThreatActivityConditions"].([]interface{})
		for _, val := range detectedThreatActivityConditions {
			isCustom := val.(map[string]interface{})["detectedThreatActivityConditionType"]
			if isCustom == "CUSTOM" {
				customDetectionType := val.(map[string]interface{})["customDetectionCondition"].(map[string]interface{})["customDetectionType"]
				threat_types = append(threat_types, customDetectionType.(string))
			} else {
				preDefinedDetectionCondition := val.(map[string]interface{})["preDefinedDetectionCondition"].(map[string]interface{})["anomalyRuleId"]
				threat_types = append(threat_types, FindThreatByCrsId(preDefinedDetectionCondition.(string)))
			}
		}
		d.Set("threat_types", schema.NewSet(schema.HashString, threat_types))
	}

	if val, ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency", val)
	}
	return nil
}

func resourceNotificationRuleLoggedThreatActivityUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId := d.Id()
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	threat_types := d.Get("threat_types").(*schema.Set).List()
	severities := d.Get("severities").(*schema.Set).List()
	impact := d.Get("impact").(*schema.Set).List()
	confidence := d.Get("confidence").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)

	severitiesString := "severities: ["
	for _, v := range severities {
		severitiesString += v.(string)
		severitiesString += ","
	}
	severitiesString = severitiesString[:len(severitiesString)-1]
	severitiesString += "]"
	if len(severities) == 4 {
		severitiesString = ""
	}

	impactString := "impactLevels: ["
	for _, v := range impact {
		impactString += v.(string)
		impactString += ","
	}
	impactString = impactString[:len(impactString)-1]
	impactString += "]"
	if len(impact) == 3 {
		impactString = ""
	}

	confidenceString := "confidenceLevels: ["
	for _, v := range confidence {
		confidenceString += v.(string)
		confidenceString += ","
	}
	confidenceString = confidenceString[:len(confidenceString)-1]
	confidenceString += "]"
	if len(confidence) == 3 {
		confidenceString = ""
	}

	threatTypesString := "["
	for _, v := range threat_types {
		str := ""
		if IsCustomThreatEvent(v.(string)) {
			str = fmt.Sprintf(`{
				detectedThreatActivityConditionType: CUSTOM
				customDetectionCondition: { customDetectionType: %s }
			}`, v)
		} else if ok, val := IsPreDefinedThreatEvent(v.(string)); ok {
			str = fmt.Sprintf(`{
				detectedThreatActivityConditionType: PRE_DEFINED
				preDefinedDetectionCondition: { anomalyRuleId: "%s" }
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
					detectedSecurityEventCondition: {
						detectedThreatActivityConditions: %s
						%s
						%s
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
	}`, ruleId, category, name, threatTypesString, severitiesString, impactString, confidenceString, channel_id, frequencyString, envString)
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

func resourceNotificationRuleLoggedThreatActivityDelete(d *schema.ResourceData, meta interface{}) error {
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
