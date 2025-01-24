package data_classification

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"log"
)

func ResourceDataClassificationOverrides() *schema.Resource {
	return &schema.Resource{
		Create: ResourceDataClassificationOverridesCreate,
		Read:   ResourceDataClassificationOverridesRead,
		Update: ResourceDataClassificationOverridesUpdate,
		Delete: ResourceDataClassificationOverridesDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the data set classification override",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the rule",
				Optional:    true,
			},
			"data_suppression_override": {
				Type:        schema.TypeString,
				Description: "Data suppression for the override (RAW/REDACT/OBFUSCATE)",
				Required:    true,
			},
			"environments": {
				Type:        schema.TypeList,
				Description: "Environments for the overrides",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"span_filter": {
				Type:        schema.TypeList,
				Description: "Span filters for the overrides",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_patterns": {
							Type:        schema.TypeList,
							Description: "Key operator and value",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "Key operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "Value for key",
										Required:    true,
									},
								},
							},
						},
						"value_patterns": {
							Type:        schema.TypeList,
							Description: "Value operator and value",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "Value operator",
										Required:    true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
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

func ResourceDataClassificationOverridesCreate(rData *schema.ResourceData, meta interface{}) error {
	name := rData.Get("name").(string)
	description := rData.Get("description").(string)
	dataSuppressionOverride := rData.Get("data_suppression_override").(string)
	environments := rData.Get("environments").(*schema.Set).List()
	spanFilter := rData.Get("span_filter").(*schema.Set).List()
	createQuery := GetOverridesCreateQuery("", name, description, dataSuppressionOverride, environments, spanFilter)
	log.Printf("This is the graphql query %s", createQuery)
	responseStr, err := common.CallExecuteQuery(createQuery, meta)
	if err != nil {
		return fmt.Errorf("error occured :%s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createDataClassificationOverride")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	rData.SetId(id)
	return nil
}

func ResourceDataClassificationOverridesRead(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(OVERRIDES_READ_QUERY, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "dataClassificationOverrideRules", id)
	if len(ruleData) == 0 {
		rData.SetId("")
		return nil
	}
	rData.Set("name", ruleData["name"].(string))
	rData.Set("description", ruleData["description"].(string))
	rData.Set("data_suppression_override", ruleData["behaviorOverride"].(map[string]interface{})["dataSuppressionOverride"].(string))
	if environmentScope, ok := ruleData["environmentScope"].(map[string]interface{}); ok {
		environmentIds := environmentScope["environmentIds"].([]interface{})
		rData.Set("environments", environmentIds)
	} else {
		rData.Set("environments", []interface{}{})
	}
	spanFilter := ruleData["spanFilter"].(map[string]interface{})
	keyValueFilter := spanFilter["keyValueFilter"].(map[string]interface{})
	keyPattern := keyValueFilter["keyPattern"].(map[string]interface{})
	keyPatternValue := keyPattern["value"].(string)
	keyPatternOp := keyPattern["operator"].(string)
	keyPatternObj := map[string]interface{}{
		"operator": keyPatternOp,
		"value":    keyPatternValue,
	}
	valuePatternObj := map[string]interface{}{}
	if valuePattern, ok := keyValueFilter["valuePattern"].(map[string]interface{}); ok {
		valuePatternValue := valuePattern["value"].(string)
		valuePatternOperator := valuePattern["operator"].(string)
		valuePatternObj = map[string]interface{}{
			"operator": valuePatternOperator,
			"value":    valuePatternValue,
		}
	}
	spanFilterObj := map[string]interface{}{
		"key_patterns":   keyPatternObj,
		"value_patterns": valuePatternObj,
	}
	rData.Set("span_filter", spanFilterObj)
	return nil
}

func ResourceDataClassificationOverridesUpdate(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	name := rData.Get("name").(string)
	description := rData.Get("description").(string)
	dataSuppressionOverride := rData.Get("data_suppression_override").(string)
	environments := rData.Get("environments").(*schema.Set).List()
	spanFilter := rData.Get("span_filter").(*schema.Set).List()
	updateQuery := GetOverridesCreateQuery(id, name, description, dataSuppressionOverride, environments, spanFilter)
	log.Printf("This is the graphql query %s", updateQuery)
	responseStr, err := common.CallExecuteQuery(updateQuery, meta)
	if err != nil {
		return fmt.Errorf("error occured :%s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateDataClassificationOverride")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	rData.SetId(updatedId)
	return nil
}

func ResourceDataClassificationOverridesDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_OVERRIDES_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId("")
	return nil
}
