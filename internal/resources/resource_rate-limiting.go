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

type RateLimitingResource struct {
	client *graphql.Client
}

func NewRateLimitingResource() resource.Resource {
	return &RateLimitingResource{}
}

func (r *RateLimitingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RateLimitingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limiting"
}

func (r *RateLimitingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.RateLimitingResourceSchema()
}

func (r *RateLimitingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Entering in Create Block")
	var data *models.RateLimitingRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleInput, err := convertRateLimitingModelToCreateInput(ctx, data)
	if ruleInput == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	id, err := getRateLimitingRuleId(ruleInput.Name, ctx, *r.client)
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

func (r *RateLimitingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.RateLimitingRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getRateLimitingRule(data.Id.ValueString(), ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	updatedData, err := convertRateLimitingRuleFieldsToModel(ctx, &response)
	tflog.Trace(ctx, "Shreyansh Gupta 3")

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedData)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	input, err := convertRateLimitingModelToUpdateInput(ctx, data, dataState.Id.ValueString())
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
	ruleName := req.ID
	id, err := getRateLimitingRuleId(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Resource Not Found", fmt.Sprintf("%s rule of this name not found", ruleName))
		return
	}
	response, err := getRateLimitingRule(id, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data, err := convertRateLimitingRuleFieldsToModel(ctx, &response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func getRateLimitingRule(id string, ctx context.Context, r graphql.Client) (generated.RateLimitingRuleFields, error) {
	rateLimitingfields := generated.RateLimitingRuleFields{}
	category := []*generated.RateLimitingRuleCategory{}
	endpointRateLimiting := generated.RateLimitingRuleCategoryEndpointRateLimiting
	category = append(category, &endpointRateLimiting)
	response, err := generated.GetRateLimitingDetails(ctx, r, category, nil)
	if err != nil {
		return rateLimitingfields, err
	}

	for _, rule := range response.RateLimitingRules.Results {
		if rule.Id == id {
			rateLimitingfields = rule.RateLimitingRuleFields
			return rateLimitingfields, nil
		}
	}

	return rateLimitingfields, nil
}

func getRateLimitingRuleId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
	category := []*generated.RateLimitingRuleCategory{}
	endpointRateLimiting := generated.RateLimitingRuleCategoryEndpointRateLimiting
	category = append(category, &endpointRateLimiting)

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
func convertRateLimitingModelToCreateInput(ctx context.Context, data *models.RateLimitingRuleModel) (*generated.InputRateLimitingRuleData, error) {
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
	category := generated.RateLimitingRuleCategoryEndpointRateLimiting
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
	conditions, err := convertToRateLimitingRuleCondition(data)
	if err != nil {
		return nil, err
	} else {
		input.Conditions = conditions
	}

	return &input, nil
}

func convertRateLimitingModelToUpdateInput(ctx context.Context, data *models.RateLimitingRuleModel, id string) (*generated.InputRateLimitingRule, error) {
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
	category := generated.RateLimitingRuleCategoryEndpointRateLimiting
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
	conditions, err := convertToRateLimitingRuleCondition(data)
	if err != nil {
		return nil, err
	} else {
		input.Conditions = conditions
	}

	return &input, nil

}

func convertRateLimitingRuleFieldsToModel(ctx context.Context, data *generated.RateLimitingRuleFields) (*models.RateLimitingRuleModel, error) {
	model := models.RateLimitingRuleModel{}

	if data.Id != "" {
		model.Id = types.StringValue(data.Id)
	}
	if data.Name != "" {
		model.Name = types.StringValue(data.Name)
	}
	if data.Description != nil {
		model.Description = types.StringValue(*data.Description)
	}
	if data.RuleConfigScope != nil && data.RuleConfigScope.EnvironmentScope != nil {
		environments, err := utils.ConvertStringPtrToTerraformSet(data.RuleConfigScope.EnvironmentScope.EnvironmentIds)
		if err != nil {
			return nil, err
		}
		model.Environments = environments
		tflog.Trace(ctx, "Shreyansh Gupta 123", map[string]interface{}{
			"environments": environments,
		})

	} else {
		model.Environments = types.SetNull(types.StringType)
	}
	model.Enabled = types.BoolValue(data.Enabled)
	endpointScope := false
	endpointLabelScope := false

	if data.Conditions != nil {
		sources := models.RateLimitingSources{}
		reqresarr := []models.RateLimitingRequestResponseCondition{}
		for _, condition := range data.GetConditions() {
			leafCondition := condition.LeafCondition.LeafConditionFields
			switch string(leafCondition.ConditionType) {
			case "KEY_VALUE":
				reqres := models.RateLimitingRequestResponseCondition{}
				if leafCondition.KeyValueCondition.GetMetadataType() != nil {
					reqres.MetadataType = types.StringValue(string(*leafCondition.KeyValueCondition.GetMetadataType()))
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

			case "SCOPE":
				if leafCondition.ScopeCondition.GetEntityScope() != nil {
					endpoints, err := utils.ConvertStringPtrToTerraformSet(leafCondition.ScopeCondition.EntityScope.GetEntityIds())
					if err != nil {
						return nil, err
					}

					sources.Endpoints = endpoints
					endpointScope = true

				}

				if leafCondition.ScopeCondition.GetLabelScope() != nil {
					endpointLabels, err := utils.ConvertStringPtrToTerraformSet(leafCondition.ScopeCondition.LabelScope.GetLabelIds())
					if err != nil {
						return nil, err
					}
					sources.EndpointLabels = endpointLabels
					endpointLabelScope = true
				}

			case "DATATYPE":

			case "IP_ADDRESS":
				ipAddressList, err := utils.ConvertStringPtrToTerraformSet(leafCondition.IpAddressCondition.GetIpAddresses())
				if err != nil {
					return nil, err
				}
				sources.IpAddress = &models.RateLimitingIpAddressSource{
					IpAddressList: ipAddressList,
					Exclude:       types.BoolValue(*leafCondition.IpAddressCondition.GetExclude()),
				}

			case "IP_LOCATION_TYPE":
				iplocationtypes, err := utils.ConvertCustomStringPtrsToTerraformSet(leafCondition.IpLocationTypeCondition.GetIpLocationTypes())
				if err != nil {

					return nil, fmt.Errorf("error converting ip location types to terraform list: %v", err)
				}
				sources.IpLocationType = &models.RateLimitingIpLocationTypeSource{
					IpLocationTypes: iplocationtypes,
					Exclude:         types.BoolValue(*leafCondition.IpLocationTypeCondition.GetExclude()),
				}

			case "IP_REPUTATION":
				sources.IpReputation = types.StringValue(string(leafCondition.IpReputationCondition.GetMinIpReputationSeverity()))

			case "REGION":
				regionIdsPointer := []*string{}
				for _, region := range leafCondition.RegionCondition.GetRegionIdentifiers() {
					regionIdsPointer = append(regionIdsPointer, &region.CountryIsoCode)
				}
				regionIds, err := utils.ConvertStringPtrToTerraformSet(regionIdsPointer)
				if err != nil {
					return nil, err
				}
				sources.Regions = &models.RateLimitingRegionsSource{
					RegionsIds: regionIds,
					Exclude:    types.BoolValue(*leafCondition.RegionCondition.GetExclude()),
				}

			case "EMAIL_DOMAIN":
				emailDomainRegexes, err := utils.ConvertStringPtrToTerraformSet(leafCondition.EmailDomainCondition.GetEmailRegexes())
				if err != nil {
					return nil, err
				}
				sources.EmailDomain = &models.RateLimitingEmailDomainSource{
					EmailDomainRegexes: emailDomainRegexes,
					Exclude:            types.BoolValue(*leafCondition.EmailDomainCondition.GetExclude()),
				}

			case "IP_CONNECTION_TYPE":

				ipConnectionTypeList, err := utils.ConvertCustomStringPtrsToTerraformSet(leafCondition.IpConnectionTypeCondition.GetIpConnectionTypes())
				if err != nil {
					return nil, fmt.Errorf("error converting ip connection types to terraform list: %v", err)
				}
				sources.IpConnectionType = &models.RateLimitingIpConnectionTypeSource{
					IpConnectionTypeList: ipConnectionTypeList,
					Exclude:              types.BoolValue(*leafCondition.IpConnectionTypeCondition.GetExclude()),
				}

			case "USER_AGENT":
				userAgentsList, err := utils.ConvertCustomStringPtrsToTerraformSet(leafCondition.UserAgentCondition.GetUserAgentRegexes())
				if err != nil {
					return nil, err
				}

				sources.UserAgents = &models.RateLimitingUserAgentsSource{
					UserAgentsList: userAgentsList,
					Exclude:        types.BoolValue(*leafCondition.UserAgentCondition.GetExclude()),
				}

			case "USER_ID":

				userIdRegexes, err := utils.ConvertStringPtrToTerraformSet(leafCondition.UserIdCondition.GetUserIdRegexes())
				if err != nil {
					return nil, err
				}

				userIds, err := utils.ConvertStringPtrToTerraformSet(leafCondition.UserIdCondition.UserIds)
				if err != nil {
					return nil, err
				}

				sources.UserId = &models.RateLimitingUserIdSource{
					UserIdRegexes: userIdRegexes,
					UserIds:       userIds,
					Exclude:       types.BoolValue(*leafCondition.UserIdCondition.GetExclude()),
				}

			case "IP_ORGANISATION":
				ipOrganisationRegexes, err := utils.ConvertStringPtrToTerraformSet(leafCondition.IpOrganisationCondition.GetIpOrganisationRegexes())
				if err != nil {
					return nil, err
				}
				sources.IpOrganisation = &models.RateLimitingIpOrganisationSource{
					IpOrganisationRegexes: ipOrganisationRegexes,
					Exclude:               types.BoolValue(*leafCondition.IpOrganisationCondition.GetExclude()),
				}

			case "IP_ASN":
				ipAsnRegexes, err := utils.ConvertStringPtrToTerraformSet(leafCondition.IpAsnCondition.GetIpAsnRegexes())
				if err != nil {
					return nil, err
				}
				sources.IpAsn = &models.RateLimitingIpAsnSource{
					IpAsnRegexes: ipAsnRegexes,
					Exclude:      types.BoolValue(*leafCondition.IpAsnCondition.GetExclude()),
				}
			}

		}
		if !endpointScope {
			sources.Endpoints = types.SetNull(types.StringType)
		}
		if !endpointLabelScope {
			sources.EndpointLabels = types.SetNull(types.StringType)
		}
		reqresset, diags := types.SetValueFrom(
			ctx,
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"metadata_type":  types.StringType,
					"value":          types.StringType,
					"key_operator":   types.StringType,
					"key_value":      types.StringType,
					"value_operator": types.StringType,
				},
			},
			reqresarr,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("Threshold config conversion failed")
		}
		sources.RequestResponse = reqresset
		model.Sources = &sources
	}
	if data.GetThresholdActionConfigs() != nil && len(data.ThresholdActionConfigs) == 1 {

		config := data.ThresholdActionConfigs[0]

		if config.Actions != nil && len(config.Actions) > 0 {
			actions := []models.RateLimitingAction{}

			for _, action := range config.GetActions() {

				actiontemp := models.RateLimitingAction{}
				switch string(action.GetActionType()) {
				case "ALERT":
					actiontemp.ActionType = types.StringValue("ALERT")
					if action.Alert != nil {
						actiontemp.EventSeverity = types.StringValue(string(action.Alert.GetEventSeverity()))
					}
					if action.Alert.AgentEffect != nil {
						headerInjections := []models.RateLimitingHeaderInjection{}
						for _, header := range action.Alert.AgentEffect.GetAgentModifications() {

							headerInj := models.RateLimitingHeaderInjection{
								Key:   types.StringValue(header.HeaderInjection.GetKey()),
								Value: types.StringValue(header.GetHeaderInjection().Value),
							}
							headerInjections = append(headerInjections, headerInj)
						}
						headerInjectionSet, diag := types.SetValueFrom(
							ctx,
							types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"key":   types.StringType,
									"value": types.StringType,
								},
							},
							headerInjections,
						)
						if diag.HasError() {
							return nil, fmt.Errorf("Threshold config conversion failed")
						}
						actiontemp.HeaderInjections = headerInjectionSet
					} else {
						actiontemp.HeaderInjections = types.SetNull(types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"key":   types.StringType,
								"value": types.StringType,
							},
						})
					}

					actions = append(actions, actiontemp)
					//alow missing
				case "BLOCK":
					actiontemp.ActionType = types.StringValue("BLOCK")
					if action.Block != nil {
						actiontemp.EventSeverity = types.StringValue(string(action.Block.GetEventSeverity()))
						actiontemp.Duration = types.StringValue(*action.Block.GetDuration())
					}
					actions = append(actions, actiontemp)

				case "MARK_FOR_TESTING":
					actiontemp.ActionType = types.StringValue("MARK_FOR_TESTING")
					actiontemp.EventSeverity = types.StringValue(string(action.GetMarkForTesting().GetEventSeverity()))
					headerInjections := []models.RateLimitingHeaderInjection{}
					for _, header := range action.Alert.AgentEffect.GetAgentModifications() {

						headerInj := models.RateLimitingHeaderInjection{
							Key:   types.StringValue(header.HeaderInjection.GetKey()),
							Value: types.StringValue(header.GetHeaderInjection().Value),
						}
						headerInjections = append(headerInjections, headerInj)
					}
					headerInjectionSet, diag := types.SetValueFrom(
						ctx,
						types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"key":   types.StringType,
								"value": types.StringType,
							},
						},
						headerInjections,
					)
					if diag.HasError() {
						return nil, fmt.Errorf("Threshold config conversion failed")
					}
					actiontemp.HeaderInjections = headerInjectionSet

				}

			}
			model.Action = actions[0]

		}

		if config.GetThresholdConfigs() != nil && len(config.ThresholdConfigs) > 0 {
			thresholdConfigs := []models.RateLimitingThresholdConfig{}

			for _, threshold := range config.GetThresholdConfigs() {

				thresholdConfig := models.RateLimitingThresholdConfig{}

				thresholdConfig.ThresholdConfigType = types.StringValue(string(threshold.GetThresholdConfigType()))
				if threshold.RollingWindowThresholdConfig != nil {
					thresholdConfig.ApiAggregateType = types.StringValue(string(threshold.GetApiAggregateType()))
					thresholdConfig.UserAggregateType = types.StringValue(string(threshold.GetUserAggregateType()))
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
			tflog.Trace(ctx, "Shreyansh Gupta powerful", map[string]interface{}{
				"thresholdConfigs": len(thresholdConfigs),
			})

			thresholdConfigsSet, diags := types.SetValueFrom(
				ctx,
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"api_aggregate_type":                       types.StringType,
						"user_aggregate_type":                      types.StringType,
						"rolling_window_count_allowed":             types.Int64Type,
						"rolling_window_duration":                  types.StringType,
						"threshold_config_type":                    types.StringType,
						"dynamic_mean_calculation_duration":        types.StringType,
						"dynamic_duration":                         types.StringType,
						"dynamic_percentage_exceding_mean_allowed": types.Int64Type,
					},
				},
				thresholdConfigs,
			)
			if diags.HasError() {
				return nil, fmt.Errorf("Threshold config conversion failed")
			}
			model.ThresholdConfigs = thresholdConfigsSet
		}

	}
	tflog.Trace(ctx, "Shreyansh Gupta 2")
	return &model, nil
}

