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
			"traceable_ip_range_rule": resourceIpRangeRule(),
			"traceable_user_attribution_rule_basic_auth": resourceUserAttributionBasicAuthRule(),
			"traceable_user_attribution_rule_req_header": resourceUserAttributionRequestHeaderRule(),
		},
		ConfigureFunc: graphqlConfigure,
	}
}

type graphqlProviderConfig struct {
	GQLServerUrl string
	ApiToken     string
}

func graphqlConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &graphqlProviderConfig{
		GQLServerUrl: d.Get("platform_url").(string),
		ApiToken:     d.Get("api_token").(string),
	}
	return config, nil
}
