package resources

import (
	"context"
	"fmt"
	"log"
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

type DataLossPreventionRequestBasedResource struct {
	client *graphql.Client
}

func NewDataLossPreventionRequestBasedResource() resource.Resource {
	return &DataLossPreventionRequestBasedResource{}
}

func (r *DataLossPreventionRequestBasedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DataLossPreventionRequestBasedResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_loss_prevention_request_based"
}

func (r *DataLossPreventionRequestBasedResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.DataLossPreventionRequestBasedResourceSchema()
}

func (r *DataLossPreventionRequestBasedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Entering in Create Block")
	var data *models.DataLossPreventionRequestBasedRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ruleInput, err := convertDLPRequestBasedModelToCreateInput(ctx, data, r.client)
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
		resp.Diagnostics.AddError("Resource already Exist", fmt.Sprintf("%s Dlp request based rule already please try with different name or import it", ruleInput.Name))
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

func (r *DataLossPreventionRequestBasedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.DataLossPreventionRequestBasedRuleModel

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

	updatedData, err := convertDLPRequestBasedRuleFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedData)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *DataLossPreventionRequestBasedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.DataLossPreventionRequestBasedRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dataState *models.DataLossPreventionRequestBasedRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)

	if resp.Diagnostics.HasError() {
		return
	}
	input, err := convertDLPRequestBasedModelToUpdateInput(ctx, data, dataState.Id.ValueString(), r.client)
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

func (r *DataLossPreventionRequestBasedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.DataLossPreventionRequestBasedRuleModel

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

func (r *DataLossPreventionRequestBasedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
	data, err := convertDLPRequestBasedRuleFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func convertDLPRequestBasedRuleFieldsToModel(ctx context.Context, data *generated.RateLimitingRuleFields) (*models.DataLossPreventionRequestBasedRuleModel, error){
	model := models.DataLossPreventionRequestBasedRuleModel{}

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

	} else {
		model.Environments = types.SetNull(types.StringType)
	}
	model.Enabled = types.BoolValue(data.Enabled)

	if data.Conditions != nil {
		sources := models.DlpRequestBasedSources{}
		requestPayloadArr := []models.DlpRequestBasedRequestPayloadScope{}
		for _, condition := range data.GetConditions() {
			leafCondition := condition.LeafCondition.LeafConditionFields
			switch string(leafCondition.ConditionType) {
			case "KEY_VALUE":
				reqres := models.DlpRequestBasedRequestPayloadScope{}
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
				requestPayloadArr = append(requestPayloadArr, reqres)

			case "SCOPE":
				if leafCondition.ScopeCondition.GetEntityScope() != nil {
					serviceIds, err := utils.ConvertStringPtrToTerraformSet(leafCondition.ScopeCondition.EntityScope.GetEntityIds())
					if err != nil {
						return nil, err
					}

					sources.ServiceScope = models.DlpRequestBasedServiceScope{
						ServiceIds: serviceIds,
					}
				}
				if leafCondition.ScopeCondition.GetUrlScope() != nil{
					urlRegexes, err := utils.ConvertStringPtrToTerraformSet(leafCondition.ScopeCondition.UrlScope.GetUrlRegexes())
					if err != nil {
						return nil, err
					}
					sources.UrlScope = models.DlpRequestBasedUrlScope{
						UrlRegexes: urlRegexes,
					}
				}

			case "DATATYPE":
				dataSetDataTypeIds := models.DataSetDataTypeIds{}
				dataTypeMatching := models.DlpReqBasedDataTypeMatching{}
				if len(leafCondition.DatatypeCondition.GetDatasetIds()) > 0 {
					dataSetsIds, err := utils.ConvertStringPtrToTerraformSet(leafCondition.DatatypeCondition.GetDatasetIds())
					if err != nil {
						return nil, err
					}
					dataSetDataTypeIds.DataSetsIds = dataSetsIds
				}else{
					dataSetDataTypeIds.DataSetsIds = types.SetNull(types.StringType)
				}
				if len(leafCondition.DatatypeCondition.GetDatatypeIds()) > 0 {
					dataTypeIds, err := utils.ConvertStringPtrToTerraformSet(leafCondition.DatatypeCondition.GetDatatypeIds())
					if err != nil {
						return nil, err
					}
					dataSetDataTypeIds.DataTypesIds = dataTypeIds
				}else{
					dataSetDataTypeIds.DataTypesIds = types.SetNull(types.StringType)
				}

				if leafCondition.DatatypeCondition.GetDatatypeMatching() != nil {
					metaDataType := leafCondition.DatatypeCondition.DatatypeMatching.RegexBasedMatching.CustomMatchingLocation.GetMetadataType()
					dataTypeMatching.MetadataType=types.StringValue(string(*metaDataType))
					if HasValue(leafCondition.DatatypeCondition.DatatypeMatching.RegexBasedMatching.CustomMatchingLocation.GetKeyCondition()) {
						operator := leafCondition.DatatypeCondition.DatatypeMatching.RegexBasedMatching.CustomMatchingLocation.KeyCondition.GetOperator()
						value := leafCondition.DatatypeCondition.DatatypeMatching.RegexBasedMatching.CustomMatchingLocation.KeyCondition.GetValue()
						dataTypeMatching.Operator=types.StringValue(string(operator))
						dataTypeMatching.Value=types.StringValue(string(value))
					}
				}
				sources.DataSetDataType = models.DlpRequestBasedDataSetDataTypeFilterSource{
					DataSetDataTypeIds: dataSetDataTypeIds,
					DataTypeMatching:   dataTypeMatching,
				}

			case "IP_ADDRESS":
				ipAddressList, err := utils.ConvertStringPtrToTerraformSet(leafCondition.IpAddressCondition.GetIpAddresses())
				if err != nil {
					return nil, err
				}
				sources.IpAddress = models.DlpRequestBasedIpAddressSource{
					IpAddressList: ipAddressList,
				}

			case "IP_LOCATION_TYPE":
				iplocationtypes, err := utils.ConvertCustomStringPtrsToTerraformSet(leafCondition.IpLocationTypeCondition.GetIpLocationTypes())
				if err != nil {

					return nil, fmt.Errorf("error converting ip location types to terraform list: %v", err)
				}
				sources.IpLocationType = models.DlpRequestBasedIpLocationTypeSource{
					IpLocationTypes: iplocationtypes,
				}

			case "REGION":
				regionIdsPointer := []*string{}
				for _, region := range leafCondition.RegionCondition.GetRegionIdentifiers() {
					regionIdsPointer = append(regionIdsPointer, &region.CountryIsoCode)
				}
				regionIds, err := utils.ConvertStringPtrToTerraformSet(regionIdsPointer)
				if err != nil {
					return nil, err
				}
				sources.Regions = models.DlpRequestBasedRegionsSource{
					RegionIds: regionIds,
				}

			}

		}
	
		requestPayloadObjType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"metadata_type":  types.StringType,
				"key_operator":   types.StringType,
				"key_value":      types.StringType,
				"value_operator": types.StringType,
				"value":          types.StringType,
			},
		}
		if len(requestPayloadArr) > 0 {
			reqresset, diags := types.SetValueFrom(
				ctx,
				requestPayloadObjType,
				requestPayloadArr,
			)
			if diags.HasError() {
				return nil, fmt.Errorf("request response conversion failed")
			}
			sources.RequestPayload = reqresset
		} else {
			sources.RequestPayload = types.SetNull(requestPayloadObjType)
		}
		
		model.Sources = sources
	}

	if data.GetTransactionActionConfigs() != nil {
		config := data.GetTransactionActionConfigs()
		actions := models.DlpAction{}

		switch string(config.Action.GetActionType()) {
		case "ALLOW":
			actions.ActionType = types.StringValue("ALLOW")
			if config.Action.Allow.GetDuration() != nil {
				actions.Duration = types.StringValue(*config.Action.Allow.GetDuration())
			}

		case "ALERT":
			actions.ActionType = types.StringValue("ALERT")
			if HasValue(config.Action.Alert.GetEventSeverity()) {
				actions.EventSeverity = types.StringValue(string(config.Action.Alert.GetEventSeverity()))
			}
		
		case "BLOCK":
			actions.ActionType = types.StringValue("BLOCK")
			if HasValue(config.Action.Block.GetEventSeverity()) {
				actions.EventSeverity = types.StringValue(string(config.Action.Block.GetEventSeverity()))
				if HasValue(config.Action.Block.GetDuration()) {
					actions.Duration = types.StringValue(*config.Action.Block.GetDuration())
				}
			}
		}
		model.Action = actions
	}

	return &model, nil
}


