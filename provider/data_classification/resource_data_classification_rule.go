package data_classification

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_suppression": {
				Type:        schema.TypeString,
				Description: "Data supression for the data type",
				Required:    true,
			},
			"sensitivity": {
				Type:        schema.TypeString,
				Description: "Sensitivity of the data type",
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
							Description: "Match or exclude",
							Required:    true,
						},
						"locations": {
							Type:        schema.TypeList,
							Description: "where to look for the data type",
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
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "key operator",
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
							Description: "key operator and value",
							Optional:    true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:        schema.TypeString,
										Description: "value operator",
										Required:    true,
									},
									"value": {
										Type:        schema.TypeString,
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

func ResourceDataClassificationCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	sensitivity := d.Get("sensitivity").(string)
	dataSuppression := d.Get("data_suppression").(string)
	enabled := d.Get("enabled").(bool)
	dataSets := d.Get("data_sets").(*schema.Set).List()
	scopedPatterns := d.Get("scoped_patterns").(*schema.Set).List()

	scopedPatternQuery := ReturnScopedPatternQuery(scopedPatterns)
	query := fmt.Sprintf(CREATE_QUERY, name, description, scopedPatternQuery, enabled, dataSuppression,sensitivity,common.InterfaceToStringSlice(dataSets))
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	id,err := common.GetIdFromResponse(responseStr,"createDataType")
	if err!=nil {
		return fmt.Errorf("error %s",err)
	}
	d.SetId(id)
	return nil
}

func ResourceDataClassificationRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(READ_QUERY, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "dataTypes", id)
	if len(ruleData) == 0 {
		d.SetId("")
		return nil
	}
	d.Set("name",ruleData["name"].(string))
	d.Set("description",ruleData["description"].(string))
	d.Set("sensitivity",ruleData["sensitivity"].(string))
	d.Set("data_suppression",ruleData["dataSuppression"].(string))
	d.Set("enabled",ruleData["enabled"].(bool))
	var dataSets []interface{}
	for _,dataSet := range ruleData["datasets"].([]interface{}) {
		dataSetData := dataSet.(map[string]interface{})
		dataSetId := dataSetData["id"].(string)
		dataSets = append(dataSets, dataSetId)
	}
	d.Set("data_sets",dataSets)
	scopedPatterns := []map[string]interface{}{}
	for _,scopedPattern := range ruleData["scopedPatterns"].([]interface{}){
		scopedPatternData := scopedPattern.(map[string]interface{})
		scopedPatternName := scopedPatternData["name"].(string)
		locations := scopedPatternData["locations"].([]interface{})
		matchType := scopedPatternData["matchType"].(string)
		urlMatchPatterns := scopedPatternData["urlMatchPatterns"].([]interface{})
		environmentIds := []interface{}{}
		if scope,ok := scopedPatternData["scope"].(map[string]interface{}); ok {
			environmentScope := scope["environmentScope"].(map[string]interface{})
			environmentIds = environmentScope["environmentIds"].([]interface{})
		}
		keyPattern := scopedPatternData["keyPattern"].(map[string]interface{})
		keyPatternValue := keyPattern["value"].(string)
		keyPatternOp := keyPattern["operator"].(string)
		keyPatternObj := map[string]interface{}{
			"operator": keyPatternOp,
			"value": keyPatternValue,
		}
		valuePatternObj := map[string]interface{}{}
		if valuePattern,ok := scopedPatternData["valuePattern"].(map[string]interface{}); ok{
			valuePatternValue := valuePattern["value"].(string)
			valuePatternOperator := valuePattern["operator"].(string)
			valuePatternObj = map[string]interface{}{
				"operator": valuePatternOperator,
				"value": valuePatternValue,
			}
		}
		scopedPatternObj := map[string]interface{}{
			"scoped_pattern_name" : scopedPatternName,
			"locations" : locations,
			"match_type" : matchType,
			"url_match_patterns" : urlMatchPatterns,
			"environments" : environmentIds,
			"key_patterns" : keyPatternObj,	
			"value_patterns": valuePatternObj,
		}
		scopedPatterns = append(scopedPatterns, scopedPatternObj)
	}
	d.Set("scoped_patterns",scopedPatterns)
	return nil
}

func ResourceDataClassificationUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	sensitivity := d.Get("sensitivity").(string)
	dataSuppression := d.Get("data_suppression").(string)
	enabled := d.Get("enabled").(bool)
	dataSets := d.Get("data_sets").(*schema.Set).List()
	scopedPatterns := d.Get("scoped_patterns").(*schema.Set).List()

	scopedPatternQuery := ReturnScopedPatternQuery(scopedPatterns)
	query := fmt.Sprintf(UPDATE_QUERY, name, description, scopedPatternQuery, enabled, dataSuppression,sensitivity,common.InterfaceToStringSlice(dataSets))
	responseStr, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	log.Printf("This is the graphql query %s", query)
	log.Printf("This is the graphql response %s", responseStr)
	id,err := common.GetIdFromResponse(responseStr,"updateDataType")
	if err!=nil {
		return fmt.Errorf("error %s",err)
	}
	d.SetId(id)
	return nil
}

func ResourceDataClassificationDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error %s", err)
	}
	d.SetId("")
	return nil
}
