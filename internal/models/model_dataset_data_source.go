package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DatasetDataModel struct {
	DataSets   types.Set `tfsdk:"data_sets"`
	DataSetIds types.Set `tfsdk:"data_set_ids"`
}
