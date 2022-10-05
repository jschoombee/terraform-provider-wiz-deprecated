package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"shell.com/terraform-provider-wiz/apiClient"
)

type resourceWizProjectType struct{}

type wizProject struct {
	provider provider
}

type wizProjectTypeData struct {
	ID                *string                    `tfsdk:"id"`
	Name              *string                    `tfsdk:"name"`
	CloudAccountLinks []CloudAccountLinkTypeData `tfsdk:"cloud_account_links"`
}

type CloudAccountLinkTypeData struct {
	GUID        string `tfsdk:"cloud_account_guid"`
	Environment string `tfsdk:"environment"`
	Shared      bool   `tfsdk:"shared"`
}

func (t resourceWizProjectType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Wiz Project configuration.",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of the Project",
				Computed:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Project Name",
				Required:            true,
				Type:                types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"cloud_account_links": {
				MarkdownDescription: "A List of cloud account ids",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(
					map[string]tfsdk.Attribute{
						"cloud_account_guid": {
							MarkdownDescription: "GUID of the cloud account",
							Required:            true,
							Type:                types.StringType,
						},
						"environment": {
							MarkdownDescription: "DTAP environment",
							Optional:            true,
							Type:                types.StringType,
						},
						"shared": {
							MarkdownDescription: "Optional event specification",
							Optional:            true,
							Type:                types.BoolType,
						},
					},
					tfsdk.ListNestedAttributesOptions{},
				),
			},
		},
	}, nil
}

func (t resourceWizProjectType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return wizProject{
		provider: provider,
	}, diags
}

func (d wizProjectTypeData) getAccountLinks(ctx context.Context) []apiClient.CloudAccountLink {
	var cloudAccountLinks []apiClient.CloudAccountLink

	for _, cl := range d.CloudAccountLinks {
		cloudAccountLinks = append(cloudAccountLinks, apiClient.CloudAccountLink{
			CloudAccount: cl.GUID,
			Environment:  cl.Environment,
			Shared:       cl.Shared,
		})
	}

	return cloudAccountLinks
}

func (d wizProjectTypeData) setAccountLinks(ctx context.Context, cloudAccountLinks apiClient.Subscriptions) {
	for _, cl := range cloudAccountLinks {
		d.CloudAccountLinks = append(d.CloudAccountLinks, CloudAccountLinkTypeData{
			GUID:        cl.SubscriptionID,
			Environment: cl.Environments[0],
			Shared:      cl.SharedAccount,
		})
	}
}

func (r wizProject) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {

	var data wizProjectTypeData

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.wizClient.CreateWizProject(ctx, apiClient.CreateProjectRequest{
		Input: apiClient.CreateProjectInput{
			Name: data.Name,
			// TODO we are passing context but it isn't being used, either use it os lose it
			CloudAccountLinks: data.getAccountLinks(ctx),
			// TODO remove hard-coded risk profile until we know what we need to do
			RiskProfile: apiClient.RiskProfile{
				BusinessImpact:      "MBI",
				HasAuthentication:   "UNKNOWN",
				HasExposedAPI:       "YES",
				IsCustomerFacing:    "NO",
				IsInternetFacing:    "YES",
				IsActivelyDeveloped: "UNKNOWN",
				IsRegulated:         "YES",
				SensitiveDataTypes:  []string{"CUSTOMER", "FINANCIAL"},
				StoresData:          "YES",
				RegulatoryStandards: []string{"SOC"},
			},
		},
	},
	)

	if err != nil {
		resp.Diagnostics.AddError("Creating Wiz Project Failed failed.",
			fmt.Sprintf("Unable to create Wiz Project, got error: %s", err))
		return
	}

	data.ID = client_resp.CreateProject.Project.ID

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}

func (r wizProject) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data wizProjectTypeData
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.provider.wizClient.UpdateWizProject(ctx, apiClient.UpdateProjectRequest{
		Input: apiClient.Input{
			ID: *data.ID,
			Override: apiClient.Override{
				Name: data.Name,
				// TODO we are passing context but it isn't being used, either use it os lose it
				CloudAccountLinks: data.getAccountLinks(ctx),
				// TODO remove hard-coded risk profile until we know what we need to do
				RiskProfile: apiClient.RiskProfile{
					BusinessImpact:      "MBI",
					HasAuthentication:   "UNKNOWN",
					HasExposedAPI:       "YES",
					IsCustomerFacing:    "NO",
					IsInternetFacing:    "YES",
					IsActivelyDeveloped: "UNKNOWN",
					IsRegulated:         "YES",
					SensitiveDataTypes:  []string{"CUSTOMER", "FINANCIAL"},
					StoresData:          "YES",
					RegulatoryStandards: []string{"SOC"},
				},
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Updating Wiz Project Failed failed.",
			fmt.Sprintf("Unable to update Wiz Project, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r wizProject) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data resourceWizProjectType
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}
func (r wizProject) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data wizProjectTypeData
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_resp, err := r.provider.wizClient.GetWizProject(ctx, apiClient.GetProjectRequest{
		First:           1,
		ProjectID:       data.ID,
		FetchTotalCount: true,
		Quick:           false,
		Query: apiClient.Query{
			Type: []string{
				"PROJECT",
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Getting Wiz project failed.",
			fmt.Sprintf("Unable to get Wiz Project, got error: %s", err))
		return
	}

	entities := client_resp.GraphSearch.Nodes[0].Entities[0]

	var subscriptions apiClient.Subscriptions
	var val []byte = []byte(client_resp.GraphSearch.Nodes[0].Entities[0].Properties.Subscriptions)

	marshal_err := json.Unmarshal(val, &subscriptions)

	if marshal_err != nil {
		resp.Diagnostics.AddError("Unmarshalling subscription data failed.",
			fmt.Sprintf("Unable to seraialise stringified subscriptions, got error: %s", err))
		panic(err)
	}

	data.Name = entities.Name
	data.ID = entities.ID
	data.setAccountLinks(ctx, subscriptions)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}

func (r wizProject) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("key"), req, resp)
}
