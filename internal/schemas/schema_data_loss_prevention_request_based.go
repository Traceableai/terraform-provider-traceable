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

func DataLossPreventionRequestBasedResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Data Loss Prevention Request Based Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the Rate Limiting Rule",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Rate Limiting Rule.",
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
				MarkdownDescription: "Description of the Rate Limiting Rule",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the Rate Limiting Rule",
				Required:            true,
			},
			"threshold_configs": schema.SetNestedAttribute{
				MarkdownDescription: "Threshold configs for the rule",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_aggregate_type": schema.StringAttribute{
							MarkdownDescription: "API aggregate type",
							Optional:            true,
						},
						"user_aggregate_type": schema.StringAttribute{
							MarkdownDescription: "User aggregate type",
							Optional:            true,
						},
						"rolling_window_count_allowed": schema.Int64Attribute{
							MarkdownDescription: "Rolling window count allowed",
							Optional:            true,
						},
						"rolling_window_duration": schema.StringAttribute{
							MarkdownDescription: "Rolling window duration",
							Optional:            true,
							Validators: []validator.String{
								validators.ValidDurationFormat(),
							},
							PlanModifiers: []planmodifier.String{
								modifiers.MatchStateIfDurationEqual(),
							},
						},
						"threshold_config_type": schema.StringAttribute{
							MarkdownDescription: "Threshold config type",
							Required:            true,
						},
						"dynamic_mean_calculation_duration": schema.StringAttribute{
							MarkdownDescription: "Dynamic mean calculation duration",
							Optional:            true,
							Validators: []validator.String{
								validators.ValidDurationFormat(),
							},
							PlanModifiers: []planmodifier.String{
								modifiers.MatchStateIfDurationEqual(),
							},
						},
						"dynamic_duration": schema.StringAttribute{
							MarkdownDescription: "Dynamic duration",
							Optional:            true,
							Validators: []validator.String{
								validators.ValidDurationFormat(),
							},
							PlanModifiers: []planmodifier.String{
								modifiers.MatchStateIfDurationEqual(),
							},
						},
						"dynamic_percentage_exceding_mean_allowed": schema.Int64Attribute{
							MarkdownDescription: "Dynamic percentage exceeding mean allowed",
							Optional:            true,
						},
						"value_type": schema.StringAttribute{
							MarkdownDescription: "Value type",
							Optional:            true,
						},
						"unique_values_allowed": schema.Int64Attribute{
							MarkdownDescription: "Unique values allowed",
							Optional:            true,
						},
						"sensitive_params_evaluation_type": schema.StringAttribute{
							MarkdownDescription: "Sensitive params evaluation type",
							Optional:            true,
						},
						"duration": schema.StringAttribute{
							MarkdownDescription: "Duration",
							Optional:            true,
							Validators: []validator.String{
								validators.ValidDurationFormat(),
							},
							PlanModifiers: []planmodifier.String{
								modifiers.MatchStateIfDurationEqual(),
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
					"header_injections": schema.SetNestedAttribute{
						MarkdownDescription: "Header fields to be injected",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									MarkdownDescription: "The header field name to inject (e.g., 'X-Custom-Header')",
									Optional:            true, // Make key optional
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "The value to set for the header field",
									Optional:            true,
								},
							},
						},
					},
				},
			},
			"sources": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{	
					"regions": schema.SingleNestedAttribute{
						MarkdownDescription: "Regions as source, It will be a list region ids (AX,DZ)",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"regions_ids": schema.SetAttribute{
								MarkdownDescription: "It will be a list of regions ids in countryIsoCode",
								Required:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"ip_location_type": schema.SingleNestedAttribute{
						MarkdownDescription: "Ip location type as source",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ip_location_types": schema.SetAttribute{
								MarkdownDescription: "Ip location type as source ([BOT,ANONYMOUS_VPN,HOSTING_PROVIDER,TOR_EXIT_NODE, PUBLIC_PROXY,SCANNER])",
								Required:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"ip_address": schema.SingleNestedAttribute{
						MarkdownDescription: "Ip address as source (LIST_OF_IP's/ALL_EXTERNAL)",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ip_address_list": schema.SetAttribute{
								MarkdownDescription: "List of ip addresses",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"request_payload": schema.SetNestedAttribute{
						MarkdownDescription: "Request payload attributes as source",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"metadata_type": schema.StringAttribute{
									MarkdownDescription: "It can be (QUERY_PARAMETER,REQUEST_HEADER,REQUEST_COOKIE,REQUEST_BODY_PARAMETER,REQUEST_BODY,HTTP_METHOD,USER_AGENT,HOST,REQUEST_BODY_SIZE,QUERY_PARAMS_COUNT,REQUEST_COOKIES_COUNT ,REQUEST_HEADERS_COUNT)",
									Required:            true,
								},
								"key_operator": schema.StringAttribute{
									MarkdownDescription: "it can be (EQUALS,NOT_EQUAL,MATCHES_REGEX,NOT_MATCH_REGEX)",
									Optional:            true,
								},
								"key_value": schema.StringAttribute{
									MarkdownDescription: "key value to match",
									Optional:            true,
								},
								"value_operator": schema.StringAttribute{
									MarkdownDescription: "it can be (EQUALS,NOT_EQUAL,MATCHES_REGEX,NOT_MATCH_REGEX,CONTAINS,NOT_CONTAIN,GREATER_THAN,LESS_THAN). All operators are not supported with certain metadata_type",
									Required:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "value to match",
									Required:            true,
								},
							},
						},
					},
					"dateset_datatype_filter": schema.SingleNestedAttribute{
						MarkdownDescription: "Datasets or datatypes source configuration",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"dateset_datatype_id": schema.SingleNestedAttribute{
								MarkdownDescription: "Datasets or datatypes ids",
								Required:            true,
								Attributes: map[string]schema.Attribute{
									"data_sets_ids": schema.SetAttribute{
										MarkdownDescription: "IDs of datasets to match",
										Optional:            true,
										ElementType:         types.StringType,
									},
									"data_types_ids": schema.SetAttribute{
										MarkdownDescription: "IDs of datatypes to match",
										Optional:            true,
										ElementType:         types.StringType,
									},
								},
							},
							"data_type_matching_metadata_type": schema.StringAttribute{
								MarkdownDescription: "It can be (QUERY_PARAMETER/REQUEST_HEADER,REQUEST_COOKIE/REQUEST_BODY_PARAMETER/REQUEST_BODY)",
								Optional:            true,
							},
							"data_type_matching_operator": schema.StringAttribute{
								MarkdownDescription: "It Can be (EQUALS/MATCHES_REGEX)",
								Optional:            true,
							},
							"data_type_matching_value": schema.StringAttribute{
								MarkdownDescription: "Value to match",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}
