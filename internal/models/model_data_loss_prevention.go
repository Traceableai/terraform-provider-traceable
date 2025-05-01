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
	ThresholdConfigs types.Set            `tfsdk:"threshold_configs"`
	Action           RateLimitingAction   `tfsdk:"action"`
	Sources          *RateLimitingSources `tfsdk:"sources"`
}
