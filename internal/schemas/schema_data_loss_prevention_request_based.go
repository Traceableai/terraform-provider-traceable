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
				MarkdownDescription: "Identifier of the DLP request based rule",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the DLP request based rule.",
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
				MarkdownDescription: "Description of the DLP request based rule",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable or Disable the DLP request based rule",
				Required:            true,
			},
			"action": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"action_type": schema.StringAttribute{
						MarkdownDescription: "ALERT , BLOCK , ALLOW",
						Required:            true,
					},
					"duration": schema.StringAttribute{
						MarkdownDescription: "how much time the action work (only allowed with ALLOW and BLOCK)",
						Optional:            true,
						Validators: []validator.String{
							validators.ValidDurationFormat(),
						},
						PlanModifiers: []planmodifier.String{
							modifiers.MatchStateIfDurationEqual(),
						},
					},
					"event_severity": schema.StringAttribute{
						MarkdownDescription: "LOW,MEDIUM,HIGH,CRITICAL (only allowed with BLOCK and ALERT)",
						Optional:            true,
					},
				},
			},
			"sources": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{	
					"service_scope": schema.SingleNestedAttribute{
						MarkdownDescription: "Service id where the rule will get applied",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"service_ids": schema.SetAttribute{
								MarkdownDescription: "It will be a list of service ids",
								Required:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"url_regex_scope": schema.SingleNestedAttribute{
						MarkdownDescription: "URL regex where the rule will get applied",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"url_regexes": schema.SetAttribute{
								MarkdownDescription: "It will be a list of url regexes",
								Required:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"regions": schema.SingleNestedAttribute{
						MarkdownDescription: "Regions as source, It will be a list region ids (AX,DZ)",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"region_ids": schema.SetAttribute{
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
							"data_type_matching": schema.SingleNestedAttribute{
								MarkdownDescription: "Datasets or datatypes ids",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"metadata_type": schema.StringAttribute{
										MarkdownDescription: "It can be (QUERY_PARAMETER/REQUEST_HEADER,REQUEST_COOKIE/REQUEST_BODY_PARAMETER/REQUEST_BODY)",
										Required:            true,
									},
									"operator": schema.StringAttribute{
										MarkdownDescription: "It Can be (EQUALS/MATCHES_REGEX) this is not required when metadata_type is REQUEST_BODY",
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
		},
	}
}
