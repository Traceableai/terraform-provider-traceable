package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomSignatureModel struct {
	Id           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Environments types.Set      `tfsdk:"environments"`
	Description  types.String   `tfsdk:"description"`
	Disabled      types.Bool     `tfsdk:"disabled"`
	PayloadCriteria PayloadCriteriaModel `tfsdk:"payload_criteria"`
	Action       ActionModel   `tfsdk:"action"`
}


type PayloadCriteriaModel struct {
	RequestResponse []RequestResponseModel `tfsdk:"request_response"`
	CustomSecRule   types.String           `tfsdk:"custom_sec_rule"`
	Attributes      []AttributeConditionModel `tfsdk:"attributes"`
}

type RequestResponseModel struct {
	KeyValueTag   *types.String `tfsdk:"key_value_tag"`
	MatchCategory    types.String `tfsdk:"match_category"`
	KeyMatchOperator       *types.String `tfsdk:"key_match_operator"`
	MatchKey  types.String `tfsdk:"match_key"`
	ValueMatchOperator          types.String `tfsdk:"value_match_operator"`
	MatchValue          types.String `tfsdk:"match_value"`
}

type AttributeConditionModel struct {
	KeyConditionOperator   types.String `tfsdk:"key_condition_operator"`
	KeyConditionValue      types.String `tfsdk:"key_condition_value"`
	ValueConditionOperator *types.String `tfsdk:"value_condition_operator"`
	ValueConditionValue    *types.String `tfsdk:"value_condition_value"`
}

type ActionModel struct {
	ActionType       types.String `tfsdk:"action_type"`
	Duration         *types.String `tfsdk:"duration"`
	EventSeverity    *types.String `tfsdk:"event_severity"`
}

