package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLabelApplicationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLabelApplicationRuleCreate,
		Read:   resourceLabelApplicationRuleRead,
		Update: resourceLabelApplicationRuleUpdate,
		Delete: resourceLabelApplicationRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the label application rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the label application rule",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the label application rule is enabled",
			},
			"condition_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_condition": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
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
						"value_condition": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value_condition_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"string_condition": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:     schema.TypeString,
													Required: true,
												},
												"string_condition_value_type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"values": {
													Type:     schema.TypeList,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"unary_condition": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Description: "Conditions that determine when the rule is applied",
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"entity_types": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dynamic_label_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"static_labels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
				Description: "Action to apply when conditions match",
			},
		},
	}
}

func resourceLabelApplicationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	conditionList := d.Get("condition_list").([]interface{})
	action := d.Get("action").([]interface{})[0].(map[string]interface{}) // Assuming only one action

	conditionsStr := parseConditions(conditionList) // A function to parse and format conditions
	actionStr := jsonifyAction(action)              // A function to convert action map into a string for the mutation

	query := fmt.Sprintf(`
mutation {
  createLabelApplicationRule(
    labelApplicationRuleData: {
      name: "%s",
      description: "%s",
      enabled: %t,
      conditionList: [%s],
      action: %s
    }
  ) {
    id
  }
}`, name, description, enabled, conditionsStr, actionStr)

	// Execute the GraphQL mutation and handle the response
	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error executing GraphQL mutation: %s", err)
	}

	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	// Set the resource ID in Terraform state
	if id, ok := response["data"].(map[string]interface{})["createLabelApplicationRule"].(map[string]interface{})["id"].(string); ok {
		d.SetId(id)
		return nil
	}

	return fmt.Errorf("Failed to create label application rule: no ID returned")
}

func resourceLabelApplicationRuleRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Get the resource ID set during creation

	// GraphQL query to read the Label Application Rule details
	query := `{
	  labelApplicationRules {
	    results {
	      id
	      labelApplicationRuleData {
	        name
	        description
	        enabled
	        conditionList {
	          leafCondition {
	            keyCondition {
	              operator
	              value
	            }
	            valueCondition {
	              valueConditionType
	              stringCondition {
	                value
	                operator
	                stringConditionValueType
	                values
	              }
	              unaryCondition {
	                operator
	              }
	            }
	          }
	        }
	        action {
	          entityTypes
	          staticLabels {
	            ids
	          }
	          operation
	          type
	          dynamicLabelKey
	        }
	      }
	    }
	  }
	}`

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	results := response["data"].(map[string]interface{})["labelApplicationRules"].(map[string]interface{})["results"].([]interface{})
	for _, item := range results {
		rule := item.(map[string]interface{})
		if rule["id"].(string) == id {
			ruleData := rule["labelApplicationRuleData"].(map[string]interface{})
			d.Set("name", ruleData["name"].(string))
			d.Set("description", ruleData["description"].(string))
			d.Set("enabled", ruleData["enabled"].(bool))

			// Parse conditions
			conditions := parseConditions(ruleData["conditionList"].([]interface{}))
			d.Set("condition_list", conditions)

			// Parse action
			action := ruleData["action"].(map[string]interface{})
			d.Set("action", jsonifyAction(action))

			break
		}
	}

	return nil
}

func resourceLabelApplicationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	conditions := jsonifyConditionList(d.Get("condition_list").([]interface{}))
	action := jsonifyAction(d.Get("action").(map[string]interface{}))

	query := fmt.Sprintf(`
mutation {
  updateLabelApplicationRule(
    labelApplicationRule: {
      id: "%s"
      labelApplicationRuleData: {
        name: "%s"
        description: "%s"
        enabled: %t
        conditionList: [%s]
        action: %s
      }
    }
  ) {
    id
  }
}`, id, name, description, enabled, conditions, action)

	log.Printf("Update query: %s", query) // Log the query for debugging purposes

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error parsing GraphQL response: %s", err)
	}

	log.Printf("GraphQL response: %s", responseStr)

	// Check if the update was successful and the correct ID was returned
	if responseData, ok := response["data"].(map[string]interface{}); ok {
		if updateResponse, ok := responseData["updateLabelApplicationRule"].(map[string]interface{}); ok {
			updatedId := updateResponse["id"].(string)
			if updatedId != id {
				return fmt.Errorf("ID mismatch after update: got %s but expected %s", updatedId, id)
			}
		} else {
			return fmt.Errorf("Update failed, no rule returned")
		}
	} else {
		return fmt.Errorf("Malformed response data")
	}

	return nil
}

func resourceLabelApplicationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Retrieve the ID of the resource to delete

	query := fmt.Sprintf(`
mutation {deleteLabelApplicationRule(id: "%s")}`, id)

	// Execute the GraphQL mutation
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	log.Printf("GraphQL response: %s", responseStr)

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return err
	}

	// The response structure needs to confirm if a delete was successful
	if response["data"] == nil || response["data"].(map[string]interface{})["deleteLabelApplicationRule"] == nil {
		return fmt.Errorf("failed to delete Label Application Rule with ID %s", id)
	}

	// If deletion was successful, remove the resource ID from the state
	d.SetId("")
	return nil
}
