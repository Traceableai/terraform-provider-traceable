package waap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/notification"
	"log"
	"strings"
)

func ResourceDetectionConfigRule() *schema.Resource {
	return &schema.Resource{
		Create:        resourceDetectionConfigCreate,
		Read:          resourceDetectionConfigRead,
		Update:        resourceDetectionConfigUpdate,
		Delete:        resourceDetectionConfigDelete,
		CustomizeDiff: validateSchema,
		Schema: map[string]*schema.Schema{
			"environment": {
				Type:        schema.TypeString,
				Description: "Environement of detection policy",
				Optional:    true,
			},
			"waap_config": {
				Type:        schema.TypeList,
				Description: "Detection Config settings for rule and subrules",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Description: "Rule id for the detection config (XSS/LFI/RFI/bola/userIdBola)",
							Required:    true,
						},
						"rule_config": {
							Type:        schema.TypeList,
							Description: "Detection Config settings for rule",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disabled": {
										Type:        schema.TypeBool,
										Description: "Flag to enable/disable the detection config(True/False)",
										Required:    true,
									},
								},
							},
						},
						"subrule_config": {
							Type:        schema.TypeList,
							Description: "Detection Config settings for subrules",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sub_rule_id": {
										Type:        schema.TypeString,
										Description: "Sub rule Id's (crs_941420/crs_941290/crs_941240)",
										Required:    true,
									},
									"sub_rule_action": {
										Type:        schema.TypeString,
										Description: "Action for sub rules (MONITOR/DISABLE/BLOCK)",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func validateSchema(ctx context.Context, rData *schema.ResourceDiff, meta interface{}) error {
	waapConfigs := rData.Get("waap_config").([]interface{})
	firstWaapConfig := waapConfigs[0].(map[string]interface{})
	ruleConfig := firstWaapConfig["rule_config"].([]interface{})
	subRuleConfig := firstWaapConfig["subrule_config"].([]interface{})
	if (len(ruleConfig) == 0 && len(subRuleConfig) == 0) || (len(ruleConfig) > 0 && len(subRuleConfig) > 0) {
		return fmt.Errorf("required atmost one rule_config or subrule_config")
	}
	return nil
}
func resourceDetectionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceDetectionConfigUpdate(d, meta)
}

func resourceDetectionConfigRead(d *schema.ResourceData, meta interface{}) error {
	crs_id := d.Id()
	environment := d.Get("environment").(string)
	anomalyScope := GetConfigScope(environment)
	readQuery := fmt.Sprintf(READ_QUERY, anomalyScope)
	var response map[string]interface{}
	log.Printf("This is the graphql query %s", readQuery)
	responseStr, err := common.CallExecuteQuery(readQuery, meta)
	if err != nil {
		return fmt.Errorf("error:%s", err)
	}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error:%s", err)
	}
	waapConfig := d.Get("waap_config").([]interface{})
	firstWaapConfig := waapConfig[0].(map[string]interface{})
	subRuleConfig := firstWaapConfig["subrule_config"].([]interface{})
	fetchedWaapConfig := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "anomalyDetectionRuleConfigs", crs_id, "ruleId", "ruleName")
	subRuleConfigArr := []map[string]interface{}{}
	if len(subRuleConfig) > 0 {
		subRuleId := subRuleConfig[0].(map[string]interface{})["sub_rule_id"].(string)
		subRuleConfigs := fetchedWaapConfig["subRuleConfigs"].([]interface{})
		anomalySubRuleAction := ""
		for _, subRuleConfig := range subRuleConfigs {
			subRuleConfigData := subRuleConfig.(map[string]interface{})
			fetchedSubRuleId := subRuleConfigData["subRuleId"].(string)
			if fetchedSubRuleId == subRuleId {
				anomalySubRuleAction = subRuleConfigData["anomalySubRuleAction"].(string)
				break
			}
		}
		subRuleConfigObj := map[string]interface{}{
			"sub_rule_id":     subRuleId,
			"sub_rule_action": anomalySubRuleAction,
		}
		subRuleConfigArr = append(subRuleConfigArr, subRuleConfigObj)
	}

	disabled := fetchedWaapConfig["configStatus"].(map[string]interface{})["disabled"].(bool)
	ruleId := notification.FindThreatByCrsId(crs_id)

	ruleConfigArr := []map[string]interface{}{}
	ruleConfigObj := map[string]interface{}{
		"disabled": disabled,
	}
	ruleConfigArr = append(ruleConfigArr, ruleConfigObj)

	waapConfigArr := []map[string]interface{}{}
	waapConfigStateObj := map[string]interface{}{
		"rule_id":        ruleId,
		"rule_config":    ruleConfigArr,
		"subrule_config": subRuleConfigArr,
	}
	waapConfigArr = append(waapConfigArr, waapConfigStateObj)
	d.Set("waap_config", waapConfigArr)
	return nil
}

func resourceDetectionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	environment := d.Get("environment").(string)
	waapConfigs := d.Get("waap_config").([]interface{})
	firstWaapConfig := waapConfigs[0].(map[string]interface{})
	ruleId := firstWaapConfig["rule_id"].(string)
	ruleConfigString := ""
	configType, err := GetConfigType(ruleId)
	ruleConfig := firstWaapConfig["rule_config"].([]interface{})
	if len(ruleConfig) > 0 {
		disabled := ruleConfig[0].(map[string]interface{})["disabled"].(bool)
		ruleConfigString, err = GetRuleConfig(ruleId, configType, disabled)
	}
	subRuleConfig := firstWaapConfig["subrule_config"].([]interface{})
	if len(subRuleConfig) > 0 {
		subRuleId := subRuleConfig[0].(map[string]interface{})["sub_rule_id"].(string)
		subRuleAction := subRuleConfig[0].(map[string]interface{})["sub_rule_action"].(string)
		ruleConfigString, err = GetSubRuleConfig(subRuleId, subRuleAction, ruleId, configType)
	}
	configScope := GetConfigScope(environment)

	if err != nil {
		return err
	}
	query := fmt.Sprintf(UPDATE_WAAP_CONFIG, configScope, ruleConfigString)
	fmt.Printf("this is the query %s",query)
	var response map[string]interface{}

	responseStr, err := common.CallExecuteQuery(query, meta)
	if strings.Contains(responseStr, "DataFetchingException") {
		_ = json.Unmarshal([]byte(responseStr), &response)
		return fmt.Errorf("error: %s", response)
	}
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error:%s", err)
	}
	if responseData, ok := response["data"].(map[string]interface{}); ok {
		rules := responseData["updateAnomalyRuleConfig"].(map[string]interface{})
		log.Println(rules)
		ruleCrsId, _ := GetRuleCrsId(ruleId)
		d.SetId(ruleCrsId)
	} else {
		return fmt.Errorf("error occurred while updating the state")
	}
	return nil
}

func resourceDetectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
