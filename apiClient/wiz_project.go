package apiClient

import (
	"context"

	"github.com/fatih/structs"
)

// #region Create Request Struct
type CreateProjectRequest struct {
	Input CreateProjectInput `structs:"input"`
}

type CreateProjectInput struct {
	Name                   *string            `structs:"name"`
	Identifiers            []interface{}      `structs:"identifiers"`
	CloudOrganizationLinks []interface{}      `structs:"cloudOrganizationLinks"`
	CloudAccountLinks      []CloudAccountLink `structs:"cloudAccountLinks"`
	RepositoryLinks        []interface{}      `structs:"repositoryLinks"`
	Description            string             `structs:"description"`
	SecurityChampions      []interface{}      `structs:"securityChampions"`
	ProjectOwners          []interface{}      `structs:"projectOwners"`
	BusinessUnit           string             `structs:"businessUnit"`
	RiskProfile            RiskProfile        `structs:"riskProfile"`
}

//#endregion

// #region Create Response Struct

type CreateProjectResponseData struct {
	CreateProject CreateProjectDetails `json:"createProject"`
}

type CreateProjectDetails struct {
	Project Project `json:"project"`
}

type Project struct {
	ID *string `json:"id"`
}

//#endregion

// #region Update Request Struct
type UpdateProjectRequest struct {
	Input Input `structs:"input"`
}

type Input struct {
	ID       string   `structs:"id"`
	Override Override `structs:"override"`
}

type Override struct {
	Name                   *string            `structs:"name"`
	Identifiers            []interface{}      `structs:"identifiers"`
	CloudOrganizationLinks []interface{}      `structs:"cloudOrganizationLinks"`
	CloudAccountLinks      []CloudAccountLink `structs:"cloudAccountLinks"`
	RepositoryLinks        []interface{}      `structs:"repositoryLinks"`
	KubernetesClusterLinks []interface{}      `structs:"kubernetesClusterLinks"`
	Description            string             `structs:"description"`
	SecurityChampions      []interface{}      `structs:"securityChampions"`
	ProjectOwners          []interface{}      `structs:"projectOwners"`
	BusinessUnit           string             `structs:"businessUnit"`
	RiskProfile            RiskProfile        `structs:"riskProfile"`
}

type CloudAccountLink struct {
	CloudAccount   string        `structs:"cloudAccount"`
	Environment    string        `structs:"environment"`
	Shared         bool          `structs:"shared"`
	ResourceTags   []interface{} `structs:"resourceTags"`
	ResourceGroups []interface{} `structs:"resourceGroups"`
}

type RiskProfile struct {
	BusinessImpact      string   `structs:"businessImpact"`
	HasAuthentication   string   `structs:"hasAuthentication"`
	HasExposedAPI       string   `structs:"hasExposedAPI"`
	IsCustomerFacing    string   `structs:"isCustomerFacing"`
	IsInternetFacing    string   `structs:"isInternetFacing"`
	IsActivelyDeveloped string   `structs:"isActivelyDeveloped"`
	IsRegulated         string   `structs:"isRegulated"`
	SensitiveDataTypes  []string `structs:"sensitiveDataTypes"`
	StoresData          string   `structs:"storesData"`
	RegulatoryStandards []string `structs:"regulatoryStandards"`
}

// #endregion

// #region Update Response Struct
type UpdateProjectResponseData struct {
	Data Data `json:"data"`
}

type Data struct {
	UpdateProject UpdateProject `json:"updateProject"`
}

type UpdateProject struct {
	Project Project `json:"project"`
}

type CloudAccount struct {
	ID string `json:"id"`
}

// #region Get Project Request Struct
type GetProjectRequest struct {
	First           int64   `structs:"first"`
	Query           Query   `structs:"query"`
	ProjectID       *string `structs:"projectId"`
	FetchTotalCount bool    `structs:"fetchTotalCount"`
	Quick           bool    `structs:"quick"`
}

type Query struct {
	Type []string `structs:"type"`
}

// #endregion

// #region Get Project Response Struct
type GetProjectResponse struct {
	Data GetProjectResponseData `json:"data"`
}

type GetProjectResponseData struct {
	GraphSearch GraphSearch `json:"graphSearch"`
}

