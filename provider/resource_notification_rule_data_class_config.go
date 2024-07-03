package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotificationRuleDataClassificationConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationRuleDataClassificationConfigCreate,
		Read:   resourceNotificationRuleDataClassificationConfigRead,
		Update: resourceNotificationRuleDataClassificationConfigUpdate,
		Delete: resourceNotificationRuleDataClassificationConfigDelete,

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
			"event_types": {
				Type:        schema.TypeSet,
				Description: "For which operation we need notification (CREATE/UPDATE/DELETE)",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sensitive_data_config_types": {
				Type:        schema.TypeSet,
				Description: "For which label config change we need notification (USER_ATTRIBUTION/DATA_TYPE/SESSION_IDENTIFICATION/DATA_SETS)",
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
				Default:     "SENSITIVE_DATA_CONFIG_CHANGE_EVENT",
			},
		},
	}
}

func resourceNotificationRuleDataClassificationConfigCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	channel_id := d.Get("channel_id").(string)
	event_types := d.Get("event_types").(*schema.Set).List()
	sensitive_data_config_types := d.Get("sensitive_data_config_types").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)
	environments := d.Get("environments").(*schema.Set).List()
	envArrayString:="["
	for _,v := range environments{
		envArrayString+=fmt.Sprintf(`"%s"`,v.(string))
		envArrayString+=","
	}
	envArrayString=envArrayString[:len(envArrayString)-1]
	envArrayString+="]"
	envString:=fmt.Sprintf(`environmentScope: { environments: %s }`,envArrayString)

	if len(environments)==0 || (len(environments)==1 && environments[0]==""){
		envString=""
	}
	frequencyString:=""
	if notification_frequency!=""{
		frequencyString=fmt.Sprintf(`rateLimitIntervalDuration: "%s"`,notification_frequency)
	}
	query:=fmt.Sprintf(`mutation {
		createNotificationRule(
			input: {
				category: %s
				ruleName: "%s"
				eventConditions: {
					sensitiveDataConfigChangeCondition: {
						sensitiveDataConfigChangeTypes: %s
						sensitiveDataConfigTypes:" %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`,category,name,event_types,sensitive_data_config_types,channel_id,frequencyString,envString)
	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
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
func resourceNotificationRuleDataClassificationConfigRead(d *schema.ResourceData, meta interface{}) error {
	id:=d.Id()
	readQuery:=`{
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
			sensitiveDataConfigChangeCondition {
			sensitiveDataConfigChangeTypes
			sensitiveDataConfigTypes
			}
		}
		rateLimitIntervalDuration
		}
	}
	}`
	var response map[string]interface{}
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		_=fmt.Errorf("Error: %s",err)
	}
	log.Printf("This is the graphql query %s", readQuery)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		_=fmt.Errorf("Error:%s", err)
	}
	ruleDetails:=getRuleDetailsFromRulesListUsingIdName(response,"notificationRules" ,id,"ruleId","ruleName")
	if len(ruleDetails)==0{
		d.SetId("")
		return nil
	}
	d.Set("name",ruleDetails["ruleName"])
	d.Set("category",ruleDetails["category"])
	d.Set("channel_id",ruleDetails["channelId"])
	eventConditions:=ruleDetails["eventConditions"]
	log.Printf("logss %s",eventConditions)
	sensitiveDataConfigChangeCondition:=eventConditions.(map[string]interface{})["sensitiveDataConfigChangeCondition"]
	if sensitiveDataConfigChangeCondition!=nil{
		sensitiveDataConfigChangeTypes:=sensitiveDataConfigChangeCondition.(map[string]interface{})["sensitiveDataConfigChangeTypes"].([]interface{})
		sensitiveDataConfigTypes:=sensitiveDataConfigChangeCondition.(map[string]interface{})["sensitiveDataConfigTypes"].([]interface{})
		d.Set("event_types",schema.NewSet(schema.HashString,sensitiveDataConfigChangeTypes))
		d.Set("sensitive_data_config_types",schema.NewSet(schema.HashString,sensitiveDataConfigTypes))
	}
	
	if val,ok := ruleDetails["rateLimitIntervalDuration"]; ok {
		d.Set("notification_frequency",val)
	}
	return nil
}

func resourceNotificationRuleDataClassificationConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId:=d.Id()
	name := d.Get("name").(string)
	channel_id := d.Get("channel_id").(string)
	event_types := d.Get("event_types").(*schema.Set).List()
	sensitive_data_config_types := d.Get("sensitive_data_config_types").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)
	category := d.Get("category").(string)
	environments := d.Get("environments").(*schema.Set).List()
	envArrayString:="["
	for _,v := range environments{
		envArrayString+=fmt.Sprintf(`"%s"`,v.(string))
		envArrayString+=","
	}
	envArrayString=envArrayString[:len(envArrayString)-1]
	envArrayString+="]"
	envString:=fmt.Sprintf(`environmentScope: { environments: %s }`,envArrayString)

	if len(environments)==0 || (len(environments)==1 && environments[0]==""){
		envString=""
	}
	frequencyString:=""
	if notification_frequency!=""{
		frequencyString=fmt.Sprintf(`rateLimitIntervalDuration: "%s"`,notification_frequency)
	}
	query:=fmt.Sprintf(`mutation {
		updateNotificationRule(
			input: {
				category: %s
				ruleId: "%s"
				ruleName: "%s"
				eventConditions: {
					sensitiveDataConfigChangeCondition: {
						sensitiveDataConfigChangeTypes: %s
						sensitiveDataConfigTypes:" %s
					}
				}
				channelId: "%s"
				%s
				%s
			}
		) {
			ruleId
		}
	}`,category,ruleId,name,event_types,sensitive_data_config_types,channel_id,frequencyString,envString)
	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
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

func resourceNotificationRuleDataClassificationConfigDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(`mutation {
		deleteNotificationRule(input: {ruleId: "%s"}) {
		  success
		}
	  }`, id)
	_, err := executeQuery(query, meta)
	if err != nil {
		return err
	}
	log.Println(query)
	d.SetId("")
	return nil
}