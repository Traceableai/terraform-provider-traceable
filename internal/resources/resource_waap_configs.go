package resources

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/traceableai/terraform-provider-traceable/internal/generated"
	"github.com/traceableai/terraform-provider-traceable/internal/models"
	"github.com/traceableai/terraform-provider-traceable/internal/schemas"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

type WaapConfigResource struct {
	client *graphql.Client
}

func NewWaapConfigResource() resource.Resource {
	return &WaapConfigResource{}
}

func (r *WaapConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WaapConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_waap_config"
}

func (r *WaapConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.WaapConfigResourceSchema()
}

func (r *WaapConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.WaapConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	err := getUpdateRules(ctx, data, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error in Updating waap config", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WaapConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.WaapConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getRuleConfigAnomalyState(data, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	updatedData, err := convertWaapConfigsFieldsToModel(data, ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedData)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WaapConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.WaapConfigModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var dataState *models.WaapConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := getUpdateRules(ctx, data, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error in Updating waap config", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WaapConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
	if resp.Diagnostics.HasError() {
		return
	}
}

func convertWaapConfigsFieldsToModel(currState *models.WaapConfigModel, ctx context.Context, data []*generated.AnomalyDetectionRuleConfigsAnomalyDetectionRuleConfigsAnomalyDetectionRuleConfigsResultSetResultsAnomalyRuleConfig) (*models.WaapConfigModel, error) {
	var waapConfigRuleConfigs []models.WaapRuleConfigModel
	var ruleConfigs []models.WaapRuleConfigModel

	conversionError := utils.ConvertElementsSet(currState.RuleConfigs, &ruleConfigs)
	if conversionError != nil {
		return nil, conversionError
	}

	for _, ruleConfig := range ruleConfigs {
		ruleName := ruleConfig.RuleName
		fetchedRuleConfig, err := findRuleConfigsWithRuleName(ruleName.ValueString(), data)
		if err != nil {
			return nil, err
		}
		disabled := fetchedRuleConfig.ConfigStatus.Disabled

		var subRules []models.WaapSubRuleConfigModel
		diags := ruleConfig.Subrules.ElementsAs(ctx, &subRules, false)
		if diags.HasError() {
			return nil, fmt.Errorf("conversion error")
		}
		var finalSubRuleModel []models.WaapSubRuleConfigModel
		for _, subRule := range subRules {
			fetchedSubruleConfig, err := findSubRuleConfigsWithSubRuleName(subRule.SubRuleName.ValueString(), fetchedRuleConfig.SubRuleConfigs)
			if err != nil {
				return nil, err
			}
			finalSubRuleModel = append(finalSubRuleModel, models.WaapSubRuleConfigModel{
				SubRuleName:   types.StringValue(subRule.SubRuleName.ValueString()),
				SubRuleAction: types.StringValue(string(fetchedSubruleConfig.AnomalySubRuleAction)),
			})
		}
		waapConfigSubRuleConfigsObjectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"sub_rule_name":   types.StringType,
				"sub_rule_action": types.StringType,
			},
		}
		if len(finalSubRuleModel) > 0 {
			subruleConfigsModel, diags := types.SetValueFrom(
				ctx,
				waapConfigSubRuleConfigsObjectType,
				finalSubRuleModel,
			)
			if diags.HasError() {
				return nil, fmt.Errorf("subrule config conversion failed")
			}
			waapConfigRuleConfigs = append(waapConfigRuleConfigs, models.WaapRuleConfigModel{
				RuleName: ruleName,
				Enabled:  types.BoolValue(!disabled),
				Subrules: subruleConfigsModel,
			})
		} else {
			waapConfigRuleConfigs = append(waapConfigRuleConfigs, models.WaapRuleConfigModel{
				RuleName: ruleName,
				Enabled:  types.BoolValue(!disabled),
				Subrules: types.SetNull(waapConfigSubRuleConfigsObjectType),
			})
		}
	}

	waapConfigModel := models.WaapConfigModel{}
	waapConfigRuleConfigsObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"rule_name": types.StringType,
			"disabled":  types.BoolType,
			"subrules": types.SetType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"sub_rule_name":   types.StringType,
					"sub_rule_action": types.StringType,
				},
			}},
		},
	}
	if len(waapConfigRuleConfigs) > 0 {
		ruleConfigsModel, diags := types.SetValueFrom(
			ctx,
			waapConfigRuleConfigsObjectType,
			waapConfigRuleConfigs,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("rule config conversion failed")
		}
		waapConfigModel.RuleConfigs = ruleConfigsModel
	} else {
		waapConfigModel.RuleConfigs = types.SetNull(waapConfigRuleConfigsObjectType)
	}
	env := currState.Environment
	if HasValue(env) {
		waapConfigModel.Environment = env
	}
	return &waapConfigModel, nil
}

