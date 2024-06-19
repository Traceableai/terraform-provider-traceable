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
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API exclusion rule",
				Required:    true,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Required:    true,
			},
			"regexes": {
				Type:        schema.TypeString,
				Description: "http path regex pattern for the rule",
				Required:    true,
			},
			"service_names": {
				Type:        schema.TypeList,
				Description: "List of service names to apply the rule",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environment_names": {
				Type:        schema.TypeList,
				Description: "List of environment names to apply the rule",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
func resourceApiExclusionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	regexes := d.Get("regexes").(string)
	serviceNames := d.Get("service_names").([]interface{})
	environmentNames := d.Get("environment_names").([]interface{})

	var spanFilters []string

	// Checking if env name list is (empty -> all env case)
	environmentNames, ok := d.Get("environment_names").([]interface{})
	if !ok {
		fmt.Println("environment_names is not a valid []interface{} type")
	}
	if len(environmentNames) > 0 {
		firstElement, ok := environmentNames[0].(string) // Proper type assertion
		if ok && firstElement != "" {
			fmt.Println("The first environment name is:", firstElement)
			spanFilters = append(spanFilters, fmt.Sprintf(`
		{relationalSpanFilter: {field: ENVIRONMENT_NAME,relationalOperator: IN,value: %s}}`, jsonifyList(environmentNames)))
		} else {
			fmt.Println("The first environment name is empty or not a string")
		}
	} else {
		fmt.Println("The environment_names list is empty")
	}

	// Checking if service name list is (empty -> all services case)
	serviceNames, ok2 := d.Get("service_names").([]interface{})
	if !ok2 {
		fmt.Println("service_names is not a valid []interface{} type")
	}
	if len(serviceNames) > 0 {
		firstElement, ok2 := serviceNames[0].(string) // Proper type assertion
		if ok2 && firstElement != "" {
			fmt.Println("The first service name is:", firstElement)
			spanFilters = append(spanFilters, fmt.Sprintf(`
		{relationalSpanFilter: {field: SERVICE_NAME,relationalOperator: IN,value: %s}}`, jsonifyList(serviceNames)))
		} else {
			fmt.Println("The first service name is empty or not a string")
		}
	} else {
		fmt.Println("The service_names list is empty")
	}

	spanFilters = append(spanFilters, fmt.Sprintf(`
        {relationalSpanFilter: {field: URL,relationalOperator: REGEX_MATCH,value: "%s"}}`, regexes))

	spanFilterQueryPart := fmt.Sprintf("spanFilters: [%s]", strings.Join(spanFilters, ","))

	query := fmt.Sprintf(`mutation{createExcludeSpanRule(input:{name:"%s" disabled:%t spanFilter:{logicalSpanFilter:{logicalOperator:AND,%s}}}){id __typename}}`, name, disabled, spanFilterQueryPart)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["createExcludeSpanRule"] != nil {
		id := response["data"].(map[string]interface{})["createExcludeSpanRule"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf(responseStr)
	}

	return nil
}

func resourceApiExclusionRuleRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	query := `{excludeSpanRules{results{id name disabled spanFilter{logicalSpanFilter{logicalOperator spanFilters{relationalSpanFilter{relationalOperator key value field}}}}}}}`

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return fmt.Errorf("error parsing JSON response: %s", err)
	}

	// Navigate through the response to find the rule with the matching ID
	results := response["data"].(map[string]interface{})["excludeSpanRules"].(map[string]interface{})["results"].([]interface{})
	for _, item := range results {
		rule := item.(map[string]interface{})
		if rule["id"].(string) == id {
			// Set the Terraform state to match the fetched data
			d.Set("name", rule["name"].(string))
			d.Set("disabled", rule["disabled"].(bool))

			// Handle spanFilter data
			if spanFilter, exists := rule["spanFilter"].(map[string]interface{}); exists {
				if logicalSpanFilter, exists := spanFilter["logicalSpanFilter"].(map[string]interface{}); exists {
					if spanFilters, exists := logicalSpanFilter["spanFilters"].([]interface{}); exists {
						var serviceNames []string
						var environmentNames []string
						for _, sf := range spanFilters {
							filter := sf.(map[string]interface{})["relationalSpanFilter"].(map[string]interface{})
							field := filter["field"].(string)
							value := filter["value"]

							switch field {
							case "URL":
								d.Set("regexes", value.(string))
							case "SERVICE_NAME":
								serviceNames = append(serviceNames, convertToStringSlicetype(value)...)
							case "ENVIRONMENT_NAME":
								environmentNames = append(environmentNames, convertToStringSlicetype(value)...)
							}
						}
						d.Set("service_names", schema.NewSet(schema.HashString, convertToInterfaceSlice(serviceNames)))
						d.Set("environment_names", schema.NewSet(schema.HashString, convertToInterfaceSlice(environmentNames)))
					}
				}
			}

			return nil
		}
	}

	return fmt.Errorf("no exclusion rule found with ID %s", d.Id())
}

func resourceApiExclusionRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Retrieve the ID of the resource
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	regexes := d.Get("regexes").(string)
	serviceNames := d.Get("service_names").([]interface{})
	environmentNames := d.Get("environment_names").([]interface{})

	var spanFilters []string

	// Checking if env name list is (empty -> all env case)
	environmentNames, ok := d.Get("environment_names").([]interface{})
	if !ok {
		fmt.Println("environment_names is not a valid []interface{} type")
	}
	if len(environmentNames) > 0 {
		firstElement, ok := environmentNames[0].(string) // Proper type assertion
		if ok && firstElement != "" {
			fmt.Println("The first environment name is:", firstElement)
			spanFilters = append(spanFilters, fmt.Sprintf(`
		{relationalSpanFilter: {field: ENVIRONMENT_NAME,relationalOperator: IN,value: %s}}`, jsonifyList(environmentNames)))
		} else {
			fmt.Println("The first environment name is empty or not a string")
		}
	} else {
		fmt.Println("The environment_names list is empty")
	}

	// Checking if service name list is (empty -> all services case)
	serviceNames, ok2 := d.Get("service_names").([]interface{})
	if !ok2 {
		fmt.Println("service_names is not a valid []interface{} type")
	}
	if len(serviceNames) > 0 {
		firstElement, ok2 := serviceNames[0].(string) // Proper type assertion
		if ok2 && firstElement != "" {
			fmt.Println("The first service name is:", firstElement)
			spanFilters = append(spanFilters, fmt.Sprintf(`
		{relationalSpanFilter: {field: SERVICE_NAME,relationalOperator: IN,value: %s}}`, jsonifyList(serviceNames)))
		} else {
			fmt.Println("The first service name is empty or not a string")
		}
	} else {
		fmt.Println("The service_names list is empty")
	}

	spanFilters = append(spanFilters, fmt.Sprintf(`
        {relationalSpanFilter: {field: URL,relationalOperator: REGEX_MATCH,value: "%s"}}`, regexes))

	spanFilterQueryPart := fmt.Sprintf("spanFilters: [%s]", strings.Join(spanFilters, ","))

	query := fmt.Sprintf(`mutation{updateExcludeSpanRule(input:{id:"%s" name:"%s" disabled:%t spanFilter:{logicalSpanFilter:{logicalOperator:AND,%s}}}){id __typename}}`, id, name, disabled, spanFilterQueryPart)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error parsing JSON response from update operation: %s", err)	
	}

	if response["data"] == nil || response["data"].(map[string]interface{})["updateExcludeSpanRule"] == nil {
		log.Printf("GraphQL update did not return expected results: %s", responseStr)
		return fmt.Errorf("update operation did not return expected results")
	}

	log.Printf("Updated API Exclusion Rule: %s", responseStr)
	if response["data"] != nil && response["data"].(map[string]interface{})["updateExcludeSpanRule"] != nil {
		updatedId := response["data"].(map[string]interface{})["updateExcludeSpanRule"].(map[string]interface{})["id"].(string)
		d.SetId(updatedId)
	} else {
		return fmt.Errorf("could not update API exclusion rule, response data is incomplete")
	}
	return nil
}
func resourceApiExclusionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Retrieve the ID of the resource to delete

	query := fmt.Sprintf(`mutation{deleteExcludeSpanRule(input:{id:%q}){success __typename}}`, id)

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL mutation for deletion: %s", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
		return fmt.Errorf("Error parsing JSON response: %s", err)
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
