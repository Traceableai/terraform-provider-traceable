package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/api"
	"github.com/traceableai/terraform-provider-traceable/shreyanshrevamp/internal/resources"
)

// Ensure provider implements provider.Provider
var _ provider.Provider = &traceableProvider{}

type traceableProvider struct {
	version string
}

type traceableProviderModel struct {
	PlatformUrl types.String `tfsdk:"platform_url"`
	ApiToken    types.String `tfsdk:"api_token"`
}

// Metadata function
func (p *traceableProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	fmt.Println("metadata intializattion")
	resp.TypeName = "traceable"
	resp.Version = p.version
}

// Schema function
func (p *traceableProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	fmt.Println("schema intializattion")
	resp.Schema = schema.Schema{
		MarkdownDescription: `
The  provider allows you to interact with Traceable Platform, managing resources on it. 
It supports creating and

Refer to the official [Traceable Documentation](https://traceable.ai) for more details.
		`,
		Attributes: map[string]schema.Attribute{
			"platform_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Url to be used to connect to Traceable Platform.",
			},
			"api_token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The API token to be used to connect to Traceable Platform.",
			},
		},
	}
}

// Configure function - Extracts Terraform version
func (p *traceableProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data traceableProviderModel
	fmt.Println("Reading provider configuration...")

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Debug: Print diagnostics
	if diags.HasError() {
		fmt.Println("Error reading config:", diags)
		return
	}
	tfVersion := req.TerraformVersion

	if tfVersion == "" {
		resp.Diagnostics.AddWarning("Terraform Version Missing", "Unable to determine Terraform version from request.")
	}
	url := data.PlatformUrl.ValueString()
	token := data.ApiToken.ValueString()

	if url == "" {
		resp.Diagnostics.AddError(
			"Missing Platform URL",
			"The platform_url attribute is required. Please provide a valid Traceable Platform URL.")
		return
	}
	if token == "" {
		resp.Diagnostics.AddError(
			"Missing API Token",
			"The api_token attribute is required. Please provide a valid Traceable API token.")
		return
	}
	client := api.NewClient(url, token, tfVersion)
	if client == nil {
		resp.Diagnostics.AddError("Failed to initialize API client", "The client could not be created. Check API URL and token.")
		return
	}
	resp.DataSourceData = &client
	resp.ResourceData = &client
	tflog.Info(ctx, fmt.Sprintf("Traceable client successfully initialized for %s with version %s", url, tfVersion))

}

// Register your resources
func (p *traceableProvider) Resources(ctx context.Context) []func() resource.Resource {
	fmt.Println("resouces intializattion")

	return []func() resource.Resource{
		resources.NewRateLimitingResource,
		resources.NewDataSetResource,
	}
}

// Register your data sources
func (p *traceableProvider) DataSources(ctx context.Context) []func() datasource.DataSource {

	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {

	return func() provider.Provider {
		return &traceableProvider{
			version: version,
		}

	}

}
