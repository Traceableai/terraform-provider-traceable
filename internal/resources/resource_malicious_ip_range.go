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

type MaliciousIpRangeResource struct {
	client *graphql.Client
}

func NewMaliciousIpRangeResource() resource.Resource {
	return &MaliciousIpRangeResource{}
}

func (r *MaliciousIpRangeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Entering in Configure Block")
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		utils.AddError(ctx, &resp.Diagnostics, fmt.Errorf("expected *graphql.Client, got: %T", req.ProviderData))
		return
	}
	r.client = client
	tflog.Trace(ctx, "Client Intialization Successfully And Existing from Configure Block")
}

func (r *MaliciousIpRangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_malicious_ip_range"
}

func (r *MaliciousIpRangeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.MaliciousIPRangeResourceSchema()
}

func (r *MaliciousIpRangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "Entering in Create Block")
	var data *models.MaliciousIpRangeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleInput, err := convertMaliciousIpRangeModelToCreateInput(ctx, data)
	if ruleInput == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	id, err := getMaliciousIpRangeId(ruleInput.RuleDetails.Name, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	if id != "" {
		resp.Diagnostics.AddError("Resource already Exist", fmt.Sprintf("%s malicious ip range rule already please try with different name or import it", ruleInput.RuleDetails.Name))
		return
	}

	response, err := generated.CreateMaliciousIpRangeRule(ctx, *r.client, *ruleInput)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	data.Id = types.StringValue(*response.CreateIpRangeRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Create Block")

}

func (r *MaliciousIpRangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "Entering in Read Block")
	var data *models.MaliciousIpRangeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := data.Name.ValueString()
	rule, err := getMaliciousIpRangeRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	ruleData, err := convertMaliciousIpRangeFieldsToModel(ctx, rule)
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

func (r *MaliciousIpRangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Entering in Update Block")
	var data *models.MaliciousIpRangeModel
	var dataState *models.MaliciousIpRangeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	input, err := convertMaliciousIpRangeModelToUpdateInput(ctx, data, dataState.Id.ValueString())
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	_, err = generated.UpdateMaliciousIpRangeRule(ctx, *r.client, *input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(*input.GetId())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Update Block")

}

func (r *MaliciousIpRangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Trace(ctx, "Entering in Delete Block")
	var data *models.MaliciousIpRangeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleid := data.Id.ValueString()
	response, err := generated.DeleteMaliciousIpRangeRule(ctx, *r.client, generated.InputIpRangeRuleDelete{Id: &ruleid})
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if response.DeleteIpRangeRule.Success != true {
		utils.AddError(ctx, &resp.Diagnostics, fmt.Errorf("failed to delete rule %s", data.Name.ValueString()))
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "Exiting in Delete Block")
}

func (r *MaliciousIpRangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "Entering in ImportState Block")
	ruleName := req.ID
	id, err := getMaliciousIpRangeId(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Resource Not Found", fmt.Sprintf("%s rule of this name not found", ruleName))
		return
	}
	response, err := getMaliciousIpRangeRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data, err := convertMaliciousIpRangeFieldsToModel(ctx, response)

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

func getMaliciousIpRangeId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
	input := &generated.InputIpRangeRulesFilter{}

	response, err := generated.GetMaliciousIpRangeRulesName(ctx, r, input)
	if err != nil {
		return "", err
	}
	for _, rule := range response.IpRangeRules.Results {
		if rule.RuleDetails.Name == ruleName {
			return *rule.Id, nil
		}
	}
	return "", nil
}
func getMaliciousIpRangeRule(ruleName string, ctx context.Context, r graphql.Client) (*generated.MaliciousIpRangeFields, error) {

	input := &generated.InputIpRangeRulesFilter{}

	response, err := generated.GetMaliciousIpRangeRuleDetails(ctx, r, input)
	if err != nil {
		return nil, err
	}
	for _, rule := range response.IpRangeRules.Results {

		if rule.RuleDetails.Name == ruleName {
			rulefields := rule.MaliciousIpRangeFields

			return &rulefields, nil
		}

	}
	return nil, nil
}

func convertMaliciousIpRangeModelToCreateInput(ctx context.Context, data *models.MaliciousIpRangeModel) (*generated.InputIpRangeRuleCreate, error) {
	var input = generated.InputIpRangeRuleCreate{}
	var ruleDetails = generated.InputIpRangeRuleDetailsRequest{}
	if HasValue(data.Name) {
		name := data.Name.ValueString()
		ruleDetails.Name = name
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		ruleDetails.Description = &description
	}
	if HasValue(data.EventSeverity) {
		eventSeverity, exist := MaliciousIpRangeEventSeverityMap[data.EventSeverity.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("EventSeverity", "Invalid EventSeverity")
		}
		ruleDetails.EventSeverity = &eventSeverity
	}
	if HasValue(data.Duration) {
		duration := data.Duration.ValueString()
		ruleDetails.ExpirationDuration = &duration
	}
	if HasValue(data.IpRange) {
		iprange, err := utils.ConvertSetToStrPointer(data.IpRange)
		if err != nil {
			return nil, fmt.Errorf("converting ip range to string pointer fails")
		}
		ruleDetails.RawIpRangeData = iprange
	} else {
		return nil, utils.NewInvalidError("ip_range", "ip range must be present and not empty")
	}
	if HasValue(data.Action) {
		action, exist := MaliciousIpRangeActionMap[data.Action.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("action", "Invalid action")
		}
		ruleDetails.RuleAction = action
	} else {
		return nil, utils.NewInvalidError("action", "action must be present and not empty")
	}

	if HasValue(data.Environments) {
		ruleScope := &generated.InputIpRangeRuleScope{}
		environments, err := utils.ConvertSetToStrPointer(data.Environments)
		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")
		}
		ruleScope.EnvironmentScope = &generated.InputIpRangeEnvironmentScope{
			EnvironmentIds: environments,
		}
		input.RuleScope = ruleScope

	}
	input.RuleDetails = ruleDetails
	return &input, nil
}

