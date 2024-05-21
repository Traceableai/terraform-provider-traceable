package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNotificationChannel() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNotificationChannel,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the channel",
				Required:    true,
			},
			"channelId": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Id of notification channel",
				Computed:    true,
			},
		},
	}
}

func dataSourceNotificationChannel(d *schema.ResourceData, meta interface{}) error {
	name=d.Get("name").(string)
	query := `{
		notificationChannels {
		  results {
			channelId
			channelName
		  }
		}
	  }`	  

	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("error while executing GraphQL query: %s", err)
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON response: %s", err)
	}

	

	return fmt.Errorf("no exclusion rule found with ID %s", d.Id())
}

