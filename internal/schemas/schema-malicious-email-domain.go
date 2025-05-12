package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MaliciousEmailDomainResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Traceable Malicious Email Domain Resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the Malicious Email Domain Rule",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Malicious Email Domain Rule.",
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
				MarkdownDescription: "Description of the Malicious Email Domain Rule",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable the Malicious Email Domain Rule",
				Required:            true,
			},
			"apply_rule_to_data_leaked_email": schema.BoolAttribute{
				MarkdownDescription: "Apply the rule to data leaked email",
				Optional:            true,
			},
			"apply_rule_to_disposable_email_domains": schema.BoolAttribute{
				MarkdownDescription: "Apply the rule to disposable email domains",
				Optional:            true,
			},
			"email_domains_list": schema.SetAttribute{
				MarkdownDescription: "Enter a list of email domains",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"email_regexes_list": schema.SetAttribute{
				MarkdownDescription: "Enter a list of email regexes",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"min_email_fraud_score_level": schema.StringAttribute{
				MarkdownDescription: "Email Fraud Score of the rule (NONE, HIGH, CRITICAL)",
				Optional:            true,
			},
			"event_severity": schema.StringAttribute{
				MarkdownDescription: "Event severity of the rule (LOW, MEDIUM, HIGH, CRITICAL)",
				Optional:            true,
			},
			"duration": schema.StringAttribute{
				MarkdownDescription: "Duration of the rule",
				Optional:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action to take when the rule is triggered (ALERT,BLOCK)",
				Required:            true,
			},
		},
	}
}
