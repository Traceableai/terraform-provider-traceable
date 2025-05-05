package resources

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/traceableai/terraform-provider-traceable/internal/generated"
	"github.com/traceableai/terraform-provider-traceable/internal/models"
	"github.com/traceableai/terraform-provider-traceable/internal/schemas"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func NewMaliciousIpTypeResource() resource.Resource {
	return &MaliciousIpTypeResource{}
}

type MaliciousIpTypeResource struct {
	client *graphql.Client
}
type MaliciousIpTypeResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IconType    types.String `tfsdk:"icon_type"`
}

func (r *MaliciousIpTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_malicious_ip_type"
}

func (r *MaliciousIpTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.MaliciousIPTypeResourceSchema()
}

func (r *MaliciousIpTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
	tflog.Trace(ctx, "Client Intialization Successfully")
}

func (r *MaliciousIpTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "Entering in Create Block")
	var data *models.MaliciousIpTypeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	input, err := convertMaliciousIpTypeModelToCreateInput(ctx, data)
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	id, err := getMaliciousIPTypeRuleId(input.RuleInfo.Name, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id != "" {
		resp.Diagnostics.AddError("Resource already Exist ", fmt.Sprintf("%s malicious ip type rule already please try with different name or import it", input.RuleInfo.Name))
		return
	}
	response, err := generated.CreateMaliciousIpTypeRule(ctx, *r.client, *input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(response.CreateMaliciousSourcesRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Create Block")

}
func (r *MaliciousIpTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "Entering in Read Block")
	var data *models.MaliciousIpTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := data.Name.ValueString()
	rule, err := getMaliciousIPTypeRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	ruleData, err := convertMaliciousIpTypeFieldsToModel(ctx, rule)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &ruleData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Read Block")

}
func (r *MaliciousIpTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Entering in Update Block")
	var dataState *models.MaliciousIpTypeModel
	var data *models.MaliciousIpTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := dataState.Name.ValueString()
	rule, err := getMaliciousIPTypeRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	input, err := convertMaliciousIpTypeModelToUpdateInput(ctx, data, dataState.Id.ValueString())
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	_, err = generated.UpdateMaliciousIpTypeRule(ctx, *r.client, *input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(input.Rule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Update Block")

}
func (r *MaliciousIpTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "Entering in ImportState Block")
	ruleName := req.ID
	id, err := getMaliciousIPTypeRuleId(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Resource Not Found", fmt.Sprintf("%s rule of this name not found", ruleName))
		return
	}
	response, err := getMaliciousIPTypeRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data, err := convertMaliciousIpTypeFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Exiting in ImportState Block")

}

func (r *MaliciousIpTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "Entering in Delete Block")
	var data *models.MaliciousIpTypeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := data.Name.ValueString()
	rule, err := getMaliciousIPTypeRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	input := generated.InputMaliciousSourcesRuleDelete{
		Id: rule.Id,
	}
	_, err = generated.DeleteMaliciousIpTypeRule(ctx, *r.client, input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "Exiting in Delete Block")

}

func getMaliciousIPTypeRuleId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
	input := &generated.InputMaliciousSourcesRulesFilter{}

	response, err := generated.GetMaliciousIpTypeRulesName(ctx, r, input)
	if err != nil {
		return "", err
	}
	for _, rule := range response.MaliciousSourcesRules.Results {
		if rule.Info.Name == ruleName {
			return rule.Id, nil
		}
	}
	return "", nil
}

func getMaliciousIPTypeRule(ruleName string, ctx context.Context, r graphql.Client) (*generated.MaliciousIpTypeRuleFields, error) {
	input := &generated.InputMaliciousSourcesRulesFilter{}

	response, err := generated.GetMaliciousIpTypeRuleDetails(ctx, r, input)
	if err != nil {
		return nil, err
	}
	for _, rule := range response.MaliciousSourcesRules.Results {
		if rule.Info.Name == ruleName {
			return &rule.MaliciousIpTypeRuleFields, nil
		}
	}
	return nil, nil
}

