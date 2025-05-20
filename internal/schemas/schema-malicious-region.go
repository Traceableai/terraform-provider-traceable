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

func MaliciousRegionResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Malicious Region Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the Malicious Region Rule",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Malicious Region Rule.",
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

			"event_severity": schema.StringAttribute{
				MarkdownDescription: "Event severity of the rule (LOW, MEDIUM, HIGH, CRITICAL)",
				Optional:            true,
			},

			"duration": schema.StringAttribute{
				MarkdownDescription: "Duration of the rule",
				Optional:            true,
				Validators: []validator.String{
					validators.ValidDurationFormat(),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.MatchStateIfDurationEqual(),
				},
			},
			"regions": schema.SetAttribute{
				MarkdownDescription: "Regions to apply the rule to(Please check documentation for the valid regions)",
				Required:            true,
				ElementType:         types.StringType,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action to take when the rule is triggered(ALERT,BLOCK,BLOCK_ALL_EXCEPT)",
				Required:            true,
			},
		},
	}
}
