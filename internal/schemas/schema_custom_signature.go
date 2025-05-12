package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CustomSignatureResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a custom signature rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the custom signature",
				Computed:            true,
				MarkdownDescription: "Identifier of the Custom Signature Rule",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the custom signature",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environments": schema.SetAttribute{
				MarkdownDescription: "Environments the rule is applicable to",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the custom signature",
				Optional:            true,
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the custom signature rule",
				Required:            true,
			},
			"payload_criteria": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"request_response": schema.SetNestedAttribute{
						MarkdownDescription: "Request/response conditions as payload match criteria",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key_value_tag": schema.StringAttribute{
									MarkdownDescription: "Which Metadatype to include",
									Optional:            true,
								},
								"match_category": schema.StringAttribute{
									MarkdownDescription: "Which operator to use",
									Required:            true,
								},
								"key_match_operator": schema.StringAttribute{
									MarkdownDescription: "Value to match",
									Optional:            true,
								},
								"match_key": schema.StringAttribute{
									MarkdownDescription: "Which operator to use",
									Required:            true,
								},
								"value_match_operator": schema.StringAttribute{
									MarkdownDescription: "Value to match",
									Required:            true,
								},
								"match_value": schema.StringAttribute{
									MarkdownDescription: "Value to match",
									Required:            true,
								},
							},
						},
					},
					"custom_sec_rule": schema.StringAttribute{
						MarkdownDescription: "custom sec rule string",
						Optional:            true,
					},
					"attributes": schema.SetNestedAttribute{
						MarkdownDescription: "Attributes conditions as payload match criteria",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key_condition_operator": schema.StringAttribute{
									MarkdownDescription: "Which operator to include",
									Required:            true,
								},
								"key_condition_value": schema.StringAttribute{
									MarkdownDescription: "Value for key operator match criteria",
									Required:            true,
								},
								"value_condition_operator": schema.StringAttribute{
									MarkdownDescription: "Value operator to use",
									Optional:            true,
								},
								"value_condition_value": schema.StringAttribute{
									MarkdownDescription: "value to use",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"action": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"action_type": schema.StringAttribute{
						MarkdownDescription: "ALERT , BLOCK , ALLOW ,MARK FOR TESTING",
						Required:            true,
					},
					"duration": schema.StringAttribute{
						MarkdownDescription: "how much time the action work",
						Optional:            true,
					},
					"event_severity": schema.StringAttribute{
						MarkdownDescription: "LOW,MEDIUM,HIGH",
						Optional:            true,
					},
				},
			},
		},
	}
}
