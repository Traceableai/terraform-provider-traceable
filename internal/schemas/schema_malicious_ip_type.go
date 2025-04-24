package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MaliciousIPTypeResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Malicious IP Type Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the Malicious IP Type Rule",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Malicious IP Range Rule.",
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
			},
			"ip_type": schema.SetAttribute{
				MarkdownDescription: "IP ranges to apply the rule to(ANONYMOUS_VPN,BOT,PUBLIC_PROXY,TOR_EXIT_NODE,HOSTING_PROVIDER)",
				Required:            true,
				ElementType:         types.StringType,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action to take when the rule is triggered(ALERT,BLOCK)",
				Required:            true,
			},
		},
	}
}
