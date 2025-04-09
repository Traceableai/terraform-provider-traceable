package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func RateLimitingResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Rate Limiting Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the Rate Limiting Rule",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Rate Limiting Rule.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environments": schema.ListAttribute{
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
			"threshold_configs": schema.ListNestedAttribute{
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
						},
						"threshold_config_type": schema.StringAttribute{
							MarkdownDescription: "Threshold config type",
							Required:            true,
						},
						"dynamic_mean_calculation_duration": schema.StringAttribute{
							MarkdownDescription: "Dynamic mean calculation duration",
							Optional:            true,
						},
						"dynamic_duration": schema.StringAttribute{
							MarkdownDescription: "Dynamic duration",
							Optional:            true,
						},
						"dynamic_percentage_exceding_mean_allowed": schema.Int64Attribute{
							MarkdownDescription: "Dynamic percentage exceeding mean allowed",
							Optional:            true,
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
					"header_injections": schema.ListNestedAttribute{
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

					"scanner": schema.SingleNestedAttribute{
						MarkdownDescription: "Scanner as source, It will be a list of scanner type",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"scanner_types_list": schema.ListAttribute{
								MarkdownDescription: "It will be a list of scanner types",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given scaner types",
								Required:            true,
							},
						},
					},
					"ip_asn": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"ip_asn_regexes": schema.ListAttribute{
								MarkdownDescription: "It will be a list of IP ASNs",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given IP ASN",
								Required:            true,
							},
						},
					},
					"ip_connection_type": schema.SingleNestedAttribute{
						MarkdownDescription: "Ip connection type as source, It will be a list of ip connection type",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ip_connection_type_list": schema.ListAttribute{
								MarkdownDescription: "It will be a list of IP connection types",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given IP connection",
								Required:            true,
							},
						},
					},
					"user_id": schema.SingleNestedAttribute{
						MarkdownDescription: "User id as source",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"user_id_regexes": schema.ListAttribute{
								MarkdownDescription: "It will be a list of user id regexes",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"user_ids": schema.ListAttribute{
								MarkdownDescription: "List of user ids",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given user id",
								Required:            true,
							},
						},
					},
					"endpoint_labels": schema.ListAttribute{
						MarkdownDescription: "Filter endpoints by labels you want to apply this rule",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"endpoints": schema.ListAttribute{
						MarkdownDescription: "List of endpoint ids",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"attribute": schema.ListNestedAttribute{
						MarkdownDescription: "Attribute based conditions for the rule",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{

							Attributes: map[string]schema.Attribute{
								"key_condition_operator": schema.StringAttribute{
									MarkdownDescription: "Key operator",
									Required:            true,
								},
								"key_condition_value": schema.StringAttribute{
									MarkdownDescription: "Key value",
									Required:            true,
								},
								"value_condition_operator": schema.StringAttribute{
									MarkdownDescription: "Value operator",
									Required:            true,
								},
								"value_condition_value": schema.StringAttribute{
									MarkdownDescription: "Value value",
									Required:            true,
								},
							},
						},
					},

					"ip_reputation": schema.StringAttribute{
						MarkdownDescription: "Ip reputation source (LOW/MEDIUM/HIGH/CRITICAL)",
						Optional:            true,
					},
					"ip_location_type": schema.SingleNestedAttribute{
						MarkdownDescription: "Ip location type as source ([BOT, TOR_EXIT_NODE, PUBLIC_PROXY])",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ip_location_types": schema.ListAttribute{
								MarkdownDescription: "Ip location type as source ([BOT, TOR_EXIT_NODE, PUBLIC_PROXY])",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given ip location types",
								Required:            true,
							},
						},
					},
					"ip_abuse_velocity": schema.StringAttribute{
						MarkdownDescription: "Ip abuse velocity as source (LOW/MEDIUM/HIGH/CRITICAL)",
						Optional:            true,
					},
					"ip_address": schema.SingleNestedAttribute{
						MarkdownDescription: "Ip address as source (LIST_OF_IP's/ALL_EXTERNAL)",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ip_address_list": schema.ListAttribute{
								MarkdownDescription: "List of ip addresses",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given ip addresses",
								Required:            true,
							},
							"ip_address_type": schema.StringAttribute{
								MarkdownDescription: "Accepts ALL_EXTERNAL",
								Optional:            true,
							},
						},
					},
					"email_domain": schema.SingleNestedAttribute{
						MarkdownDescription: "Email domain as source, It will be a list of email domain regexes",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"email_domain_regexes": schema.ListAttribute{
								MarkdownDescription: "It will be a list of email domain regexes",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given email domains regexes",
								Required:            true,
							},
						},
					},
					"user_agents": schema.SingleNestedAttribute{
						MarkdownDescription: "User agents as source, It will be a list of user agents",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"user_agents_list": schema.ListAttribute{
								MarkdownDescription: "It will be a list of user agents",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given user agents",
								Required:            true,
							},
						},
					},
					"regions": schema.SingleNestedAttribute{
						MarkdownDescription: "Regions as source, It will be a list region ids (AX,DZ)",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"regions_ids": schema.ListAttribute{
								MarkdownDescription: "It will be a list of regions ids in countryIsoCode",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given regions",
								Required:            true,
							},
						},
					},
					"ip_organisation": schema.SingleNestedAttribute{
						MarkdownDescription: "Ip organisation as source, It will be a list of ip organisation",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"ip_organisation_regexes": schema.ListAttribute{
								MarkdownDescription: "It will be a list of ip organisations",
								Required:            true,
								ElementType:         types.StringType,
							},
							"exclude": schema.BoolAttribute{
								MarkdownDescription: "Set it to true to exclude given ip organisation",
								Required:            true,
							},
						},
					},
					"request_response": schema.ListNestedAttribute{
						MarkdownDescription: "Request/response attributes as source",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"metadata_type": schema.StringAttribute{
									MarkdownDescription: "Which Metadatype to include",
									Required:            true,
								},
								"key_operator": schema.StringAttribute{
									MarkdownDescription: "Which operator to use",
									Optional:            true,
								},
								"key_value": schema.StringAttribute{
									MarkdownDescription: "Value to match",
									Optional:            true,
								},
								"value_operator": schema.StringAttribute{
									MarkdownDescription: "Which operator to use",
									Optional:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Value to match",
									Optional:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
