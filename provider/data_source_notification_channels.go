package provider

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )

// func DataSourceNotificationChannel() *schema.Resource {
// 	return &schema.Resource{
// 		Read: DataSourceNotificationChannelRead,

// 		Schema: map[string]*schema.Schema{
// 			"name": {
// 				Type:        schema.TypeString,
// 				Description: "Name of the channel",
// 				Required:    true,
// 			},
// 			"channel_id": {
// 				Type:        schema.TypeString,
// 				Description: "Id of notification channel",
// 				Computed:    true,
// 			},
// 		},
// 	}
// }

// func DataSourceNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
// 	name := d.Get("name").(string)
// 	query := `{
// 		notificationChannels {
// 		  results {
// 			channelId
// 			channelName
// 		  }
// 		}
// 	  }`

// 	responseStr, err := ExecuteQuery(query, meta)
// 	if err != nil {
// 		return fmt.Errorf("error while executing GraphQL query: %s", err)
// 	}

// 	var response map[string]interface{}
// 	err = json.Unmarshal([]byte(responseStr), &response)
// 	if err != nil {
// 		return fmt.Errorf("error parsing JSON response: %s", err)
// 	}
// 	log.Printf("this is the gql response %s", response)
// 	ruleDetails := GetRuleDetailsFromRulesListUsingIdName(response, "notificationChannels", name, "channelId", "channelName")
// 	if len(ruleDetails) == 0 {
// 		return fmt.Errorf("no rules found with name %s", name)
// 	}
// 	channelId := ruleDetails["channelId"].(string)
// 	log.Printf("Rule found with name %s", channelId)
// 	d.Set("channel_id", channelId)
// 	d.SetId(channelId)
// 	return nil
// }