func convertMaliciousIpTypeModelToCreateInput(ctx context.Context, data *models.MaliciousIpTypeModel) (*generated.InputMaliciousSourcesRuleCreate, error) {
	var input = generated.InputMaliciousSourcesRuleCreate{}
	var ruleInfo = generated.InputMaliciousSourcesRuleInfo{}
	var ruleAction = generated.InputMaliciousSourcesRuleAction{}

	if HasValue(data.Name) {
		ruleInfo.Name = data.Name.ValueString()
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		ruleInfo.Description = &description
	}

	if HasValue(data.EventSeverity) {
		eventSeverity, exist := MaliciousIpTypeEventSeverityMap[data.EventSeverity.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("event_severity", "Invalid EventSeverity")
		}
		ruleAction.EventSeverity = &eventSeverity
	}
	if HasValue(data.Action) {
		action, exist := MaliciousIpTypeActionMap[data.Action.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("action", "Invalid action")
		}
		ruleAction.RuleActionType = action
	} else {
		return nil, utils.NewInvalidError("action", "Action field must be present and must not be empty")
	}

	if HasValue(data.Duration) {
		duration := data.Duration.ValueString()
		expirationDetails := generated.InputMaliciousSourcesRuleExpirationDetails{
			ExpirationDuration: duration,
		}
		ruleAction.ExpirationDetails = &expirationDetails
	}

	if HasValue(data.Environments) {
		var ruleScope = generated.InputMaliciousSourcesRuleScope{}
		var environmentScope = generated.InputMaliciousSourcesRuleEnvironmentScope{}
		environments, err := utils.ConvertSetToStrPointer(data.Environments)
		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")
		}
		environmentScope.EnvironmentIds = environments
		ruleScope.EnvironmentScope = &environmentScope
		input.RuleScope = &ruleScope
	} else {
		return nil, utils.NewInvalidError("environments", "Environments field must be present and must not be empty")
	}

	if HasValue(data.Enabled) {
		if !data.Enabled.ValueBool() {
			return nil, utils.NewInvalidError("enabled", "during creation enabled field must be true")
		}
	}
	if HasValue(data.IpType) {

		ruleCondition := generated.InputMaliciousSourcesRuleCondition{}
		ipLocationCondition := generated.InputMaliciousSourcesRuleIpLocationTypeCondition{}
		ipLocationTypes := []*generated.MaliciousSourcesRuleIpLocationType{}

		for _, ipType := range data.IpType.Elements() {

			if ipType, ok := ipType.(types.String); ok {
				ipType, exist := MaliciousIpTypeMap[ipType.ValueString()]
				if !exist {
					return nil, utils.NewInvalidError("ip_type", fmt.Sprintf("%s Invalid Ip Type", ipType))
				}
				ipLocationTypes = append(ipLocationTypes, &ipType)
			}
		}

		ipLocationCondition.IpLocationTypes = ipLocationTypes
		ruleCondition.ConditionType = generated.MaliciousSourcesRuleConditionTypeIpLocationType
		ruleCondition.IpLocationTypeCondition = &ipLocationCondition
		ruleInfo.Conditions = append(ruleInfo.Conditions, &ruleCondition)
	} else {
		return nil, utils.NewInvalidError("ip_type", "IpType field must be present and must not be empty")
	}

	ruleInfo.Action = ruleAction
	input.RuleInfo = ruleInfo

	return &input, nil

}

