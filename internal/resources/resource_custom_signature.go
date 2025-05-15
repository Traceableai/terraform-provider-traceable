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

type CustomSignatureResource struct {
	client *graphql.Client
}

func NewCustomSignatureResource() resource.Resource {
	return &CustomSignatureResource{}
}

func (r *CustomSignatureResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomSignatureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_signature"
}

func (r *CustomSignatureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.CustomSignatureResourceSchema()
}

func (r *CustomSignatureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Entering Create Block")
	var data *models.CustomSignatureModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, err := convertCustomSignatureModelToCreateInput(ctx, data)
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	id, err := getCustomSignatureId(input.Name, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	if id != "" {
		resp.Diagnostics.AddError("Resource already exists",
			fmt.Sprintf("%s custom signature already exists, please try with a different name or import it", input.Name))
		return
	}

	signature, err := generated.CreateCustomSignature(ctx, *r.client, *input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	data.Id = types.StringValue(signature.CreateCustomSignatureRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Exiting Create Block")
}

func (r *CustomSignatureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.CustomSignatureModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getCustomSignatureRule(data.Id.ValueString(), ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	updatedData, err := convertCustomSignatureFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedData)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomSignatureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.CustomSignatureModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var dataState *models.CustomSignatureModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)

	if resp.Diagnostics.HasError() {
		return
	}
	
	
	input, err := convertCustomSignatureModelToUpdateInput(ctx, data, dataState.Id.ValueString())
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp1, err2 := generated.UpdateCustomSignature(ctx, *r.client, *input)
	if err2 != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(resp1.UpdateCustomSignatureRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CustomSignatureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.CustomSignatureModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := generated.DeleteCustomSignature(ctx, *r.client, generated.InputCustomSignatureRuleDelete{Id: data.Id.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting custom signature", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func getCustomSignatureRule(id string, ctx context.Context, r graphql.Client) (*generated.CustomSignatureFields, error) {
	customSignatureFeilds := generated.CustomSignatureFields{}
	filter := &generated.InputCustomSignatureRulesFilter{}
	response, err := generated.GetCustomSignature(ctx, r, *filter)
	if err != nil {
		return nil, err
	}

	for _, rule := range response.CustomSignatureRules.Results {
		if rule.Id == id {
			customSignatureFeilds = rule.CustomSignatureFields
			return &customSignatureFeilds, nil
		}
	}

	return nil, nil
}

func convertCustomSignatureFieldsToModel(ctx context.Context, data *generated.CustomSignatureFields) (*models.CustomSignatureModel, error) {
	customSignatureModel := models.CustomSignatureModel{}
	if data.Id != "" {
		customSignatureModel.Id = types.StringValue(data.Id)
	}
	if data.Name != "" {
		customSignatureModel.Name = types.StringValue(data.Name)
	}
	if data.Description != "" {
		customSignatureModel.Description = types.StringValue(data.Description)
	}
	if data.RuleScope != nil && data.RuleScope.EnvironmentScope != nil {
		environments, err := utils.ConvertStringPtrToTerraformSet(data.RuleScope.EnvironmentScope.EnvironmentIds)
		if err != nil {
			return nil, err
		}
		customSignatureModel.Environments = environments
		tflog.Trace(ctx, "Shreyansh Gupta 123", map[string]interface{}{
			"environments": environments,
		})

	} else {
		customSignatureModel.Environments = types.SetNull(types.StringType)
	}
	customSignatureModel.Disabled = types.BoolValue(*data.Disabled)
	clauseGroup := data.RuleDefinition.ClauseGroup
	clauses := clauseGroup.Clauses
	requestReponseModel := []models.RequestResponseModel{}
	attributeConditionModel := []models.AttributeConditionModel{}
	for _, clause := range clauses {
		clauseType := clause.GetClauseType()
		if clauseType == "MATCH_EXPRESSION" || clauseType == "KEY_VALUE_EXPRESSION" {
			if clauseType == "KEY_VALUE_EXPRESSION" {
				matchCategory := clause.GetKeyValueExpression().GetMatchCategory()
				matchKey := clause.GetKeyValueExpression().GetMatchKey()
				valueMatchOperator := clause.GetKeyValueExpression().GetValueMatchOperator()
				matchValue := clause.GetKeyValueExpression().GetMatchValue()
				keyValueTag := string(clause.GetKeyValueExpression().GetKeyValueTag())
				KeyMatchOperator := string(clause.GetKeyValueExpression().GetKeyMatchOperator())
				kvTag := types.StringValue(keyValueTag)
				keyMatchOp := types.StringValue(KeyMatchOperator)
				requestReponseModel = append(requestReponseModel, models.RequestResponseModel{
					MatchCategory:      types.StringValue(string(*matchCategory)),
					MatchKey:           types.StringValue(string(matchKey)),
					ValueMatchOperator: types.StringValue(string(valueMatchOperator)),
					MatchValue:         types.StringValue(matchValue),
					KeyValueTag:        kvTag,
					KeyMatchOperator:   keyMatchOp,
				})
			} else {
				matchCategory := clause.GetMatchExpression().GetMatchCategory()
				matchKey := clause.GetMatchExpression().GetMatchKey()
				valueMatchOperator := clause.GetMatchExpression().GetMatchOperator()
				matchValue := clause.GetMatchExpression().GetMatchValue()
				requestReponseModel = append(requestReponseModel, models.RequestResponseModel{
					MatchCategory:      types.StringValue(string(*matchCategory)),
					MatchKey:           types.StringValue(string(matchKey)),
					ValueMatchOperator: types.StringValue(string(valueMatchOperator)),
					MatchValue:         types.StringValue(*matchValue),
					KeyValueTag:        types.StringNull(),
					KeyMatchOperator:   types.StringNull(),
				})
			}
		} else if clauseType == "ATTRIBUTE_KEY_VALUE_EXPRESSION" {
			keyConditionOperator := clause.GetAttributeKeyValueExpression().GetKeyCondition().GetOperator()
			keyConditionValue := clause.GetAttributeKeyValueExpression().GetKeyCondition().GetValue()
			valueConditionOperator := clause.GetAttributeKeyValueExpression().GetValueCondition().GetOperator()
			valueConditionValue := clause.GetAttributeKeyValueExpression().GetValueCondition().GetValue()
			valConditionOperator := types.StringValue(string(valueConditionOperator))
			valConditionValue := types.StringValue(string(valueConditionValue))
			attributeConditionModel = append(attributeConditionModel, models.AttributeConditionModel{
				KeyConditionOperator:   types.StringValue(string(keyConditionOperator)),
				KeyConditionValue:      types.StringValue(string(keyConditionValue)),
				ValueConditionOperator: valConditionOperator,
				ValueConditionValue:    valConditionValue,
			})
		} else if clauseType == "CUSTOM_SEC_RULE" {
			customSignatureModel.PayloadCriteria.CustomSecRule = types.StringValue(*clause.GetCustomSecRule().GetInputSecRuleString())
		}
	}
	requestResponseObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key_value_tag":        types.StringType,
			"match_category":       types.StringType,
			"key_match_operator":   types.StringType,
			"match_key":            types.StringType,
			"value_match_operator": types.StringType,
			"match_value":          types.StringType,
		},
	}
	if len(requestReponseModel) > 0 {
		reqresset, diags := types.SetValueFrom(
			ctx,
			requestResponseObjectType,
			requestReponseModel,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("request response conversion failed")
		}
		customSignatureModel.PayloadCriteria.RequestResponse = reqresset
	} else {
		customSignatureModel.PayloadCriteria.RequestResponse = types.SetNull(requestResponseObjectType)
	}
	attributesObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key_condition_operator":   types.StringType,
			"key_condition_value":      types.StringType,
			"value_condition_operator": types.StringType,
			"value_condition_value":    types.StringType,
		},
	}
	if len(attributeConditionModel) > 0 {
		attributesSet, diags := types.SetValueFrom(
			ctx,
			attributesObjectType,
			attributeConditionModel,
		)
		if diags.HasError() {
			return nil, fmt.Errorf("request response conversion failed")
		}
		customSignatureModel.PayloadCriteria.Attributes = attributesSet
	} else {
		customSignatureModel.PayloadCriteria.Attributes = types.SetNull(attributesObjectType)
	}
	ruleEffect := data.RuleEffect
	if ruleEffect.GetEventSeverity() != nil {
		eventSev := types.StringValue(string(*ruleEffect.GetEventSeverity()))
		customSignatureModel.Action.EventSeverity = eventSev
	} else {
		customSignatureModel.Action.EventSeverity = types.StringNull()
	}
	eventType := types.StringValue(string(ruleEffect.GetEventType()))
	customSignatureModel.Action.ActionType = eventType
	if data.GetBlockingExpirationDuration() != nil {
		duration := types.StringValue(string(*data.GetBlockingExpirationDuration()))
		customSignatureModel.Action.Duration = duration
	} else {
		customSignatureModel.Action.Duration = types.StringNull()
	}
	return &customSignatureModel, nil
}

func convertCustomSignatureModelToUpdateInput(ctx context.Context, data *models.CustomSignatureModel, id string) (*generated.InputCustomSignatureRuleUpdate, error) {
	input := generated.InputCustomSignatureRuleUpdate{}
	err := checkInputCondition(data)
	if err != nil {
		return nil, err
	}
	if id != "" {
		input.Id = id
	} else {
		return nil, fmt.Errorf("id can not be empty")
	}
	if HasValue(data.Name) {
		name := data.Name.ValueString()
		input.Name = name
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		input.Description = description
	}
	if HasValue(data.Disabled) {
		disabled := data.Disabled.ValueBool()
		input.Disabled = &disabled
	}
	internal := false
	input.Internal = &internal
	scope, err := convertToCustomSignatureRuleConfigScope(data.Environments)
	if err != nil {
		return nil, err
	} else {
		input.RuleScope = scope
	}
	RuleDefinitionScope, err := convertToCustomSignatureRuleDefination(data)
	if err != nil {
		return nil, err
	} else {
		input.RuleDefinition = *RuleDefinitionScope
	}
	ruleEffect, err := convertToCustomSignatureRuleEffect(data)
	if err != nil {
		return nil, err
	} else {
		input.RuleEffect = *ruleEffect
	}
	if HasValue(data.Action.Duration) {
		duration := data.Action.Duration.ValueString()
		input.BlockingExpirationDuration = &duration
	}
	if HasValue(data.PayloadCriteria.Attributes) && data.Action.ActionType == types.StringValue("ALLOW") {
		return nil, utils.NewInvalidError("action_type", "action_type ALLOW is not valid with payload_criteria attributes")
	}
	return &input, nil
}

func getCustomSignatureId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
	filter := &generated.InputCustomSignatureRulesFilter{}

	response, err := generated.GetCustomSignatureId(ctx, r, *filter)
	if err != nil {
		return "", err
	}
	for _, rule := range response.CustomSignatureRules.Results {
		if rule.Name == ruleName {
			return rule.GetId(), nil
		}
	}
	return "", nil
}

