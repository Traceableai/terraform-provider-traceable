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

type DataLossPreventionUserBasedResource struct {
	client *graphql.Client
}

func NewDataLossPreventionUserBasedResource() resource.Resource {
	return &DataLossPreventionUserBasedResource{}
}

func (r *DataLossPreventionUserBasedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Entering in Configure Block")
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
	tflog.Trace(ctx, "Client Intialization Successfully And Existing from Configure Block")
}

func (r *DataLossPreventionUserBasedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_loss_prevention_user_based"
}

func (r *DataLossPreventionUserBasedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.RateLimitingResourceSchema()
}

func (r *DataLossPreventionUserBasedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Entering in Create Block")
	var data *models.RateLimitingRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleInput, err := convertDataLossPreventionRateLimitingModelToCreateInput(ctx, data, r.client)
	if ruleInput == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	id, err := getDataLossPreventionRateLimitingRuleId(ruleInput.Name, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	if id != "" {
		resp.Diagnostics.AddError("Resource already Exist", fmt.Sprintf("%s rate limiting rule already please try with different name or import it", ruleInput.Name))
		return
	}

	rule, err := generated.CreateRateLimitingRule(ctx, *r.client, *ruleInput)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	data.Id = types.StringValue(rule.CreateRateLimitingRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Exiting in Create Block")

}

func (r *DataLossPreventionUserBasedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.RateLimitingRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getDataLossPreventionRateLimitingRule(data.Id.ValueString(), ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	updatedData, err := convertRateLimitingRuleFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedData)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *DataLossPreventionUserBasedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.RateLimitingRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataState *models.RateLimitingRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)

	if resp.Diagnostics.HasError() {
		return
	}
	input, err := convertDataLossPreventionRateLimitingModelToUpdateInput(ctx, data, dataState.Id.ValueString(), r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error in Updating rate limiting rule", err.Error())
		return

	}

	resp1, err2 := generated.UpdateRateLimitingRule(ctx, *r.client, *input)
	if err2 != nil {
		resp.Diagnostics.AddError("Error in Updating rate limiting rule", err.Error())
		return
	}
	data.Id = types.StringValue(resp1.UpdateRateLimitingRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataLossPreventionUserBasedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.RateLimitingRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := generated.DeleteRateLimitingRule(ctx, *r.client, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting rate limiting rule", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *DataLossPreventionUserBasedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ruleName := req.ID
	id, err := getDataLossPreventionRateLimitingRuleId(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Resource Not Found", fmt.Sprintf("%s rule of this name not found", ruleName))
		return
	}
	response, err := getDataLossPreventionRateLimitingRule(id, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data, err := convertRateLimitingRuleFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func getDataLossPreventionRateLimitingRule(id string, ctx context.Context, r graphql.Client) (*generated.RateLimitingRuleFields, error) {
	rateLimitingfields := generated.RateLimitingRuleFields{}
	category := []*generated.RateLimitingRuleCategory{}
	enumeration := generated.RateLimitingRuleCategoryDataExfiltration
	category = append(category, &enumeration)
	response, err := generated.GetRateLimitingDetails(ctx, r, category, nil)
	if err != nil {
		return nil, err
	}

	for _, rule := range response.RateLimitingRules.Results {
		if rule.Id == id {
			rateLimitingfields = rule.RateLimitingRuleFields
			return &rateLimitingfields, nil
		}
	}

	return nil, nil
}

func getDataLossPreventionRateLimitingRuleId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
	category := []*generated.RateLimitingRuleCategory{}
	enumeration := generated.RateLimitingRuleCategoryDataExfiltration
	category = append(category, &enumeration)

	response, err := generated.GetRateLimitingRulesName(ctx, r, category, nil)
	if err != nil {
		return "", err
	}
	for _, rule := range response.RateLimitingRules.Results {
		if rule.Name == ruleName {
			return rule.GetId(), nil
		}
	}
	return "", nil

}
func convertDataLossPreventionRateLimitingModelToCreateInput(ctx context.Context, data *models.RateLimitingRuleModel, client *graphql.Client) (*generated.InputRateLimitingRuleData, error) {
	var input = generated.InputRateLimitingRuleData{}
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
	if HasValue(data.Enabled) {
		enabled := data.Enabled.ValueBool()
		input.Enabled = enabled
	}
	category := generated.RateLimitingRuleCategoryDataExfiltration
	input.Category = category
	scope, err := convertToRuleConfigScope(data.Environments)
	if err != nil {
		return nil, err
	} else {
		input.RuleConfigScope = scope
	}
	status, err := convertToRateLimitingRuleStatus(data)
	if err != nil {
		return nil, err
	} else {
		input.RuleStatus = status
	}
	thresholdActionConfigs, err := convertToRateLimitingRuleThresholdActionConfigType(data)
	if err != nil {
		return nil, err
	} else {
		input.ThresholdActionConfigs = thresholdActionConfigs
	}
	conditions, err := convertToRateLimitingRuleCondition(ctx, data, client)
	if err != nil {
		return nil, err
	} else {
		input.Conditions = conditions
	}

	if HasValue(data.Sources.EndpointLabels) && HasValue(data.Sources.Endpoints) {
		return nil, utils.NewInvalidError("sources.endpoint", "endpoint_labels field must not be present at same time ")
	}

	return &input, nil
}

func convertDataLossPreventionRateLimitingModelToUpdateInput(ctx context.Context, data *models.RateLimitingRuleModel, id string, client *graphql.Client) (*generated.InputRateLimitingRule, error) {
	input := generated.InputRateLimitingRule{}

	if id != "" {
		input.Id = id
	} else {
		return nil, fmt.Errorf("Id can not be empty")
	}
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
	if HasValue(data.Enabled) {
		enabled := data.Enabled.ValueBool()
		input.Enabled = enabled
	}
	category := generated.RateLimitingRuleCategoryDataExfiltration
	input.Category = category
	scope, err := convertToRuleConfigScope(data.Environments)
	if err != nil {
		return nil, err
	} else {
		input.RuleConfigScope = scope
	}
	status, err := convertToRateLimitingRuleStatus(data)
	if err != nil {
		return nil, err
	} else {
		input.RuleStatus = status
	}
	thresholdActionConfigs, err := convertToRateLimitingRuleThresholdActionConfigType(data)
	if err != nil {
		return nil, err
	} else {
		input.ThresholdActionConfigs = thresholdActionConfigs
	}
	conditions, err := convertToRateLimitingRuleCondition(ctx, data, client)
	if err != nil {
		return nil, err
	} else {
		input.Conditions = conditions
	}

	if HasValue(data.Sources.EndpointLabels) && HasValue(data.Sources.Endpoints) {
		return nil, utils.NewInvalidError("sources.endpoint", "endpoint_labels field must not be present at same time ")
	}

	return &input, nil

}
