package label_management

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ResourceLabelCreationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLabelCreationRuleCreate,
		Read:   resourceLabelCreationRuleRead,
		Update: resourceLabelCreationRuleUpdate,
		Delete: resourceLabelCreationRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the label",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A description of the label",
				Optional:    true,
				Default:     "",
			},
			"color": {
				Type:        schema.TypeString,
				Description: "The color code of the label",
				Optional:    true,
				Default:     "",
			},
		},
	}
}

func resourceLabelCreationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	color := d.Get("color").(string)

	mutation := fmt.Sprintf(CREATE_LABEL_QUERY, name, description, color)

	// Execute the GraphQL mutation
	responseStr, err := common.CallExecuteQuery(mutation, meta)
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
			d.Set("name", createResponse["name"])
			d.Set("description", createResponse["description"])
			d.Set("color", createResponse["color"])
			return nil
		}
	}

	return fmt.Errorf("label exist with same name")
}

func resourceLabelCreationRuleRead(d *schema.ResourceData, meta interface{}) error {
	responseStr, err := common.CallExecuteQuery(READ_LABELS_QUERY, meta)
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
		if labelID, ok := label["id"].(string); ok && labelID == d.Id() {
			// Set the Terraform state to match the fetched data
			if name, ok := label["name"].(string); ok {
				d.Set("name", name)
			}
			if description, ok := label["description"].(string); ok {
				d.Set("description", description)
			}
			if color, ok := label["color"].(string); ok {
				d.Set("color", color)
			}
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceLabelCreationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	color := d.Get("color").(string)

	mutation := fmt.Sprintf(UPDATE_LEBEL_QUERY, id, name, description, color)

	responseStr, err := common.CallExecuteQuery(mutation, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL update mutation: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response from update mutation: %s", err)
	}

	if responseData, ok := response["data"].(map[string]interface{})["updateLabel"].(map[string]interface{}); ok {
		if id, ok := responseData["id"].(string); ok {
			d.SetId(id)
		}
		if name, ok := responseData["name"].(string); ok {
			d.Set("name", name)
		}
		if description, ok := responseData["description"].(string); ok {
			d.Set("description", description)
		}
		if color, ok := responseData["color"].(string); ok {
			d.Set("color", color)
		}
	} else {
		return fmt.Errorf("error updating label: no label returned from API")
	}

	return nil
}

func resourceLabelCreationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("") // This tells Terraform that the resource no longer exists.
	log.Printf("[INFO] Resource %s effectively removed from Terraform state but not deleted from backend", d.Id())
	return nil
}
