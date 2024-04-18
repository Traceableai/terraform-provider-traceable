package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"platform_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Platform url where we need to create rules",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "platform api token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"example_ip_range_rule": resourceIpRangeRule(),
			// "example_api_abuse_rule": resourceApiAbuseRule(),
			// Add more resources here as needed
		},
		ConfigureFunc: graphqlConfigure,
	}
}

type graphqlProviderConfig struct {
	GQLServerUrl string
	ApiToken     string
}

func graphqlConfigure( d *schema.ResourceData) (interface{}, error) {
	config := &graphqlProviderConfig{
		GQLServerUrl: d.Get("platform_url").(string),
		ApiToken:     d.Get("api_token").(string),
	}
	return config, nil
}
