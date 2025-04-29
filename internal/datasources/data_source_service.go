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

func NewServiceDataSource() datasource.DataSource {
	return &ServiceDataSource{}
}

type ServiceDataSource struct {
	client *graphql.Client
}

func (d *ServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (d *ServiceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.ServiceDataSourceSchema()
}

func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.ServiceDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !utils.HasValue(config.Services) {
		resp.Diagnostics.AddError("services field must be present and must not be empty", "")
	}
	servicePtr, err := utils.ConvertSetToStrPointer(config.Services)
	if err != nil {
		resp.Diagnostics.AddError("Error converting set to string pointer", err.Error())
		return
	}
	serviceIds, err := GetServiceIds(servicePtr, ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting service labels", err.Error())
		return
	}
	serviceIdSet, err := utils.ConvertStringPtrToTerraformSet(serviceIds)
	if err != nil {
		resp.Diagnostics.AddError("Error converting string pointer to terraform set", err.Error())
		return
	}
	config.ServiceIds = serviceIdSet
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func GetServiceIds(serviceNames []*string, ctx context.Context, r graphql.Client) ([]*string, error) {
	serviceIds := []*string{}
	entityType := generated.EntityTypeService

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

	nameKeyExpression := generated.InputAttributeExpression{
		Key: "name",
	}

	nameFilter := generated.InputFilter{
		KeyExpression: &nameKeyExpression,
		Operator:      generated.FilterOperatorTypeLike,
		Value:         "",
		Type:          generated.FilterTypeAttribute,
	}

	filterBy := []*generated.InputFilter{
		&nameFilter,
	}
	scope := "SERVICE"

	response, err := generated.GetEntitiesIds(ctx, r, &entityType, &scope, between, nil, filterBy, nil, &limit, &offset, &includeInactive)

	if err != nil {
		return nil, err
	}

	servicePresent := map[string]bool{}
	for _, serviceName := range serviceNames {
		if serviceName != nil {
			servicePresent[*serviceName] = true
		}
	}

	for _, service := range response.GetEntities().Results {
		if service.Name != nil {
			nameInterface := *service.Name
			serviceName, okName := nameInterface.(string)
			serviceId := service.EntityId

			if okName && servicePresent[serviceName] {
				idCopy := serviceId
				serviceIds = append(serviceIds, &idCopy)
				servicePresent[serviceName] = false
			}
		}
	}

	for _, serviceName := range serviceNames {
		if serviceName != nil && servicePresent[*serviceName] {
			return nil, utils.NewInvalidError("service", fmt.Sprintf("%s is not a supported service", *serviceName))
		}
	}

	return serviceIds, nil
}
