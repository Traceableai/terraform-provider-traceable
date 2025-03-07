package data_classification

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"log"
)

func DataSourceDatTypeId() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceDatTypeIdRead,

		Schema: map[string]*schema.Schema{
			"data_type_name": {
				Type:        schema.TypeString,
				Description: "Name of the data type",
				Required:    true,
			},
			"data_type_id": {
				Type:        schema.TypeString,
				Description: "Id of the data type",
				Computed:    true,
			},
		},
	}
}

func DataSourceDatTypeIdRead(rData *schema.ResourceData, meta interface{}) error {
	dataTypeName := rData.Get("data_type_name").(string)
	responseStr, err := common.CallExecuteQuery(DATA_TYPE_READ_QUERY, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response: %s", err)
	}
	log.Printf("this is the gql response %s", response)
	ruleDetails := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "dataTypes", dataTypeName)
	if len(ruleDetails) == 0 {
		return fmt.Errorf("no data types found with name %s", dataTypeName)
	}
	dataTypeId := ruleDetails["id"].(string)
	log.Printf("data type found with name %s", dataTypeId)
	rData.Set("data_type_id", dataTypeId)
	rData.SetId(dataTypeId)
	return nil
}