func convertCustomSignatureModelToCreateInput(ctx context.Context, data *models.CustomSignatureModel) (*generated.InputCustomSignatureRuleDescriptor, error) {
	errInInput := checkInputCondition(data)
	if errInInput != nil {
		return nil, errInInput
	}
	if data == nil {
		return nil, fmt.Errorf("data model is nil")
	}
	var input = generated.InputCustomSignatureRuleDescriptor{}
	if HasValue(data.Name) {
		name := data.Name.ValueString()
		input.Name = name
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Description) {
		description := data.Description.ValueString()
		input.Description = description
	}
	if HasValue(data.Disabled) {
		disabled := data.Disabled.ValueBool()
		input.Disabled = &disabled
	}
	internal := false
	input.Internal = &internal
	scope, err := convertToCustomSignatureRuleConfigScope(data.Environments)
	if err != nil {
		return nil, err
	} else {
		input.RuleScope = scope
	}
	RuleDefinitionScope, err := convertToCustomSignatureRuleDefination(data)
	if err != nil {
		return nil, err
	} else {
		input.RuleDefinition = *RuleDefinitionScope
	}
	ruleEffect, err := convertToCustomSignatureRuleEffect(data)
	if err != nil {
		return nil, err
	} else {
		input.RuleEffect = *ruleEffect
	}
	if HasValue(data.Action.Duration) {
		duration := data.Action.Duration.ValueString()
		input.BlockingExpirationDuration = &duration
	}
	if HasValue(data.PayloadCriteria.Attributes) && data.Action.ActionType == types.StringValue("ALLOW") {
		return nil, utils.NewInvalidError("action_type", "action_type ALLOW is not valid with payload_criteria attributes")
	}
	return &input, nil
}

