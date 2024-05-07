package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiNamingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiNamingRuleCreate,
		Read:   resourceApiNamingRuleRead,
		Update: resourceApiNamingRuleUpdate,
		Delete: resourceApiNamingRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the API naming rule",
				Required:    true,
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Required:    true,
			},
			"regexes": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of regex patterns for the rule",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"values": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Corresponding values for the regex patterns",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_names": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of service names to apply the rule",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environment_names": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of environment names to apply the rule",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceApiNamingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	regexes := d.Get("regexes").([]interface{})
	values := d.Get("values").([]interface{})
	serviceNames := d.Get("service_names").([]interface{})
	environmentNames := d.Get("environment_names").([]interface{})

	var spanFilters []string

	// Checking if service name list is (empty -> all services case)
	if len(serviceNames) > 0 && serviceNames[0] != "" {
		spanFilters = append(spanFilters, fmt.Sprintf(`
		{
			relationalSpanFilter: {
				field: SERVICE_NAME
				relationalOperator: IN
				value: %s
			}
		}`, jsonifyList(serviceNames)))
	}

	// Checling if env name list is (empty -> all env case)
	if len(environmentNames) > 0 && environmentNames[0] != "" {
		spanFilters = append(spanFilters, fmt.Sprintf(`
		{
			relationalSpanFilter: {
				field: ENVIRONMENT_NAME
				relationalOperator: IN
				value: %s
			}
		}`, jsonifyList(environmentNames)))
	}

	// spanFilter part of the query willbe empty in case of all env and all services
	spanFilterQueryPart := ""
	if len(spanFilters) > 0 {
		spanFilterQueryPart = fmt.Sprintf("spanFilters: [%s]", strings.Join(spanFilters, ","))
	} else {
		spanFilterQueryPart = "spanFilters: []" // No filters to apply
	}

	query := fmt.Sprintf(`
mutation {
  createApiNamingRule(
    input: {
      name: "%s"
      disabled: %t
      apiNamingRuleConfig: {
        apiNamingRuleConfigType: SEGMENT_MATCHING
        segmentMatchingBasedRuleConfig: {
          regexes: %s
          values: %s
        }
      }
      spanFilter: {
        logicalSpanFilter: {
          logicalOperator: AND
          %s
        }
      }
    }
  ) {
    id
  }
}`, name, disabled, jsonifyList(regexes), jsonifyList(values), spanFilterQueryPart)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Printf("Error while executing GraphQL query: %s", err)
		return err
	}

	log.Printf("GraphQL response: %s", responseStr)

	if response["data"] != nil && response["data"].(map[string]interface{})["createApiNamingRule"] != nil {
		id := response["data"].(map[string]interface{})["createApiNamingRule"].(map[string]interface{})["id"].(string)
		d.SetId(id)
		log.Printf(query)
	} else {
		return fmt.Errorf("could not create API naming rule, no ID returned")
	}

	return nil
}

// function to convert a list of strings to a GraphQL-compatible string list
func jsonifyList(list []interface{}) string {
	var strList []string
	for _, item := range list {
		strList = append(strList, fmt.Sprintf(`"%s"`, item))
	}
	return "[" + strings.Join(strList, ", ") + "]"
}

func resourceApiNamingRuleRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Get the resource ID set during creation

	// GraphQL query to read the API Naming Rule details
	query := `{
        apiNamingRules {
            results {
                id
                name
                disabled
                apiNamingRuleConfig {
                    apiNamingRuleConfigType
                    segmentMatchingBasedRuleConfig {
                        regexes
                        values
                    }
                }
                spanFilter {
                    logicalSpanFilter {
                        logicalOperator
                        spanFilters {
                            relationalSpanFilter {
                                field
                                relationalOperator
                                value
                            }
                        }
                    }
                }
            }
        }
    }`

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		log.Printf("Error while executing GraphQL query: %s", err)
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	results := response["data"].(map[string]interface{})["apiNamingRules"].(map[string]interface{})["results"].([]interface{})
	for _, item := range results {
		rule := item.(map[string]interface{})
		if rule["id"].(string) == id {
			d.Set("name", rule["name"].(string))
			d.Set("disabled", rule["disabled"].(bool))

			config := rule["apiNamingRuleConfig"].(map[string]interface{})
			if configType, ok := config["segmentMatchingBasedRuleConfig"].(map[string]interface{}); ok {
				d.Set("regexes", configType["regexes"].([]interface{}))
				d.Set("values", configType["values"].([]interface{}))
			}

			// Optionally handle other fields like spanFilter, etc.
			// ...
			break
		}
	}

	return nil
}

func resourceApiNamingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	regexes := jsonifyList(d.Get("regexes").([]interface{}))
	values := jsonifyList(d.Get("values").([]interface{}))
	serviceNames := jsonifyList(d.Get("service_names").([]interface{}))
	environmentNames := jsonifyList(d.Get("environment_names").([]interface{}))

	query := fmt.Sprintf(`
mutation {
  updateApiNamingRule(
    input: {
      id: "%s"
      name: "%s"
      disabled: %t
      apiNamingRuleConfig: {
        apiNamingRuleConfigType: SEGMENT_MATCHING
        segmentMatchingBasedRuleConfig: {
          regexes: %s
          values: %s
        }
      }
      spanFilter: {
        logicalSpanFilter: {
          logicalOperator: AND
          spanFilters: [
            {
              relationalSpanFilter: {
                field: SERVICE_NAME
                relationalOperator: IN
                value: %s
              }
            },
            {
              relationalSpanFilter: {
                field: ENVIRONMENT_NAME
                relationalOperator: IN
                value: %s
              }
            }
          ]
        }
      }
    }
  ) {
    id
  }
}`, id, name, disabled, regexes, values, serviceNames, environmentNames)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Printf("Error while executing GraphQL query: %s", err)
		return err
	}

	log.Printf("GraphQL response: %s", responseStr)
	if response["data"] != nil && response["data"].(map[string]interface{})["updateApiNamingRule"] != nil {
		updatedId := response["data"].(map[string]interface{})["updateApiNamingRule"].(map[string]interface{})["id"].(string)
		d.SetId(updatedId)
	} else {
		return fmt.Errorf("could not update API naming rule, response data is incomplete")
	}

	return nil
}

func resourceApiNamingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Retrieve the ID of the resource to delete

	query := fmt.Sprintf(`
mutation {
  deleteApiNamingRule(input: { id: "%s" }) {
    success
    __typename
  }
}`, id)

	// Execute the GraphQL mutation
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		log.Printf("Error while executing GraphQL query: %s", err)
		return err
	}

	log.Printf("GraphQL response: %s", responseStr)

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	// Check the success field from the response to ensure the deletion was processed correctly
	successResponse := response["data"].(map[string]interface{})["deleteApiNamingRule"].(map[string]interface{})
	success := successResponse["success"].(bool)
	if !success {
		return fmt.Errorf("failed to delete API Naming Rule with ID %s", id)
	}

	// If deletion was successful, remove the resource ID from the state
	d.SetId("")
	return nil
}
