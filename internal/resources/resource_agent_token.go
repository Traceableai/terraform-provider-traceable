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

type AgentTokenResource struct {
	client *graphql.Client
}

func NewAgentTokenResource() resource.Resource {
	return &AgentTokenResource{}
}

func (r *AgentTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AgentTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_token"
}

func (r *AgentTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.AgentTokenResourceSchema()
}

func (r *AgentTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Entering Create Block")
	var data *models.AgentTokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, err := convertAgentTokenModelToCreateInput(ctx, data)
	if input == nil || err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	agentToken, err := generated.CreateAgentToken(ctx, *r.client, *input)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	data.Id = types.StringValue(*agentToken.CreateAgentToken.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Exiting Create Block")
}

func (r *AgentTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.AgentTokenModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := getAgentToken(data.Id.ValueString(), ctx, r.client)
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	updatedData, err := convertAgentTokenFieldsToModel(ctx, response)

	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedData)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AgentTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.AgentTokenModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var dataState *models.AgentTokenModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)

	if resp.Diagnostics.HasError() {
		return
	}
	
	input, err := convertAgentTokenModelToUpdateInput(ctx, data, dataState.Id.ValueString())
	if err != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}

	resp1, err2 := generated.UpdateAgentToken(ctx, *r.client, *input)
	if err2 != nil {
		utils.AddError(ctx, &resp.Diagnostics, err)
		return
	}
	data.Id = types.StringValue(*resp1.UpdateAgentTokenMetadata.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AgentTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.AgentTokenModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := generated.DeleteAgentToken(ctx, *r.client, generated.InputDeleteAgentTokenInput{Id: data.Id.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting agent token", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func convertAgentTokenModelToCreateInput(ctx context.Context,data *models.AgentTokenModel) (*generated.InputCreateAgentTokenInput,error){
	if data == nil {
		return nil, fmt.Errorf("data model is nil")
	}
	var input = generated.InputCreateAgentTokenInput{}
	if HasValue(data.Name) {
		name := data.Name.ValueString()
		input.Name = name
	} else {
		return nil, utils.NewInvalidError("Name", "Name field must be present and must not be empty")
	}
	return &input,nil
}

func convertAgentTokenModelToUpdateInput(ctx context.Context,data *models.AgentTokenModel,id string) (*generated.InputUpdateAgentTokenMetadataInput,error){
	if data == nil {
		return nil, fmt.Errorf("data model is nil")
	}
	var input = generated.InputUpdateAgentTokenMetadataInput{}
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
	return &input,nil
}

func getAgentToken(id string, ctx context.Context, client *graphql.Client) (*generated.GetAgentTokenAgentTokenMetadataAgentTokenMetadataResultSetResultsAgentTokenMetadata, error) {
	agentTokenFeilds := generated.GetAgentTokenAgentTokenMetadataAgentTokenMetadataResultSetResultsAgentTokenMetadata{}
	response, err := generated.GetAgentToken(ctx, *client)
	if err != nil {
		return nil, err
	}

	for _, rule := range response.AgentTokenMetadata.Results {
		if *rule.Id == id {
			agentTokenFeilds = *rule
			return &agentTokenFeilds, nil
		}
	}

	return nil, nil
}

func convertAgentTokenFieldsToModel(ctx context.Context, fields *generated.GetAgentTokenAgentTokenMetadataAgentTokenMetadataResultSetResultsAgentTokenMetadata) (*models.AgentTokenModel, error) {
	if fields == nil {
		return nil, fmt.Errorf("fields model is nil")
	}
	var data = models.AgentTokenModel{}
	data.Id = types.StringValue(*fields.Id)
	data.Name = types.StringValue(fields.Name)
	return &data, nil
}