package data_classification

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"log"
)

func DataSourceDatSetId() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceDatSetIdRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the data set",
				Required:    true,
			},
			"data_set_id": {
				Type:        schema.TypeString,
				Description: "Id of the data set",
				Computed:    true,
			},
		},
	}
}

func DataSourceDatSetIdRead(rData *schema.ResourceData, meta interface{}) error {
	name := rData.Get("name").(string)
	responseStr, err := common.CallExecuteQuery(DATA_SET_READ_QUERY, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response: %s", err)
	}
	log.Printf("this is the gql response %s", response)
	ruleDetails := common.CallGetRuleDetailsFromRulesListUsingIdName(response, "dataSets", name)
	if len(ruleDetails) == 0 {
		return fmt.Errorf("no data sets found with name %s", name)
	}
	data_set_id := ruleDetails["id"].(string)
	log.Printf("data set found with name %s", data_set_id)
	rData.Set("data_set_id", data_set_id)
	rData.SetId(data_set_id)
	return nil
}