func convertMaliciousIpTypeModelToUpdateInput(ctx context.Context, data *models.MaliciousIpTypeModel, id string) (*generated.InputMaliciousSourcesRuleUpdate, error) {
	var input = generated.InputMaliciousSourcesRuleUpdate{}
	rule := generated.InputMaliciousSourcesRule{}

	var ruleInfo = generated.InputMaliciousSourcesRuleInfo{}
	var ruleAction = generated.InputMaliciousSourcesRuleAction{}

	rule.Id = id

	if HasValue(data.Name) {
		ruleInfo.Name = data.Name.ValueString()
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		ruleInfo.Description = &description
	}

	if HasValue(data.EventSeverity) {
		eventSeverity, exist := MaliciousIpTypeEventSeverityMap[data.EventSeverity.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("event_severity", "Invalid EventSeverity")
		}
		ruleAction.EventSeverity = &eventSeverity
	}
	if HasValue(data.Action) {
		action, exist := MaliciousIpTypeActionMap[data.Action.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("action", "Invalid action")
		}
		ruleAction.RuleActionType = action
	} else {
		return nil, utils.NewInvalidError("action", "Action field must be present and must not be empty")
	}

	if HasValue(data.Duration) {
		duration := data.Duration.ValueString()
		expirationDetails := generated.InputMaliciousSourcesRuleExpirationDetails{
			ExpirationDuration: duration,
		}
		ruleAction.ExpirationDetails = &expirationDetails
	}

	if HasValue(data.Environments) {
		ruleScope := generated.InputMaliciousSourcesRuleScope{}
		environmentScope := generated.InputMaliciousSourcesRuleEnvironmentScope{}
		environments, err := utils.ConvertSetToStrPointer(data.Environments)
		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")
		}
		environmentScope.EnvironmentIds = environments
		ruleScope.EnvironmentScope = &environmentScope
		rule.Scope = &ruleScope
	} else {
		return nil, utils.NewInvalidError("environments", "Environments field must be present and must not be empty")
	}

	if HasValue(data.Enabled) {
		rule.Status = generated.InputMaliciousSourcesRuleStatus{
			Disabled: !data.Enabled.ValueBool(),
		}

	}
	if HasValue(data.IpType) {
		ruleCondition := generated.InputMaliciousSourcesRuleCondition{}
		ipLocationCondition := generated.InputMaliciousSourcesRuleIpLocationTypeCondition{}
		ipLocationTypes := []*generated.MaliciousSourcesRuleIpLocationType{}

		for _, ipType := range data.IpType.Elements() {

			if ipType, ok := ipType.(types.String); ok {
				ipType, exist := MaliciousIpTypeMap[ipType.ValueString()]
				if !exist {
					return nil, utils.NewInvalidError("ip_type", fmt.Sprintf("%s Invalid Ip Type", ipType))
				}
				ipLocationTypes = append(ipLocationTypes, &ipType)
			}
		}

		ipLocationCondition.IpLocationTypes = ipLocationTypes
		ruleCondition.ConditionType = generated.MaliciousSourcesRuleConditionTypeIpLocationType
		ruleCondition.IpLocationTypeCondition = &ipLocationCondition
		ruleInfo.Conditions = append(ruleInfo.Conditions, &ruleCondition)
	} else {
		return nil, utils.NewInvalidError("ip_type", "IpType field must be present and must not be empty")
	}

	ruleInfo.Action = ruleAction
	rule.Info = ruleInfo
	input.Rule = rule

	return &input, nil

}

func convertMaliciousIpTypeFieldsToModel(ctx context.Context, rule *generated.MaliciousIpTypeRuleFields) (*models.MaliciousIpTypeModel, error) {
	var data = models.MaliciousIpTypeModel{}

	data.Id = types.StringValue(rule.Id)
	data.Name = types.StringValue(rule.Info.Name)
	if rule.Info.Description != nil {
		data.Description = types.StringValue(*rule.Info.Description)
	}
	if rule.Info.Action.EventSeverity != nil {
		data.EventSeverity = types.StringValue(string(*rule.Info.Action.EventSeverity))
	}

	data.Action = types.StringValue(string(rule.Info.Action.RuleActionType))
	if rule.Info.Action.ExpirationDetails != nil {
		data.Duration = types.StringValue(rule.Info.Action.ExpirationDetails.ExpirationDuration)
	}
	data.Enabled = types.BoolValue(!rule.Status.Disabled)

	if len(rule.Info.Conditions) > 0 {
		ipTypes := []*string{}
		for _, condition := range rule.Info.Conditions {
			if condition.ConditionType == generated.MaliciousSourcesRuleConditionTypeIpLocationType {
				if condition.IpLocationTypeCondition != nil {
					if condition.IpLocationTypeCondition.IpLocationTypes != nil {
						for _, ipLocationType := range condition.IpLocationTypeCondition.IpLocationTypes {
							iptype := string(*ipLocationType)
							ipTypes = append(ipTypes, &iptype)
						}
					}
				}
			}
		}
		ipTypeSet, err := utils.ConvertStringPtrToTerraformSet(ipTypes)
		if err != nil {
			return nil, fmt.Errorf("converting ip types to string pointer fails")
		}
		data.IpType = ipTypeSet
	} else {
		data.IpType = types.SetNull(types.StringType)
	}
	if rule.Scope != nil && rule.Scope.EnvironmentScope != nil {
		environments, err := utils.ConvertStringPtrToTerraformSet(rule.Scope.EnvironmentScope.EnvironmentIds)
		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")
		}
		data.Environments = environments
	} else {
		data.Environments = types.SetNull(types.StringType)
	}
	return &data, nil
}