func convertToCustomSignatureRuleEffect(data *models.CustomSignatureModel) (*generated.InputCustomSignatureRuleEffect, error) {
	if !HasValue(data) {
		return nil, nil
	}
	var ruleEffect *generated.InputCustomSignatureRuleEffect
	eventType := data.Action.ActionType.ValueString()
	_, isValidEventType := CustomSignatureRuleEventTypeMap[eventType]
	if !isValidEventType {
		return nil, utils.NewInvalidError("sources action_type", fmt.Sprintf(" %s Invalid action_type", eventType))
	}
	ruleEventType := generated.CustomSignatureRuleEventType(eventType)
	if HasValue(data.Action.EventSeverity) {
		eventSeverity := data.Action.EventSeverity.ValueString()
		_, isValidEventSev := RateLimitingRuleEventSeverityMap[eventSeverity]
		if !isValidEventSev {
			return nil, utils.NewInvalidError("sources action_type event_severity", fmt.Sprintf(" %s Invalid event_severity", eventSeverity))
		}
		ruleEventSeverity := generated.CustomSignatureRuleEventSeverity(eventSeverity)
		ruleEffect = &generated.InputCustomSignatureRuleEffect{
			Effects:       []*generated.InputCustomSignatureRuleEffectWithModification{},
			EventSeverity: &ruleEventSeverity,
			EventType:     ruleEventType,
		}
	}else{
		ruleEffect = &generated.InputCustomSignatureRuleEffect{
			Effects:       []*generated.InputCustomSignatureRuleEffectWithModification{},
			EventType:     ruleEventType,
		}
	}
	return ruleEffect, nil
}