func findSubRuleConfigsWithSubRuleName(subRuleName string, subRuleConfigs []*generated.AnomalyRuleConfigFieldsSubRuleConfigsAnomalySubRuleConfig) (generated.AnomalyRuleConfigFieldsSubRuleConfigsAnomalySubRuleConfig, error) {
	for _, subRuleConfig := range subRuleConfigs {
		if subRuleConfig.SubRuleName == subRuleName {
			return *subRuleConfig, nil
		}
	}
	return generated.AnomalyRuleConfigFieldsSubRuleConfigsAnomalySubRuleConfig{}, utils.NewInvalidError("subrules sub_rule_name", fmt.Sprintf("sub_rule_name %s not found", subRuleName))
}

func findRuleConfigsWithRuleName(ruleName string, data []*generated.AnomalyDetectionRuleConfigsAnomalyDetectionRuleConfigsAnomalyDetectionRuleConfigsResultSetResultsAnomalyRuleConfig) (*generated.AnomalyRuleConfigFields, error) {
	ruleId, err := GetRuleId(ruleName)
	if err != nil {
		return nil, err

	}
	for _, ruleConfig := range data {
		if ruleConfig.RuleId == ruleId {
			return &ruleConfig.AnomalyRuleConfigFields, nil
		}
	}
	return nil, fmt.Errorf("rule config not found")
}

func getRuleConfigAnomalyState(data *models.WaapConfigModel, ctx context.Context, r graphql.Client) ([]*generated.AnomalyDetectionRuleConfigsAnomalyDetectionRuleConfigsAnomalyDetectionRuleConfigsResultSetResultsAnomalyRuleConfig, error) {
	inputFilter := generated.InputAnomalyScope{}
	inputFilter.ScopeType = generated.AnomalyScopeTypeEnvironment
	if HasValue(data.Environment) {
		environment := data.Environment.ValueString()
		inputFilter.EnvironmentScope = &generated.InputAnomalyEnvironmentScope{
			Id: environment,
		}
	} else {
		inputFilter.ScopeType = generated.AnomalyScopeTypeCustomer
	}
	response, err := generated.AnomalyDetectionRuleConfigs(ctx, r, inputFilter)
	if err != nil {
		return nil, err
	}

	return response.AnomalyDetectionRuleConfigs.Results, nil
}

