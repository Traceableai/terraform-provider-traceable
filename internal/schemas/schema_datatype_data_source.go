package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DataTypeDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Data Type Data Source",
		Attributes: map[string]schema.Attribute{
			"data_types": schema.SetAttribute{
				MarkdownDescription: "Identifier of the Data Type Rule",
				Required:            true,
				ElementType:         types.StringType,
			},
			"data_type_ids": schema.SetAttribute{
				MarkdownDescription: "ID of the Data Type Rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
