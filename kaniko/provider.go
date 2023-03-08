package kaniko

import (
	"context"
	"github.com/gitlawr/terraform-provider-kaniko/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure kanikoProvider satisfies various provider interfaces.
var _ provider.Provider = &kanikoProvider{}

// kanikoProvider defines the provider implementation.
type kanikoProvider struct {
	version string
}

// kanikoProviderModel describes the provider data model.
type kanikoProviderModel struct {
	ConfigPath types.String `tfsdk:"config_path"`
}

func (p *kanikoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kaniko"
	resp.Version = p.version
}

func (p *kanikoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config_path": schema.StringAttribute{
				Description: "Path to the kube config file.",
				Optional:    true,
			},
		},
	}
}

func (p *kanikoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring kaniko client")

	var config kanikoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}
	restConfig, err := utils.GetConfig(config.ConfigPath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get rest.Config", err.Error())
		return
	}

	resp.DataSourceData = restConfig
	resp.ResourceData = restConfig
}

func (p *kanikoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewImageResource,
	}
}

func (p *kanikoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kanikoProvider{
			version: version,
		}
	}
}