func convertToDlpRequestBasedRuleStatus() (*generated.InputRateLimitingRuleStatus, error) {
	var internal = false
	status := &generated.InputRateLimitingRuleStatus{
		Internal: &internal,
	}
	return status, nil
}

func convertToDlpRequestBasedTransactionActionConfigType(data *models.DataLossPreventionRequestBasedRuleModel) (*generated.InputRateLimitingTransactionActionConfig, error) {
	configTypes := generated.InputRateLimitingTransactionActionConfig{}
	actions := generated.InputRateLimitingRuleAction{}
	if HasValue(data.Action) {
		if HasValue(data.Action.ActionType) {
			switch data.Action.ActionType.ValueString() {
			case "ALERT":
				if HasValue(data.Action.Duration) {
					return nil, utils.NewInvalidError("action duration", "duration not required with action_type alert")
				}

				if !HasValue(data.Action.EventSeverity) {
					return nil, utils.NewInvalidError("action event_severity", "event_severity must present and must not be empty")
				}
				eventSeverity, ok := RateLimitingRuleEventSeverityMap[data.Action.EventSeverity.ValueString()]
				if !ok {
					return nil, utils.NewInvalidError("action event_severity", fmt.Sprintf("%s, is not a valid type of event_severity", data.Action.EventSeverity.ValueString()))
				}
				actions = generated.InputRateLimitingRuleAction{
					ActionType: generated.RateLimitingRuleActionTypeAlert,
					Alert: &generated.InputRateLimitingRuleAlertAction{
						EventSeverity: eventSeverity,
					},
				}
			case "BLOCK":
				if !HasValue(data.Action.EventSeverity) {
					return nil, utils.NewInvalidError("Action EventSeverity", "EventSeverity must present and must not be empty")
				}
				duration := data.Action.Duration.ValueString()
				if duration!=""{
					actions = generated.InputRateLimitingRuleAction{
						ActionType: generated.RateLimitingRuleActionTypeBlock,
						Block: &generated.InputRateLimitingRuleBlockAction{
							EventSeverity: RateLimitingRuleEventSeverityMap[data.Action.EventSeverity.ValueString()],
							Duration:      &duration,
						},
					}
				}else{
					actions = generated.InputRateLimitingRuleAction{
						ActionType: generated.RateLimitingRuleActionTypeBlock,
						Block: &generated.InputRateLimitingRuleBlockAction{
							EventSeverity: RateLimitingRuleEventSeverityMap[data.Action.EventSeverity.ValueString()],
						},
					}
				}
			case "ALLOW":
				if HasValue(data.Action.EventSeverity) {
					return nil, utils.NewInvalidError("Action EventSeverity", "EventSeverity not required with action_type ALLOW")
				}
				duration := data.Action.Duration.ValueString()
				if duration!=""{
					actions = generated.InputRateLimitingRuleAction{
						ActionType: generated.RateLimitingRuleActionTypeAllow,
						Allow: &generated.InputRateLimitingRuleAllowAction{
							Duration: &duration,
						},
					}
				}else{
					actions = generated.InputRateLimitingRuleAction{
						ActionType: generated.RateLimitingRuleActionTypeAllow,
					}
				}
			default:
				return nil, utils.NewInvalidError("Action ActionType", fmt.Sprintf("%s is not a valid action datatype", data.Action.ActionType.ValueString()))
			}

		} else {
			return nil, utils.NewInvalidError("Action ActionType", "must be present and must not be empty")
		}

	} else {
		return nil, utils.NewInvalidError("Action ", "Action must be present and not be empty")
	}

	configTypes = generated.InputRateLimitingTransactionActionConfig{
		Action:          actions,
	}

	return &configTypes, nil
}

