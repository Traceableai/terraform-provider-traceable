package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNotificationRuleLoggedThreatActivity() *schema.Resource {
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
			},
			"severities": {
				Type:        schema.TypeSet,
				Description: "Severites of threat events you want to notify",
				Required:    true,
			},
			"impact": {
				Type:        schema.TypeSet,
				Description: "Impact of threat events you want to notify",
				Required:    true,
			},
			"confidence": {
				Type:        schema.TypeSet,
				Description: "Confidence of threat events you want to notify",
				Required:    true,
			},
			"notification_frequency": {
				Type:        schema.TypeString,
				Description: "No more than one notification every configured notification_frequency",
				Optional:    true,
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

	severitiesString:="["
	for _,v := range severities{
		severitiesString+=v.(string)
		severitiesString+=","
	}
	severitiesString=severitiesString[:len(severitiesString)-1]
	severitiesString+="]"
	if len(severities)==4{
		severitiesString=""
	}

	impactString:="["
	for _,v := range impact{
		impactString+=v.(string)
		impactString+=","
	}
	impactString=impactString[:len(impactString)-1]
	impactString+="]"
	if len(impact)==3{
		impactString=""
	}

	confidenceString:="["
	for _,v := range confidence{
		confidenceString+=v.(string)
		confidenceString+=","
	}
	confidenceString=confidenceString[:len(confidenceString)-1]
	confidenceString+="]"
	if len(confidence)==3{
		confidenceString=""
	}

	threatTypesString:="["
	for _,v := range threat_types{
		str:=""
		if isCustomThreatEvent(v.(string))==true{
			str=fmt.Sprintf(`{
				detectedThreatActivityConditionType: CUSTOM
				customDetectionCondition: { customDetectionType: %s }
			}`,v)
		}else if ok,val:=isPreDefinedThreatEvent(v.(string));ok{
			str=fmt.Sprintf(`{
				detectedThreatActivityConditionType: PRE_DEFINED
				'
				': { anomalyRuleId: %s }
			}`,val)
		}
		if str==""{
			return fmt.Errorf("Threat type %s not expected",v)
		}
		threatTypesString+=str
	}
	threatTypesString+="]"

	frequencyString:=""
	if notification_frequency!=""{
		frequencyString=fmt.Sprintf(`rateLimitIntervalDuration: "%s"`,notification_frequency)
	}
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
	query:=fmt.Sprintf(`mutation {
		createNotificationRule(
			input: {
				category: DETECTED_SECURITY_EVENT
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
	}`,name,threatTypesString,severitiesString,impactString,confidenceString,channel_id,frequencyString,envString)
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
func resourceNotificationRuleLoggedThreatActivityRead(d *schema.ResourceData, meta interface{}) error {
	id:=d.Id()
	readQuery:=`{
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
			  detectedSecurityEventCondition {
				detectedThreatActivityConditions {
				  detectedThreatActivityConditionType
				  customDetectionCondition {
					customDetectionType
				  }
				  preDefinedDetectionCondition {
					anomalyRuleId
				  }
				}
				severities
				impactLevels
				confidenceLevels
			  }
			}
			rateLimitIntervalDuration
		  }
		}
	  }`
	var response map[string]interface{}
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		_=fmt.Errorf("Error:%s",err)
	}
	log.Printf("This is the graphql query %s", readQuery)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		_=fmt.Errorf("Error:%s", err)
	}
	ruleDetails:=getRuleDetailsFromRulesListUsingIdName(response,"notificationRules" ,id,"ruleId","ruleName")
	d.Set("name",ruleDetails["ruleName"])
	d.Set("channel_id",ruleDetails["channelId"])
	envs:=ruleDetails["environmentScope"].(map[string]interface{})["environments"]
	d.Set("environments",schema.NewSet(schema.HashString,envs.([]interface{})))

	return nil
}

func resourceNotificationRuleLoggedThreatActivityUpdate(d *schema.ResourceData, meta interface{}) error {
	ruleId:=d.Id()
	name := d.Get("name").(string)
	environments := d.Get("environments").(*schema.Set).List()
	channel_id := d.Get("channel_id").(string)
	threat_types := d.Get("threat_types").(*schema.Set).List()
	severities := d.Get("severities").(*schema.Set).List()
	impact := d.Get("impact").(*schema.Set).List()
	confidence := d.Get("confidence").(*schema.Set).List()
	notification_frequency := d.Get("notification_frequency").(string)

	severitiesString:="["
	for _,v := range severities{
		severitiesString+=v.(string)
		severitiesString+=","
	}
	severitiesString=severitiesString[:len(severitiesString)-1]
	severitiesString+="]"
	if len(severities)==4{
		severitiesString=""
	}

	impactString:="["
	for _,v := range impact{
		impactString+=v.(string)
		impactString+=","
	}
	impactString=impactString[:len(impactString)-1]
	impactString+="]"
	if len(impact)==3{
		impactString=""
	}

	confidenceString:="["
	for _,v := range confidence{
		confidenceString+=v.(string)
		confidenceString+=","
	}
	confidenceString=confidenceString[:len(confidenceString)-1]
	confidenceString+="]"
	if len(confidence)==3{
		confidenceString=""
	}

	threatTypesString:="["
	for _,v := range threat_types{
		str:=""
		if isCustomThreatEvent(v.(string))==true{
			str=fmt.Sprintf(`{
				detectedThreatActivityConditionType: CUSTOM
				customDetectionCondition: { customDetectionType: %s }
			}`,v)
		}else if ok,val:=isPreDefinedThreatEvent(v.(string));ok{
			str=fmt.Sprintf(`{
				detectedThreatActivityConditionType: PRE_DEFINED
				preDefinedDetectionCondition: { anomalyRuleId: %s }
			}`,val)
		}
		if str==""{
			return fmt.Errorf("Threat type %s not expected",v)
		}
		threatTypesString+=str
	}
	threatTypesString+="]"

	frequencyString:=""
	if notification_frequency!=""{
		frequencyString=fmt.Sprintf(`rateLimitIntervalDuration: "%s"`,notification_frequency)
	}
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
	query:=fmt.Sprintf(`mutation {
		updateNotificationRule(
			input: {
				ruleId: "%s"
				category: DETECTED_SECURITY_EVENT
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
	}`,ruleId,name,threatTypesString,severitiesString,impactString,confidenceString,channel_id,frequencyString,envString)
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

func resourceNotificationRuleLoggedThreatActivityDelete(d *schema.ResourceData, meta interface{}) error {
	
	return nil
}