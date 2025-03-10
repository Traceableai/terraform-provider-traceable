package provider

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// func dataSourceServiceId() *schema.Resource {
// 	return &schema.Resource{
// 		Read: dataSourceServiceIdRead,

// 		Schema: map[string]*schema.Schema{
// 			"service_name": {
// 				Type:        schema.TypeString,
// 				Description: "name of the service",
// 				Required:    true,
// 			},
// 			"enviroment_name": {
// 				Type:        schema.TypeString,
// 				Description: "Environement of the service name",
// 				Required:    true,
// 			},
// 			"service_id": {
// 				Type:        schema.TypeString,
// 				Description: "Id of the service",
// 				Computed:    true,
// 			},
// 		},
// 	}
// }

// func dataSourceServiceIdRead(d *schema.ResourceData, meta interface{}) error {
// 	service_name := d.Get("service_name").(string)
// 	enviroment_name := d.Get("enviroment_name").(string)

// 	currentTime := time.Now().UTC()
// 	endTime := currentTime.Format("2006-01-02T15:04:05.000Z")
// 	lastWeekTime := currentTime.AddDate(0, 0, -7)
// 	stTime := lastWeekTime.Format("2006-01-02T15:04:05.000Z")

// 	query := fmt.Sprintf(`{
// 		entities(
// 		  scope: "SERVICE"
// 		  limit: 1
// 		  between: {startTime: "%s", endTime: "%s"}
// 		  includeInactive: true
// 		  filterBy: [
// 					  {
// 						  keyExpression: { key: "name" }
// 						  operator: EQUALS
// 						  value: "%s"
// 						  type: ATTRIBUTE
// 					  }
// 					  {
// 						  keyExpression: { key: "environment" }
// 						  operator: IN
// 						  value: ["%s"]
// 						  type: ATTRIBUTE
// 					  }
// 				  ]
// 		) {
// 		  results {
// 			entityId: id
// 			id: attribute(expression: {key: "id"})
// 			name: attribute(expression: {key: "name"})
// 		  }
// 		}
// 	  }
// 	  `, stTime, endTime, service_name, enviroment_name)

// 	responseStr, err := ExecuteQuery(query, meta)
// 	if err != nil {
// 		return fmt.Errorf("error while executing GraphQL query: %s", err)
// 	}

// 	var response map[string]interface{}
// 	err = json.Unmarshal([]byte(responseStr), &response)
// 	if err != nil {
// 		return fmt.Errorf("error parsing JSON response: %s", err)
// 	}
// 	log.Printf("this is the gql response %s", response)
// 	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "entities", service_name)
// 	if len(ruleDetails) == 0 {
// 		return fmt.Errorf("no services found with name %s", service_name)
// 	}
// 	service_id := ruleDetails["id"].(string)
// 	log.Printf("Service found with name %s", service_id)
// 	d.Set("service_id", service_id)
// 	d.SetId(service_id)
// 	return nil
// }