func convertToCustomSignatureRuleDefination(data *models.CustomSignatureModel) (*generated.InputCustomSignatureRuleDefinition, error) {
	if !HasValue(data) {
		return nil, nil
	}
	var ruleDefinition *generated.InputCustomSignatureRuleDefinition
	var clauseGroup *generated.InputCustomSignatureRuleClauseGroup
	var customSignatureRuleClauseRequest []*generated.InputCustomSignatureRuleClauseRequest
	if HasValue(data.PayloadCriteria) {
		if HasValue(data.PayloadCriteria.RequestResponse) {
			var customSigReqRes []*models.RequestResponseModel
			err := utils.ConvertElementsSet(data.PayloadCriteria.RequestResponse, &customSigReqRes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert elements: %v", err)
			}
			for _, reqResp := range customSigReqRes {
				_, isValidMatchCat := CustomSignatureRuleMatchCategoryMap[reqResp.MatchCategory.ValueString()]
				if !isValidMatchCat {
					return nil, utils.NewInvalidError("sources request_response match_category", fmt.Sprintf(" %s Invalid match_category", reqResp.MatchCategory.ValueString()))
				}

				_, isValidMatchKey := CustomSignatureRuleMatchKeyMap[reqResp.MatchKey.ValueString()]
				if !isValidMatchKey {
					return nil, utils.NewInvalidError("sources request_response match_key", fmt.Sprintf(" %s Invalid match_key", reqResp.MatchKey.ValueString()))
				}

				_, isValidValueMatchOp := RateLimitingKeyValueMatchOperatorMap[reqResp.ValueMatchOperator.ValueString()]
				if !isValidValueMatchOp {
					return nil, utils.NewInvalidError("sources request_response value_match_operator", fmt.Sprintf(" %s Invalid value_match_operator", reqResp.ValueMatchOperator.ValueString()))
				}

				var clause *generated.InputCustomSignatureRuleClauseRequest
				if HasValue(reqResp.KeyValueTag) && HasValue(reqResp.KeyMatchOperator) {
					_, isValidKeyValTag := CustomSignatureKeyValuesExpressionMap[reqResp.KeyValueTag.ValueString()]
					if isValidKeyValTag {
						_, KeyMatchOperator := RateLimitingKeyValueMatchOperatorMap[reqResp.KeyMatchOperator.ValueString()]
						if !KeyMatchOperator {
							return nil, utils.NewInvalidError("sources request_response key_match_operator", fmt.Sprintf(" %s Invalid key_match_operator", reqResp.KeyMatchOperator.ValueString()))
						}

						matchKey := reqResp.MatchKey.ValueString()
						matchCategory := generated.CustomSignatureRuleMatchCategory(reqResp.MatchCategory.ValueString())
						valueMatchOperator := generated.CustomSignatureRuleMatchOperator(reqResp.ValueMatchOperator.ValueString())
						keyValueTag := generated.CustomSignatureRuleKeyValueTag(reqResp.KeyValueTag.ValueString())
						keyMatchOperator := generated.CustomSignatureRuleMatchOperator(reqResp.KeyMatchOperator.ValueString())
						matchValue := reqResp.MatchValue.ValueString()
						clause = &generated.InputCustomSignatureRuleClauseRequest{
							ClauseType: generated.CustomSignatureRuleClauseTypeKeyValueExpression,
							KeyValueExpression: &generated.InputCustomSignatureRuleKeyValueExpression{
								KeyMatchOperator:   keyMatchOperator,
								KeyValueTag:        keyValueTag,
								MatchCategory:      &matchCategory,
								MatchKey:           matchKey,
								MatchValue:         matchValue,
								ValueMatchOperator: valueMatchOperator,
							},
						}
					} else {
						return nil, utils.NewInvalidError("key_value_tag", fmt.Sprintf("%s not a valid key_value_tag", reqResp.KeyValueTag.ValueString()))
					}
				} else if !HasValue(reqResp.KeyValueTag) && !HasValue(reqResp.KeyMatchOperator) {
					matchKey := generated.CustomSignatureRuleMatchKey(reqResp.MatchKey.ValueString())
					matchCategory := generated.CustomSignatureRuleMatchCategory(reqResp.MatchCategory.ValueString())
					valueMatchOperator := generated.CustomSignatureRuleMatchOperator(reqResp.ValueMatchOperator.ValueString())
					matchValue := reqResp.MatchValue.ValueString()
					clause = &generated.InputCustomSignatureRuleClauseRequest{
						ClauseType: generated.CustomSignatureRuleClauseTypeMatchExpression,
						MatchExpression: &generated.InputCustomSignatureRuleMatchExpression{
							MatchCategory: &matchCategory,
							MatchValue:    &matchValue,
							MatchKey:      matchKey,
							MatchOperator: valueMatchOperator,
						},
					}
				} else {
					return nil, utils.NewInvalidError("key_value_tag or key_match_operator", "either both or none of them is required")
				}
				customSignatureRuleClauseRequest = append(customSignatureRuleClauseRequest, clause)
			}
		}
		if HasValue(data.PayloadCriteria.Attributes) {
			var customSigAttributes []*models.AttributeConditionModel
			err := utils.ConvertElementsSet(data.PayloadCriteria.Attributes, &customSigAttributes)
			if err != nil {
				return nil, fmt.Errorf("failed to convert elements: %v", err)
			}
			for _, attr := range customSigAttributes {
				_, keyCondtionOperatorExist := RateLimitingKeyValueMatchOperatorMap[attr.KeyConditionOperator.ValueString()]

				if !keyCondtionOperatorExist {
					return nil, utils.NewInvalidError("sources attributes key_condition_operator", fmt.Sprintf(" %s Invalid keyOperator", attr.KeyConditionOperator.ValueString()))
				}
				_, valCondtionOperatorExist := RateLimitingKeyValueMatchOperatorMap[attr.ValueConditionOperator.ValueString()]

				if !valCondtionOperatorExist {
					return nil, utils.NewInvalidError("sources attributes value_condition_operator", fmt.Sprintf(" %s Invalid valueOperator", attr.ValueConditionOperator.ValueString()))
				}

				keyCondition := &generated.InputCustomSignatureStringMatchCondition{
					Operator: generated.CustomSignatureRuleMatchOperator(attr.KeyConditionOperator.ValueString()),
					Value:    attr.KeyConditionValue.ValueString(),
				}
				valueCondition := &generated.InputCustomSignatureStringMatchCondition{
					Operator: generated.CustomSignatureRuleMatchOperator(attr.ValueConditionOperator.ValueString()),
					Value:    attr.ValueConditionValue.ValueString(),
				}

				clause := &generated.InputCustomSignatureRuleClauseRequest{
					ClauseType: generated.CustomSignatureRuleClauseTypeAttributeKeyValueExpression,
					AttributeKeyValueExpression: &generated.InputCustomSignatureRuleAttributeKeyValueExpression{
						KeyCondition:   keyCondition,
						ValueCondition: valueCondition,
					},
				}
				customSignatureRuleClauseRequest = append(customSignatureRuleClauseRequest, clause)
			}
		}
		if HasValue(data.PayloadCriteria.CustomSecRule) {
			customSecRuleString := data.PayloadCriteria.CustomSecRule.ValueString()
			clause := &generated.InputCustomSignatureRuleClauseRequest{
				ClauseType: generated.CustomSignatureRuleClauseTypeCustomSecRule,
				CustomSecRule: &generated.InputCustomSignatureSecRule{
					InputSecRuleString: &customSecRuleString,
				},
			}
			customSignatureRuleClauseRequest = append(customSignatureRuleClauseRequest, clause)
		}
	}
	clauseGroup = &generated.InputCustomSignatureRuleClauseGroup{
		Clauses:        customSignatureRuleClauseRequest,
		ClauseOperator: generated.CustomSignatureRuleClauseOperatorAnd,
	}
	ruleDefinition = &generated.InputCustomSignatureRuleDefinition{
		ClauseGroup: *clauseGroup,
		Labels:      []*generated.InputCustomSignatureRuleLabel{},
	}
	return ruleDefinition, nil
}

