package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EndpointDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Endpoint Data Source",
		Attributes: map[string]schema.Attribute{
			"endpoints": schema.SetAttribute{
				MarkdownDescription: "Identifier of the Endpoint Rule",
				Required:            true,
				ElementType:         types.StringType,
			},
			"endpoint_ids": schema.SetAttribute{
				MarkdownDescription: "ID of the Endpoint Rule.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}
