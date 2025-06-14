package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/traceableai/terraform-provider-traceable/internal/modifiers"
	"github.com/traceableai/terraform-provider-traceable/internal/validators"
)

func CustomSignatureResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a custom signature rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the custom signature",
				Computed:            true,
				MarkdownDescription: "Identifier of the Custom Signature Rule",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
									MarkdownDescription: "Which multi valued api attribute to include it can be (HEADER/PARAMETER/COOKIE)",
									Optional:            true,
								},
								"match_category": schema.StringAttribute{
									MarkdownDescription: "Where to find api attribute (REQUEST/RESPONSE)",
									Required:            true,
								},
								"key_match_operator": schema.StringAttribute{
									MarkdownDescription: "Operator to use for key match with key_value_tag (EQUALS/NOT_EQUAL/MATCHES_REGEX/NOT_MATCH_REGEX/CONTAINS/NOT_CONTAIN/GREATER_THAN/LESS_THAN). All operarots are not valid with certain metadeta type.",
									Optional:            true,
								},
								"match_key": schema.StringAttribute{
									MarkdownDescription: "Which single valued api attribute to include (URL/HEADER_NAME/HEADER_VALUE/PARAMETER_NAME/PARAMETER_VALUE/HTTP_METHOD/HOST/USER_AGENT/STATUS_CODE/BODY/BODY_SIZE/COOKIE_NAME/COOKIE_VALUE/QUERY_PARAMS_COUNT/HEADERS_COUNT/COOKIES_COUNT)",
									Required:            true,
								},
								"value_match_operator": schema.StringAttribute{
									MarkdownDescription: "Operator to use for value match. Accepts same operators as key_match_operator",
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
						MarkdownDescription: "Custom security rule string",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							modifiers.SuppressDiffIfCustomSecRuleEqual(),
						},
					},
					"attributes": schema.SetNestedAttribute{
						MarkdownDescription: "Attributes conditions as payload match criteria",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key_condition_operator": schema.StringAttribute{
									MarkdownDescription: "Which operator to include. Same operators as other feilds.",
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
						MarkdownDescription: "how much time the action work (PT1M,PT5M)",
						Optional:            true,
						Validators: []validator.String{
							validators.ValidDurationFormat(),
						},
						PlanModifiers: []planmodifier.String{
							modifiers.MatchStateIfDurationEqual(),
						},
					},
					"event_severity": schema.StringAttribute{
						MarkdownDescription: "LOW,MEDIUM,HIGH,CRITICAL",
						Optional:            true,
					},
				},
			},
		},
	}
}
