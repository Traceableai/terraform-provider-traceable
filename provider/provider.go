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
			"traceable_ip_range_rule":                            resourceIpRangeRule(),
			"traceable_user_attribution_rule_basic_auth":         resourceUserAttributionBasicAuthRule(),
			"traceable_user_attribution_rule_req_header":         resourceUserAttributionRequestHeaderRule(),
			"traceable_user_attribution_rule_jwt_authentication": resourceUserAttributionJwtAuthRule(),
			"traceable_user_attribution_rule_response_body": resourceUserAttributionResponseBodyRule(),
			"traceable_user_attribution_rule_custom_json": resourceUserAttributionCustomJsonRule(),
			"traceable_user_attribution_rule_custom_token": resourceUserAttributionCustomTokenRule(),
			"traceable_notification_channel": resourceNotificationChannelRule(),
			"traceable_notification_rule_logged_threat_activity": resourceNotificationRuleLoggedThreatActivity(),
			"traceable_notification_rule_blocked_threat_activity": resourceNotificationRuleBlockedThreatActivity(),
			"traceable_notification_rule_threat_actor_status": resourceNotificationRuleThreatActorStatusChange(),
			"traceable_notification_rule_actor_severity_change": resourceNotificationRuleActorSeverityChange(),
			"traceable_notification_rule_api_documentation": resourceNotificationRuleApiDocumentation(),
			"traceable_notification_rule_posture_events": resourceNotificationRulePostureEvents(),
			"traceable_api_naming_rule":                          resourceApiNamingRule(),
			// "traceable_api_exclusion_rule":                       resourceApiExclusionRule(),
			"traceable_label_creation_rule":                      resourceLabelCreationRule(),
			"traceable_agent_token":                              resourceAgentToken(),

		},
		DataSourcesMap: map[string]*schema.Resource{
			"traceable_notification_channels": dataSourceNotificationChannel(),
			"traceable_splunk_integration":    dataSourceSplunkIntegration(),
			"traceable_syslog_integration":    dataSourceSyslogIntegration(),
			"traceable_endpoint_id":           dataSourceEndpointId(),
			"traceable_service_id":            dataSourceServiceId(),
			"traceable_agent_token":           dataSourceAgentToken(),
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
