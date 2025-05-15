package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServiceDataModel struct {
	Services   types.Set `tfsdk:"services"`
	ServiceIds types.Set `tfsdk:"service_ids"`
}
