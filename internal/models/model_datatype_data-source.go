package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataTypeDataModel struct {
	DataTypes   types.Set `tfsdk:"data_types"`
	DataTypeIds types.Set `tfsdk:"data_type_ids"`
}
