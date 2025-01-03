package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/rate_limiting"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"github.com/traceableai/terraform-provider-traceable/provider/malicious_sources"
	// "github.com/traceableai/terraform-provider-traceable/provider/data_classification"
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
			"traceable_ip_range_rule_block": malicious_sources.ResourceIpRangeRuleBlock(),
			"traceable_ip_range_rule_allow": malicious_sources.ResourceIpRangeRuleAllow(),
			"traceable_ip_range_rule_alert": malicious_sources.ResourceIpRangeRuleAlert(),
			"traceable_region_rule_block": malicious_sources.ResourceRegionRuleBlock(),
			"traceable_region_rule_alert": malicious_sources.ResourceRegionRuleAlert(),
			"traceable_email_domain_block": malicious_sources.ResourceEmailDomainBlock(),
			"traceable_email_domain_alert": malicious_sources.ResourceEmailDomainAlert(),
			"traceable_ip_type_rule_alert": malicious_sources.ResourceIpTypeRuleAlert(),
			"traceable_ip_type_rule_block": malicious_sources.ResourceIpTypeRuleBlock(),
			//"traceable_user_attribution_rule_basic_auth":         resourceUserAttributionBasicAuthRule(),
			//"traceable_user_attribution_rule_req_header":         resourceUserAttributionRequestHeaderRule(),
			//"traceable_user_attribution_rule_jwt_authentication": resourceUserAttributionJwtAuthRule(),
			//"traceable_user_attribution_rule_response_body": resourceUserAttributionResponseBodyRule(),
			//"traceable_user_attribution_rule_custom_json": resourceUserAttributionCustomJsonRule(),
			//"traceable_user_attribution_rule_custom_token": resourceUserAttributionCustomTokenRule(),
			// "traceable_notification_channel": resourceNotificationChannelRule(),
			// "traceable_notification_rule_logged_threat_activity": resourceNotificationRuleLoggedThreatActivity(),
			// "traceable_notification_rule_blocked_threat_activity": resourceNotificationRuleBlockedThreatActivity(),
			// "traceable_notification_rule_threat_actor_status": resourceNotificationRuleThreatActorStatusChange(),
			// "traceable_notification_rule_actor_severity_change": resourceNotificationRuleActorSeverityChange(),
			"traceable_api_naming_rule": resourceApiNamingRule(),
			// "traceable_api_exclusion_rule":                       resourceApiExclusionRule(),
			"traceable_label_creation_rule": resourceLabelCreationRule(),
			"traceable_rate_limiting_block": rate_limiting.ResourceRateLimitingRuleBlock(),
			"traceable_detection_policies":  resourceDetectionConfigRule(),
			// "traceable_agent_token":                              resourceAgentToken(),
			"traceable_custom_signature_allow":                   custom_signature.ResourceCustomSignatureAllowRule(),
			"traceable_custom_signature_block":                   custom_signature.ResourceCustomSignatureBlockRule(),
			"traceable_custom_signature_alert":                   custom_signature.ResourceCustomSignatureAlertRule(),
			"traceable_custom_signature_testing":                   custom_signature.ResourceCustomSignatureTestingRule(),
			// "traceable_data_classification_rule":                   data_classification.ResourceDataClassification(),

		},
		DataSourcesMap: map[string]*schema.Resource{
			// "traceable_notification_channels": dataSourceNotificationChannel(),
			//"traceable_splunk_integration":    dataSourceSplunkIntegration(),
			//"traceable_syslog_integration":    dataSourceSyslogIntegration(),
			"traceable_endpoint_id": dataSourceEndpointId(),
			"traceable_service_id":  dataSourceServiceId(),
			// "traceable_agent_token":           dataSourceAgentToken(),
		},
		ConfigureFunc: graphqlConfigure,
	}
}

func graphqlConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &common.GraphqlProviderConfig{
		GQLServerUrl: d.Get("platform_url").(string),
		ApiToken:     d.Get("api_token").(string),
	}
	return config, nil
}