func convertToCustomSignatureRuleConfigScope(environments types.Set) (*generated.InputCustomSignatureRuleScope, error) {
	if !HasValue(environments) {
		return nil, nil
	}

	var scope *generated.InputCustomSignatureRuleScope
	var envIds []*string
	for _, env := range environments.Elements() {
		if env, ok := env.(types.String); ok {
			env1 := env.ValueString()
			envIds = append(envIds, &env1)
		}
	}
	scope = &generated.InputCustomSignatureRuleScope{
		EnvironmentScope: &generated.InputCustomSignatureEnvironmentScope{
			EnvironmentIds: envIds,
		},
	}
	return scope, nil
}

func (r *CustomSignatureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ruleName := req.ID
	id, err := getCustomSignatureId(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Resource Not Found", fmt.Sprintf("%s rule of this name not found", ruleName))
		return
	}
	response, err := getCustomSignatureRule(id, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data, err := convertCustomSignatureFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func checkInputCondition(data *models.CustomSignatureModel) (error) {
	if data.Action.ActionType.ValueString() == "NORMAL_DETECTION" && HasValue(data.Action.Duration) {
		return utils.NewInvalidError("duration not required with action_type", fmt.Sprintf("duration not required with action_type %s", data.Action.ActionType.ValueString()))
	}
	if data.Action.ActionType.ValueString() == "TESTING_DETECTION" && HasValue(data.Action.Duration) {
		return utils.NewInvalidError("duration not required with action_type", fmt.Sprintf("duration not required with action_type %s", data.Action.ActionType.ValueString()))
	}
	if data.Action.ActionType.ValueString() == "ALLOW" && HasValue(data.Action.EventSeverity) {
		return utils.NewInvalidError("event_severity not required with action_type", fmt.Sprintf("event_severity not required with action_type %s", data.Action.ActionType.ValueString()))
	}
	if data.Action.ActionType.ValueString() == "TESTING_DETECTION" && HasValue(data.Action.EventSeverity) {
		return utils.NewInvalidError("event_severity not required with action_type", fmt.Sprintf("event_severity not required with action_type %s", data.Action.ActionType.ValueString()))
	}
	if data.Action.ActionType.ValueString() == "BLOCK" && !HasValue(data.Action.EventSeverity) {
		return utils.NewInvalidError("event_severity required with action_type", fmt.Sprintf("event_severity required with action_type %s", data.Action.ActionType.ValueString()))
	}
	if data.Action.ActionType.ValueString() == "NORMAL_DETECTION" && !HasValue(data.Action.EventSeverity) {
		return utils.NewInvalidError("event_severity required with action_type", fmt.Sprintf("event_severity required with action_type %s", data.Action.ActionType.ValueString()))
	}
	return nil
}
