package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLabelCreationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLabelCreationRuleCreate,
		Read:   resourceLabelCreationRuleRead,
		Update: resourceLabelCreationRuleUpdate,
		Delete: resourceLabelCreationRuleDelete,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The key of the label",
				Required:    true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A description of the label",
				Optional:    true, // Optional as per requirements
				Default:     "",
			},
			"color": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The color code of the label",
				Optional:    true, // Optional as per requirements
				Default:     "",
			},
		},
	}
}

func resourceLabelCreationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	key := d.Get("key").(string)
	description := d.Get("description").(string)
	color := d.Get("color").(string)

	mutation := fmt.Sprintf(`mutation{createLabel(label:{key:"%s",description:"%s",color:"%s"}){id key description color}}`, key, description, color)

	// Execute the GraphQL mutation
	responseStr, err := executeQuery(mutation, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response: %s", err)
	}

	// Handle the response to set the resource ID and possibly other attributes
	if responseData, ok := response["data"].(map[string]interface{}); ok {
		if createResponse, ok := responseData["createLabel"].(map[string]interface{}); ok {
			d.SetId(createResponse["id"].(string))
			d.Set("key", createResponse["key"])
			d.Set("description", createResponse["description"])
			d.Set("color", createResponse["color"])
			return nil
		}
	}

	return fmt.Errorf("could not create label Creation rule, no ID returned")
}

func resourceLabelCreationRuleRead(d *schema.ResourceData, meta interface{}) error {

	query := `query{labels{results{id key description color}}}`

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response: %s", err)
	}

	// Navigate through the response to find the label with the matching ID
	labelsResponse := response["data"].(map[string]interface{})["labels"].(map[string]interface{})["results"].([]interface{})
	for _, item := range labelsResponse {
		label := item.(map[string]interface{})
		if label["id"].(string) == d.Id() {
			// Set the Terraform state to match the fetched data
			d.Set("key", label["key"].(string))
			d.Set("description", label["description"].(string))
			d.Set("color", label["color"].(string))
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceLabelCreationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	// Retrieve the current ID and attributes from the schema
	id := d.Id()
	key := d.Get("key").(string)
	description := d.Get("description").(string)
	color := d.Get("color").(string)

	mutation := fmt.Sprintf(`mutation{updateLabel(label:{id:"%s",key:"%s",description:"%s",color:"%s"}){id key description color}}`, id, key, description, color)

	responseStr, err := executeQuery(mutation, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL update mutation: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response from update mutation: %s", err)
	}

	if responseData, ok := response["data"].(map[string]interface{})["updateLabel"].(map[string]interface{}); ok {
		d.SetId(responseData["id"].(string))
		d.Set("key", responseData["key"].(string))
		d.Set("description", responseData["description"].(string))
		d.Set("color", responseData["color"].(string))
	} else {
		return fmt.Errorf("error updating label: no label returned from API")
	}

	return nil
}

func resourceLabelCreationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// Since actual deletion is not supported by the API, we simply remove the resource from the state
	d.SetId("") // This tells Terraform that the resource no longer exists.

	// Optionally, you could log this operation to inform that it's a dummy operation
	log.Printf("[INFO] Resource %s effectively removed from Terraform state but not deleted from backend", d.Id())

	return nil
}
