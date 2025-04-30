package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ServiceDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Service Data Source",
		Attributes: map[string]schema.Attribute{
			"services": schema.SetAttribute{
				MarkdownDescription: "Identifier of the Service",
				Required:            true,
				ElementType:         types.StringType,
			},
			"service_ids": schema.SetAttribute{
				MarkdownDescription: "ID of the Service.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
