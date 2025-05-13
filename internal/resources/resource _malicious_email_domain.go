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

// Value needs to be set to MIN_SEVERITY for user to set MinEmailFraudScoreLevel, since value for EmailFraudScoreType will be constant, not exposing to the end user
const DefaultEmailFraudScoreType = "MIN_SEVERITY"

func NewMaliciousEmailDomainResource() resource.Resource {
	return &MaliciousEmailDomainResource{}
}

type MaliciousEmailDomainResource struct {
	client *graphql.Client
}

func (r *MaliciousEmailDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_malicious_email_domain"
}

func (r *MaliciousEmailDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.MaliciousEmailDomainResourceSchema()
}

func (r *MaliciousEmailDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MaliciousEmailDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "Entering in Create Block")
	var data *models.MaliciousEmailDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	input, err := convertMaliciousEmailDomainModelToCreateInput(ctx, data)
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	id, err := getMaliciousEmailDomainRuleId(input.RuleInfo.Name, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id != "" {
		resp.Diagnostics.AddError("Resource already Exist ", fmt.Sprintf("%s malicious email domain rule already exists please try with different name or import it", input.RuleInfo.Name))
		return
	}
	response, err := generated.CreateMaliciousEmailDomainRule(ctx, *r.client, *input)
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

func (r *MaliciousEmailDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "Entering in Read Block")
	var data *models.MaliciousEmailDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := data.Name.ValueString()
	rule, err := getMaliciousEmailDomainRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	ruleData, err := convertMaliciousEmailDomainFieldsToModel(ctx, rule)
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

func (r *MaliciousEmailDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Entering in Update Block")
	var dataState *models.MaliciousEmailDomainModel
	var data *models.MaliciousEmailDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := dataState.Name.ValueString()
	rule, err := getMaliciousEmailDomainRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	input, err := convertMaliciousEmailDomainModelToUpdateInput(ctx, data, dataState.Id.ValueString())
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	_, err = generated.UpdateMaliciousEmailDomainRule(ctx, *r.client, *input)
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
func (r *MaliciousEmailDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Trace(ctx, "Entering in ImportState Block")
	ruleName := req.ID
	id, err := getMaliciousEmailDomainRuleId(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Resource Not Found", fmt.Sprintf("%s rule of this name not found", ruleName))
		return
	}
	response, err := getMaliciousEmailDomainRule(ruleName, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data, err := convertMaliciousEmailDomainFieldsToModel(ctx, response)

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

func (r *MaliciousEmailDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "Entering in Delete Block")
	var data *models.MaliciousEmailDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleName := data.Name.ValueString()
	rule, err := getMaliciousEmailDomainRule(ruleName, ctx, *r.client)
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
	_, err = generated.DeleteMaliciousEmailDomainRule(ctx, *r.client, input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "Exiting in Delete Block")

}

func getMaliciousEmailDomainRuleId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
	input := &generated.InputMaliciousSourcesRulesFilter{}

	response, err := generated.GetMaliciousEmailDomainRulesName(ctx, r, input)
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

func getMaliciousEmailDomainRule(ruleName string, ctx context.Context, r graphql.Client) (*generated.MaliciousEmailDomainRuleFields, error) {
	input := &generated.InputMaliciousSourcesRulesFilter{}

	response, err := generated.GetMaliciousEmailDomainRuleDetails(ctx, r, input)
	if err != nil {
		return nil, err
	}
	for _, rule := range response.MaliciousSourcesRules.Results {
		if rule.Info.Name == ruleName {
			return &rule.MaliciousEmailDomainRuleFields, nil
		}
	}
	return nil, nil
}

func convertMaliciousEmailDomainModelToCreateInput(ctx context.Context, data *models.MaliciousEmailDomainModel) (*generated.InputMaliciousSourcesRuleCreate, error) {
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

	// Atleast any one of the following fields needs to be present
	if HasValue(data.EmailDomainsList) ||
		HasValue(data.EmailRegexesList) ||
		HasValue(data.MinEmailFraudScoreLevel) ||
		(HasValue(data.ApplyRuleToDataLeakedEmail) && data.ApplyRuleToDataLeakedEmail.ValueBool()) ||
		(HasValue(data.ApplyRuleToDisposableEmailDomains) && data.ApplyRuleToDisposableEmailDomains.ValueBool()) {

		ruleCondition := generated.InputMaliciousSourcesRuleCondition{}
		emailDomainCondition := generated.InputMaliciousSourcesRuleEmailDomainCondition{}

		if HasValue(data.EmailDomainsList) {
			emailDomainsList, err := utils.ConvertSetToStrPointer(data.EmailDomainsList)
			if err != nil {
				return nil, fmt.Errorf("converting email domains to string pointer fails")
			}
			emailDomainCondition.EmailDomains = emailDomainsList
		}
		if HasValue(data.EmailRegexesList) {
			emailRegexesList, err := utils.ConvertSetToStrPointer(data.EmailRegexesList)
			if err != nil {
				return nil, fmt.Errorf("converting email regexes to string pointer fails")
			}
			emailDomainCondition.EmailRegexes = emailRegexesList
		}
		if HasValue(data.MinEmailFraudScoreLevel) {
			minEmailFraudScoreLevel, exist := MaliciousEmailDomainMinEmailFraudScoreLevel[data.MinEmailFraudScoreLevel.ValueString()]
			if !exist {
				return nil, utils.NewInvalidError("min_email_fraud_score_level", "Invalid min email fraud score level")
			}
			if emailDomainCondition.EmailFraudScore == nil {
				emailDomainCondition.EmailFraudScore = &generated.InputMaliciousSourcesRuleEmailFraudScore{}
			}
			emailDomainCondition.EmailFraudScore.MinEmailFraudScoreLevel = &minEmailFraudScoreLevel
			emailDomainCondition.EmailFraudScore.EmailFraudScoreType = DefaultEmailFraudScoreType
		}
		if HasValue(data.ApplyRuleToDataLeakedEmail) {
			val := data.ApplyRuleToDataLeakedEmail.ValueBool()
			emailDomainCondition.DataLeakedEmail = &val
		}
		if HasValue(data.ApplyRuleToDisposableEmailDomains) {
			val := data.ApplyRuleToDisposableEmailDomains.ValueBool()
			emailDomainCondition.DisposableEmailDomain = &val
		}
		ruleCondition.ConditionType = generated.MaliciousSourcesRuleConditionTypeEmailDomain
		ruleCondition.EmailDomainCondition = &emailDomainCondition
		ruleInfo.Conditions = append(ruleInfo.Conditions, &ruleCondition)
	} else {
		return nil, utils.NewInvalidError("email_domains_list, email_regexes_list, min_email_fraud_score_level, apply_rule_to_data_leaked_email, apply_rule_to_disposable_email_domains", "All fields cannot be empty, at least one of them must be present and cannot contain only default value. For eg: ONLY having apply_rule_to_data_leaked_email = false and apply_rule_to_disposable_email_domains = false is not allowed")
	}

	ruleInfo.Action = ruleAction
	input.RuleInfo = ruleInfo

	return &input, nil

}

func convertMaliciousEmailDomainModelToUpdateInput(ctx context.Context, data *models.MaliciousEmailDomainModel, id string) (*generated.InputMaliciousSourcesRuleUpdate, error) {
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

	// Atleast any one of the following fields needs to be present
	if HasValue(data.EmailDomainsList) ||
		HasValue(data.EmailRegexesList) ||
		HasValue(data.MinEmailFraudScoreLevel) ||
		(HasValue(data.ApplyRuleToDataLeakedEmail) && data.ApplyRuleToDataLeakedEmail.ValueBool()) ||
		(HasValue(data.ApplyRuleToDisposableEmailDomains) && data.ApplyRuleToDisposableEmailDomains.ValueBool()) {

		ruleCondition := generated.InputMaliciousSourcesRuleCondition{}
		emailDomainCondition := generated.InputMaliciousSourcesRuleEmailDomainCondition{}

		if HasValue(data.EmailDomainsList) {
			emailDomainsList, err := utils.ConvertSetToStrPointer(data.EmailDomainsList)
			if err != nil {
				return nil, fmt.Errorf("converting email domains to string pointer fails")
			}
			emailDomainCondition.EmailDomains = emailDomainsList
		}
		if HasValue(data.EmailRegexesList) {
			emailRegexesList, err := utils.ConvertSetToStrPointer(data.EmailRegexesList)
			if err != nil {
				return nil, fmt.Errorf("converting email regexes to string pointer fails")
			}
			emailDomainCondition.EmailRegexes = emailRegexesList
		}
		if HasValue(data.MinEmailFraudScoreLevel) {
			minEmailFraudScoreLevel, exist := MaliciousEmailDomainMinEmailFraudScoreLevel[data.MinEmailFraudScoreLevel.ValueString()]
			if !exist {
				return nil, utils.NewInvalidError("min_email_fraud_score_level", "Invalid min email fraud score level")
			}
			if emailDomainCondition.EmailFraudScore == nil {
				emailDomainCondition.EmailFraudScore = &generated.InputMaliciousSourcesRuleEmailFraudScore{}
			}
			emailDomainCondition.EmailFraudScore.MinEmailFraudScoreLevel = &minEmailFraudScoreLevel
			emailDomainCondition.EmailFraudScore.EmailFraudScoreType = DefaultEmailFraudScoreType
		}
		if HasValue(data.ApplyRuleToDataLeakedEmail) {
			val := data.ApplyRuleToDataLeakedEmail.ValueBool()
			emailDomainCondition.DataLeakedEmail = &val
		}
		if HasValue(data.ApplyRuleToDisposableEmailDomains) {
			val := data.ApplyRuleToDisposableEmailDomains.ValueBool()
			emailDomainCondition.DisposableEmailDomain = &val
		}
		ruleCondition.ConditionType = generated.MaliciousSourcesRuleConditionTypeEmailDomain
		ruleCondition.EmailDomainCondition = &emailDomainCondition
		ruleInfo.Conditions = append(ruleInfo.Conditions, &ruleCondition)
	} else {
		return nil, utils.NewInvalidError("email_domains_list, email_regexes_list, min_email_fraud_score_level, apply_rule_to_data_leaked_email, apply_rule_to_disposable_email_domains", "All fields cannot be empty, at least one of them must be present and cannot contain only default value. For eg: ONLY having apply_rule_to_data_leaked_email = false and apply_rule_to_disposable_email_domains = false is not allowed")
	}

	ruleInfo.Action = ruleAction
	rule.Info = ruleInfo
	input.Rule = rule

	return &input, nil

}

func convertMaliciousEmailDomainFieldsToModel(ctx context.Context, rule *generated.MaliciousEmailDomainRuleFields) (*models.MaliciousEmailDomainModel, error) {
	var data = models.MaliciousEmailDomainModel{}

	var emailDomains []*string
	var emailRegexes []*string

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
		for _, condition := range rule.Info.Conditions {
			if condition.ConditionType == generated.MaliciousSourcesRuleConditionTypeEmailDomain {
				if condition.EmailDomainCondition != nil {
					if condition.EmailDomainCondition.EmailDomains != nil {
						for _, emailDomain := range condition.EmailDomainCondition.EmailDomains {
							emailDomain := string(*emailDomain)
							emailDomains = append(emailDomains, &emailDomain)
						}
						emailDomainsSet, err := utils.ConvertCustomStringPtrsToTerraformSet(emailDomains)
						if err != nil {
							return nil, fmt.Errorf("Error converting email domains to terraform set")
						}
						data.EmailDomainsList = emailDomainsSet
					}
					if condition.EmailDomainCondition.EmailRegexes != nil {
						for _, emailRegex := range condition.EmailDomainCondition.EmailRegexes {
							emailRegex := string(*emailRegex)
							emailRegexes = append(emailRegexes, &emailRegex)
						}
						emailRegexesSet, err := utils.ConvertCustomStringPtrsToTerraformSet(emailRegexes)
						if err != nil {
							return nil, fmt.Errorf("Error converting email regexes to terraform set")
						}
						data.EmailRegexesList = emailRegexesSet
					}
					if condition.EmailDomainCondition.DataLeakedEmail != nil {
						data.ApplyRuleToDataLeakedEmail = types.BoolValue(*condition.EmailDomainCondition.DataLeakedEmail)
					}
					if condition.EmailDomainCondition.DisposableEmailDomain != nil {
						data.ApplyRuleToDisposableEmailDomains = types.BoolValue(*condition.EmailDomainCondition.DisposableEmailDomain)
					}
					if condition.EmailDomainCondition.EmailFraudScore != nil {
						if condition.EmailDomainCondition.EmailFraudScore.MinEmailFraudScoreLevel != nil {
							data.MinEmailFraudScoreLevel = types.StringValue(string(*condition.EmailDomainCondition.EmailFraudScore.MinEmailFraudScoreLevel))
						}
					}
				}
			}
		}
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