func getUpdateRules(ctx context.Context, data *models.WaapConfigModel, r *graphql.Client) error {
	inputArr := []generated.InputScopedAnomalyRuleConfigUpdate{}
	var ruleConfigs []models.WaapRuleConfigModel
	var subruleConfigs []models.WaapSubRuleConfigModel

	conversionError := utils.ConvertElementsSet(data.RuleConfigs, &ruleConfigs)
	if conversionError != nil {
		return conversionError
	}
	for _, ruleConfig := range ruleConfigs {
		diags := ruleConfig.Subrules.ElementsAs(ctx, &subruleConfigs, false)
		if diags.HasError() {
			return fmt.Errorf("conversion error")
		}
		if !ruleConfig.Enabled.ValueBool() && len(subruleConfigs) > 0 {
			return fmt.Errorf("rule must be set to enabled to modify subrules. either remove subrules block to disable the rule")
		}
		if len(subruleConfigs) > 0 {
			for _, subruleConfig := range subruleConfigs {
				if !HasValue(subruleConfig.SubRuleAction) || !HasValue(subruleConfig.SubRuleName) {
					return utils.NewInvalidError("subrules", "field must be present and must not be empty")
				}
				input, err := convertWaapConfigModelToUpdateInput(ctx, data, &ruleConfig, subruleConfig.SubRuleName.ValueString(), subruleConfig.SubRuleAction.ValueString())
				if err != nil {
					return err
				}
				inputArr = append(inputArr, *input)
			}
		} else {
			input, err := convertWaapConfigModelToUpdateInput(ctx, data, &ruleConfig, "", "")
			if err != nil {
				return err
			}
			inputArr = append(inputArr, *input)
		}
		if len(subruleConfigs) > 0 {
			firstInput, err := convertWaapConfigModelToUpdateInput(ctx, data, &ruleConfig, "", "")
			if err != nil {
				return err
			}
			inputArr = append([]generated.InputScopedAnomalyRuleConfigUpdate{
				*firstInput,
			}, inputArr...)
		}
	}
	for _, input := range inputArr {
		_, err := generated.UpdateAnomalyRuleConfig(ctx, *r, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildInputAnomalyScope(data *models.WaapConfigModel) (generated.InputAnomalyScope, error) {
	inputAnomalyScope := generated.InputAnomalyScope{}
	if HasValue(data.Environment) {
		inputAnomalyScope.ScopeType = generated.AnomalyScopeTypeEnvironment
		inputAnomalyScope.EnvironmentScope = &generated.InputAnomalyEnvironmentScope{
			Id: data.Environment.ValueString(),
		}
		return inputAnomalyScope, nil
	}
	inputAnomalyScope.ScopeType = generated.AnomalyScopeTypeCustomer
	return inputAnomalyScope, nil
}

func buildInputAnomalyRuleConfig(ruleConfig *models.WaapRuleConfigModel, subRuleName string, subruleAction string) (generated.InputAnomalyRuleConfigUpdate, error) {
	inputAnomalyRuleConfigUpdate := generated.InputAnomalyRuleConfigUpdate{}
	enabled := !ruleConfig.Enabled.ValueBool()
	internal := false
	configType, err := GetConfigType(ruleConfig.RuleName.ValueString())
	if err != nil {
		return generated.InputAnomalyRuleConfigUpdate{}, err
	}
	inputAnomalyRuleConfigUpdate.ConfigType = generated.AnomalyDetectionConfigType(configType)
	mainRuleName := ruleConfig.RuleName.ValueString()
	mainRuleId, errInGettingMainRuleId := GetRuleId(mainRuleName)
	if errInGettingMainRuleId != nil {
		return generated.InputAnomalyRuleConfigUpdate{}, errInGettingMainRuleId
	}
	inputAnomalyRuleConfigUpdate.RuleId = mainRuleId
	if subRuleName == "" && subruleAction == "" {
		inputAnomalyRuleConfigUpdate.ConfigStatus = &generated.InputAnomalyConfigStatusChange{
			Disabled: &enabled,
			Internal: &internal,
		}
	} else {
		if _, ok := WaapConfigSubRuleActionMap[subruleAction]; !ok {
			return generated.InputAnomalyRuleConfigUpdate{}, utils.NewInvalidError("rule_configs subrules", fmt.Sprintf("%s is not a valid subrule action", subruleAction))
		}
		anomalySubruleAction := generated.AnomalySubRuleAction(subruleAction)

		subruleId, errInGettingSubRuleId := GetSubRuleId(mainRuleName, subRuleName)
		if errInGettingSubRuleId != nil {
			return generated.InputAnomalyRuleConfigUpdate{}, errInGettingSubRuleId
		}
		inputAnomalyRuleConfigUpdate.SubRuleConfigs = []*generated.InputAnomalySubRuleConfigUpdate{
			{
				SubRuleId:            subruleId,
				AnomalySubRuleAction: &anomalySubruleAction,
			},
		}
	}
	return inputAnomalyRuleConfigUpdate, nil
}

func convertWaapConfigModelToUpdateInput(ctx context.Context, data *models.WaapConfigModel, ruleConfig *models.WaapRuleConfigModel, subRuleName string, subruleAction string) (*generated.InputScopedAnomalyRuleConfigUpdate, error) {
	input := generated.InputScopedAnomalyRuleConfigUpdate{}
	inputAnomalyScope, anomalyScopeBuildErr := buildInputAnomalyScope(data)
	if anomalyScopeBuildErr != nil {
		return nil, anomalyScopeBuildErr
	}
	input.AnomalyScope = inputAnomalyScope
	inputRuleConfig, ruleConfigBuildError := buildInputAnomalyRuleConfig(ruleConfig, subRuleName, subruleAction)
	if ruleConfigBuildError != nil {
		return nil, ruleConfigBuildError
	}
	input.RuleConfig = inputRuleConfig
	return &input, nil
}
