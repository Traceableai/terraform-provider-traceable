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

func NewDataTypeDataSource() datasource.DataSource {
	return &DataTypeDataSource{}
}

type DataTypeDataSource struct {
	client *graphql.Client
}

func (d *DataTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_types"
}

func (d *DataTypeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.DataTypeDataSourceSchema()
}

func (d *DataTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.DataTypeDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !utils.HasValue(config.DataTypes) {
		resp.Diagnostics.AddError("Data types field must be present and must not be empty", "")
	}
	dataTypePtr, err := utils.ConvertSetToStrPointer(config.DataTypes)
	if err != nil {
		resp.Diagnostics.AddError("Error converting set to string pointer", err.Error())
		return
	}
	dataTypeIds, err := GetDataTypeId(dataTypePtr, ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting endpoint labels", err.Error())
		return
	}

	dataTypeId, err := utils.ConvertStringPtrToTerraformSet(dataTypeIds)
	if err != nil {
		resp.Diagnostics.AddError("Error converting string pointer to terraform set", err.Error())
		return
	}
	config.DataTypeIds = dataTypeId
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func GetDataTypeId(dataTypeNames []*string, ctx context.Context, r graphql.Client) ([]*string, error) {
	dataTypeIds := []*string{}
	response, err := generated.GetDataTypesId(ctx, r)
	if err != nil {
		return nil, err
	}
	dataTypeNamesPresent := map[string]bool{}
	for _, key := range dataTypeNames {
		dataTypeNamesPresent[*key] = true
	}
	for _, dataType := range response.DataTypes.Results {
		if dataTypeNamesPresent[dataType.Name] {
			dataTypeIds = append(dataTypeIds, &dataType.Id)
			dataTypeNamesPresent[dataType.Name] = false
		}
	}
	for _, key := range dataTypeNames {
		if dataTypeNamesPresent[*key] {
			return nil, utils.NewInvalidError("data_type", fmt.Sprintf("%s is not a supported data type", *key))
		}
	}
	return dataTypeIds, nil
}
