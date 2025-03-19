package data_classification

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"log"
)

func ResourceDataSetsRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourceDataSetsRuleCreate,
		Read:   ResourceDataSetsRuleRead,
		Update: ResourceDataSetsRuleUpdate,
		Delete: ResourceDataSetsRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the data set",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the data set",
				Optional:    true,
			},
			"icon_type": {
				Type:        schema.TypeString,
				Description: "Icon for the data set",
				Optional:    true,
			},
		},
	}
}

func ResourceDataSetsRuleCreate(rData *schema.ResourceData, meta interface{}) error {
	name := rData.Get("name").(string)
	if name == "" {
		return fmt.Errorf("non empty string required")
	}
	description := rData.Get("description").(string)
	iconType := rData.Get("icon_type").(string)
	createQuery := GetDatSetQuery("", name, description, iconType)
	log.Printf("This is the graphql query %s", createQuery)
	responseStr, err := common.CallExecuteQuery(createQuery, meta)
	if err != nil {
		return fmt.Errorf("error occured :%s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	id, err := common.GetIdFromResponse(responseStr, "createDataSet")
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	rData.SetId(id)
	return nil
}

func ResourceDataSetsRuleRead(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	log.Println("Id from read ", id)
	responseStr, err := common.CallExecuteQuery(DATA_SET_READ_QUERY, meta)
	if err != nil {
		return err
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	ruleData := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "dataSets", id)
	if len(ruleData) == 0 {
		rData.SetId("")
		return nil
	}
	
	rData.Set("name", ruleData["name"].(string))
	rData.Set("description", ruleData["description"].(string))
	rData.Set("iconType", ruleData["iconType"].(string))
	return nil
}

func ResourceDataSetsRuleUpdate(rData *schema.ResourceData, meta interface{}) error {
	id := rData.Id()
	name := rData.Get("name").(string)
	description := rData.Get("description").(string)
	iconType := rData.Get("icon_type").(string)
	updateQuery := GetDatSetQuery(id, name, description, iconType)
	log.Printf("This is the graphql query %s", updateQuery)
	responseStr, err := common.CallExecuteQuery(updateQuery, meta)
	if err != nil {
		return fmt.Errorf("error occured :%s", err)
	}
	log.Printf("This is the graphql response %s", responseStr)
	updatedId, err := common.GetIdFromResponse(responseStr, "updateDataSet")
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	rData.SetId(updatedId)
	return nil
}

func ResourceDataSetsRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	query := fmt.Sprintf(DELETE_DATA_SET_QUERY, id)
	_, err := common.CallExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	d.SetId("")
	return nil
}