type GraphSearch struct {
	TotalCount      int64    `json:"totalCount"`
	MaxCountReached bool     `json:"maxCountReached"`
	PageInfo        PageInfo `json:"pageInfo"`
	Nodes           []Node   `json:"nodes"`
}

type Node struct {
	Entities []Entity `json:"entities"`
}

type Entity struct {
	ID             *string     `json:"id"`
	Name           *string     `json:"name"`
	Type           string      `json:"type"`
	Properties     Properties  `json:"properties"`
	OriginalObject interface{} `json:"originalObject"`
}

type Properties struct {
	VertexID           string      `json:"_vertexID"`
	BusinessImpact     string      `json:"businessImpact"`
	BusinessUnit       interface{} `json:"businessUnit"`
	Description        interface{} `json:"description"`
	ExternalID         string      `json:"externalId"`
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	ProductCategory    string      `json:"productCategory"`
	ProductSubCategory string      `json:"productSubCategory"`
	Slug               string      `json:"slug"`
	Subscriptions      string      `json:"subscriptions"`
	UpdatedAt          string      `json:"updatedAt"`
}

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type Subscriptions []struct {
	Environments         []string      `json:"environments"`
	SharedAccount        bool          `json:"sharedAccount"`
	SharedResourceGroups []interface{} `json:"sharedResourceGroups"`
	SharedTags           struct {
	} `json:"sharedTags"`
	Status                 string `json:"status"`
	SubscriptionExternalID string `json:"subscriptionExternalId"`
	SubscriptionID         string `json:"subscriptionId"`
}

//#endregion

func (c *Client) CreateWizProject(ctx context.Context, req CreateProjectRequest) (*CreateProjectResponseData, error) {
	create_req := `
	  mutation CreateProject($input: CreateProjectInput!) {
		  createProject(input: $input) {
			project {
			  id
			}
		  }
		}`

	s := structs.New(req)
	request_mapped := s.Map()
	response := &CreateProjectResponseData{}

	c.doRequest(create_req, request_mapped, response)

	return response, nil
}

func (c *Client) UpdateWizProject(ctx context.Context, req UpdateProjectRequest) (*UpdateProjectResponseData, error) {
	update_req := `
		mutation UpdateProject($input: UpdateProjectInput!) {
			updateProject(input: $input) {
			project {
				id
				name
				identifiers
				description
				businessUnit
				projectOwners {
				id
				name
				email
				}
				securityChampions {
				id
				name
				email
				}
				cloudOrganizationLinks {
				cloudOrganization {
					id
				}
				environment
				resourceTags {
					key
					value
				}
				shared
				resourceGroups
				}
				cloudAccountLinks {
				cloudAccount {
					id
				}
				environment
				resourceTags {
					key
					value
				}
				shared
				resourceGroups
				}
				kubernetesClustersLinks {
				kubernetesCluster {
					id
				}
				environment
				namespaces
				shared
				}
				repositoryLinks {
				repository {
					id
				}
				}
				riskProfile {
				businessImpact
				hasAuthentication
				isInternetFacing
				hasExposedAPI
				storesData
				sensitiveDataTypes
				regulatoryStandards
				isCustomerFacing
				isRegulated
				}
			}
			}
		}
	  `

	s := structs.New(req)
	request_mapped := s.Map()
	response := &UpdateProjectResponseData{}

	c.doRequest(update_req, request_mapped, response)

	return response, nil
}

func (c *Client) GetWizProject(ctx context.Context, req GetProjectRequest) (*GetProjectResponseData, error) {
	get_req := `
	query GraphSearch(
		$query: GraphEntityQueryInput
		$controlId: ID
		$projectId: String!
		$first: Int
		$after: String
		$fetchTotalCount: Boolean!
		$quick: Boolean
	) {
		graphSearch(
		query: $query
		controlId: $controlId
		projectId: $projectId
		first: $first
		after: $after
		quick: $quick
		) {
		totalCount @include(if: $fetchTotalCount)
		maxCountReached @include(if: $fetchTotalCount)
		pageInfo {
			endCursor
			hasNextPage
		}
		nodes {
			entities {
			id
			name
			type
			properties
			originalObject
			}
		}
		}
	}
	  `
	s := structs.New(req)
	request_mapped := s.Map()
	response := &GetProjectResponseData{}
	c.doRequest(get_req, request_mapped, response)

	return response, nil
}
