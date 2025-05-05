package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EndpointLabelsDataModel struct {
	Labels   types.Set `tfsdk:"labels"`
	LabelIds types.Set `tfsdk:"label_ids"`
}
