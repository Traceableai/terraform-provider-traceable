package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func MatchStateIfDurationEqual() planmodifier.String {
	return matchStateIfDurationEqual{}
}

type matchStateIfDurationEqual struct{}

func (m matchStateIfDurationEqual) Description(_ context.Context) string {
	return "Suppress diff if plan and state duration are equivalent."
}

func (m matchStateIfDurationEqual) MarkdownDescription(_ context.Context) string {
	return m.Description(context.Background())
}

func (m matchStateIfDurationEqual) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {

	if req.PlanValue.IsUnknown() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.StateValue.IsNull() {
		return
	}

	planDur, err := utils.ConvertDurationToSeconds(req.PlanValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid format of giving time", err.Error())
		return
	}
	stateDur, _ := utils.ConvertDurationToSeconds(req.StateValue.ValueString())

	if planDur == stateDur {
		resp.PlanValue = req.StateValue
	}
}
