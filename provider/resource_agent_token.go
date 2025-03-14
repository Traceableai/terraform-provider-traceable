package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAgentToken() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAgentTokenCreate,
		Read:   ResourceAgentTokenRead,
		Update: ResourceAgentTokenUpdate,
		Delete: ResourceAgentTokenDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the agent token",
				Required:    true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "The agent token value",
				Computed:    true,
			},
			"created_by": {
				Type:        schema.TypeString,
				Description: "The creator of the agent token",
				Computed:    true,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Description: "The creation timestamp of the agent token",
				Computed:    true,
			},
			"last_used_timestamp": {
				Type:        schema.TypeString,
				Description: "The last used timestamp of the agent token",
				Computed:    true,
			},
		},
	}
}

func ResourceAgentTokenCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	query := fmt.Sprintf(`mutation {
		createAgentToken(input: {name: "%s"}) {
			id
			name
			token
			createdBy
			creationTimestamp
			__typename
		}
	}`, name)

	var response map[string]interface{}
	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	log.Printf("GraphQL response: %s", responseStr)

	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
	}

	if data, ok := response["data"].(map[string]interface{})["createAgentToken"].(map[string]interface{}); ok {
		d.SetId(data["id"].(string))
		d.Set("name", data["name"].(string))
		d.Set("token", data["token"].(string))
		d.Set("created_by", data["createdBy"].(string))
		d.Set("creation_timestamp", data["creationTimestamp"].(string))
	} else {
		return fmt.Errorf("could not create agent token, no data returned")
	}

	return nil
}

func ResourceAgentTokenRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	query := `{agentTokenMetadata {results {id name createdBy creationTimestamp lastUsedTimestamp __typename}}}`

	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
	}

	if results, ok := response["data"].(map[string]interface{})["agentTokenMetadata"].(map[string]interface{})["results"].([]interface{}); ok {
		for _, result := range results {
			if token := result.(map[string]interface{}); token["id"].(string) == id {
				d.Set("name", token["name"].(string))
				// Preserve the token value in state
				if v, ok := d.GetOk("token"); ok {
					d.Set("token", v)
				}
				d.Set("created_by", token["createdBy"].(string))
				d.Set("creation_timestamp", token["creationTimestamp"].(string))
				d.Set("last_used_timestamp", token["lastUsedTimestamp"].(string))
				return nil
			}
		}
		d.SetId("")
		return nil
	}

	return nil
}

func ResourceAgentTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)

	query := fmt.Sprintf(`mutation {
		updateAgentTokenMetadata(
			input: {id: "%s", name: "%s"}
		) {
			id
			name
			__typename
		}
	}`, id, name)

	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["updateAgentTokenMetadata"] != nil {
		updatedId := response["data"].(map[string]interface{})["updateAgentTokenMetadata"].(map[string]interface{})["id"].(string)
		d.SetId(updatedId)
	} else {
		return fmt.Errorf("could not update Agent Token, no data returned")
	}

	return nil
}

func ResourceAgentTokenDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	query := fmt.Sprintf(`mutation {
		deleteAgentToken(
			input: {id: "%s", forceDelete: false}
		) {
			id
			__typename
		}
	}`, id)

	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error while parsing GraphQL response: %s", err)
	}

	d.SetId("")
	return nil
}
