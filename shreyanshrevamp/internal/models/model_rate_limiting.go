package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RateLimitingRuleModel struct {
	Id               types.String      `tfsdk:"id"`
	Name             types.String      `tfsdk:"name"`
	Environments     types.List        `tfsdk:"environments"`
	Description      types.String      `tfsdk:"description"`
	Enabled          types.Bool        `tfsdk:"enabled"`
	ThresholdConfigs []ThresholdConfig `tfsdk:"threshold_configs"`
	Action           Action            `tfsdk:"action"`
	Sources          *Sources          `tfsdk:"sources"`
}

type ThresholdConfig struct {
	ApiAggregateType                     types.String `tfsdk:"api_aggregate_type"`
	UserAggregateType                    types.String `tfsdk:"user_aggregate_type"`
	RollingWindowCountAllowed            types.Int64  `tfsdk:"rolling_window_count_allowed"`
	RollingWindowDuration                types.String `tfsdk:"rolling_window_duration"`
	ThresholdConfigType                  types.String `tfsdk:"threshold_config_type"`
	DynamicMeanCalculationDuration       types.String `tfsdk:"dynamic_mean_calculation_duration"`
	DynamicDuration                      types.String `tfsdk:"dynamic_duration"`
	DynamicPercentageExcedingMeanAllowed types.Int64  `tfsdk:"dynamic_percentage_exceding_mean_allowed"`
}

type Action struct {
	ActionType       types.String      `tfsdk:"action_type"`
	Duration         types.String      `tfsdk:"duration"`
	EventSeverity    types.String      `tfsdk:"event_severity"`
	HeaderInjections []HeaderInjection `tfsdk:"header_injections"`
}

type HeaderInjection struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type Sources struct {
	Scanner          *ScannerSource             `tfsdk:"scanner"`
	IpAsn            *IpAsnSource               `tfsdk:"ip_asn"`
	IpConnectionType *IpConnectionTypeSource    `tfsdk:"ip_connection_type"`
	UserId           *UserIdSource              `tfsdk:"user_id"`
	EndpointLabels   types.List                 `tfsdk:"endpoint_labels"`
	Endpoints        types.List                 `tfsdk:"endpoints"`
	Attributes       []AttributeCondition       `tfsdk:"attribute"`
	IpReputation     types.String               `tfsdk:"ip_reputation"`
	IpLocationType   *IpLocationTypeSource      `tfsdk:"ip_location_type"`
	IpAbuseVelocity  types.String               `tfsdk:"ip_abuse_velocity"`
	IpAddress        *IpAddressSource           `tfsdk:"ip_address"`
	EmailDomain      *EmailDomainSource         `tfsdk:"email_domain"`
	UserAgents       *UserAgentsSource          `tfsdk:"user_agents"`
	Regions          *RegionsSource             `tfsdk:"regions"`
	IpOrganisation   *IpOrganisationSource      `tfsdk:"ip_organisation"`
	RequestResponse  []RequestResponseCondition `tfsdk:"request_response"`
}

type ScannerSource struct {
	ScannerTypesList types.List `tfsdk:"scanner_types_list"`
	Exclude          types.Bool `tfsdk:"exclude"`
}

type IpAsnSource struct {
	IpAsnRegexes types.List `tfsdk:"ip_asn_regexes"`
	Exclude      types.Bool `tfsdk:"exclude"`
}

type IpConnectionTypeSource struct {
	IpConnectionTypeList types.List `tfsdk:"ip_connection_type_list"`
	Exclude              types.Bool `tfsdk:"exclude"`
}

type UserIdSource struct {
	UserIdRegexes types.List `tfsdk:"user_id_regexes"`
	UserIds       types.List `tfsdk:"user_ids"`
	Exclude       types.Bool `tfsdk:"exclude"`
}

type AttributeCondition struct {
	KeyConditionOperator   types.String `tfsdk:"key_condition_operator"`
	KeyConditionValue      types.String `tfsdk:"key_condition_value"`
	ValueConditionOperator types.String `tfsdk:"value_condition_operator"`
	ValueConditionValue    types.String `tfsdk:"value_condition_value"`
}

type IpLocationTypeSource struct {
	IpLocationTypes types.List `tfsdk:"ip_location_types"`
	Exclude         types.Bool `tfsdk:"exclude"`
}

type IpAddressSource struct {
	IpAddressList types.List   `tfsdk:"ip_address_list"`
	Exclude       types.Bool   `tfsdk:"exclude"`
	IpAddressType types.String `tfsdk:"ip_address_type"`
}

type EmailDomainSource struct {
	EmailDomainRegexes types.List `tfsdk:"email_domain_regexes"`
	Exclude            types.Bool `tfsdk:"exclude"`
}

type UserAgentsSource struct {
	UserAgentsList types.List `tfsdk:"user_agents_list"`
	Exclude        types.Bool `tfsdk:"exclude"`
}

type RegionsSource struct {
	RegionsIds types.List `tfsdk:"regions_ids"`
	Exclude    types.Bool `tfsdk:"exclude"`
}

type IpOrganisationSource struct {
	IpOrganisationRegexes types.List `tfsdk:"ip_organisation_regexes"`
	Exclude               types.Bool `tfsdk:"exclude"`
}

type RequestResponseCondition struct {
	MetadataType  types.String `tfsdk:"metadata_type"`
	Value         types.String `tfsdk:"value"`
	KeyOperator   types.String `tfsdk:"key_operator"`
	KeyValue      types.String `tfsdk:"key_value"`
	ValueOperator types.String `tfsdk:"value_operator"`
}
