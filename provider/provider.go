//"shell.com/terraform-provider-wiz/apiClient"

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"shell.com/terraform-provider-wiz/apiClient"
)

type provider struct {
	wizClient apiClient.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type providerData struct {
	Endpoint     types.String `tfsdk:"endpoint"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ClientId.Null || data.ClientSecret.Null {
		resp.Diagnostics.AddError("No valid credentials provided.", "Both client_id and client_secret are needed.")
		return
	}

	config := apiClient.ClientConfig{
		Credentials: apiClient.ClientCredentials{
			ClientID:     data.ClientId.Value,
			ClientSecret: data.ClientSecret.Value,
			Endpoint:     data.Endpoint.Value,
		},
	}

	client, err := apiClient.CreateClient(config)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("failed to create client: %s", err.Error()))
		resp.Diagnostics.AddError("failed to create client.", err.Error())
	}

	p.wizClient = *client

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"wiz_project": resourceWizProjectType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		//register tfsdk datasources here
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				MarkdownDescription: "The base URL of the Wiz API.",
				Optional:            true,
				Type:                types.StringType,
			},
			"client_id": {
				MarkdownDescription: "service account client_id used to log in to Wiz.",
				Optional:            true,
				Type:                types.StringType,
			},
			"client_secret": {
				MarkdownDescription: "service account secret (client_secret) used to log in to Wiz.",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
