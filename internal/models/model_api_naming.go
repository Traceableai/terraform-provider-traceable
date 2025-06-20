package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ApiNamingModel struct {
	Id              types.String         `tfsdk:"id"`
	Name            types.String         `tfsdk:"name"`
	Disabled        types.Bool           `tfsdk:"disabled"`
	ServiceNames    types.Set            `tfsdk:"service_names"`
	EnvironmentNames types.Set            `tfsdk:"environment_names"`
	Regexes         types.Set            `tfsdk:"regexes"`
	Values          types.Set            `tfsdk:"values"`
}