func convertToRateLimitingRuleStatus(data *models.RateLimitingRuleModel) (*generated.InputRateLimitingRuleStatus, error) {
	var internal = false
	var status *generated.InputRateLimitingRuleStatus
	status = &generated.InputRateLimitingRuleStatus{
		Internal: &internal,
	}
	return status, nil
}

func convertToRateLimitingRuleThresholdActionConfigType(data *models.RateLimitingRuleModel) ([]*generated.InputRateLimitingRuleThresholdActionConfig, error) {
	configTypes := []*generated.InputRateLimitingRuleThresholdActionConfig{}
	actions := []*generated.InputRateLimitingRuleAction{}
	thresholdConfigs := []*generated.InputRateLimitingRuleThresholdConfig{}
	if HasValue(data.Action) {
		if HasValue(data.Action.ActionType) {
			switch data.Action.ActionType.ValueString() {
			case "ALERT":
				if !HasValue(data.Action.EventSeverity) {
					return nil, utils.NewInvalidError("Action EventSeverity", "EventSeverity must present and must not be empty")
				}
				agentEffect, err := convertToRateLimitingRuleAgentEffect(data)
				if err != nil {
					return nil, err
				}
				eventSeverity, ok := RateLimitingRuleEventSeverityMap[data.Action.EventSeverity.ValueString()]
				if !ok {
					return nil, utils.NewInvalidError("Action EventSeverity", fmt.Sprintf("%s, is not a valid type of Event Severity", data.Action.EventSeverity.ValueString()))

				}
				action := &generated.InputRateLimitingRuleAction{
					ActionType: generated.RateLimitingRuleActionTypeAlert,
					Alert: &generated.InputRateLimitingRuleAlertAction{
						EventSeverity: eventSeverity,
						AgentEffect:   agentEffect,
					},
				}
				actions = append(actions, action)
			case "BLOCK":
				if !HasValue(data.Action.EventSeverity) {
					return nil, utils.NewInvalidError("Action EventSeverity", "EventSeverity must present and must not be empty")
				}
				if !HasValue(data.Action.Duration) {
					return nil, utils.NewInvalidError("Action Duration", "Duration must be present and must not be empty")
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
				if !HasValue(data.Action.Duration) {
					return nil, utils.NewInvalidError("Action Duration", "Duration must be present and must not be empty")
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

				if !HasValue(data.Action.EventSeverity) {
					return configTypes, fmt.Errorf("Event Severity should present")
				}
				eventSeverity, ok := RateLimitingRuleEventSeverityMap[data.Action.EventSeverity.ValueString()]
				if !ok {
					return nil, utils.NewInvalidError("Action EventSeverity", fmt.Sprintf("%s, is not a valid type of Event Severity", data.Action.EventSeverity.ValueString()))

				}
				agentEffect, err := convertToRateLimitingRuleAgentEffect(data)
				if err != nil {
					return configTypes, err
				}
				action := &generated.InputRateLimitingRuleAction{
					ActionType: generated.RateLimitingRuleActionTypeMarkForTesting,
					MarkForTesting: &generated.InputRateLimitingRuleMarkForTestingAction{
						EventSeverity: eventSeverity,
						AgentEffect:   agentEffect,
					},
				}
				actions = append(actions, action)

			default:
				return nil, utils.NewInvalidError("Action ActionType", fmt.Sprintf("%s is not a valid action datatype", data.Action.ActionType.ValueString()))
			}

		} else {
			return nil, utils.NewInvalidError("Action ActionType", "must be present and must not be empty")
		}

	} else {
		return nil, utils.NewInvalidError("Action ", "Action must be present and not be empty")
	}

	if HasValue(data.ThresholdConfigs) {
		var thresholdConfigsModel []models.RateLimitingThresholdConfig
		err := utils.ConvertElementsSet(data.ThresholdConfigs, &thresholdConfigsModel)
		if err != nil {
			return nil, fmt.Errorf("failed to convert threshold configs: %v", err)
		}
		for _, config := range thresholdConfigsModel {

			switch config.ThresholdConfigType.ValueString() {
			case "ROLLING_WINDOW":
				if !HasValue(config.ApiAggregateType) && !(config.ApiAggregateType.ValueString() == "PER_ENDPOINT" || config.ApiAggregateType.ValueString() == "ACROSS_ENDPOINT") {
					return nil, utils.NewInvalidError("threshold_configs api_aggregate_type", "ApiAggregateType must be present,not empty and of valid type")
				}
				if !HasValue(config.UserAggregateType) && !(config.UserAggregateType.ValueString() == "PER_USER" || config.UserAggregateType.ValueString() == "ACROSS_USER") {
					return nil, utils.NewInvalidError("threshold_configs user_aggregate_type", "ApiAggregateType must be present,not empty and of valid type")
				}
				if !HasValue(config.RollingWindowCountAllowed) {
					return nil, utils.NewInvalidError("threshold_configs rolling_window_count_allowed", "RollingWindowCountAllowed must be present")
				}
				if !HasValue(config.RollingWindowDuration) {
					return nil, utils.NewInvalidError("threshold_configs rolling_window_duration", "RollingWindowDuration must be present")
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

			case "DYNAMIC":
				if !HasValue(config.DynamicDuration) {
					return nil, utils.NewInvalidError("threshold_configs dynamic_duration", "DynamicDuration must be present and not empty")
				}
				if !HasValue(config.DynamicMeanCalculationDuration) {
					return nil, utils.NewInvalidError("threshold_configs dynamic_mean_calculation_duration", "DynamicMeanCalculationDuration must be present and not empty")
				}
				if !HasValue(config.DynamicPercentageExcedingMeanAllowed) {
					return nil, utils.NewInvalidError("threshold_configs dynamic_percentage_exceding_mean_allowed", "DynamicPercentageExcedingMeanAllowed must be present and not empty")
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
				return nil, utils.NewInvalidError("threshold_configs threshold_config_type", fmt.Sprintf("%s is not a vaidl thresholdConfigType", config.ThresholdConfigType.ValueString()))

			}

		}

	} else {
		return nil, utils.NewInvalidError("threhold_config", "Must be present can not empty")
	}

	configType := &generated.InputRateLimitingRuleThresholdActionConfig{
		Actions:          actions,
		ThresholdConfigs: thresholdConfigs,
	}

	configTypes = append(configTypes, configType)
	return configTypes, nil

}

func convertToRateLimitingRuleAgentEffect(data *models.RateLimitingRuleModel) (*generated.InputRateLimitingRuleAgentEffect, error) {
	var agentEffect *generated.InputRateLimitingRuleAgentEffect
	agentModifications := []*generated.InputRateLimitingRuleAgentModification{}

	if HasValue(data.Action.HeaderInjections) {
		headerInjectionsList := []models.RateLimitingHeaderInjection{}
		err := utils.ConvertElementsSet(data.Action.HeaderInjections, &headerInjectionsList)
		if err != nil {
			return nil, fmt.Errorf("failed to convert header injections: %v", err)
		}
		for _, injection := range headerInjectionsList {
			if !HasValue(injection.Key) {
				return nil, utils.NewInvalidError("action header_injections key", "key must be present and not empty")
			}
			if !HasValue(injection.Value) {
				return nil, utils.NewInvalidError("action header_injections value", "value must be present and not empty")
			}
			key := injection.Key.ValueString()
			value := injection.Value.ValueString()
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

func convertToRateLimitingRuleCondition(data *models.RateLimitingRuleModel) ([]*generated.InputRateLimitingRuleCondition, error) {

	conditions := []*generated.InputRateLimitingRuleCondition{}

	if HasValue(data.Sources) {

		if HasValue(data.Sources.Scanner) {

			if !HasValue(data.Sources.Scanner.ScannerTypesList) {
				return nil, utils.NewInvalidError("sources scanner scanner_types_list", " Must be present and not empty")
			}
			if !HasValue(data.Sources.Scanner.Exclude) {
				return nil, utils.NewInvalidError("sources scanner exclude", " Must be present and not empty")
			}
			var scannerTypes []*string

			for _, scanner := range data.Sources.Scanner.ScannerTypesList.Elements() {

				if scanner, ok := scanner.(types.String); ok {
					if !RateLimitingRuleScannerMap[scanner.ValueString()] {
						return nil, utils.NewInvalidError("sources scanner scanner_types_list", fmt.Sprintf("Scanner %s is not a valid scanner type", scanner.ValueString()))
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
				return nil, utils.NewInvalidError("sources ip_asn ip_asn_regexes", " Must be present and not empty")
			}
			if !HasValue(data.Sources.IpAsn.Exclude) {
				return nil, utils.NewInvalidError("sources ip_asn exclude", " Must be present and not empty")
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
				return nil, utils.NewInvalidError("sources ip_connection_type ip_connection_type_list", " Must be present and not empty")
			}
			if !HasValue(data.Sources.IpConnectionType.Exclude) {
				return nil, utils.NewInvalidError("sources  ip_connection_type exclude", " Must be present and not empty")
			}
			var ipConnectionTypes []*generated.RateLimitingRuleIpConnectionType

			for _, ipConnectionType := range data.Sources.IpConnectionType.IpConnectionTypeList.Elements() {

				if ipConnectionType, ok := ipConnectionType.(types.String); ok {
					connection := ipConnectionType.ValueString()
					val, exist := RateLimitingRuleIpConnectionTypeMap[connection]
					if !exist {
						return nil, utils.NewInvalidError("sources ip_connection_type ip_connection_type_list", fmt.Sprintf("Ip connection type %s is not a valid ip connection type", connection))
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
				return nil, utils.NewInvalidError("sources user_id user_id_regexes", " Must be present and not empty")
			}
			if !HasValue(data.Sources.UserId.Exclude) {
				return nil, utils.NewInvalidError("sources user_id exclude", " Must be present and not empty")
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

		if HasValue(data.Sources.IpReputation) {
			minIpReputationSeverity, exist := RateLimitingRuleIpReputationSeverityMap[data.Sources.IpReputation.ValueString()]
			if !exist {
				return nil, utils.NewInvalidError("sources ip_reputation", fmt.Sprintf(" %s Invalid Ip Reputation Severity", data.Sources.IpReputation.ValueString()))
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
				return nil, utils.NewInvalidError("sources ip_location_type ip_location_types", " Must be present and not empty")
			}
			if !HasValue(data.Sources.IpLocationType.Exclude) {
				return nil, utils.NewInvalidError("sources ip_location_type exclude", " Must be present and not empty")
			}
			ipLocationTypes := []*generated.RateLimitingRuleIpLocationType{}
			for _, ipLocationType := range data.Sources.IpLocationType.IpLocationTypes.Elements() {
				if locationType, ok := ipLocationType.(types.String); ok {

					ipLocationType, exist := RateLimitingRuleIpLocationTypeMap[locationType.ValueString()]
					if !exist {
						return nil, utils.NewInvalidError("sources ip_location_types", fmt.Sprintf("%s Invalid Ip location Type", locationType.ValueString()))
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
				return nil, utils.NewInvalidError("sources ip_abuse_velocity", fmt.Sprintf(" %s Invalid Ip Abuse Velocity", data.Sources.IpAbuseVelocity.ValueString()))
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
				return nil, utils.NewInvalidError("sources ip_address ip_address_list", " Must be present and not empty")
			}

			if !HasValue(data.Sources.IpAddress.Exclude) {
				return nil, utils.NewInvalidError("sources ip_address exclude", " Must be present and not empty")
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
				return nil, utils.NewInvalidError("sources email_domain email_domain_regexes", " Must be present and not empty")
			}
			if !HasValue(data.Sources.EmailDomain.EmailDomainRegexes) {
				return nil, utils.NewInvalidError("sources email_domain exclude", " Must be present and not empty")
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
				return nil, utils.NewInvalidError("sources user_agents exclude", " Must be present and not empty")
			}
			if !HasValue(data.Sources.UserAgents.UserAgentsList) {
				return nil, utils.NewInvalidError("sources user_agents user_agents_list", " Must be present and not empty")
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
				return nil, utils.NewInvalidError("sources regions exclude", " Must be present and not empty")
			}
			if !HasValue(data.Sources.Regions.RegionsIds) {
				return nil, utils.NewInvalidError("sources regions region_ids", " Must be present and not empty")
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
				return nil, utils.NewInvalidError("sources ip_organisation exclude", " Must be present and not empty")
			}
			if !HasValue(data.Sources.IpOrganisation.IpOrganisationRegexes) {
				return nil, utils.NewInvalidError("sources ip_organisation ip_organisation_regexes", " Must be present and not empty")
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
			requestResponseElement := []models.RateLimitingRequestResponseCondition{}
			err := utils.ConvertElementsSet(data.Sources.RequestResponse, &requestResponseElement)
			if err != nil {
				return nil, fmt.Errorf("converting request response set to slice fails")
			}

			for _, requestResponse := range requestResponseElement {

				keyValueCondition := generated.InputRateLimitingRuleKeyValueCondition{}
				if !HasValue(requestResponse.MetadataType) {
					return nil, utils.NewInvalidError("sources request_response metadata_type", " Must be present and not empty")
				}
				metadataType, exists := RateLimitingRuleKeyValueConditionMetadataTypeMap[requestResponse.MetadataType.ValueString()]
				if !exists {
					return nil, utils.NewInvalidError("sources request_response metadata_type", fmt.Sprintf(" %s Inavlid MetadataType", requestResponse.MetadataType.ValueString()))
				}
				keyValueCondition.MetadataType = &metadataType

				if HasValue(requestResponse.KeyValue) && HasValue(requestResponse.KeyOperator) {
					keyConditionValue := requestResponse.KeyValue.ValueString()
					keyConditionOperator, exist := RateLimitingKeyValueMatchOperatorMap[requestResponse.KeyOperator.ValueString()]

					if !exist {
						return nil, utils.NewInvalidError("sources request_response key_operator", fmt.Sprintf(" %s Inavlid keyOperator", requestResponse.KeyOperator.ValueString()))
					}

					keyValueCondition.KeyCondition = &generated.InputRateLimitingRuleStringCondition{
						Operator: keyConditionOperator,
						Value:    keyConditionValue,
					}

				}

				if HasValue(requestResponse.ValueOperator) && HasValue(requestResponse.Value) {
					valueConditionValue := requestResponse.Value.ValueString()
					valueConditionOperator, exist := RateLimitingKeyValueMatchOperatorMap[requestResponse.ValueOperator.ValueString()]

					if !exist {
						return nil, utils.NewInvalidError("sources request_response value_operator", fmt.Sprintf(" %s Inavlid keyOperator", requestResponse.KeyOperator.ValueString()))
					}

					keyValueCondition.ValueCondition = &generated.InputRateLimitingRuleStringCondition{
						Operator: valueConditionOperator,
						Value:    valueConditionValue,
					}
				}
				var input = generated.InputRateLimitingRuleCondition{
					LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
						ConditionType:     generated.RateLimitingRuleLeafConditionTypeKeyValue,
						KeyValueCondition: &keyValueCondition,
					},
				}

				conditions = append(conditions, &input)

			}

		}

	}
	return conditions, nil

}
