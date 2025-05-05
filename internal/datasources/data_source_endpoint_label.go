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

func NewEndpointLabelDataSource() datasource.DataSource {
	return &EndpointLabelDataSource{}
}

type EndpointLabelDataSource struct {
	client *graphql.Client
}

func (d *EndpointLabelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EndpointLabelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_labels"
}

func (d *EndpointLabelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.EndpointLabelsDataSourceSchema()
}

func (d *EndpointLabelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.EndpointLabelsDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !utils.HasValue(config.Labels) {
		resp.Diagnostics.AddError("Labels field must be present and must not be empty", "")
	}
	labelptr, err := utils.ConvertSetToStrPointer(config.Labels)
	if err != nil {
		resp.Diagnostics.AddError("Error converting set to string pointer", err.Error())
		return
	}
	labelIds, err := GetEndpointLabeslId(labelptr, ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError("Error getting endpoint labels", err.Error())
		return
	}
	labelsId, err := utils.ConvertStringPtrToTerraformSet(labelIds)
	if err != nil {
		resp.Diagnostics.AddError("Error converting string pointer to terraform set", err.Error())
		return
	}
	config.LabelIds = labelsId
	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func GetEndpointLabeslId(endpointLabels []*string, ctx context.Context, r graphql.Client) ([]*string, error) {
	endpointLabelIds := []*string{}
	response, err := generated.GetEndpointLabelsId(ctx, r)
	if err != nil {
		return nil, err
	}
	endpointLabelsPresent := map[string]bool{}
	for _, key := range endpointLabels {
		endpointLabelsPresent[*key] = true
	}
	for _, label := range response.Labels.Results {
		if endpointLabelsPresent[label.Key] {
			endpointLabelIds = append(endpointLabelIds, &label.Id)
			endpointLabelsPresent[label.Key] = false
		}
	}

	for _, key := range endpointLabels {
		if endpointLabelsPresent[*key] {
			return nil, utils.NewInvalidError("endpoint_labels", fmt.Sprintf("%s is not a supported endpoint label", *key))
		}
	}
	return endpointLabelIds, nil
}
