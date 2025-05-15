package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EndpointLabelsDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Endpoint Labels Data Source",
		Attributes: map[string]schema.Attribute{
			"labels": schema.SetAttribute{
				MarkdownDescription: "Identifier of the Endpoint Label Rule",
				Required:            true,
				ElementType:         types.StringType,
			},
			"label_ids": schema.SetAttribute{
				MarkdownDescription: "ID of the Endpoint Label Rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
