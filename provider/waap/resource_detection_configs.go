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
							Description: "Rule id for the detection config",
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
										Description: "Flag to enable/disable the detection config",
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
										Description: "Sub rule Id",
										Required:    true,
									},
									"sub_rule_action": {
										Type:        schema.TypeString,
										Description: "Action for sub rules",
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

	waapConfig := d.Get("waap_config").([]interface{})
	firstWaapConfig := waapConfig[0].(map[string]interface{})
	subRuleConfig := firstWaapConfig["subrule_config"].([]interface{})
	subRuleId := subRuleConfig[0].(map[string]interface{})["sub_rule_id"].(string)

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
	fetchedWaapConfig := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "anomalyDetectionRuleConfigs", crs_id, "ruleId", "ruleName")

	disabled := fetchedWaapConfig["configStatus"].(map[string]interface{})["disabled"].(bool)
	ruleId := notification.FindThreatByCrsId(crs_id)
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
	ruleConfigArr := []map[string]interface{}{}
	ruleConfigObj := map[string]interface{}{
		"disabled": disabled,
	}
	ruleConfigArr = append(ruleConfigArr, ruleConfigObj)

	subRuleConfigArr := []map[string]interface{}{}
	subRuleConfigObj := map[string]interface{}{
		"sub_rule_id":     subRuleId,
		"sub_rule_action": anomalySubRuleAction,
	}
	subRuleConfigArr = append(subRuleConfigArr, subRuleConfigObj)
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

	ruleConfig := firstWaapConfig["rule_config"].([]interface{})
	disabled := ruleConfig[0].(map[string]interface{})["disabled"].(bool)

	subRuleConfig := firstWaapConfig["subrule_config"].([]interface{})
	subRuleId := subRuleConfig[0].(map[string]interface{})["sub_rule_id"].(string)
	subRuleAction := subRuleConfig[0].(map[string]interface{})["sub_rule_action"].(string)

	configScope := GetConfigScope(environment)
	configType, err := GetConfigType(ruleId)
	if err != nil {
		return err
	}
	ruleConfigs, err := GetRuleConfig(ruleId, configType, disabled, subRuleId, subRuleAction)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(UPDATE_WAAP_CONFIG, configScope, ruleConfigs)

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
