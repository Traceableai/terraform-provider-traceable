package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/traceableai/terraform-provider-traceable/provider/common"
	"github.com/traceableai/terraform-provider-traceable/provider/custom_signature"
	"github.com/traceableai/terraform-provider-traceable/provider/enumeration"
	"github.com/traceableai/terraform-provider-traceable/provider/malicious_sources"
	"github.com/traceableai/terraform-provider-traceable/provider/rate_limiting"
	"github.com/traceableai/terraform-provider-traceable/provider/waap"
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
			"traceable_region_rule_block":   malicious_sources.ResourceRegionRuleBlock(),
			"traceable_region_rule_alert":   malicious_sources.ResourceRegionRuleAlert(),
			"traceable_email_domain_block":  malicious_sources.ResourceEmailDomainBlock(),
			"traceable_email_domain_alert":  malicious_sources.ResourceEmailDomainAlert(),
			"traceable_ip_type_rule_alert":  malicious_sources.ResourceIpTypeRuleAlert(),
			"traceable_ip_type_rule_block":  malicious_sources.ResourceIpTypeRuleBlock(),
			// "traceable_user_attribution_rule_basic_auth":                  ResourceUserAttributionBasicAuthRule(),
			// "traceable_user_attribution_rule_req_header":                  ResourceUserAttributionRequestHeaderRule(),
			// "traceable_user_attribution_rule_jwt_authentication":          ResourceUserAttributionJwtAuthRule(),
			// "traceable_user_attribution_rule_response_body":               ResourceUserAttributionResponseBodyRule(),
			// "traceable_user_attribution_rule_custom_json":                 ResourceUserAttributionCustomJsonRule(),
			// "traceable_user_attribution_rule_custom_token":                ResourceUserAttributionCustomTokenRule(),
			// "traceable_notification_channel":                              notification.ResourceNotificationChannelRule(),
			// "traceable_notification_rule_logged_threat_activity":          notification.ResourceNotificationRuleLoggedThreatActivity(),
			// "traceable_notification_rule_protection_configuration_change": notification.ResourceNotificationRuleProtectionConfig(),
			// "traceable_notification_rule_team_activity":                   notification.ResourceNotificationRuleTeamActivity(),
			// "traceable_notification_rule_api_naming":                      notification.ResourceNotificationRuleApiNaming(),
			// "traceable_notification_rule_api_documentation":               notification.ResourceNotificationRuleApiDocumentation(),
			// "traceable_notification_rule_data_collection":                 notification.ResourceNotificationRuleDataCollection(),
			// "traceable_notification_rule_risk_scoring":                    notification.ResourceNotificationRuleRiskScoring(),
			// "traceable_notification_rule_exclude_rule":                    notification.ResourceNotificationRuleExcludeRule(),
			// "traceable_notification_rule_blocked_threat_activity":         notification.ResourceNotificationRuleBlockedThreatActivity(),
			// "traceable_notification_rule_threat_actor_status":             notification.ResourceNotificationRuleThreatActorStatusChange(),
			// "traceable_notification_rule_actor_severity_change":           notification.ResourceNotificationRuleActorSeverityChange(),
			// "traceable_notification_rule_label_configuration":             notification.ResourceNotificationRuleLabelConfiguration(),
			// "traceable_notification_rule_notification_configuration":      notification.ResourceNotificationRuleNotificationConfiguration(),
			// "traceable_notification_rule_data_class_configuration":        notification.ResourceNotificationRuleDataClassificationConfig(),
			// "traceable_notification_rule_posture_events":                  notification.ResourceNotificationRulePostureEvents(),
			"traceable_api_naming_rule":          ResourceApiNamingRule(),
			"traceable_label_creation_rule":      resourceLabelCreationRule(),
			"traceable_rate_limiting_block":      rate_limiting.ResourceRateLimitingRuleBlock(),
			"traceable_enumeration_rule":         enumeration.ResourceEnumerationRule(),
			"traceable_waap_policies":            waap.ResourceDetectionConfigRule(),
			"traceable_custom_signature_allow":   custom_signature.ResourceCustomSignatureAllowRule(),
			"traceable_custom_signature_block":   custom_signature.ResourceCustomSignatureBlockRule(),
			"traceable_custom_signature_alert":   custom_signature.ResourceCustomSignatureAlertRule(),
			"traceable_custom_signature_testing": custom_signature.ResourceCustomSignatureTestingRule(),
			// "traceable_api_exclusion_rule":                       ResourceApiExclusionRule(),
			// "traceable_agent_token":                              ResourceAgentToken(),
			// "traceable_data_classification_rule":                   data_classification.ResourceDataClassification(),
			// "traceable_data_classification_overrides":                   data_classification.ResourceDataClassificationOverrides(),
			// "traceable_data_sets":                   data_classification.ResourceDataSetsRule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "traceable_notification_channels": DataSourceNotificationChannel(),
			//"traceable_splunk_integration":    DataSourceSplunkIntegration(),
			//"traceable_syslog_integration":    DataSourceSyslogIntegration(),
			"traceable_endpoint_id": dataSourceEndpointId(),
			"traceable_service_id":  dataSourceServiceId(),
			// "traceable_data_set_id":  data_classification.DataSourceDatSetId(),
			// "traceable_agent_token":           DataSourceAgentToken(),
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
