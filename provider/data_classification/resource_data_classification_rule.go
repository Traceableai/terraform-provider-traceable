package data_classification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceDataClassification() *schema.Resource {
	return &schema.Resource{
		Create: ResourceDataClassificationCreate,
		Read:   ResourceDataClassificationRead,
		Update: ResourceDataClassificationUpdate,
		Delete: ResourceDataClassificationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name of the data type",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the rule",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable the data type with this flag",
				Required:    true,
			},
			"data_sets": {
				Type:        schema.TypeList,
				Description: "List of Data sets  Id for this data type",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_suppression": {
				Type:        schema.TypeString,
				Description: "Data supression for the data type (RAW/REDACT/OBFUSCATE)",
				Required:    true,
			},
			"sensitivity": {
				Type:        schema.TypeString,
				Description: "Sensitivity of the data type (LOW/HIGH/MEDIUM/CRITICAL)",
				Required:    true,
			},
			"scoped_patterns": {
				Type:        schema.TypeList,
				Description: "Scope of the data type of rule",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scoped_pattern_name": {
							Type:        schema.TypeString,
							Description: "name of scoped pattern",
							Optional:    true,
						},
						"environments": {
							Type:        schema.TypeList,
							Description: "Environment where it will be applied",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"match_type": {
							Type:        schema.TypeString,
							Description: "Match or exclude (IGNORE/MATCH)",
							Required:    true,
						},
						"locations": {
							Type:        schema.TypeList,
							Description: "where to look for the data type (ANY/QUERY/REQUEST_HEADER/RESPONSE_HEADER/REQUEST_BODY/RESPONSE_BODY/REQUEST_COOKIE)",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"url_match_patterns": {
							Type:        schema.TypeList,
							Description: "url matching regex",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"key_patterns": {
							Type:        schema.TypeList,
							Description: "key operator and value",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "key operator (MATCHES_REGEX/EQUALS)",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "value for key",
										Required:    true,
									},
								},
							},
						},
						"value_patterns": {
							Type:        schema.TypeList,
							Description: "key operator and value ",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "value operator (MATCHES_REGEX/EQUALS)",
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

func ResourceDataClassificationCreate(rData *schema.ResourceData, meta interface{}) error {
	name := rData.Get("name").(string)
	description := rData.Get("description").(string)
	sensitivity := rData.Get("sensitivity").(string)
	dataSuppression := rData.Get("data_suppression").(string)
	enabled := rData.Get("enabled").(bool)
	dataSets := rData.Get("data_sets").([]interface{})
	scopedPatterns := rData.Get("scoped_patterns").([]interface{})

	scopedPatternQuery := GetScopedPatternQuery(scopedPatterns)
	query := fmt.Sprintf(CREATE_QUERY, name, description, scopedPatternQuery, enabled, dataSuppression, sensitivity, common.InterfaceToStringSlice(dataSets))
	log.Printf("This is the graphql query %s", query)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createDataType")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	rData.SetId(id)
	return nil
}

func ResourceDataClassificationRead(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(DATA_TYPE_READ_QUERY, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "dataTypes", id)
	if len(ruleData) == 0 {
		rData.SetId("")
		return nil
	}
	rData.Set("name", ruleData["name"].(string))
	rData.Set("description", ruleData["description"].(string))
	rData.Set("sensitivity", ruleData["sensitivity"].(string))
	rData.Set("data_suppression", ruleData["dataSuppression"].(string))
	rData.Set("enabled", ruleData["enabled"].(bool))
	var dataSets []interface{}
	for _, dataSet := range ruleData["datasets"].([]interface{}) {
		dataSetData := dataSet.(map[string]interface{})
		dataSetId := dataSetData["id"].(string)
		dataSets = append(dataSets, dataSetId)
	}
	rData.Set("data_sets", dataSets)
	scopedPatterns := []map[string]interface{}{}
	for _, scopedPattern := range ruleData["scopedPatterns"].([]interface{}) {
		scopedPatternData := scopedPattern.(map[string]interface{})
		scopedPatternName := scopedPatternData["name"].(string)
		locations := scopedPatternData["locations"].([]interface{})
		matchType := scopedPatternData["matchType"].(string)
		urlMatchPatterns := scopedPatternData["urlMatchPatterns"].([]interface{})
		environmentIds := []interface{}{}
		if scope, ok := scopedPatternData["scope"].(map[string]interface{}); ok {
			environmentScope := scope["environmentScope"].(map[string]interface{})
			environmentIds = environmentScope["environmentIds"].([]interface{})
		}
		keyPattern := scopedPatternData["keyPattern"].(map[string]interface{})
		keyPatternValue := keyPattern["value"].(string)
		keyPatternOp := keyPattern["operator"].(string)
		keyPatternObj := map[string]interface{}{
			"operator": keyPatternOp,
			"value":    keyPatternValue,
		}
		valuePatternObj := map[string]interface{}{}
		if valuePattern, ok := scopedPatternData["valuePattern"].(map[string]interface{}); ok {
			valuePatternValue := valuePattern["value"].(string)
			valuePatternOperator := valuePattern["operator"].(string)
			valuePatternObj = map[string]interface{}{
				"operator": valuePatternOperator,
				"value":    valuePatternValue,
			}
		}
		scopedPatternObj := map[string]interface{}{
			"scoped_pattern_name": scopedPatternName,
			"locations":           locations,
			"match_type":          matchType,
			"url_match_patterns":  urlMatchPatterns,
			"environments":        environmentIds,
			"key_patterns":        keyPatternObj,
			"value_patterns":      valuePatternObj,
		}
		scopedPatterns = append(scopedPatterns, scopedPatternObj)
	}
	rData.Set("scoped_patterns", scopedPatterns)
	return nil
}

func ResourceDataClassificationUpdate(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	name := rData.Get("name").(string)
	description := rData.Get("description").(string)
	sensitivity := rData.Get("sensitivity").(string)
	dataSuppression := rData.Get("data_suppression").(string)
	enabled := rData.Get("enabled").(bool)
	dataSets := rData.Get("data_sets").([]interface{})
	scopedPatterns := rData.Get("scoped_patterns").([]interface{})

	scopedPatternQuery := GetScopedPatternQuery(scopedPatterns)
	query := fmt.Sprintf(UPDATE_QUERY, id, name, description, scopedPatternQuery, enabled, dataSuppression, sensitivity, common.InterfaceToStringSlice(dataSets))
	log.Printf("This is the graphql query %s", query)
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateDataType")
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	rData.SetId(updatedId)
	return nil
}

func ResourceDataClassificationDelete(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	query := fmt.Sprintf(DELETE_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	rData.SetId("")
	return nil
}
