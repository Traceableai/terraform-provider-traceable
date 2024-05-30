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
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the API naming rule",
				Required:    true,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable the rule",
				Required:    true,
			},
			"regexes": {
				Type:        schema.TypeList,
				Description: "List of regex patterns for the rule",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"values": {
				Type:        schema.TypeList,
				Description: "Corresponding values for the regex patterns",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceApiNamingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	disabled := d.Get("disabled").(bool)
	regexes := d.Get("regexes").([]interface{})
	values := d.Get("values").([]interface{})
	serviceNames := d.Get("service_names").([]interface{})
	environmentNames := d.Get("environment_names").([]interface{})

	if len(regexes) != len(values) {
		return fmt.Errorf("the number of regexes (%d) does not match the number of values (%d)", len(regexes), len(values))
	}

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

	// spanFilter part of the query willbe empty in case of all env and all services
	spanFilterQueryPart := ""
	if len(spanFilters) > 0 {
		spanFilterQueryPart = fmt.Sprintf("spanFilters: [%s]", strings.Join(spanFilters, ","))
	} else {
		spanFilterQueryPart = "spanFilters: []" // No filters to apply
	}

	query := fmt.Sprintf(`mutation{createApiNamingRule(input:{name:"%s" disabled:%t apiNamingRuleConfig:{apiNamingRuleConfigType:SEGMENT_MATCHING segmentMatchingBasedRuleConfig:{regexes:%s,values:%s}}spanFilter:{logicalSpanFilter:{logicalOperator:AND,%s}}}){id}}`, name, disabled, jsonifyList(regexes), jsonifyList(values), spanFilterQueryPart)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)

	}

	log.Printf("GraphQL response: %s", responseStr)

	if response["data"] != nil && response["data"].(map[string]interface{})["createApiNamingRule"] != nil {
		id := response["data"].(map[string]interface{})["createApiNamingRule"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf("could not create API naming rule, no ID returned")
	}

	return nil
}

func resourceApiNamingRuleRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id() // Get the resource ID set during creation

	// GraphQL query to read the API Naming Rule details
	query := `{apiNamingRules{results{id name disabled apiNamingRuleConfig{apiNamingRuleConfigType segmentMatchingBasedRuleConfig{regexes values}}spanFilter{logicalSpanFilter{logicalOperator spanFilters{relationalSpanFilter{field relationalOperator value}}}}}}}`

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
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
	regexes := d.Get("regexes").([]interface{})
	values := d.Get("values").([]interface{})
	serviceNames := d.Get("service_names").([]interface{})
	environmentNames := d.Get("environment_names").([]interface{})

	if len(regexes) != len(values) {
		return fmt.Errorf("the number of regexes (%d) does not match the number of values (%d)", len(regexes), len(values))
	}

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

	spanFilterQueryPart := ""
	if len(spanFilters) > 0 {
		spanFilterQueryPart = fmt.Sprintf("spanFilters: [%s]", strings.Join(spanFilters, ","))
	} else {
		spanFilterQueryPart = "spanFilters: []" // No filters to apply
	}

	query := fmt.Sprintf(`mutation{updateApiNamingRule(input:{id:"%s" name:"%s" disabled:%t apiNamingRuleConfig:{apiNamingRuleConfigType:SEGMENT_MATCHING segmentMatchingBasedRuleConfig:{regexes:%s,values:%s}}spanFilter:{logicalSpanFilter:{logicalOperator:AND,%s}}}){id}}`, id, name, disabled, jsonifyList(regexes), jsonifyList(values), spanFilterQueryPart)

	log.Printf((query))
	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
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
mutation {deleteApiNamingRule(input: { id: "%s" }) {success,__typename}}`, id)

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
