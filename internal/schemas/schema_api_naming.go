package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ApiNamingResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Manages a API naming rule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the API naming rule",
				Computed:            true,
				MarkdownDescription: "Identifier of the API Naming Rule",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the API naming rule",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: "Flag to enable or disable the rule",
				Required:    true,
			},
			"regexes": schema.SetAttribute{
				Description: "List of regex patterns for the rule",
				Required:    true,
				ElementType: types.StringType,
			},
			"values": schema.SetAttribute{
				Description: "Corresponding values for the regex patterns",
				Required:    true,
				ElementType: types.StringType,
			},
			"service_names": schema.SetAttribute{
				Description: "List of service names to apply the rule",
				Required:    true,
				ElementType: types.StringType,
			},
			"environment_names": schema.SetAttribute{
				Description: "List of environment names to apply the rule",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}
