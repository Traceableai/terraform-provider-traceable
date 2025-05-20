package modifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func SuppressDiffIfSanitizedEqual() planmodifier.String {
	return suppressDiffIfSanitizedEqual{}
}

type suppressDiffIfSanitizedEqual struct{}

func (m suppressDiffIfSanitizedEqual) Description(_ context.Context) string {
	return "Suppress diff if sanitized plan and state values are equal."
}

func (m suppressDiffIfSanitizedEqual) MarkdownDescription(_ context.Context) string {
	return m.Description(context.Background())
}

func (m suppressDiffIfSanitizedEqual) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.StateValue.IsNull() {
		return
	}

	sanitizedPlan := strings.TrimSpace(utils.EscapeString(req.PlanValue.ValueString()))
	sanitizedState := strings.TrimSpace(utils.EscapeString(req.StateValue.ValueString()))
	tflog.Info(ctx, "sanitizedPlan", map[string]interface{}{
		"sanitizedPlan": sanitizedPlan,
		"sanitizedState": sanitizedState,
	})
	if sanitizedPlan == sanitizedState {
		// Suppress diff by setting the planned value to the state value
		resp.PlanValue = req.StateValue
	}
}