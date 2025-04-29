package datasources

import (
	"context"
	"fmt"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/traceableai/terraform-provider-traceable/internal/generated"
	"github.com/traceableai/terraform-provider-traceable/internal/models"
	"github.com/traceableai/terraform-provider-traceable/internal/schemas"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func NewEndpointDataSource() datasource.DataSource {
	return &EndpointDataSource{}
}

type EndpointDataSource struct {
	client *graphql.Client
}

func (d *EndpointDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *EndpointDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoints"
}

func (d *EndpointDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.EndpointDataSourceSchema()
}

func (d *EndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.EndpointDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !utils.HasValue(config.Endpoints) {
		resp.Diagnostics.AddError("endpoints field must be present and must not be empty", "")
	}
	endpointPtr, err := utils.ConvertSetToStrPointer(config.Endpoints)
	if err != nil {
		resp.Diagnostics.AddError("Error converting set to string pointer", err.Error())
		return
	}
	endpointIds, err := GetEndpointIds(endpointPtr, ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting endpoint labels", err.Error())
		return
	}
	endpointIdSet, err := utils.ConvertStringPtrToTerraformSet(endpointIds)
	if err != nil {
		resp.Diagnostics.AddError("Error converting string pointer to terraform set", err.Error())
		return
	}
	config.EndpointIds = endpointIdSet
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func GetEndpointIds(endpointNames []*string, ctx context.Context, r graphql.Client) ([]*string, error) {
	endpointIds := []*string{}
	entityType := generated.EntityTypeApi

	currentTime := time.Now().UTC()
	endTime := currentTime.Format("2006-01-02T15:04:05.000Z")
	lastWeekTime := currentTime.AddDate(0, 0, -7)
	startTime := lastWeekTime.Format("2006-01-02T15:04:05.000Z")
	between := generated.InputTimeRange{
		StartTime: startTime,
		EndTime:   endTime,
	}

	includeInactive := true
	offset := int64(0)
	limit := int64(100)

	apiDiscoveryKeyExpression := generated.InputAttributeExpression{
		Key: "apiDiscoveryState",
	}
	discoveryKeyExpression := generated.InputAttributeExpression{
		Key: "discoverySources",
	}

	apiDiscoveryFilter := generated.InputFilter{
		KeyExpression: &apiDiscoveryKeyExpression,
		Operator:      generated.FilterOperatorTypeIn,
		Value:         []string{"DISCOVERED", "UNDER_DISCOVERY"},
		Type:          generated.FilterTypeAttribute,
	}

	discoveryFilter := generated.InputFilter{
		KeyExpression: &discoveryKeyExpression,
		Operator:      generated.FilterOperatorTypeIn,
		Value:         []string{"Live Traffic"},
		Type:          generated.FilterTypeAttribute,
	}

	filterBy := []*generated.InputFilter{
		&apiDiscoveryFilter,
		&discoveryFilter,
	}
	scope := "API"

	response, err := generated.GetEndpointIds(ctx, r, &entityType, &scope, between, nil, filterBy, nil, &limit, &offset, &includeInactive)

	if err != nil {
		return nil, err
	}

	endpointPresent := map[string]bool{}
	for _, endpointName := range endpointNames {
		if endpointName != nil {
			endpointPresent[*endpointName] = true
		}
	}

	for _, edp := range response.GetEntities().Results {
		if edp.Name != nil {
			nameInterface := *edp.Name
			endpointName, okName := nameInterface.(string)
			endpointId := edp.EntityId

			if okName && endpointPresent[endpointName] {
				idCopy := endpointId
				endpointIds = append(endpointIds, &idCopy)
				endpointPresent[endpointName] = false
			}
		}
	}

	for _, endpointName := range endpointNames {
		if endpointName != nil && endpointPresent[*endpointName] {
			return nil, utils.NewInvalidError("endpoint", fmt.Sprintf("%s is not a supported endpoint", *endpointName))
		}
	}

	return endpointIds, nil
}
