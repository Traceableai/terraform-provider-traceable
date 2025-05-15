package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WaapConfigModel struct {
	Environment types.String `tfsdk:"environment"`
	RuleConfigs types.Set    `tfsdk:"rule_configs"`
}

type WaapRuleConfigModel struct {
	RuleName types.String `tfsdk:"rule_name"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	Subrules types.Set    `tfsdk:"subrules"`
}

type WaapSubRuleConfigModel struct {
	SubRuleName   types.String `tfsdk:"name"`
	SubRuleAction types.String `tfsdk:"action"`
}
