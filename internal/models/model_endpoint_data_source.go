package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EndpointDataModel struct {
	Endpoints   types.Set `tfsdk:"endpoints"`
	EndpointIds types.Set `tfsdk:"endpoint_ids"`
}
