package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AgentTokenModel struct {
	Id              types.String         `tfsdk:"id"`
	Name            types.String         `tfsdk:"name"`
}