func convertDLPRequestBasedModelToCreateInput(ctx context.Context, data *models.DataLossPreventionRequestBasedRuleModel, client *graphql.Client) (*generated.InputRateLimitingRuleData, error) {
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
	status, err := convertToDlpRequestBasedRuleStatus()
	if err != nil {
		return nil, err
	} else {
		input.RuleStatus = status
	}
	transactionActionConfigs, err := convertToDlpRequestBasedTransactionActionConfigType(data)
	if err != nil {
		return nil, err
	} else {
		input.TransactionActionConfigs = transactionActionConfigs
	}
	conditions, err := convertToDlpRequestBasedCondition(ctx, data, client)
	if err != nil {
		return nil, err
	} else {
		input.Conditions = conditions
	}

	return &input, nil
}

func convertToDlpRequestBasedCondition(ctx context.Context, data *models.DataLossPreventionRequestBasedRuleModel, client *graphql.Client) ([]*generated.InputRateLimitingRuleCondition, error) {
	conditions := []*generated.InputRateLimitingRuleCondition{}

	if HasValue(data.Sources.IpLocationType) {
		if !HasValue(data.Sources.IpLocationType.IpLocationTypes) {
			return nil, utils.NewInvalidError("sources ip_location_type ip_location_types", " Must be present and not empty")
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
		exclude := false

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

	if HasValue(data.Sources.IpAddress) {

		if !HasValue(data.Sources.IpAddress.IpAddressList) {
			return nil, utils.NewInvalidError("sources ip_address ip_address_list", " Must be present and not empty")
		}
		ipAddresses := []*string{}
		for _, ipAddress := range data.Sources.IpAddress.IpAddressList.Elements() {
			if ip, ok := ipAddress.(types.String); ok {
				ipAddr := ip.ValueString()
				ipAddresses = append(ipAddresses, &ipAddr)
			}
		}
		exclude := false
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

	if HasValue(data.Sources.Regions) {
		
		if !HasValue(data.Sources.Regions.RegionIds) {
			return nil, utils.NewInvalidError("sources regions region_ids", " Must be present and not empty")
		}
		regionIdentifieres := []*generated.InputRateLimitingRegionIdentifier{}

		regions, err := utils.ConvertSetToStrPointer(data.Sources.Regions.RegionIds)
		if err != nil {
			return nil, fmt.Errorf("converting regions to string pointer fails")
		}
		_, err = GetCountriesId(regions, ctx, *client)
		if err != nil {
			return nil, err
		}
		for _, region := range regions {
			identifiers := &generated.InputRateLimitingRegionIdentifier{
				CountryIsoCode: *region,
			}
			regionIdentifieres = append(regionIdentifieres, identifiers)
		}
		exclude := false
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

	if HasValue(data.Sources.ServiceScope){
	
		if !HasValue(data.Sources.ServiceScope.ServiceIds) {
			return nil, utils.NewInvalidError("sources service_scope service_ids", " Must be present and not empty")
		}
		serviceIds := []*string{}
		for _, serviceId := range data.Sources.ServiceScope.ServiceIds.Elements() {
			if serviceId, ok := serviceId.(types.String); ok {
				serviceIdStr := serviceId.ValueString()
				serviceIds = append(serviceIds, &serviceIdStr)
			}
		}
		scopeType := generated.RateLimitingRuleScopeConditionTypeEntity
		var input = generated.InputRateLimitingRuleCondition{
			LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
				ConditionType: generated.RateLimitingRuleLeafConditionTypeScope,
				ScopeCondition: &generated.InputRateLimitingRuleScopeCondition{
					ScopeType: scopeType,
					EntityScope: &generated.InputRateLimitingRuleEntityScope{
						EntityIds: serviceIds,
						EntityType: generated.RateLimitingRuleEntityTypeService,
					},
				},
			},
		}
		conditions = append(conditions, &input)
	}

	if HasValue(data.Sources.UrlScope) {
		if !HasValue(data.Sources.UrlScope.UrlRegexes) {
			return nil, utils.NewInvalidError("sources url_regex_scope url_regexes", " Must be present and not empty")
		}
		urlRegexes := []*string{}
		for _, urlRegex := range data.Sources.UrlScope.UrlRegexes.Elements() {
			if urlRegex, ok := urlRegex.(types.String); ok {
				urlRegexStr := urlRegex.ValueString()
				urlRegexes = append(urlRegexes, &urlRegexStr)
			}
		}
		scopeType := generated.RateLimitingRuleScopeConditionTypeUrl
		var input = generated.InputRateLimitingRuleCondition{
			LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
				ConditionType: generated.RateLimitingRuleLeafConditionTypeScope,
				ScopeCondition: &generated.InputRateLimitingRuleScopeCondition{
					ScopeType: scopeType,
					UrlScope: &generated.InputRateLimitingRuleUrlScope{
						UrlRegexes: urlRegexes,
					},
				},
			},
		}
		conditions = append(conditions, &input)
	}

	if HasValue(data.Sources.RequestPayload) {
		requestPayloadElement := []models.DlpRequestBasedRequestPayloadScope{}
		err := utils.ConvertElementsSet(data.Sources.RequestPayload, &requestPayloadElement)
		if err != nil {
			return nil, fmt.Errorf("converting request payload set to slice fails")
		}

		for _, requestPayload := range requestPayloadElement {
			keyValueCondition := generated.InputRateLimitingRuleKeyValueCondition{}
				if !HasValue(requestPayload.MetadataType) {
					return nil, utils.NewInvalidError("sources request_payload metadata_type", " Must be present and not empty")
				}
				metadataType, exists := DlpRequestBasedRequestPayloadMetadataTypeMap[requestPayload.MetadataType.ValueString()]
				if !exists {
					return nil, utils.NewInvalidError("sources request_payload metadata_type", fmt.Sprintf(" %s Invalid MetadataType", requestPayload.MetadataType.ValueString()))
				}
				keyValueCondition.MetadataType = &metadataType

				if HasValue(requestPayload.KeyValue) && HasValue(requestPayload.KeyOperator) {
					keyConditionValue := requestPayload.KeyValue.ValueString()
					keyConditionOperator, exist := DlpRequestBasedRequestPayloadKeyOperatorMap[requestPayload.KeyOperator.ValueString()]

					if !exist {
						return nil, utils.NewInvalidError("sources request_payload key_operator", fmt.Sprintf(" %s Invalid keyOperator", requestPayload.KeyOperator.ValueString()))
					}

					keyValueCondition.KeyCondition = &generated.InputRateLimitingRuleStringCondition{
						Operator: keyConditionOperator,
						Value:    keyConditionValue,
					}
				}

				if HasValue(requestPayload.ValueOperator) && HasValue(requestPayload.Value) {
					valueConditionValue := requestPayload.Value.ValueString()
					valueConditionOperator, exist := RateLimitingKeyValueMatchOperatorMap[requestPayload.ValueOperator.ValueString()]

					if !exist {
						return nil, utils.NewInvalidError("sources request_payload value_operator", fmt.Sprintf(" %s Invalid keyOperator", requestPayload.ValueOperator.ValueString()))
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

	if HasValue(data.Sources.DataSetDataType) {
		var input = generated.InputRateLimitingRuleCondition{}
		dataSetIds := []*string{}
		if HasValue(data.Sources.DataSetDataType.DataSetDataTypeIds.DataSetsIds) {
			for _, dataSetId := range data.Sources.DataSetDataType.DataSetDataTypeIds.DataSetsIds.Elements() {
				if dataSetId, ok := dataSetId.(types.String); ok {
					dataSetIdStr := dataSetId.ValueString()
					dataSetIds = append(dataSetIds, &dataSetIdStr)
				}
			}
		}
		dataTypeIds := []*string{}
		if HasValue(data.Sources.DataSetDataType.DataSetDataTypeIds.DataTypesIds) {
			for _, dataTypeId := range data.Sources.DataSetDataType.DataSetDataTypeIds.DataTypesIds.Elements() {
				if dataTypeId, ok := dataTypeId.(types.String); ok {
					dataTypeIdStr := dataTypeId.ValueString()
					dataTypeIds = append(dataTypeIds, &dataTypeIdStr)
				}
			}
		}
		// log.Printf("datasets %s,datatypes %s",dataSetIds,dataTypeIds)
		if (len(dataSetIds)==0 && len(dataTypeIds)==0) || (len(dataSetIds)>0 && len(dataTypeIds)>0){
			return nil,utils.NewInvalidError("data_sets_ids or data_types_ids", " Must be present and not empty")
		}
		dataLocation := generated.RateLimitingRuleDataLocationRequest
		if HasValue(data.Sources.DataSetDataType.DataTypeMatching) {
			metadataType,ok := DlpRequestBasedDatatypeMatchingMetadataTypeMap[data.Sources.DataSetDataType.DataTypeMatching.MetadataType.ValueString()]
			if !ok {
				return nil,utils.NewInvalidError("dateset_datatype_filter data_type_matching metadata_type",fmt.Sprintf("Invalid metadata_type %s", data.Sources.DataSetDataType.DataTypeMatching.MetadataType.ValueString()))
			}
			datatypeMatchingType := generated.RateLimitingRuleDatatypeMatchingTypeRegexBasedMatching
			log.Printf("oppp %s val %s",data.Sources.DataSetDataType.DataTypeMatching.Operator.ValueString(),data.Sources.DataSetDataType.DataTypeMatching.Value.ValueString())
			if metadataType == generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestBody && (data.Sources.DataSetDataType.DataTypeMatching.Operator.ValueString()!="" || data.Sources.DataSetDataType.DataTypeMatching.Value.ValueString()!=""){
				return nil,utils.NewInvalidError("dateset_datatype_filter data_type_matching operator value","Operator and Value not required for metadata_type REQUEST_BODY")
			}
			if metadataType != generated.RateLimitingRuleKeyValueConditionMetadataTypeRequestBody && (data.Sources.DataSetDataType.DataTypeMatching.Operator.ValueString()=="" || data.Sources.DataSetDataType.DataTypeMatching.Value.ValueString()==""){
				return nil,utils.NewInvalidError("dateset_datatype_filter data_type_matching operator value",fmt.Sprintf("Operator and Value required for metadata_type %s", data.Sources.DataSetDataType.DataTypeMatching.MetadataType.ValueString()))
			}
			if HasValue(data.Sources.DataSetDataType.DataTypeMatching.Operator) && HasValue(data.Sources.DataSetDataType.DataTypeMatching.Value) {
				operator,ok := DlpRequestBasedDatatypeMatchingKeyOperatorMap[data.Sources.DataSetDataType.DataTypeMatching.Operator.ValueString()]
				if !ok {
					return nil,utils.NewInvalidError("dateset_datatype_filter data_type_matching operator",fmt.Sprintf("Invalid operator %s", data.Sources.DataSetDataType.DataTypeMatching.Operator.ValueString()))
				}
				value := data.Sources.DataSetDataType.DataTypeMatching.Value.ValueString()
				input = generated.InputRateLimitingRuleCondition{
					LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
						ConditionType: generated.RateLimitingRuleLeafConditionTypeDatatype,
						DatatypeCondition: &generated.InputRateLimitingRuleDatatypeCondition{
							DataLocation: &dataLocation,
							DatasetIds: dataSetIds,
							DatatypeIds: dataTypeIds,
							DatatypeMatching: &generated.InputRateLimitingRuleDatatypeMatching{
								DatatypeMatchingType : &datatypeMatchingType,
								RegexBasedMatching : &generated.InputRateLimitingRuleRegexBasedMatching{
									CustomMatchingLocation : &generated.InputRateLimitingRuleKeyValueCondition{
										MetadataType : &metadataType,
										KeyCondition: &generated.InputRateLimitingRuleStringCondition{
											Operator: operator,
											Value: value,
										},
									},
								},
							},
						},
					},
				}
			}else{
				input = generated.InputRateLimitingRuleCondition{
					LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
						ConditionType: generated.RateLimitingRuleLeafConditionTypeDatatype,
						DatatypeCondition: &generated.InputRateLimitingRuleDatatypeCondition{
							DataLocation: &dataLocation,
							DatasetIds: dataSetIds,
							DatatypeIds: dataTypeIds,
							DatatypeMatching: &generated.InputRateLimitingRuleDatatypeMatching{
								DatatypeMatchingType : &datatypeMatchingType,
								RegexBasedMatching : &generated.InputRateLimitingRuleRegexBasedMatching{
									CustomMatchingLocation : &generated.InputRateLimitingRuleKeyValueCondition{
										MetadataType : &metadataType,
									},
								},
							},
						},
					},
				}
			}
		}else{
			input = generated.InputRateLimitingRuleCondition{
				LeafCondition: &generated.InputRateLimitingRuleLeafCondition{
					ConditionType: generated.RateLimitingRuleLeafConditionTypeDatatype,
					DatatypeCondition: &generated.InputRateLimitingRuleDatatypeCondition{
						DataLocation: &dataLocation,
						DatasetIds: dataSetIds,
						DatatypeIds: dataTypeIds,
					},
				},
			}
		}
		conditions = append(conditions, &input)
	}
	return conditions, nil
}


func convertDLPRequestBasedModelToUpdateInput(ctx context.Context, data *models.DataLossPreventionRequestBasedRuleModel, id string, client *graphql.Client) (*generated.InputRateLimitingRule, error) {
	var input = generated.InputRateLimitingRule{}
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
	status, err := convertToDlpRequestBasedRuleStatus()
	if err != nil {
		return nil, err
	} else {
		input.RuleStatus = status
	}
	transactionActionConfigs, err := convertToDlpRequestBasedTransactionActionConfigType(data)
	if err != nil {
		return nil, err
	} else {
		input.TransactionActionConfigs = transactionActionConfigs
	}
	conditions, err := convertToDlpRequestBasedCondition(ctx, data, client)
	if err != nil {
		return nil, err
	} else {
		input.Conditions = conditions
	}

	return &input, nil
}