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

type MaliciousRegionResource struct {
	client *graphql.Client
}

func NewMaliciousRegionResource() resource.Resource {
	return &MaliciousRegionResource{}
}

func (r *MaliciousRegionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MaliciousRegionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_malicious_region"
}

func (r *MaliciousRegionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.MaliciousRegionResourceSchema()
}

func (r *MaliciousRegionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "Entering in Create Block")
	var data *models.MaliciousRegionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleInput, err := convertMaliciousRegionModelToCreateInput(ctx, data, *r.client)
	if ruleInput == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	id, err := getMaliciousRegionId(ruleInput.Name, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	if id != "" {
		resp.Diagnostics.AddError("Resource already Exist", fmt.Sprintf("%s malicious region rule already please try with different name or import it", ruleInput.Name))
		return
	}

	response, err := generated.CreateMaliciousRegionRule(ctx, *r.client, *ruleInput)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(*&response.CreateRegionRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Create Block")

}

func (r *MaliciousRegionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "Entering in Read Block")
	var data *models.MaliciousRegionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := data.Name.ValueString()
	rule, err := getMaliciousRegionRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	ruleData, err := convertMaliciousRegionFieldsToModel(ctx, rule)
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

func (r *MaliciousRegionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Entering in Update Block")
	var data *models.MaliciousRegionModel
	var dataState *models.MaliciousRegionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	input, err := convertMaliciousRegionModelToUpdateInput(ctx, data, dataState.Id.ValueString(), *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	_, err = generated.UpdateMaliciousRegionRule(ctx, *r.client, *input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(input.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Update Block")

}

func (r *MaliciousRegionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "Entering in Delete Block")
	var data *models.MaliciousRegionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleid := data.Id.ValueString()
	response, err := generated.DeleteMaliciousRegionRule(ctx, *r.client, generated.InputRegionRuleDelete{Id: ruleid})
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if response.DeleteRegionRule.Success != true {
		utils.AddError(ctx, &resp.Diagnostics, fmt.Errorf("failed to delete rule %s", data.Name.ValueString()))
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "Exiting in Delete Block")

}

func (r *MaliciousRegionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "Entering in ImportState Block")
	ruleName := req.ID
	id, err := getMaliciousRegionId(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Resource Not Found", fmt.Sprintf("%s rule of this name not found", ruleName))
		return
	}
	response, err := getMaliciousRegionRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data, err := convertMaliciousRegionFieldsToModel(ctx, response)

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

func convertMaliciousRegionModelToCreateInput(ctx context.Context, data *models.MaliciousRegionModel, r graphql.Client) (*generated.InputRegionRuleCreate, error) {
	var input = generated.InputRegionRuleCreate{}
	if HasValue(data.Name) {
		name := data.Name.ValueString()
		input.Name = name
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		input.Description = &description
	}
	if HasValue(data.EventSeverity) {
		eventSeverity, exist := MaliciousRegionEventSeverityMap[data.EventSeverity.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("event_severity", "Invalid EventSeverity")
		}
		input.EventSeverity = &eventSeverity
	}
	if HasValue(data.Duration) {
		duration := data.Duration.ValueString()
		input.Duration = &duration
	}
	if HasValue(data.Regions) {
		regions, err := utils.ConvertSetToStrPointer(data.Regions)
		if err != nil {
			return nil, fmt.Errorf("converting regions to string pointer fails")
		}
		regionsId, err := GetCountriesId(regions, ctx, r)
		if err != nil {
			return nil, err
		}
		input.RegionIds = regionsId
	} else {
		return nil, utils.NewInvalidError("regions", "regions field must be present and not empty")
	}

	if HasValue(data.Environments) {
		ruleScope := &generated.InputRegionRuleScope{}
		environments, err := utils.ConvertSetToStrPointer(data.Environments)
		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")
		}
		ruleScope.EnvironmentScope = &generated.InputRegionEnvironmentScope{
			EnvironmentIds: environments,
		}
		input.RuleScope = ruleScope

	}
	if HasValue(data.Action) {
		action, exist := MaliciousRegionActionMap[data.Action.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("action", "Invalid action")
		}
		input.Type = action
	} else {
		return nil, utils.NewInvalidError("action", "action must be present and not empty")
	}
	return &input, nil
}

func convertMaliciousRegionModelToUpdateInput(ctx context.Context, data *models.MaliciousRegionModel, id string, r graphql.Client) (*generated.InputRegionRuleUpdate, error) {
	var input = generated.InputRegionRuleUpdate{}
	input.Id = id
	if HasValue(data.Name) {
		name := data.Name.ValueString()
		input.Name = name
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		input.Description = &description
	}
	if HasValue(data.EventSeverity) {
		eventSeverity, exist := MaliciousRegionEventSeverityMap[data.EventSeverity.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("event_severity", "Invalid EventSeverity")
		}
		input.EventSeverity = &eventSeverity
	}
	if HasValue(data.Duration) {
		duration := data.Duration.ValueString()
		input.Duration = &duration
	}
	if HasValue(data.Regions) {
		regions, err := utils.ConvertSetToStrPointer(data.Regions)
		if err != nil {
			return nil, fmt.Errorf("converting regions to string pointer fails")
		}
		regionsId, err := GetCountriesId(regions, ctx, r)
		if err != nil {
			return nil, err
		}
		input.RegionIds = regionsId
	} else {
		return nil, utils.NewInvalidError("regions", "regions field must be present and not empty")
	}

	if HasValue(data.Environments) {
		ruleScope := &generated.InputRegionRuleScope{}
		environments, err := utils.ConvertSetToStrPointer(data.Environments)
		if err != nil {
			return nil, fmt.Errorf("converting environments to string pointer fails")
		}
		ruleScope.EnvironmentScope = &generated.InputRegionEnvironmentScope{
			EnvironmentIds: environments,
		}
		input.RuleScope = ruleScope

	}
	if HasValue(data.Enabled) {
		disabled := !data.Enabled.ValueBool()
		input.Disabled = &disabled
	}
	if HasValue(data.Action) {
		action, exist := MaliciousRegionActionMap[data.Action.ValueString()]
		if !exist {
			return nil, utils.NewInvalidError("action", "Invalid action")
		}
		input.Type = action
	} else {
		return nil, utils.NewInvalidError("action", "action must be present and not empty")
	}
	return &input, nil

}

func getMaliciousRegionId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
	input := &generated.InputRegionRulesFilter{}

	response, err := generated.GetMaliciousRegionRulesName(ctx, r, input)
	if err != nil {
		return "", err
	}
	for _, rule := range response.RegionRules.Results {
		if rule.Name == ruleName {
			return rule.Id, nil
		}
	}
	return "", nil
}

func getMaliciousRegionRule(ruleName string, ctx context.Context, r graphql.Client) (*generated.MaliciousRegionRuleFields, error) {
	input := &generated.InputRegionRulesFilter{}

	response, err := generated.GetMaliciousRegionRuleDetails(ctx, r, input)
	if err != nil {
		return nil, err
	}
	for _, rule := range response.RegionRules.Results {
		if rule.Name == ruleName {
			rulefields := rule.MaliciousRegionRuleFields
			return &rulefields, nil
		}

	}
	return nil, nil

}

func convertMaliciousRegionFieldsToModel(ctx context.Context, data *generated.MaliciousRegionRuleFields) (*models.MaliciousRegionModel, error) {
	model := models.MaliciousRegionModel{}
	model.Id = types.StringValue(data.Id)
	model.Name = types.StringValue(data.Name)
	if data.Description != nil {
		model.Description = types.StringValue(*data.Description)
	}
	model.Action = types.StringValue(string(data.Type))
	if data.EventSeverity != nil {
		model.EventSeverity = types.StringValue(string(*data.EventSeverity))
	}
	model.Enabled = types.BoolValue(!data.Disabled)
	if data.Expiration != nil {
		model.Duration = types.StringValue(data.Expiration.Duration)
	}
	if len(data.Regions) > 0 {
		isoCodes := []*string{}
		for _, region := range data.Regions {
			if region.Country != nil {
				isoCode := region.Country.IsoCode
				isoCodes = append(isoCodes, &isoCode)
			}
		}
		regions, err := utils.ConvertStringPtrToTerraformSet(isoCodes)
		if err != nil {
			return nil, fmt.Errorf("converting regions to string pointer fails")
		}
		model.Regions = regions
	} else {
		model.Regions = types.SetNull(types.StringType)
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
