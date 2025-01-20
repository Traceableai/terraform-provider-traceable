package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEndpointId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEndpointIdRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the endpoint",
				Required:    true,
			},
			"service_name": {
				Type:        schema.TypeString,
				Description: "service of the endpoint",
				Required:    true,
			},
			"enviroment_name": {
				Type:        schema.TypeString,
				Description: "environment of the endpoint",
				Required:    true,
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Description: "Id of the endpoint",
				Computed:    true,
			},
		},
	}
}

func dataSourceEndpointIdRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	service_name := d.Get("service_name").(string)
	enviroment_name := d.Get("enviroment_name").(string)

	currentTime := time.Now().UTC()
	endTime := currentTime.Format("2006-01-02T15:04:05.000Z")
	lastWeekTime := currentTime.AddDate(0, 0, -7)
	stTime := lastWeekTime.Format("2006-01-02T15:04:05.000Z")

	query := fmt.Sprintf(`{
		entities(
			scope: "API"
			limit: 1
			between: {
				startTime: "%s"
				endTime: "%s"
			}
			filterBy: [
				{
					keyExpression: { key: "name" }
					operator: EQUALS
					value: "%s"
					type: ATTRIBUTE
				}
				{
					keyExpression: { key: "environment" }
					operator: IN
					value: ["%s"]
					type: ATTRIBUTE
				}
				{
					keyExpression: { key: "serviceName" }
					operator: EQUALS
					value: "%s"
					type: ATTRIBUTE
				}
			]
			includeInactive: true
		) {
			results {
				entityId: id
				id: attribute(expression: { key: "id" })
				name: attribute(expression: { key: "name" })
			}
		}
	}`, stTime, endTime, name, enviroment_name, service_name)

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
	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "entities", name)
	if len(ruleDetails) == 0 {
		return fmt.Errorf("no endpoints found with name %s", name)
	}
	endpoint_id := ruleDetails["id"].(string)
	log.Printf("endpoint found with name %s", endpoint_id)
	d.Set("endpoint_id", endpoint_id)
	d.SetId(endpoint_id)
	return nil
}
