package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MaliciousEmailDomainModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Environments                      types.Set    `tfsdk:"environments"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
	EventSeverity                     types.String `tfsdk:"event_severity"`
	Duration                          types.String `tfsdk:"duration"`
	EmailDomainsList                  types.Set    `tfsdk:"email_domains_list"`
	EmailRegexesList                  types.Set    `tfsdk:"email_regexes_list"`
	MinEmailFraudScoreLevel           types.String `tfsdk:"min_email_fraud_score_level"`
	Action                            types.String `tfsdk:"action"`
	ApplyRuleToDataLeakedEmail        types.Bool   `tfsdk:"apply_rule_to_data_leaked_email"`
	ApplyRuleToDisposableEmailDomains types.Bool   `tfsdk:"apply_rule_to_disposable_email_domains"`
}
