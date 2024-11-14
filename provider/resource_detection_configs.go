package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDetectionConfigRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceDetectionConfigCreate,
		Read:   resourceDetectionConfigRead,
		Update: resourceDetectionConfigUpdate,
		Delete: resourceDetectionConfigDelete,

		Schema: map[string]*schema.Schema{
			"environment": {
				Type:        schema.TypeString,
				Description: "Environement of detection policy",
				Optional:    true,
			},
			"config_name": {
				Type:        schema.TypeString,
				Description: "Detection Config name",
				Required:    true,
			},
            "disabled": {
                Type:        schema.TypeString,
                Description: "Flag to enable/disable the detection config",
                Optional:    true,
            },
            "sub_rule_id": {
                Type:        schema.TypeString,
                Description: "Sub rule id to enable/disable sub rule blocking",
                Optional:    true,
            },
            "sub_rule_blocking_enabled": {
                Type:        schema.TypeString,
                Description: "Enable/Disable blocking for sub rules",
                Optional:    true,
            },
		},
	}
}

func resourceDetectionConfigCreate(d *schema.ResourceData, meta interface{}) error {
    return resourceDetectionConfigUpdate(d, meta)
}

func resourceDetectionConfigRead(d *schema.ResourceData, meta interface{}) error {
    crs_id:=d.Id()
    environment := d.Get("environment").(string)
    anomalyScope:=`anomalyScope: { scopeType: CUSTOMER }`
    if environment!=""{
        anomalyScope=fmt.Sprintf(`anomalyScope: {scopeType: ENVIRONMENT, environmentScope: {id: "%s"}}`,environment)
    }
    readQuery:=fmt.Sprintf(`{
                  anomalyDetectionRuleConfigs(
                    %s
                  ) {
                    count
                    total
                    results {
                      configStatus {
                        disabled
                        internal
                        __typename
                      }
                      hidden
                      eventFamily
                      configType
                      ruleCategory
                      ruleId
                      ruleName

                      subRuleConfigs {
                        blockingEnabled
                        configStatus {
                          disabled
                          internal
                          __typename
                        }
                        subRuleId
                        subRuleName


                        anomalyProtectionType
                        anomalyRuleSeverity
                        __typename
                      }
                      anomalyRuleSeverity
                      __typename
                    }
                    __typename
                  }
                }`,anomalyScope)
    var response map[string]interface{}
    responseStr, err := executeQuery(readQuery, meta)
    if err != nil {
        return fmt.Errorf("Error:%s", err)
    }
    log.Printf("This is the graphql query %s", readQuery)
    err = json.Unmarshal([]byte(responseStr), &response)
    if err != nil {
        return fmt.Errorf("Error:%s", err)
    }
    detection_config:=getRuleDetailsFromRulesListUsingIdName(response,"anomalyDetectionRuleConfigs" ,crs_id,"ruleId","ruleName")

    sub_rule_id := d.Get("sub_rule_id").(string)
    if sub_rule_id!=""{
        log.Println("going inside")
        sub_rules_array:=detection_config["subRuleConfigs"].([]interface{})
        for _, sub_rule := range sub_rules_array {
            sub_rule_config:=sub_rule.(map[string]interface{})
            if sub_rule_config["subRuleId"]==sub_rule_id{
                log.Printf("found a match %s",sub_rule_config)
                is_blocking_enabled:=sub_rule_config["blockingEnabled"].(bool)
                d.Set("sub_rule_blocking_enabled",is_blocking_enabled)
                break
            }
        }
    }
    is_rule_disabled:=detection_config["configStatus"].(map[string]interface{})["disabled"].(bool)
    d.Set("disabled",is_rule_disabled)
	return nil
}

func resourceDetectionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
    environment := d.Get("environment").(string)
    config_name := d.Get("config_name").(string)
    disabled := d.Get("disabled").(string)
    sub_rule_id := d.Get("sub_rule_id").(string)
    sub_rule_blocking_enabled := d.Get("sub_rule_blocking_enabled").(string)


    if (sub_rule_blocking_enabled=="" && disabled=="") || (disabled!="" && sub_rule_blocking_enabled!=""){
        return fmt.Errorf("Required one of `sub_rule_blocking_enable` or `disabled`")
    }

    configScope:="anomalyScope: { scopeType: CUSTOMER }"
    if environment!=""{
        configScope=fmt.Sprintf(`anomalyScope: {
                                                 scopeType: ENVIRONMENT
                                                 environmentScope: { id: "%s" }
                                             }`,environment)
    }
    _,rule_crs_id:=isPreDefinedThreatEvent(config_name)
    configType:="API_DEFINITION_METADATA"
    if strings.Contains(rule_crs_id,"crs"){
        configType="MODSECURITY"
    }

    ruleConfigs:=fmt.Sprintf(`ruleConfig: {
                      ruleId: "%s"
                      configType: %s
                      configStatus: { disabled: %s }
                  }`,rule_crs_id,configType,disabled)
    if sub_rule_id!=""{
        ruleConfigs=fmt.Sprintf(`ruleConfig: {
                      ruleId: "%s"
                      configType: %s
                      subRuleConfigs: [{ subRuleId: "%s", blockingEnabled: %s }]
                  }`,rule_crs_id,configType,sub_rule_id,sub_rule_blocking_enabled)
    }


    query:=fmt.Sprintf(`mutation {
                            updateAnomalyRuleConfig(
                                update: {
                                    %s
                                    %s
                                }
                            ) {
                                ruleId
                                configStatus { disabled }
                                subRuleConfigs{
                                            subRuleId
                                            blockingEnabled
                                        }
                                __typename
                            }
                        }`,configScope,ruleConfigs)

    var response map[string]interface{}
    responseStr, err := executeQuery(query, meta)
    if strings.Contains(responseStr,"DataFetchingException"){
        _ = json.Unmarshal([]byte(responseStr), &response)
        return fmt.Errorf("Error: %s", response)
    }
    if err != nil {
        return fmt.Errorf("Error: %s", err)
    }
    log.Printf("This is the graphql query %s", query)
    log.Printf("This is the graphql response %s", responseStr)
    err = json.Unmarshal([]byte(responseStr), &response)
    if err != nil {
        return fmt.Errorf("Error:%s", err)
    }
    rules := response["data"].(map[string]interface{})["updateAnomalyRuleConfig"].(map[string]interface{})
    log.Println(rules)
    d.SetId(rule_crs_id)
    return nil
}

func resourceDetectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}