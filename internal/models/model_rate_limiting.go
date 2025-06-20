package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RateLimitingRuleModel struct {
	Id               types.String         `tfsdk:"id"`
	Name             types.String         `tfsdk:"name"`
	Environments     types.Set            `tfsdk:"environments"`
	Description      types.String         `tfsdk:"description"`
	Enabled          types.Bool           `tfsdk:"enabled"`
	ThresholdConfigs types.Set            `tfsdk:"threshold_configs"`
	Action           RateLimitingAction   `tfsdk:"action"`
	Sources          *RateLimitingSources `tfsdk:"sources"`
}

type RateLimitingThresholdConfig struct {
	ApiAggregateType                     types.String `tfsdk:"api_aggregate_type"`
	UserAggregateType                    types.String `tfsdk:"user_aggregate_type"`
	RollingWindowCountAllowed            types.Int64  `tfsdk:"rolling_window_count_allowed"`
	RollingWindowDuration                types.String `tfsdk:"rolling_window_duration"`
	ThresholdConfigType                  types.String `tfsdk:"threshold_config_type"`
	DynamicMeanCalculationDuration       types.String `tfsdk:"dynamic_mean_calculation_duration"`
	DynamicDuration                      types.String `tfsdk:"dynamic_duration"`
	DynamicPercentageExcedingMeanAllowed types.Int64  `tfsdk:"dynamic_percentage_exceding_mean_allowed"`
	ValueType                            types.String `tfsdk:"value_type"`
	UniqueValuesAllowed                  types.Int64  `tfsdk:"unique_values_allowed"`
	Duration                             types.String `tfsdk:"duration"`
	SensitiveParamsEvaluationType        types.String `tfsdk:"sensitive_params_evaluation_type"`
}

type RateLimitingAction struct {
	ActionType       types.String `tfsdk:"action_type"`
	Duration         types.String `tfsdk:"duration"`
	EventSeverity    types.String `tfsdk:"event_severity"`
	HeaderInjections types.Set    `tfsdk:"header_injections"`
}

type RateLimitingHeaderInjection struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type RateLimitingSources struct {
	Scanner          *RateLimitingScannerSource          `tfsdk:"scanner"`
	IpAsn            *RateLimitingIpAsnSource            `tfsdk:"ip_asn"`
	IpConnectionType *RateLimitingIpConnectionTypeSource `tfsdk:"ip_connection_type"`
	UserId           *RateLimitingUserIdSource           `tfsdk:"user_id"`
	EndpointLabels   types.Set                           `tfsdk:"endpoint_labels"`
	Endpoints        types.Set                           `tfsdk:"endpoints"`
	IpReputation     types.String                        `tfsdk:"ip_reputation"`
	IpLocationType   *RateLimitingIpLocationTypeSource   `tfsdk:"ip_location_type"`
	IpAbuseVelocity  types.String                        `tfsdk:"ip_abuse_velocity"`
	IpAddress        *RateLimitingIpAddressSource        `tfsdk:"ip_address"`
	EmailDomain      *RateLimitingEmailDomainSource      `tfsdk:"email_domain"`
	UserAgents       *RateLimitingUserAgentsSource       `tfsdk:"user_agents"`
	Regions          *RateLimitingRegionsSource          `tfsdk:"regions"`
	IpOrganisation   *RateLimitingIpOrganisationSource   `tfsdk:"ip_organisation"`
	RequestResponse  types.Set                           `tfsdk:"request_response"`
	DataSet          types.Set                           `tfsdk:"data_set"`
	DataType         types.Set                           `tfsdk:"data_type"`
}

type RateLimitingScannerSource struct {
	ScannerTypesList types.Set  `tfsdk:"scanner_types_list"`
	Exclude          types.Bool `tfsdk:"exclude"`
}

type RateLimitingIpAsnSource struct {
	IpAsnRegexes types.Set  `tfsdk:"ip_asn_regexes"`
	Exclude      types.Bool `tfsdk:"exclude"`
}

type RateLimitingIpConnectionTypeSource struct {
	IpConnectionTypeList types.Set  `tfsdk:"ip_connection_type_list"`
	Exclude              types.Bool `tfsdk:"exclude"`
}

type RateLimitingUserIdSource struct {
	UserIdRegexes types.Set  `tfsdk:"user_id_regexes"`
	UserIds       types.Set  `tfsdk:"user_ids"`
	Exclude       types.Bool `tfsdk:"exclude"`
}

type RateLimitingIpLocationTypeSource struct {
	IpLocationTypes types.Set  `tfsdk:"ip_location_types"`
	Exclude         types.Bool `tfsdk:"exclude"`
}

type RateLimitingIpAddressSource struct {
	IpAddressList types.Set    `tfsdk:"ip_address_list"`
	Exclude       types.Bool   `tfsdk:"exclude"`
	IpAddressType types.String `tfsdk:"ip_address_type"`
}

type RateLimitingEmailDomainSource struct {
	EmailDomainRegexes types.Set  `tfsdk:"email_domain_regexes"`
	Exclude            types.Bool `tfsdk:"exclude"`
}

type RateLimitingUserAgentsSource struct {
	UserAgentsList types.Set  `tfsdk:"user_agents_list"`
	Exclude        types.Bool `tfsdk:"exclude"`
}

type RateLimitingRegionsSource struct {
	RegionsIds types.Set  `tfsdk:"regions_ids"`
	Exclude    types.Bool `tfsdk:"exclude"`
}

type RateLimitingIpOrganisationSource struct {
	IpOrganisationRegexes types.Set  `tfsdk:"ip_organisation_regexes"`
	Exclude               types.Bool `tfsdk:"exclude"`
}

type RateLimitingRequestResponseCondition struct {
	MetadataType  types.String `tfsdk:"metadata_type"`
	Value         types.String `tfsdk:"value"`
	KeyOperator   types.String `tfsdk:"key_operator"`
	KeyValue      types.String `tfsdk:"key_value"`
	ValueOperator types.String `tfsdk:"value_operator"`
}

type RateLimitingDataType struct {
	DataLocation types.String `tfsdk:"data_location"`
	DataTypesIds types.Set    `tfsdk:"data_types_ids"`
}

type RateLimitingDataSet struct {
	DataLocation types.String `tfsdk:"data_location"`
	DataSetsIds  types.Set    `tfsdk:"data_sets_ids"`
}
