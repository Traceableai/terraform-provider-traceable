package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DatasetDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Dataset Data Source",
		Attributes: map[string]schema.Attribute{
			"data_sets": schema.SetAttribute{
				MarkdownDescription: "Identifier of the Dataset Rule",
				Required:            true,
				ElementType:         types.StringType,
			},
			"data_set_ids": schema.SetAttribute{
				MarkdownDescription: "ID of the Dataset Rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
