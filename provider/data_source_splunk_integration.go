package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSplunkIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSplunkIntegrationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the integration",
				Required:    true,
			},
			"splunk_id": {
				Type:        schema.TypeString,
				Description: "Id of splunk integration",
				Computed:    true,
			},
		},
	}
}

func dataSourceSplunkIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	query := `{
		splunkIntegrations {
		  results {
			id
			name
			description
			httpEventCollectorUrl
		  }
		}
	  }`

	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response: %s", err)
	}
	log.Printf("this is the gql response %s", response)
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "splunkIntegrations", name)
	if len(ruleDetails) == 0 {
		return fmt.Errorf("No rules found with name %s", name)
	}
	splunkId := ruleDetails["id"].(string)
	log.Printf("Rule found with name %s", splunkId)
	d.Set("splunk_id", splunkId)
	d.SetId(splunkId)
	return nil
}
