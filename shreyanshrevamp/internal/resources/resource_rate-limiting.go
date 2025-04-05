package resources

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/generated"
	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/models"
	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/schemas"
	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/utils"
)

type RateLimitingResource struct {
	client *graphql.Client
}

func NewRateLimitingResource() resource.Resource {
	return &RateLimitingResource{}
}

func (r *RateLimitingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	fmt.Println("hello inside resource")
	tflog.Trace(ctx, "Client Intialization Successfully")
}

func (r *RateLimitingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limiting"
}

func (r *RateLimitingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schemas.RateLimitingResourceSchema()
}

func (r *RateLimitingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.RateLimitingRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Entering in Create Block", map[string]any{
		"shreyanshdata": data,
	})
	ruleInput, err := convertModelToCreateInput(ctx, data)
	if ruleInput == nil || err != nil {
		resp.Diagnostics.AddError("Error converting model to input", err.Error())
		return
	}

	tflog.Trace(ctx, "Entering in Create Block", map[string]any{
		"input": ruleInput,
	})

	resp1, err2 := generated.CreateRateLimitingRule(ctx, *r.client, *ruleInput)
	if err2 != nil {
		resp.Diagnostics.AddError("Error creating rate limiting rule", err2.Error())
		return
	}
	tflog.Trace(ctx, "Entering in Create Block", map[string]any{
		"response": resp1,
	})
	// data1, err3 := convertCreateResponseToModel(ctx, resp1.CreateRateLimitingRule.RateLimitingRuleFields)
	// if err3 != nil {
	// 	resp.Diagnostics.AddError("Error converting create response to model", err3.Error())
	// 	return
	// }
	data.Id = types.StringValue(resp1.CreateRateLimitingRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *RateLimitingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.RateLimitingRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RateLimitingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	input, err := convertModelToUpdateInput(ctx, data, dataState.Id.ValueString())
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

func (r *RateLimitingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

func (r *RateLimitingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

func convertRateLimitingRuleFieldsToModel(ctx context.Context, data *generated.RateLimitingRuleFields) (*models.RateLimitingRuleModel, error) {
	var model *models.RateLimitingRuleModel
	sources := models.Sources{}
	reqresarr := []models.RequestResponseCondition{}
	for _, condition := range data.GetConditions() {
		leafCondition := condition.LeafCondition.LeafConditionFields
		switch string(leafCondition.ConditionType) {
		case "KEY_VALUE":
			reqres := models.RequestResponseCondition{}
			if leafCondition.KeyValueCondition.GetMetadataType() != "" {
				reqres.MetadataType = types.StringValue(string(leafCondition.KeyValueCondition.GetMetadataType()))
			}
			if leafCondition.KeyValueCondition.GetValueCondition() != nil {
				reqres.Value = types.StringValue(leafCondition.KeyValueCondition.GetValueCondition().GetValue())
				reqres.ValueOperator = types.StringValue(string(leafCondition.KeyValueCondition.GetValueCondition().GetOperator()))
			}
			if leafCondition.KeyValueCondition.GetKeyCondition() != nil {
				reqres.KeyOperator = types.StringValue(string(leafCondition.KeyValueCondition.GetKeyCondition().GetOperator()))
				reqres.KeyValue = types.StringValue(leafCondition.KeyValueCondition.GetKeyCondition().GetValue())
			}
			reqresarr = append(reqresarr, reqres)

		case "SCOPE_VALUE":
			if leafCondition.ScopeCondition.GetEntityScope() != nil {
				sources.Endpoints = utils.ConvertStringPtrSliceToTerraformList(leafCondition.ScopeCondition.EntityScope.GetEntityIds())
			}
			if leafCondition.ScopeCondition.GetLabelScope() != nil {
				sources.EndpointLabels = utils.ConvertStringPtrSliceToTerraformList(leafCondition.ScopeCondition.LabelScope.GetLabelIds())
			}

		case "DATATYPE":

		case "IP_ADDRESS":
			sources.IpAddress = &models.IpAddressSource{
				IpAddressList: utils.ConvertStringPtrSliceToTerraformList(leafCondition.IpAddressCondition.GetIpAddresses()),
				Exclude:       types.BoolValue(*leafCondition.IpAddressCondition.GetExclude()),
			}

		case "IP_LOCATION_TYPE":
			iplocationtypes, err := utils.ConvertCustomStringPtrsToTerraformList(leafCondition.IpLocationTypeCondition.GetIpLocationTypes())
			if err != nil {
				fmt.Errorf("error converting ip location types to terraform list: %v", err)
				return nil, err
			}
			sources.IpLocationType = &models.IpLocationTypeSource{
				IpLocationTypes: iplocationtypes,
				Exclude:         types.BoolValue(*leafCondition.IpLocationTypeCondition.GetExclude()),
			}

		case "IP_REPUTATION":
			sources.IpReputation = types.StringValue(string(leafCondition.IpReputationCondition.GetMinIpReputationSeverity()))

		case "REGION":
			regionIds, err := utils.ConvertCustomStringPtrsToTerraformList(leafCondition.RegionCondition.GetRegionIdentifiers())
			if err != nil {
				fmt.Errorf("error converting region identifiers to terraform list: %v", err)
				return nil, err
			}
			sources.Regions = &models.RegionsSource{
				RegionsIds: regionIds,
				Exclude:    types.BoolValue(*leafCondition.RegionCondition.GetExclude()),
			}

		case "EMAIL_DOMAIN":
			sources.EmailDomain = &models.EmailDomainSource{
				EmailDomainRegexes: utils.ConvertStringPtrSliceToTerraformList(leafCondition.EmailDomainCondition.GetEmailRegexes()),
				Exclude:            types.BoolValue(*leafCondition.EmailDomainCondition.GetExclude()),
			}

		case "IP_CONNECTION_TYPE":
			ipConnectionTypeList, err := utils.ConvertCustomStringPtrsToTerraformList(leafCondition.IpConnectionTypeCondition.GetIpConnectionTypes())
			if err != nil {
				fmt.Errorf("error converting ip connection types to terraform list: %v", err)
				return nil, err
			}
			sources.IpConnectionType = &models.IpConnectionTypeSource{
				IpConnectionTypeList: ipConnectionTypeList,
				Exclude:              types.BoolValue(*leafCondition.IpConnectionTypeCondition.GetExclude()),
			}

		case "USER_AGENT":
			sources.UserAgents = &models.UserAgentsSource{
				UserAgentsList: utils.ConvertStringPtrSliceToTerraformList(leafCondition.UserAgentCondition.GetUserAgentRegexes()),
				Exclude:        types.BoolValue(*leafCondition.UserAgentCondition.GetExclude()),
			}

		case "USER_ID":
			sources.UserId = &models.UserIdSource{
				UserIdRegexes: utils.ConvertStringPtrSliceToTerraformList(leafCondition.UserIdCondition.GetUserIdRegexes()),
				Exclude:       types.BoolValue(*leafCondition.UserIdCondition.GetExclude()),
			}

		case "IP_ORGANISATION":
			sources.IpOrganisation = &models.IpOrganisationSource{
				IpOrganisationRegexes: utils.ConvertStringPtrSliceToTerraformList(leafCondition.IpOrganisationCondition.GetIpOrganisationRegexes()),
				Exclude:               types.BoolValue(*leafCondition.IpOrganisationCondition.GetExclude()),
			}

		}

	}
	sources.RequestResponse = reqresarr

	thresholdConfigs := []models.ThresholdConfig{}
	actions := []models.Action{}

	for _, config := range data.GetThresholdActionConfigs() {
		for _, action := range config.GetActions() {
			actiontemp := models.Action{}
			switch string(action.GetActionType()) {
			case "ALERT":
				actiontemp.ActionType = types.StringValue("ALERT")
				actiontemp.EventSeverity = types.StringValue(string(action.Alert.GetEventSeverity()))
				actiontemp.HeaderInjections = []models.HeaderInjection{}
				for _, header := range action.Alert.AgentEffect.GetAgentModifications() {
					headerInj := models.HeaderInjection{
						Key:   types.StringValue(header.HeaderInjection.GetKey()),
						Value: types.StringValue(header.GetHeaderInjection().Value),
					}
					actiontemp.HeaderInjections = append(actiontemp.HeaderInjections, headerInj)
				}
				actions = append(actions, actiontemp)
				//alow missing
			case "MARK_FOR_TESTING":
				actiontemp.ActionType = types.StringValue("MARK_FOR_TESTING")
				actiontemp.EventSeverity = types.StringValue(string(action.GetMarkForTesting().GetEventSeverity()))
				actiontemp.HeaderInjections = []models.HeaderInjection{}
				for _, header := range action.GetMarkForTesting().GetAgentEffect().GetAgentModifications() {
					headerInj := models.HeaderInjection{
						Key:   types.StringValue(header.GetHeaderInjection().Key),
						Value: types.StringValue(header.GetHeaderInjection().Value),
					}
					actiontemp.HeaderInjections = append(actiontemp.HeaderInjections, headerInj)
				}
				actions = append(actions, actiontemp)

			}

		}

		for _, threshold := range config.GetThresholdConfigs() {
			thresholdConfig := models.ThresholdConfig{}

			thresholdConfig.ApiAggregateType = types.StringValue(string(threshold.GetApiAggregateType()))
			thresholdConfig.UserAggregateType = types.StringValue(string(threshold.GetUserAggregateType()))
			thresholdConfig.ThresholdConfigType = types.StringValue(string(threshold.GetThresholdConfigType()))
			if threshold.RollingWindowThresholdConfig != nil {
				thresholdConfig.RollingWindowCountAllowed = types.Int64Value(threshold.GetRollingWindowThresholdConfig().GetCountAllowed())
				thresholdConfig.RollingWindowDuration = types.StringValue(threshold.GetRollingWindowThresholdConfig().GetDuration())
			}

			if threshold.GetDynamicThresholdConfig() != nil {

				thresholdConfig.DynamicDuration = types.StringValue(threshold.GetDynamicThresholdConfig().GetDuration())
				thresholdConfig.DynamicMeanCalculationDuration = types.StringValue(threshold.GetDynamicThresholdConfig().GetMeanCalculationDuration())
				thresholdConfig.DynamicPercentageExcedingMeanAllowed = types.Int64Value(threshold.GetDynamicThresholdConfig().GetPercentageExceedingMeanAllowed())
			}

			thresholdConfigs = append(thresholdConfigs, thresholdConfig)

		}

	}

	model = &models.RateLimitingRuleModel{
		Id:               types.StringValue(data.GetId()),
		Name:             types.StringValue(data.GetName()),
		Enabled:          types.BoolValue(data.GetEnabled()),
		Description:      types.StringValue(*data.GetDescription()),
		Environments:     utils.ConvertStringPtrSliceToTerraformList(data.GetRuleConfigScope().GetEnvironmentScope().GetEnvironmentIds()),
		Sources:          &sources,
		Action:           actions[0],
		ThresholdConfigs: thresholdConfigs,
	}

	return model, nil
}

func convertModelToCreateInput(ctx context.Context, data *models.RateLimitingRuleModel) (*generated.InputRateLimitingRuleData, error) {
	var input *generated.InputRateLimitingRuleData
	name := data.Name.ValueString()
	enabled := data.Enabled.ValueBool()
	category := generated.RateLimitingRuleCategoryEndpointRateLimiting
	tflog.Trace(ctx, "why not category", map[string]any{
		"category": category,
	})

	description := data.Description.ValueString()
	scope, err1 := convertRuleConfigScope(data)
	if err1 != nil {
		return nil, err1
	}
	status, err2 := convertRuleStatus(data)
	if err2 != nil {
		return nil, err2
	}
	thresholdActionConfigs, err3 := convertThresholdActionConfigType(data)
	if err3 != nil {
		return nil, err3
	}
	conditions, err4 := convertToInputRateLimitingCondition(data)
	if err4 != nil {
		return nil, err4
	}

	input = &generated.InputRateLimitingRuleData{
		Category:               category,
		Conditions:             conditions,
		Description:            &description,
		Enabled:                enabled,
		Name:                   name,
		RuleConfigScope:        scope,
		RuleStatus:             status,
		ThresholdActionConfigs: thresholdActionConfigs,
	}
	return input, nil

}

func convertModelToUpdateInput(ctx context.Context, data *models.RateLimitingRuleModel, id string) (*generated.InputRateLimitingRule, error) {
	name := data.Name.ValueString()
	enabled := data.Enabled.ValueBool()
	category := generated.RateLimitingRuleCategoryEndpointRateLimiting
	tflog.Trace(ctx, "why not category", map[string]any{
		"category": category,
	})

	description := data.Description.ValueString()
	scope, err1 := convertRuleConfigScope(data)
	if err1 != nil {
		return nil, err1
	}
	status, err2 := convertRuleStatus(data)
	if err2 != nil {
		return nil, err2
	}
	thresholdActionConfigs, err3 := convertThresholdActionConfigType(data)
	if err3 != nil {
		return nil, err3
	}
	conditions, err4 := convertToInputRateLimitingCondition(data)
	if err4 != nil {
		return nil, err4
	}

	input := &generated.InputRateLimitingRule{
		Id:                     id,
		Category:               category,
		Conditions:             conditions,
		Description:            &description,
		Enabled:                enabled,
		Name:                   name,
		RuleConfigScope:        scope,
		RuleStatus:             status,
		ThresholdActionConfigs: thresholdActionConfigs,
	}
	return input, nil

}

func convertRuleConfigScope(data *models.RateLimitingRuleModel) (*generated.InputRuleConfigScope, error) {
	var scope *generated.InputRuleConfigScope

	if !HasValue(data.Environments) {
		return nil, nil
	}

	var environments []*string

	for _, env := range data.Environments.Elements() {
		if env, ok := env.(types.String); ok {
			env1 := env.ValueString()
			environments = append(environments, &env1)
		}

	}
	scope = &generated.InputRuleConfigScope{
		EnvironmentScope: &generated.InputEnvironmentScope{
			EnvironmentIds: environments,
		},
	}
	return scope, nil
}

func convertRuleStatus(data *models.RateLimitingRuleModel) (*generated.InputRateLimitingRuleStatus, error) {
	var internal = false
	var status *generated.InputRateLimitingRuleStatus
	status = &generated.InputRateLimitingRuleStatus{
		Internal: &internal,
	}
	return status, nil
}

func convertThresholdActionConfigType(data *models.RateLimitingRuleModel) ([]*generated.InputRateLimitingRuleThresholdActionConfig, error) {
	var configTypes []*generated.InputRateLimitingRuleThresholdActionConfig
	var actions []*generated.InputRateLimitingRuleAction
	var thresholdConfigs []*generated.InputRateLimitingRuleThresholdConfig
	if HasValue(data.Action) {
		if HasValue(data.Action.ActionType) {

			switch data.Action.ActionType.ValueString() {

			case "ALERT":
				if HasValue(data.Action.Duration) {
					return nil, fmt.Errorf("Duration should not be in alert")
				}
				if !HasValue(data.Action.EventSeverity) {
					return nil, fmt.Errorf("Event Severity should present")
				}
				agentEffect, err := getRateLimitingRuleAgentEffect(data)
				if err != nil {
					return nil, err
				}

				eventSeverity := data.Action.EventSeverity.ValueString()

				action := &generated.InputRateLimitingRuleAction{
					ActionType: generated.RateLimitingRuleActionTypeAlert,
					Alert: &generated.InputRateLimitingRuleAlertAction{
						EventSeverity: RateLimitingRuleEventSeverityMap[eventSeverity],
						AgentEffect:   agentEffect,
					},
				}
				actions = append(actions, action)

			case "BLOCK":
				if HasValue(data.Action.HeaderInjections) {
					return nil, fmt.Errorf("No Header Injection allowed in Block")
				}
				if !HasValue(data.Action.EventSeverity) {
					return nil, fmt.Errorf("Event Severity should present")
				}
				if !HasValue(data.Action.Duration) {
					return nil, fmt.Errorf("Duration is required")
				}

				duration := data.Action.Duration.ValueString()
				action := &generated.InputRateLimitingRuleAction{
					ActionType: generated.RateLimitingRuleActionTypeBlock,
					Block: &generated.InputRateLimitingRuleBlockAction{
						EventSeverity: RateLimitingRuleEventSeverityMap[data.Action.EventSeverity.ValueString()],
						Duration:      &duration,
					},
				}
				actions = append(actions, action)

			case "ALLOW":
				if HasValue(data.Action.HeaderInjections) {
					return nil, fmt.Errorf("No Header Injection allowed in Block")
				}
				if HasValue(data.Action.EventSeverity) {
					return nil, fmt.Errorf("Event Severity should not  present")
				}
				if !HasValue(data.Action.Duration) {
					return nil, fmt.Errorf("Duration is required")
				}
				duration := data.Action.Duration.ValueString()
				action := &generated.InputRateLimitingRuleAction{
					ActionType: generated.RateLimitingRuleActionTypeAllow,
					Allow: &generated.InputRateLimitingRuleAllowAction{
						Duration: &duration,
					},
				}
				actions = append(actions, action)

			case "MARK_FOR_TESTING":
				if HasValue(data.Action.Duration) {
					return configTypes, fmt.Errorf("Duration should not be in alert")
				}
				if !HasValue(data.Action.EventSeverity) {
					return configTypes, fmt.Errorf("Event Severity should present")
				}
				agentEffect, err := getRateLimitingRuleAgentEffect(data)
				if err != nil {
					return configTypes, err
				}
				eventSeverity := RateLimitingRuleEventSeverityMap[data.Action.EventSeverity.ValueString()]
				action := &generated.InputRateLimitingRuleAction{
					ActionType: generated.RateLimitingRuleActionTypeMarkForTesting,
					MarkForTesting: &generated.InputRateLimitingRuleMarkForTestingAction{
						EventSeverity: eventSeverity,
						AgentEffect:   agentEffect,
					},
				}
				actions = append(actions, action)

			default:
				return nil, fmt.Errorf("Invalid action type: %s", data.Action.ActionType.ValueString())
			}

		}

	}

	if HasValue(data.ThresholdConfigs) {
		for _, config := range data.ThresholdConfigs {

			switch config.ThresholdConfigType.ValueString() {
			case "ROLLING_WINDOW":
				if !HasValue(config.ApiAggregateType) && !(config.ApiAggregateType.ValueString() == "PER_ENDPOINT" || config.ApiAggregateType.ValueString() == "ACROSS_ENDPOINT") {
					return nil, fmt.Errorf("API Aggregate Type Is Not Correct")
				}
				if !HasValue(config.UserAggregateType) && !(config.UserAggregateType.ValueString() == "PER_USER" || config.UserAggregateType.ValueString() == "ACROSS_USER") {
					return nil, fmt.Errorf("User  Aggregate Type Is Not Correct")
				}
				if !HasValue(config.RollingWindowCountAllowed) {
					return nil, fmt.Errorf("Rolling Window Count must present")
				}
				if !HasValue(config.RollingWindowDuration) {
					return nil, fmt.Errorf("Rolling Window Duration must present")
				}

				apiAggregateType := RateLimitingApiAggregateMap[config.ApiAggregateType.ValueString()]
				userAggregateType := RateLimitingUserAggregateMap[config.UserAggregateType.ValueString()]
				cntAllowed := config.RollingWindowCountAllowed.ValueInt64()
				duration := config.RollingWindowDuration.ValueString()

				thresholdConfig := &generated.InputRateLimitingRuleThresholdConfig{
					ThresholdConfigType: generated.RateLimitingRuleThresholdConfigTypeRollingWindow,
					ApiAggregateType:    apiAggregateType,
					UserAggregateType:   userAggregateType,
					RollingWindowThresholdConfig: &generated.InputRollingWindowThresholdConfig{
						CountAllowed: cntAllowed,
						Duration:     duration,
					},
				}
				thresholdConfigs = append(thresholdConfigs, thresholdConfig)

			case "DYNMAIC":
				if HasValue(config.ApiAggregateType) {
					return nil, fmt.Errorf("API Aggregate Type should  not Present")
				}
				if HasValue(config.UserAggregateType) {
					return nil, fmt.Errorf("user  aggregate Type should  Not")
				}
				if !HasValue(config.DynamicDuration) {
					return nil, fmt.Errorf("dynamic duration must present")
				}
				if !HasValue(config.DynamicMeanCalculationDuration) {
					return nil, fmt.Errorf("dynamic mean calculation duration should present")
				}
				if !HasValue(config.DynamicPercentageExcedingMeanAllowed) {
					return nil, fmt.Errorf("dynamic percentage exceding mean allowed should present")
				}
				percentageExceedingMeanAllowed := config.DynamicPercentageExcedingMeanAllowed.ValueInt64()
				meanCalculationDuration := config.DynamicMeanCalculationDuration.ValueString()
				duration := config.DynamicDuration.ValueString()

				thresholdConfig := &generated.InputRateLimitingRuleThresholdConfig{
					ThresholdConfigType: generated.RateLimitingRuleThresholdConfigTypeDynamic,
					ApiAggregateType:    generated.RateLimitingRuleApiAggregateTypePerEndpoint,
					UserAggregateType:   generated.RateLimitingRuleUserAggregateTypePerUser,
					DynamicThresholdConfig: &generated.InputDynamicThresholdConfig{
						PercentageExceedingMeanAllowed: percentageExceedingMeanAllowed,
						MeanCalculationDuration:        meanCalculationDuration,
						Duration:                       duration,
					},
				}
				thresholdConfigs = append(thresholdConfigs, thresholdConfig)
			default:
				return nil, fmt.Errorf("Invalid Threshold Config Type %s", config.ThresholdConfigType.ValueString())

			}

		}

	}

	configType := &generated.InputRateLimitingRuleThresholdActionConfig{
		Actions:          actions,
		ThresholdConfigs: thresholdConfigs,
	}

	configTypes = append(configTypes, configType)
	return configTypes, nil

}

func getRateLimitingRuleAgentEffect(data *models.RateLimitingRuleModel) (*generated.InputRateLimitingRuleAgentEffect, error) {
	var agentEffect *generated.InputRateLimitingRuleAgentEffect
	agentModifications := []*generated.InputRateLimitingRuleAgentModification{}

	if HasValue(data.Action.HeaderInjections) {
		for _, injection := range data.Action.HeaderInjections {
			if !HasValue(injection.Key) {
				return nil, fmt.Errorf("Key should be present")
			}
			if !HasValue(injection.Value) {
				return nil, fmt.Errorf("Key should be present")
			}
			key := injection.Key.ValueString()
			value := injection.Value.ValueString()
			if key == "" && value == "" {
				return nil, fmt.Errorf("In Header Injection key and value can not be empty string ")
			}
			temp := &generated.InputRateLimitingRuleAgentModification{
				AgentModificationType: generated.RateLimitingRuleAgentModificationTypeHeaderInjection,
				HeaderInjection: generated.InputRateLimitingRuleHeaderInjection{
					HeaderCategory: generated.RateLimitingRuleMatchCategoryRequest,
					Key:            key,
					Value:          value,
				},
			}
			agentModifications = append(agentModifications, temp)
		}
	}
	agentEffect = &generated.InputRateLimitingRuleAgentEffect{
		AgentModifications: agentModifications,
	}
	return agentEffect, nil

}

func convertToInputRateLimitingCondition(data *models.RateLimitingRuleModel) ([]*generated.InputRateLimitingRuleCondition, error) {

	var conditions []*generated.InputRateLimitingRuleCondition

	if HasValue(data.Sources) {

		if HasValue(data.Sources.Scanner) {

			if !HasValue(data.Sources.Scanner.ScannerTypesList) {
				return nil, fmt.Errorf("scanners type list can not be empty")
			}
			if !HasValue(data.Sources.Scanner.Exclude) {
				return nil, fmt.Errorf("exclude should be present")
			}
			var scannerTypes []*string

			for _, scanner := range data.Sources.Scanner.ScannerTypesList.Elements() {

				if scanner, ok := scanner.(types.String); ok {
					if !RateLimitingRuleScannerMap[scanner.ValueString()] {
						return nil, fmt.Errorf("Scanner %s is not a valid scanner type", scanner.ValueString())
					}
					sc := scanner.ValueString()
					scannerTypes = append(scannerTypes, &sc)
				}
			}
			exclude := data.Sources.Scanner.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{

				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{

					ConditionType: generated.RateLimitingRuleLeafConditionTypeRequestScannerType,
					RequestScannerTypeCondition: &generated.InputRateLimitingRuleRequestScannerTypeCondition{
						Exclude:      &exclude,
						ScannerTypes: scannerTypes,
					},
				},
			}

			conditions = append(conditions, &input)

		}
		if HasValue(data.Sources.IpAsn) {
			if !HasValue(data.Sources.IpAsn.IpAsnRegexes) {
				return nil, fmt.Errorf("ip asn regexes can not be empty")
			}
			if !HasValue(data.Sources.IpAsn.Exclude) {
				return nil, fmt.Errorf("exclude should be present")
			}
			var ipAsnRegexes []*string

			for _, ipAsnRegex := range data.Sources.IpAsn.IpAsnRegexes.Elements() {

				if ipAsnRegex, ok := ipAsnRegex.(types.String); ok {
					sc := ipAsnRegex.ValueString()
					ipAsnRegexes = append(ipAsnRegexes, &sc)
				}
			}
			exclude := data.Sources.IpAsn.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{

				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{

					ConditionType: generated.RateLimitingRuleLeafConditionTypeIpAsn,
					IpAsnCondition: &generated.InputRateLimitingRuleIpAsnCondition{
						Exclude:      &exclude,
						IpAsnRegexes: ipAsnRegexes,
					},
				},
			}

			conditions = append(conditions, &input)

		}

		if HasValue(data.Sources.IpConnectionType) {
			if !HasValue(data.Sources.IpConnectionType.IpConnectionTypeList) {
				return nil, fmt.Errorf("ip connection type list can not be empty")
			}
			if !HasValue(data.Sources.IpConnectionType.Exclude) {
				return nil, fmt.Errorf("exclude should be present")
			}
			var ipConnectionTypes []*generated.RateLimitingRuleIpConnectionType

			for _, ipConnectionType := range data.Sources.IpConnectionType.IpConnectionTypeList.Elements() {

				if ipConnectionType, ok := ipConnectionType.(types.String); ok {
					connection := ipConnectionType.ValueString()
					val, exist := RateLimitingRuleIpConnectionTypeMap[connection]
					if !exist {
						return nil, fmt.Errorf("Ip connection type %s is not a valid ip connection type", connection)
					}
					ipConnectionTypes = append(ipConnectionTypes, &val)
				}
			}
			exclude := data.Sources.IpConnectionType.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{

				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{

					ConditionType: generated.RateLimitingRuleLeafConditionTypeIpConnectionType,
					IpConnectionTypeCondition: &generated.InputRateLimitingRuleIpConnectionTypeCondition{
						Exclude:           &exclude,
						IpConnectionTypes: ipConnectionTypes,
					},
				},
			}

			conditions = append(conditions, &input)
		}

		if HasValue(data.Sources.UserId) {

			if HasValue(data.Sources.UserId.UserIdRegexes) && HasValue(data.Sources.UserId.UserIds) {
				return nil, fmt.Errorf("user id regexes and user ids can not be present together")
			}
			if !HasValue(data.Sources.UserId.Exclude) {
				return nil, fmt.Errorf("exclude should be present")
			}
			exclude := data.Sources.UserId.Exclude.ValueBool()
			if HasValue(data.Sources.UserId.UserIds) {
				userIds := []*string{}
				for _, userId := range data.Sources.UserId.UserIds.Elements() {
					if userId, ok := userId.(types.String); ok {
						id := userId.ValueString()
						userIds = append(userIds, &id)

					}
				}
				var input = generated.InputRateLimitingRuleCondition{

					LeafCondition: &generated.InputRateLimitingRuleLeafCondition{

						ConditionType: generated.RateLimitingRuleLeafConditionTypeUserId,
						UserIdCondition: &generated.InputRateLimitingRuleUserIdCondition{
							Exclude: &exclude,
							UserIds: userIds,
						},
					},
				}
				conditions = append(conditions, &input)
			}
			if HasValue(data.Sources.UserId.UserIdRegexes) {
				userIdRegexes := []*string{}
				for _, userIdRegex := range data.Sources.UserId.UserIdRegexes.Elements() {
					if userIdRegex, ok := userIdRegex.(types.String); ok {
						id := userIdRegex.ValueString()
						userIdRegexes = append(userIdRegexes, &id)

					}
				}
				var input = generated.InputRateLimitingRuleCondition{
					LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
						ConditionType: generated.RateLimitingRuleLeafConditionTypeUserId,
						UserIdCondition: &generated.InputRateLimitingRuleUserIdCondition{
							Exclude:       &exclude,
							UserIdRegexes: userIdRegexes,
						},
					},
				}
				conditions = append(conditions, &input)
			}
		}
		// labels ki id chaiye
		if HasValue(data.Sources.EndpointLabels) {
			labelIds := []*string{}
			for _, label := range data.Sources.EndpointLabels.Elements() {
				if label, ok := label.(types.String); ok {
					id := label.ValueString()
					labelIds = append(labelIds, &id)

				}
			}
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeScope,
					ScopeCondition: &generated.InputRateLimitingRuleScopeCondition{
						ScopeType: generated.RateLimitingRuleScopeConditionTypeLabel,
						LabelScope: &generated.InputRateLimitingRuleLabelScope{
							LabelIds:  labelIds,
							LabelType: generated.RateLimitingRuleLabelTypeApi,
						},
					},
				},
			}
			conditions = append(conditions, &input)
		}

		//endpoint id chaiye
		if HasValue(data.Sources.Endpoints) {

			endpointIds := []*string{}
			for _, endpoint := range data.Sources.Endpoints.Elements() {
				if endpoint, ok := endpoint.(types.String); ok {
					id := endpoint.ValueString()
					endpointIds = append(endpointIds, &id)

				}
			}
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeScope,
					ScopeCondition: &generated.InputRateLimitingRuleScopeCondition{
						ScopeType: generated.RateLimitingRuleScopeConditionTypeEntity,
						EntityScope: &generated.InputRateLimitingRuleEntityScope{
							EntityIds:  endpointIds,
							EntityType: generated.RateLimitingRuleEntityTypeApi,
						},
					},
				},
			}
			conditions = append(conditions, &input)

		}

		if HasValue(data.Sources.Attributes) {

			for _, attribute := range data.Sources.Attributes {

				if !HasValue(attribute.KeyConditionOperator) {
					return nil, fmt.Errorf("KeyConditionOperator should be present")
				}
				if !HasValue(attribute.ValueConditionOperator) {
					return nil, fmt.Errorf("ValueConditionOperator should be present")
				}
				keyConditionOperator, exist := RateLimitingKeyValueMatchOperatorMap[attribute.KeyConditionOperator.ValueString()]
				if !exist {
					return nil, fmt.Errorf("Invalid KeyConditionOperator")
				}
				valueConditionOperator, exist := RateLimitingKeyValueMatchOperatorMap[attribute.ValueConditionOperator.ValueString()]
				if !exist {
					return nil, fmt.Errorf("Invalid ValueConditionOperator")
				}
				keyConditionValue := attribute.KeyConditionValue.ValueString()
				valueConditionValue := attribute.ValueConditionValue.ValueString()

				var input = generated.InputRateLimitingRuleCondition{
					LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
						ConditionType: generated.RateLimitingRuleLeafConditionTypeKeyValue,
						KeyValueCondition: &generated.InputRateLimitingRuleKeyValueCondition{
							MetadataType: generated.RateLimitingRuleKeyValueConditionMetadataTypeTag,
							KeyCondition: &generated.InputRateLimitingRuleStringCondition{
								Operator: keyConditionOperator,
								Value:    keyConditionValue,
							},
							ValueCondition: &generated.InputRateLimitingRuleStringCondition{
								Operator: valueConditionOperator,
								Value:    valueConditionValue,
							},
						},
					},
				}

				conditions = append(conditions, &input)
			}

		}

		if HasValue(data.Sources.IpReputation) {
			minIpReputationSeverity, exist := RateLimitingRuleIpReputationSeverityMap[data.Sources.IpReputation.ValueString()]
			if !exist {
				return nil, fmt.Errorf("Invalid Ip Reputation Severity")
			}
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeIpReputation,
					IpReputationCondition: &generated.InputRateLimitingRuleIpReputationCondition{
						MinIpReputationSeverity: minIpReputationSeverity,
					},
				},
			}
			conditions = append(conditions, &input)

		}

		if HasValue(data.Sources.IpLocationType) {
			if !HasValue(data.Sources.IpLocationType.IpLocationTypes) {
				return nil, fmt.Errorf("IpLocationTypes should be present")
			}
			if !HasValue(data.Sources.IpLocationType.Exclude) {
				return nil, fmt.Errorf("Exclude should be present")
			}
			ipLocationTypes := []*generated.RateLimitingRuleIpLocationType{}
			for _, ipLocationType := range data.Sources.IpLocationType.IpLocationTypes.Elements() {
				if locationType, ok := ipLocationType.(types.String); ok {

					ipLocationType, exist := RateLimitingRuleIpLocationTypeMap[locationType.ValueString()]
					if !exist {
						return nil, fmt.Errorf("Invalid Ip Location Type")
					}
					ipLocationTypes = append(ipLocationTypes, &ipLocationType)
				}
			}
			exclude := data.Sources.IpLocationType.Exclude.ValueBool()

			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeIpLocationType,
					IpLocationTypeCondition: &generated.InputRateLimitingRuleIpLocationTypeCondition{
						IpLocationTypes: ipLocationTypes,
						Exclude:         &exclude,
					},
				},
			}
			conditions = append(conditions, &input)

		}

		if HasValue(data.Sources.IpAbuseVelocity) {
			minIpAbuseVelocity, exist := RateLimitingRuleIpAbuseVelocityMap[data.Sources.IpAbuseVelocity.ValueString()]
			if !exist {
				return nil, fmt.Errorf("Invalid Ip Abuse Velocity")
			}
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeIpAbuseVelocity,
					IpAbuseVelocityCondition: &generated.InputRateLimitingRuleIpAbuseVelocityCondition{
						MinIpAbuseVelocity: minIpAbuseVelocity,
					},
				},
			}
			conditions = append(conditions, &input)

		}

		if HasValue(data.Sources.IpAddress) {

			if !HasValue(data.Sources.IpAddress.IpAddressList) {
				return nil, fmt.Errorf("Ip Address List should be present")
			}
			if !HasValue(data.Sources.IpAddress.Exclude) {
				return nil, fmt.Errorf("Exclude should be present")
			}
			ipAddresses := []*string{}
			for _, ipAddress := range data.Sources.IpAddress.IpAddressList.Elements() {
				if ip, ok := ipAddress.(types.String); ok {
					ipAddr := ip.ValueString()
					ipAddresses = append(ipAddresses, &ipAddr)
				}
			}
			exclude := data.Sources.IpAddress.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeIpAddress,
					IpAddressCondition: &generated.InputRateLimitingRuleIpAddressCondition{
						RawInputIpData: ipAddresses,
						Exclude:        &exclude,
					},
				},
			}
			conditions = append(conditions, &input)

		}

		if HasValue(data.Sources.EmailDomain) {
			if !HasValue(data.Sources.EmailDomain.Exclude) {
				return nil, fmt.Errorf("Exclude should be present")
			}
			if !HasValue(data.Sources.EmailDomain.EmailDomainRegexes) {
				return nil, fmt.Errorf("Email Domain List should be present")
			}
			emailRegexes := []*string{}
			for _, emailRegex := range data.Sources.EmailDomain.EmailDomainRegexes.Elements() {
				if emailRegex, ok := emailRegex.(types.String); ok {
					emailRegexStr := emailRegex.ValueString()
					emailRegexes = append(emailRegexes, &emailRegexStr)
				}
			}
			exclude := data.Sources.EmailDomain.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeEmailDomain,
					EmailDomainCondition: &generated.InputRateLimitingRuleEmailDomainCondition{
						EmailRegexes: emailRegexes,
						Exclude:      &exclude,
					},
				},
			}
			conditions = append(conditions, &input)

		}

		if HasValue(data.Sources.UserAgents) {
			if !HasValue(data.Sources.UserAgents.Exclude) {
				return nil, fmt.Errorf("Exclude should be present")
			}
			if !HasValue(data.Sources.UserAgents.UserAgentsList) {
				return nil, fmt.Errorf("User Agents List should be present")
			}
			userAgents := []*string{}
			for _, userAgent := range data.Sources.UserAgents.UserAgentsList.Elements() {
				if userAgent, ok := userAgent.(types.String); ok {
					userAgentStr := userAgent.ValueString()
					userAgents = append(userAgents, &userAgentStr)
				}
			}
			exclude := data.Sources.UserAgents.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeUserAgent,
					UserAgentCondition: &generated.InputRateLimitingRuleUserAgentCondition{
						UserAgentRegexes: userAgents,
						Exclude:          &exclude,
					},
				},
			}
			conditions = append(conditions, &input)
		}
		// region alpha 2 iso code required
		if HasValue(data.Sources.Regions) {
			if !HasValue(data.Sources.Regions.Exclude) {
				return nil, fmt.Errorf("Exclude should be present")
			}
			if !HasValue(data.Sources.Regions.RegionsIds) {
				return nil, fmt.Errorf("Regions List should be present")
			}
			regionIdentifieres := []*generated.InputRateLimitingRegionIdentifier{}
			for _, region := range data.Sources.Regions.RegionsIds.Elements() {
				if region, ok := region.(types.String); ok {
					identifiers := &generated.InputRateLimitingRegionIdentifier{
						CountryIsoCode: region.ValueString(),
					}
					regionIdentifieres = append(regionIdentifieres, identifiers)
				}
			}
			exclude := data.Sources.Regions.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeRegion,
					RegionCondition: &generated.InputRateLimitingRuleRegionCondition{
						RegionIdentifiers: regionIdentifieres,
						Exclude:           &exclude,
					},
				},
			}
			conditions = append(conditions, &input)
		}

		if HasValue(data.Sources.IpOrganisation) {
			if !HasValue(data.Sources.IpOrganisation.Exclude) {
				return nil, fmt.Errorf("Exclude should be present")
			}
			if !HasValue(data.Sources.IpOrganisation.IpOrganisationRegexes) {
				return nil, fmt.Errorf("Ip Organisation Regexes should be present")
			}
			ipOrganisationRegexes := []*string{}
			for _, ipOrganisationRegex := range data.Sources.IpOrganisation.IpOrganisationRegexes.Elements() {
				if ipOrganisationRegex, ok := ipOrganisationRegex.(types.String); ok {
					ipOrganisationRegexStr := ipOrganisationRegex.ValueString()
					ipOrganisationRegexes = append(ipOrganisationRegexes, &ipOrganisationRegexStr)
				}
			}
			exclude := data.Sources.IpOrganisation.Exclude.ValueBool()
			var input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeIpOrganisation,
					IpOrganisationCondition: &generated.InputRateLimitingRuleIpOrganisationCondition{
						IpOrganisationRegexes: ipOrganisationRegexes,
						Exclude:               &exclude,
					},
				},
			}
			conditions = append(conditions, &input)
		}

		if HasValue(data.Sources.RequestResponse) {

			for _, requestResponse := range data.Sources.RequestResponse {

				if !HasValue(requestResponse.MetadataType) {
					return nil, fmt.Errorf("Metadata Type should be present")
				}
				metadataType, exists := RateLimitingRuleKeyValueConditionMetadataTypeMap[requestResponse.MetadataType.ValueString()]
				if !exists {
					return nil, fmt.Errorf("Invalid Metadata Type")
				}
				var keyConditionOperator generated.RateLimitingRuleKeyValueMatchOperator
				var keyConditionValue string
				var valueConditionOperator generated.RateLimitingRuleKeyValueMatchOperator
				var valueConditionValue string
				if HasValue(requestResponse.KeyOperator) {
					keyConditionOperator = RateLimitingKeyValueMatchOperatorMap[requestResponse.KeyOperator.ValueString()]
				}
				if HasValue(requestResponse.ValueOperator) {
					valueConditionOperator = RateLimitingKeyValueMatchOperatorMap[requestResponse.ValueOperator.ValueString()]
				}

				if HasValue(requestResponse.KeyValue) {
					keyConditionValue = requestResponse.KeyValue.ValueString()
				}
				if HasValue(requestResponse.Value) {
					valueConditionValue = requestResponse.Value.ValueString()
				}

				var input = generated.InputRateLimitingRuleCondition{
					LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
						ConditionType: generated.RateLimitingRuleLeafConditionTypeKeyValue,
						KeyValueCondition: &generated.InputRateLimitingRuleKeyValueCondition{
							MetadataType: metadataType,
							KeyCondition: &generated.InputRateLimitingRuleStringCondition{
								Operator: keyConditionOperator,
								Value:    keyConditionValue,
							},
							ValueCondition: &generated.InputRateLimitingRuleStringCondition{
								Operator: valueConditionOperator,
								Value:    valueConditionValue,
							},
						},
					},
				}

				conditions = append(conditions, &input)
			}
		}

	}

	return conditions, nil

}
