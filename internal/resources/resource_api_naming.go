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

func NewApiNamingResource() resource.Resource {
	return &ApiNamingResource{}
}

type ApiNamingResource struct {
	client *graphql.Client
}

func (r *ApiNamingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_naming"
}

func (r *ApiNamingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ApiNamingResourceSchema()
}

func (r *ApiNamingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApiNamingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "Entering in Create Block")
	var data *models.ApiNamingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	input, err := convertApiNamingModelToCreateInput(ctx, data)
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	//commented this as we can create multiple rules with same name for api naming now
	// id, err := getApiNamingRuleId(input.Name, ctx, *r.client)
	// if err != nil {
	// 	utils.AddError(ctx, &resp.Diagnostics, err)
	// 	return
	// }
	// if id != "" {
	// 	resp.Diagnostics.AddError("Resource already Exist ", fmt.Sprintf("%s api naming rule already please try with different name or import it", input.Name))
	// 	return
	// }
	response, err := generated.CreateApiNamingRule(ctx, *r.client, *input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(response.CreateApiNamingRule.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Exiting in Create Block")
}

func (r *ApiNamingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "Entering in Read Block")
	var data *models.ApiNamingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleID := data.Id.ValueString()
	rule, err := getApiNamingRuleById(ruleID, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	ruleData, err := convertApiNamingFieldsToModel(ctx, rule)
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
func (r *ApiNamingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "Entering in Update Block")
	var dataState *models.ApiNamingModel
	var data *models.ApiNamingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleId := dataState.Id.ValueString()
	rule, err := getApiNamingRuleById(ruleId, ctx, *r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if rule == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	input, err := convertApiNamingModelToUpdateInput(ctx, data, dataState.Id.ValueString())
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	_, err = generated.UpdateApiNamingRule(ctx, *r.client, *input)
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

func (r *ApiNamingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.ApiNamingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := generated.DeleteApiNamingRule(ctx, *r.client, generated.InputApiNamingRuleDelete{Id: data.Id.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting api naming rule", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)

}

// func getApiNamingRuleId(ruleName string, ctx context.Context, r graphql.Client) (string, error) {
// 	response, err := generated.GetApiNamingRule(ctx, r)
// 	if err != nil {
// 		return "", err
// 	}
// 	for _, rule := range response.ApiNamingRules.Results {
// 		if rule.Name == ruleName {
// 			return rule.Id, nil
// 		}
// 	}
// 	return "", nil
// }

func convertApiNamingModelToCreateInput(ctx context.Context, data *models.ApiNamingModel) (*generated.InputApiNamingRuleCreate, error) {
	var input = generated.InputApiNamingRuleCreate{}

	if HasValue(data.Name) {
		input.Name = data.Name.ValueString()
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Disabled) {
		input.Disabled = data.Disabled.ValueBool()
	}

	spanFilter,errInSpanFilter := buildSpanFilter(ctx,data)
	if errInSpanFilter != nil {
		return nil, errInSpanFilter
	}
	input.SpanFilter = *spanFilter

	apiNamingConfig,errInApiNamingConfig := buildApiNamingRuleConfig(ctx,data)
	if errInApiNamingConfig != nil {
		return nil, errInApiNamingConfig
	}
	input.ApiNamingRuleConfig = *apiNamingConfig

	return &input, nil
}


func buildSpanFilter(ctx context.Context, data *models.ApiNamingModel) (*generated.InputTraceableSpanProcessingRuleFilter,error){
	var spanFilter = generated.InputTraceableSpanProcessingRuleFilter{}
	if HasValue(data.EnvironmentNames) && HasValue(data.ServiceNames) {
		environments, _ := utils.ConvertSetToStrPointer(data.EnvironmentNames)
		services, _ := utils.ConvertSetToStrPointer(data.ServiceNames)
		spanFilters:=[]*generated.InputTraceableSpanProcessingRuleFilter{}
		relationalOperator:=generated.TraceableSpanProcessingRelationalOperatorIn
		if (len(environments)) > 0 {
			field:=generated.TraceableSpanProcessingFilterFieldEnvironmentName
			relationalSpanFilter:=generated.InputTraceableSpanProcessingRelationalFilter{
				Field: &field,
				Value: environments,
				RelationalOperator: relationalOperator,
			}
			envNameSpanFilter := generated.InputTraceableSpanProcessingRuleFilter{
				RelationalSpanFilter: &relationalSpanFilter,
			}
			spanFilters = append(spanFilters, &envNameSpanFilter)
		}
		if (len(services)) > 0 {
			field:=generated.TraceableSpanProcessingFilterFieldServiceName
			relationalSpanFilter:=generated.InputTraceableSpanProcessingRelationalFilter{
				Field: &field,
				Value: services,
				RelationalOperator: relationalOperator,
			}
			servNameSpanFilter := generated.InputTraceableSpanProcessingRuleFilter{
				RelationalSpanFilter: &relationalSpanFilter,
			}
			spanFilters = append(spanFilters, &servNameSpanFilter)
		}
		spanFilter.LogicalSpanFilter = &generated.InputTraceableSpanProcessingLogicalFilter{
			LogicalOperator: generated.LogicalOperatorAnd,
			SpanFilters: spanFilters,
		}
		return &spanFilter, nil
	}
	return nil, nil
}

func convertApiNamingModelToUpdateInput(ctx context.Context, data *models.ApiNamingModel, id string) (*generated.InputApiNamingRuleUpdate, error) {
	var input = generated.InputApiNamingRuleUpdate{}
	
	if id != "" {
		input.Id = id
	} else {
		return nil, fmt.Errorf("id can not be empty")
	}
	if HasValue(data.Name) {
		input.Name = data.Name.ValueString()
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	if HasValue(data.Disabled) {
		input.Disabled = data.Disabled.ValueBool()
	}

	spanFilter,errInSpanFilter := buildSpanFilter(ctx,data)
	if errInSpanFilter != nil {
		return nil, errInSpanFilter
	}
	input.SpanFilter = *spanFilter

	apiNamingConfig,errInApiNamingConfig := buildApiNamingRuleConfig(ctx,data)
	if errInApiNamingConfig != nil {
		return nil, errInApiNamingConfig
	}
	input.ApiNamingRuleConfig = *apiNamingConfig
	return &input, nil
}

func buildApiNamingRuleConfig(ctx context.Context, data *models.ApiNamingModel) (*generated.InputApiNamingRuleConfig, error) {
	if HasValue(data.Regexes) && HasValue(data.Values) {
		regexes, _ := utils.ConvertSetToStrPointer(data.Regexes)
		values, _ := utils.ConvertSetToStrPointer(data.Values)
		if len(regexes)==0 || len(values)==0{
			return nil,utils.NewInvalidError("regexes values","both must be non empty")
		}
		if len(regexes) != len(values) {
			return nil,utils.NewInvalidError("regexes values","both must be of equal lengths")
		}
		ruleConfigType := generated.ApiNamingRuleConfigTypeSegmentMatching
		return &generated.InputApiNamingRuleConfig{
			ApiNamingRuleConfigType: ruleConfigType,
			SegmentMatchingBasedRuleConfig : &generated.InputSegmentMatchingBasedRuleConfig{
				Regexes: regexes,
				Values: values,
			},
		}, nil
	}
	return nil, nil
}

func getApiNamingRuleById(ruleId string, ctx context.Context, r graphql.Client) (*generated.ApiNamingRuleFeilds, error) {
	response, err := generated.GetApiNamingRule(ctx, r)
	if err != nil {
		return nil, err
	}
	for _, rule := range response.ApiNamingRules.Results {
		if rule.Id == ruleId {
			return &rule.ApiNamingRuleFeilds, nil
		}
	}
	return nil, nil
}

func convertApiNamingFieldsToModel(ctx context.Context,data *generated.ApiNamingRuleFeilds) (*models.ApiNamingModel, error) {
	finalModel := &models.ApiNamingModel{}
	fmt.Printf("this is data %s",data)
	finalModel.Name=types.StringValue(data.Name)
	finalModel.Id=types.StringValue(data.Id)
	finalModel.Disabled=types.BoolValue(data.Disabled)
	if HasValue(data.SpanFilter){
		serviceNameFeild := generated.TraceableSpanProcessingFilterFieldServiceName
		envNameFeild := generated.TraceableSpanProcessingFilterFieldEnvironmentName
		spanFilters:=data.SpanFilter.LogicalSpanFilter.GetSpanFilters()
		for _, spanFilterObj := range spanFilters{
			field := spanFilterObj.RelationalSpanFilter.GetField()
			if field == &serviceNameFeild {
				interfaceSlice := spanFilterObj.RelationalSpanFilter.Value.([]interface{})
				fmt.Printf("interace %v",interfaceSlice)

				var stringPtrs []*string
				for _, item := range interfaceSlice {
					strPtr, ok := item.(*string)
					if !ok {
						return nil, fmt.Errorf("expected *string in slice but got %T", item)
					}
					stringPtrs = append(stringPtrs, strPtr)
				}
				fmt.Printf("this is svc name %v",stringPtrs)
				svc, err := utils.ConvertStringPtrToTerraformSet(stringPtrs)
				if err != nil {
					return nil, fmt.Errorf("failed to convert service names: %w", err)
				}
				finalModel.ServiceNames=svc
			}else if field == &envNameFeild {
				interfaceSlice := spanFilterObj.RelationalSpanFilter.Value.([]interface{})
				fmt.Printf("interace %v",interfaceSlice)

				var stringPtrs []*string
				for _, item := range interfaceSlice {
					strPtr, ok := item.(*string)
					if !ok {
						return nil, fmt.Errorf("expected *string in slice but got %T", item)
					}
					stringPtrs = append(stringPtrs, strPtr)
				}
				env, err := utils.ConvertStringPtrToTerraformSet(stringPtrs)
				if err != nil {
					return nil, fmt.Errorf("failed to convert service names: %w", err)
				}
				finalModel.EnvironmentNames=env
			}
		}
	}else{
		finalModel.ServiceNames=types.SetNull(types.StringType)
		finalModel.EnvironmentNames=types.SetNull(types.StringType)
	}
	if HasValue(data.ApiNamingRuleConfig){
		regex := data.ApiNamingRuleConfig.SegmentMatchingBasedRuleConfig.GetRegexes()
		values := data.ApiNamingRuleConfig.SegmentMatchingBasedRuleConfig.GetValues()
		finalModel.Regexes, _ = utils.ConvertStringPtrToTerraformSet(regex)
		finalModel.Values, _ = utils.ConvertStringPtrToTerraformSet(values)
	}
	return finalModel, nil
}