package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiExclusionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiExclusionRuleCreate,
		Read:   resourceApiExclusionRuleRead,
		Update: resourceApiExclusionRuleUpdate,
		Delete: resourceApiExclusionRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the API exclusion rule",
				Required:    true,
			},
			"disabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Required:    true,
			},
			"span_filters": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of span filters for the exclusion rule",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"relational_operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}
func resourceApiExclusionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	spanFilters := parseSpanFilters(d.Get("span_filters").([]interface{}))

	// Build the spanFilters GraphQL part based on provided filter data
	var graphqlSpanFilters string
	if len(spanFilters) > 0 {
		graphqlSpanFilters = fmt.Sprintf("spanFilters: [%s]", strings.Join(spanFilters, ", "))
	} else {
		// Default filter focusing on URL if no specific environment or service filters are provided
		graphqlSpanFilters = `spanFilters: [
            {
              relationalSpanFilter: {
                field: URL
                relationalOperator: REGEX_MATCH
                value: "hello\\test"
              }
            }
        ]`
	}

	query := fmt.Sprintf(`
mutation {
  createExcludeSpanRule(
    input: {
      name: "%s"
      disabled: %t
      spanFilter: {
        logicalSpanFilter: {
          logicalOperator: AND
          %s
        }
      }
    }
  ) {
    id
    __typename
  }
}`, name, disabled, graphqlSpanFilters)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Printf("Error while executing GraphQL query: %s", err)
		return err
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["createExcludeSpanRule"] != nil {
		id := response["data"].(map[string]interface{})["createExcludeSpanRule"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf("could not create API exclusion rule, no ID returned")
	}

	return nil
}

func parseSpanFilters(filters []interface{}) []string {
	var filterStrs []string
	for _, item := range filters {
		f := item.(map[string]interface{})
		// Check if value is not just an empty string or an array with an empty string
		values, ok := f["value"].([]interface{})
		if ok && len(values) > 0 && values[0] != "" {
			filterStr := fmt.Sprintf(`{
				relationalSpanFilter: {
					field: "%s"
					relationalOperator: %s
					value: %s
				}
			}`, f["field"].(string), f["relational_operator"].(string), jsonifyList(values))
			filterStrs = append(filterStrs, filterStr)
		}
	}
	return filterStrs
}

func jsonifySpanFilters(filters []interface{}) string {
	var filterStrs []string
	for _, item := range filters {
		f := item.(map[string]interface{})
		filterStr := fmt.Sprintf(`{
			relationalSpanFilter: {
				field: "%s"
				relationalOperator: %s
				value: "%s"
			}
		}`, f["field"].(string), f["relational_operator"].(string), f["value"].(string))
		filterStrs = append(filterStrs, filterStr)
	}
	return "[" + strings.Join(filterStrs, ", ") + "]"
}

func resourceApiExclusionRuleRead(d *schema.ResourceData, meta interface{}) error {

	query := `
{
  excludeSpanRules {
    results {
      id
      name
      creationTime
      lastUpdatedTime
      disabled
      spanFilter {
        logicalSpanFilter {
          logicalOperator
          spanFilters {
            relationalSpanFilter {
              relationalOperator
              key
              value
              field
              __typename
            }
            __typename
          }
          __typename
        }
        __typename
      }
      __typename
    }
    __typename
  }
}
`
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		log.Printf("Error while executing GraphQL query: %s", err)
		return err
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Printf("Error parsing JSON response: %s", err)
		return err
	}

	// Navigate through the response to find the rule with the matching ID
	results := response["excludeSpanRules"].(map[string]interface{})["results"].([]interface{})
	for _, item := range results {
		rule := item.(map[string]interface{})
		if rule["id"].(string) == d.Id() {
			d.Set("name", rule["name"].(string))
			d.Set("disabled", rule["disabled"].(bool))

			// Extract spanFilter details if they exist
			if spanFilter, ok := rule["spanFilter"].(map[string]interface{}); ok {
				logicalSpanFilter := spanFilter["logicalSpanFilter"].(map[string]interface{})
				d.Set("logical_operator", logicalSpanFilter["logicalOperator"].(string))

				// Process each relationalSpanFilter and collect span filter details
				var spanFilters []interface{}
				for _, span := range logicalSpanFilter["spanFilters"].([]interface{}) {
					filter := span.(map[string]interface{})["relationalSpanFilter"].(map[string]interface{})
					spanFilters = append(spanFilters, map[string]interface{}{
						"relational_operator": filter["relationalOperator"].(string),
						"key":                 filter["key"].(string),
						"value":               filter["value"].(string),
						"field":               filter["field"].(string),
					})
				}
				if err := d.Set("span_filters", spanFilters); err != nil {
					log.Printf("Error setting span_filters: %s", err)
					return err
				}
			}

			return nil
		}
	}

	log.Printf("No rule found with ID %s", d.Id())
	return fmt.Errorf("no rule found with ID %s", d.Id())
}

func resourceApiExclusionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Retrieve the ID of the resource
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	spanFiltersInterface := d.Get("span_filters").([]interface{})

	spanFilters := make([]string, 0)
	for _, filter := range spanFiltersInterface {
		f := filter.(map[string]interface{})
		field := f["field"].(string)
		relationalOperator := f["relational_operator"].(string)
		value := f["value"].(string)

		// Handling empty lists for specific filters like service or environment
		if field == "SERVICE_NAME" || field == "ENVIRONMENT_NAME" {
			values := strings.Split(value, ",")
			if len(values) == 1 && values[0] == "" {
				continue // Skip adding this filter if the list is effectively empty
			}
		}

		spanFilters = append(spanFilters, fmt.Sprintf(`{
            relationalSpanFilter: {
                field: %q
                relationalOperator: %q
                value: [%q]
            }
        }`, field, relationalOperator, value))
	}

	// Handle case where no valid filters are added
	if len(spanFilters) == 0 {
		spanFilters = append(spanFilters, `{
            relationalSpanFilter: {
                field: "URL"
                relationalOperator: "REGEX_MATCH"
                value: ["hello\\test"]
            }
        }`)
	}

	// Build the spanFilters part of the GraphQL query
	spanFilterQueryPart := fmt.Sprintf("[%s]", strings.Join(spanFilters, ", "))

	query := fmt.Sprintf(`
mutation {
    updateExcludeSpanRule(
        input: {
            id: %q
            name: %q
            disabled: %t
            spanFilter: {
                logicalSpanFilter: {
                    logicalOperator: AND
                    spanFilters: %s
                }
            }
        }
    ) {
        id
        __typename
    }
}`, id, name, disabled, spanFilterQueryPart)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		log.Printf("Error while executing GraphQL query: %s", err)
		return err
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		log.Printf("Error parsing JSON response from update operation: %s", err)
		return err
	}

	if response["data"] == nil || response["data"].(map[string]interface{})["updateExcludeSpanRule"] == nil {
		log.Printf("GraphQL update did not return expected results: %s", responseStr)
		return fmt.Errorf("update operation did not return expected results")
	}

	log.Printf("Updated API Exclusion Rule: %s", responseStr)
	return nil
}
func resourceApiExclusionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Retrieve the ID of the resource to delete

	query := fmt.Sprintf(`
mutation {
  deleteExcludeSpanRule(input: { id: %q }) {
    success
    __typename
  }
}`, id)

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		log.Printf("Error while executing GraphQL mutation for deletion: %s", err)
		return err
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		log.Printf("Error parsing JSON response: %s", err)
		return err
	}

	// Check the success field from the response to ensure the deletion was processed correctly
	successResponse := response["data"].(map[string]interface{})["deleteExcludeSpanRule"].(map[string]interface{})
	success := successResponse["success"].(bool)
	if !success {
		return fmt.Errorf("failed to delete API Exclusion Rule with ID %s", id)
	}

	// If deletion was successful, remove the resource ID from the state
	d.SetId("")
	return nil
}
