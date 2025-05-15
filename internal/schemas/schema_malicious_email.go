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

func MaliciousEmailResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Malicious Email Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the Malicious Email Rule",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Malicious Email Rule.",
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
			"email_domain": schema.SetAttribute{
				MarkdownDescription: "Email domains to apply the rule to",
				Required:            true,
				ElementType:         types.StringType,
			},
			"email_domain_regexes": schema.SetAttribute{
				MarkdownDescription: "Email domain regexes to apply the rule to",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"email_fraud_score": schema.StringAttribute{
				MarkdownDescription: "Email fraud score(None/HIGH/CRITICAL)",
				Required:            true,
			},
			"data_leaked_email": schema.BoolAttribute{
				MarkdownDescription: "Users from leaked email domain",
				Required:            true,
			},
			"disposable_email_domain": schema.BoolAttribute{
				MarkdownDescription: "Users from disposable email domain",
				Required:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action to take when the rule is triggered(ALERT,BLOCK)",
				Required:            true,
			},
		},
	}
}
