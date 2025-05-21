package datasources

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/traceableai/terraform-provider-traceable/internal/generated"
	"github.com/traceableai/terraform-provider-traceable/internal/models"
	"github.com/traceableai/terraform-provider-traceable/internal/schemas"
	"github.com/traceableai/terraform-provider-traceable/internal/utils"
)

func NewDatasetDataSource() datasource.DataSource {
	return &DatasetDataSource{}
}

type DatasetDataSource struct {
	client *graphql.Client
}

func (d *DatasetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DatasetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datasets"
}

func (d *DatasetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.DatasetDataSourceSchema()
}

func (d *DatasetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.DatasetDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !utils.HasValue(config.DataSets) {
		resp.Diagnostics.AddError("Data set field must be present and must not be empty", "")
	}
	dataSetPtr, err := utils.ConvertSetToStrPointer(config.DataSets)
	if err != nil {
		resp.Diagnostics.AddError("Error converting set to string pointer", err.Error())
		return
	}
	dataSetIds, err := GetDataSetId(dataSetPtr, ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting endpoint labels", err.Error())
		return
	}
	dataSetId, err := utils.ConvertStringPtrToTerraformSet(dataSetIds)
	if err != nil {
		resp.Diagnostics.AddError("Error converting string pointer to terraform set", err.Error())
		return
	}
	config.DataSetIds = dataSetId
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func GetDataSetId(dataSetNames []*string, ctx context.Context, r graphql.Client) ([]*string, error) {
	dataSetIds := []*string{}
	response, err := generated.GetDataSetsId(ctx, r)
	if err != nil {
		return nil, err
	}
	dataSetNamesPresent := map[string]bool{}
	for _, key := range dataSetNames {
		dataSetNamesPresent[*key] = true
	}
	for _, dataSet := range response.DataSets.Results {
		if dataSetNamesPresent[dataSet.Name] {
			dataSetIds = append(dataSetIds, &dataSet.Id)
			dataSetNamesPresent[dataSet.Name] = false
		}
	}
	for _, key := range dataSetNames {
		if dataSetNamesPresent[*key] {
			return nil, utils.NewInvalidError("data_set", fmt.Sprintf("%s is not a supported data set", *key))
		}
	}
	return dataSetIds, nil
}
