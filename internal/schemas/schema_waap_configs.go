package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func WaapConfigResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a waap config to enable/disable",
		Attributes: map[string]schema.Attribute{
			"environment": schema.StringAttribute{
				MarkdownDescription: "Environment the rule is applicable to",
				Optional:            true,
			},
			"rule_configs": schema.SetNestedAttribute{
				MarkdownDescription: "List of WAAF rule configurations.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_name": schema.StringAttribute{
							MarkdownDescription: "Identifier of the underlying WAAF rule.",
							Required:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the rule is enabled (true) or disabled (false).",
							Required:            true,
						},
						"subrules": schema.SetNestedAttribute{
							MarkdownDescription: "List of sub rule configurations.",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: "Identifier of the underlying sub rule.",
										Required:            true,
									},
									"action": schema.StringAttribute{
										MarkdownDescription: "Whether the sub rule is MONITOR/DISABLED/BLOCK.",
										Required:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
