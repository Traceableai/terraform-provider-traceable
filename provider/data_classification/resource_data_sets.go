package data_classification

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

)

func NewDataSetResource() resource.Resource {
	return &DataSetResource{}
}

type DataSetResource struct {
	client *graphql.Client
}
type DataSetResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	IconType    types.String `tfsdk:"icon_type"`
}

func (r *DataSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_set"
}

func (r *DataSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Traceable DataSet",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the Data Set",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the DataSet.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the DataSet",
				Optional:            true,
			},
			"icon_type": schema.StringAttribute{
				MarkdownDescription: "Icon Type of the DataSet",
				Optional:            true,
			},
		},
}


}



func (r *DataSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	fmt.Println("hello inside resource")
	tflog.Trace(ctx,"Client Intialization Successfully")
}



func(r * DataSetResource)Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse){
	tflog.Trace(ctx,"Entering in Create Block")
	var data * DataSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existDataSet,err:=getDataSet(ctx,data.Name.ValueString(),r.client)

	if(err!=nil){
		resp.Diagnostics.AddError("Failed to Connect to Platform",
			fmt.Sprintf("Some internal error occurs %s check your token and url once and retry",err.Error()),
		)
		return
	}

	
	if(existDataSet!=nil){
		tflog.Trace(ctx,"DataSet is already present")
		resp.Diagnostics.AddError("Data Set already Exist",
		fmt.Sprintf("DataSet with name %s, already exist please import (in case of import import id is name of Data Set)or try using other name", data.Name.ValueString()),
		
	)
	return
	}

	 input:= InputDataSetCreate{
		  Name: data.Name.ValueString(),
			Description: data.Description.ValueString(),
			IconType: data.IconType.ValueString(),
		
		}
   
   dataSet, err1:=createDataSet(ctx,*r.client,input)
	 tflog.Debug(ctx,"create Data Set response",map[string]interface{}{
		"resp":dataSet.CreateDataSet,
	 })

	 if(err1!=nil){
		resp.Diagnostics.AddError("Failed to Connect to Platform",
			fmt.Sprintf("Some internal error occurs %s check your token and url once and retry",err.Error()),
		)
	return
	 }

	 
   if(dataSet.CreateDataSet.Id == ""){
		resp.Diagnostics.AddError("Some Internal Error Occur","Error Message: After creating dataset Id is empty" )
	 }
	 data.Id=types.StringValue(dataSet.CreateDataSet.Id)
	 resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}


func (r * DataSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse){
	var data *DataSetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}


	dataSet,err:=getDataSet(ctx,data.Name.ValueString(),r.client)

	if(err ==nil){
     
	}
	if(dataSet==nil){
    
	}

	data.Id = types.StringValue(dataSet.Id)
	data.Name = types.StringValue(dataSet.Name)
  data.IconType=types.StringValue(dataSet.IconType)
	data.Description=types.StringValue(dataSet.Description)
  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)


}
func (r *DataSetResource)Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse){
	var data *DataSetResourceModel 

	var dataState * DataSetResourceModel 

	resp.Diagnostics.Append(req.State.Get(ctx, &dataState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	input:=InputDataSetUpdate{
		Id: dataState.Id.ValueString(),
		Name:data.Name.ValueString(),
		IconType: data.IconType.ValueString(),
	  Description: data.Description.ValueString(),
	}
	tflog.Trace(ctx,"Not able to update 1 ")
  tflog.Trace(ctx,"updatedDataSetResPonse",map[string]interface{}{
		"updated shreyansh 1 ":input,
		"client":*r.client,
		"ctx":ctx,
	})
	updateDataSetResponse,err :=updateDataSet(ctx,*r.client,input)

	
	if(err != nil){
		resp.Diagnostics.AddError("Failed to Connect to Platform",
		fmt.Sprintf("Some internal error occurs %s check your token and url once and retry",err.Error()),
	)
		return
	} 
	


	data.Description=types.StringValue(updateDataSetResponse.UpdateDataSet.Description)
	data.IconType=types.StringValue(updateDataSetResponse.UpdateDataSet.IconType)
	data.Name=types.StringValue(updateDataSetResponse.UpdateDataSet.Name)
	data.Id=types.StringValue(updateDataSetResponse.UpdateDataSet.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)


}
func( r *DataSetResource)Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse){

	var data *DataSetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_,err:=deleteDataSet(ctx,*r.client,data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Failed to Connect to Platform",
		fmt.Sprintf("Some internal error occurs %s check your token and url once and retry",err.Error()),
	)
		return
	}
	resp.State.RemoveResource(ctx)
	
}
func (r *DataSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
     dataSetName:=req.ID

		 dataSet,err:=getDataSet(ctx,dataSetName,r.client)
		 if err !=nil{
			resp.Diagnostics.AddError("Failed to Connect to Platform",
		fmt.Sprintf("Some internal error occurs %s check your token and url once and retry",err.Error()),	)
		  return
		 }

		 if dataSet==nil{
			resp.Diagnostics.AddError("Failed to Connect to Platform",
			fmt.Sprintf("Some internal error occurs %s check your token and url once and retry",err.Error()),
		)
			return
		 }

		 resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(dataSet.Id))...)
		 resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"),types.StringValue(dataSet.Name) )...)
		 resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"),types.StringValue(dataSet.Description) )...)
		 resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("icon_type"),types.StringValue(dataSet.IconType) )...)
}


func getDataSet(ctx context.Context,dataSetName string,client *graphql.Client ) (*dataSetsDataSetsDataSetResultSetResultsDataSet,error) {
	dataSetsResponse,err := dataSets(ctx,*client)
	if err !=nil {
		return nil,err
	}
	  for _ , dataSet:= range dataSetsResponse.DataSets.Results{
			if(dataSet.Name==dataSetName){
           return &dataSet,nil
			}			   
		}
		return nil,nil
}
