func convertMaliciousIpRangeModelToUpdateInput(ctx context.Context, data *models.MaliciousIpRangeModel, id string) (*generated.InputIpRangeRuleUpdate, error) {
	var input = generated.InputIpRangeRuleUpdate{}
	var ruleDetails = generated.InputIpRangeRuleDetailsRequest{}
	if HasValue(data.Name) {
		name := data.Name.ValueString()
		ruleDetails.Name = name
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		ruleDetails.Description = &description
	}
	if HasValue(data.EventSeverity) {
		eventSeverity, exist := MaliciousIpRangeEventSeverityMap[data.EventSeverity.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("EventSeverity", "Invalid EventSeverity")
		}
		ruleDetails.EventSeverity = &eventSeverity
	}
	if HasValue(data.Duration) {
		duration := data.Duration.ValueString()
		ruleDetails.ExpirationDuration = &duration
	}
	if HasValue(data.IpRange) {
		iprange, err := utils.ConvertSetToStrPointer(data.IpRange)
		if err != nil {
			return nil, fmt.Errorf("converting ip range to string pointer fails")
		}
		ruleDetails.RawIpRangeData = iprange
	} else {
		return nil, utils.NewInvalidError("ip_range", "ip range must be present and not empty")
	}
	if HasValue(data.Action) {
		action, exist := MaliciousIpRangeActionMap[data.Action.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("action", "Invalid action")
		}
		ruleDetails.RuleAction = action
	} else {
		return nil, utils.NewInvalidError("action", "action must be present and not empty")
	}

	if HasValue(data.Environments) {
		ruleScope := &generated.InputIpRangeRuleScope{}
		environments, err := utils.ConvertSetToStrPointer(data.Environments)
		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")
		}
		ruleScope.EnvironmentScope = &generated.InputIpRangeEnvironmentScope{
			EnvironmentIds: environments,
		}
		input.RuleScope = ruleScope

	}
	if HasValue(data.Enabled) {
		enabled := !data.Enabled.ValueBool()
		input.Disabled = &enabled
	}
	input.Id = &id
	input.RuleDetails = ruleDetails
	return &input, nil
}

func convertMaliciousIpRangeFieldsToModel(ctx context.Context, data *generated.MaliciousIpRangeFields) (*models.MaliciousIpRangeModel, error) {
	model := models.MaliciousIpRangeModel{}
	if data.Id != nil {
		model.Id = types.StringValue(*data.Id)
	}
	model.Enabled = types.BoolValue(!data.Disabled)

	model.Name = types.StringValue(data.RuleDetails.Name)
	if data.RuleDetails.Description != nil {
		model.Description = types.StringValue(*data.RuleDetails.Description)
	}

	if data.RuleDetails.EventSeverity != nil {
		model.EventSeverity = types.StringValue(string(*data.RuleDetails.EventSeverity))
	}

	model.Action = types.StringValue(MaliciousIpRangeIpRangeMapResponse[string(data.RuleDetails.RuleAction)])
	if data.RuleDetails.ExpirationDetails != nil {
		model.Duration = types.StringValue(data.RuleDetails.ExpirationDetails.ExpirationDuration)
	}

	if data.RuleDetails.RawIpRangeData != nil {
		iprange, err := utils.ConvertStringPtrToTerraformSet(data.RuleDetails.RawIpRangeData)
		if err != nil {
			return nil, fmt.Errorf("converting ip range to string pointer fails")
		}
		model.IpRange = iprange
	} else {
		model.IpRange = types.SetNull(types.StringType)
	}

	if data.RuleScope != nil && data.RuleScope.EnvironmentScope != nil {

		environments, err := utils.ConvertStringPtrToTerraformSet(data.RuleScope.EnvironmentScope.EnvironmentIds)

		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")

		}
		model.Environments = environments
	} else {
		model.Environments = types.SetNull(types.StringType)
	}

	return &model, nil

}
