package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MaliciousIpTypeModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Environments  types.Set    `tfsdk:"environments"`
	Description   types.String `tfsdk:"description"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	EventSeverity types.String `tfsdk:"event_severity"`
	Duration      types.String `tfsdk:"duration"`
	IpType        types.Set    `tfsdk:"ip_type"`
	Action        types.String `tfsdk:"action"`
}
