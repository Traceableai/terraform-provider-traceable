package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataLossPreventionRequestBasedRuleModel struct {
	Id               types.String         `tfsdk:"id"`
	Name             types.String         `tfsdk:"name"`
	Environments     types.Set            `tfsdk:"environments"`
	Description      types.String         `tfsdk:"description"`
	Enabled          types.Bool           `tfsdk:"enabled"`
	Action           DlpAction   `tfsdk:"action"`
	Sources          DlpRequestBasedSources `tfsdk:"sources"`
}

type DlpAction struct {
	ActionType       types.String `tfsdk:"action_type"`
	Duration         types.String `tfsdk:"duration"`
	EventSeverity    types.String `tfsdk:"event_severity"`
}

type DlpRequestBasedSources struct {
	IpLocationType   DlpRequestBasedIpLocationTypeSource   `tfsdk:"ip_location_type"`
	Regions          DlpRequestBasedRegionsSource          `tfsdk:"regions"`
	IpAddress        DlpRequestBasedIpAddressSource        `tfsdk:"ip_address"`
	RequestPayload   types.Set                           `tfsdk:"request_payload"`
	DataSetDataType  DlpRequestBasedDataSetDataTypeFilterSource        `tfsdk:"dateset_datatype_filter"`
	ServiceScope     DlpRequestBasedServiceScope           `tfsdk:"service_scope"`
	UrlScope         DlpRequestBasedUrlScope               `tfsdk:"url_regex_scope"`
}

type DlpRequestBasedRequestPayloadScope struct {
	MetadataType types.String `tfsdk:"metadata_type"`
	KeyOperator types.String `tfsdk:"key_operator"`
	KeyValue types.String `tfsdk:"key_value"`
	ValueOperator types.String `tfsdk:"value_operator"`
	Value types.String `tfsdk:"value"`
}
type DlpRequestBasedServiceScope struct {
	ServiceIds types.Set    `tfsdk:"service_ids"`
}
type DlpRequestBasedUrlScope struct {
	UrlRegexes types.Set    `tfsdk:"url_regexes"`
}

type DlpRequestBasedIpAddressSource struct {
	IpAddressList types.Set  `tfsdk:"ip_address_list"`
}

type DlpRequestBasedRegionsSource struct {
	RegionIds types.Set  `tfsdk:"region_ids"`
}

type DlpRequestBasedIpLocationTypeSource struct {
	IpLocationTypes types.Set  `tfsdk:"ip_location_types"`
}

type DlpRequestBasedDataSetDataTypeFilterSource struct {
	DataSetDataTypeIds DataSetDataTypeIds  `tfsdk:"dateset_datatype_id"`
	DataTypeMatching DlpReqBasedDataTypeMatching  `tfsdk:"data_type_matching"`
}

type DlpReqBasedDataTypeMatching struct {
	MetadataType types.String `tfsdk:"metadata_type"`
	Operator types.String `tfsdk:"operator"`
	Value types.String `tfsdk:"value"`
}

type DataSetDataTypeIds struct {
	DataSetsIds types.Set `tfsdk:"data_sets_ids"`
	DataTypesIds types.Set `tfsdk:"data_types_ids"`
}