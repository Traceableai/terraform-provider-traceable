package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type WaapConfigModel struct {
	Environment types.String `tfsdk:"environment"`
	RuleConfigs types.Set    `tfsdk:"rule_configs"`
}

type WaapRuleConfigModel struct {
	RuleName   types.String `tfsdk:"rule_name"`
	Disabled types.Bool   `tfsdk:"disabled"`
	Subrules types.Set    `tfsdk:"subrules"`
}

type WaapSubRuleConfigModel struct {
	SubRuleName types.String `tfsdk:"sub_rule_name"`
	SubRuleAction types.String `tfsdk:"sub_rule_action"`
}