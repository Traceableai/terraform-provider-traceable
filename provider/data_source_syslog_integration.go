package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSyslogIntegration() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceSyslogIntegrationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the integration",
				Required:    true,
			},
			"syslog_id": {
				Type:        schema.TypeString,
				Description: "Id of syslog integration",
				Computed:    true,
			},
		},
	}
}

func DataSourceSyslogIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	query := `{
		syslogServerIntegrations {
		  results {
			id
			name
			description 
		  }
		}
	  }
	  `

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
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "syslogServerIntegrations", name)
	if len(ruleDetails) == 0 {
		return fmt.Errorf("no rules found with name %s", name)
	}
	syslogId := ruleDetails["id"].(string)
	log.Printf("Rule found with name %s", syslogId)
	d.Set("syslog_id", syslogId)
	d.SetId(syslogId)
	return nil
}
