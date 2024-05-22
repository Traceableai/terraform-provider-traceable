package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNotificationChannel() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNotificationChannelRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the channel",
				Required:    true,
			},
			"channel_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Id of notification channel",
				Computed:    true,
			},
		},
	}
}

func dataSourceNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	name:=d.Get("name").(string)
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
	log.Println("this is the gql response %s",response)
	ruleDetails:=getRuleDetailsFromRulesListUsingIdName(response,"notificationChannels" ,name,"channelId","channelName")
	if len(ruleDetails)==0{
		return fmt.Errorf("No rules found with name %s",name)
	}
	channelId:=ruleDetails["channelId"].(string)
	log.Println("Rule found with name %s",channelId)
	d.Set("channel_id",channelId)
	d.SetId(channelId)
	return nil
}

