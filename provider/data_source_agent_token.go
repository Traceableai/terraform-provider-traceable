package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAgentToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentTokenRead,

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

func dataSourceAgentTokenRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	query := `{agentTokenMetadata {results {id name createdBy creationTimestamp lastUsedTimestamp __typename}}}`

	responseStr, err := ExecuteQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while parsing GraphQL response: %s", err)
	}

	if results, ok := response["data"].(map[string]interface{})["agentTokenMetadata"].(map[string]interface{})["results"].([]interface{}); ok {
		for _, result := range results {
			if token := result.(map[string]interface{}); token["name"].(string) == name {
				d.SetId(token["id"].(string))
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
		return fmt.Errorf("No agent token found with name %s", name)
	}

	return nil
}
