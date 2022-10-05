package apiClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/machinebox/graphql"
	errorsHandler "shell.com/terraform-provider-wiz/errors"
)

type Client struct {
	AccessToken string
	Graphql     *graphql.Client
}

func UnmarshalTokenResponse(data []byte) (TokenResponse, error) {
	var r TokenResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func AuthLoginRequest(credentials ClientCredentials) *http.Request {
	apiUrl := "https://auth.wiz.io"
	resource := "/oauth/token"
	data := url.Values{}

	data.Set("grant_type", "client_credentials")
	data.Set("client_id", credentials.ClientID)
	data.Set("client_secret", credentials.ClientSecret)
	data.Set("audience", "beyond-api")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	return r

}

func DoLogin(authLoginRequest *http.Request) string {
	client := &http.Client{}
	authLoginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(authLoginRequest)

	if err != nil {
		panic(fmt.Sprintf("Ran into with initial authentication: %s", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	tokenResponse, err := UnmarshalTokenResponse(body)
	if err != nil {
		panic(fmt.Sprintf("Ran into when unmarshalling token response: %s", err))
	}

	return tokenResponse.AccessToken

}

func CreateClient(config ClientConfig) (*Client, error) {

	credentials, err := GetCredentials(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials, error: %s", err.Error())
	}

	authLoginRequest := AuthLoginRequest(credentials)

	return &Client{
		AccessToken: DoLogin(authLoginRequest),
		Graphql:     graphql.NewClient(credentials.Endpoint),
	}, nil
}

func GetCredentials(config ClientConfig) (ClientCredentials, error) {
	credentials := config.Credentials

	return credentials, nil
}

func (client *Client) doRequest(query string, vars map[string]interface{}, responseData interface{}) error {
	req := graphql.NewRequest(query)

	if vars != nil {
		for k, v := range vars {
			req.Var(k, v)
		}
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+client.AccessToken)

	// define a Context for the request
	ctx := context.Background()

	// run and capture the response
	if err := client.Graphql.Run(ctx, req, &responseData); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return err
	}
	return nil
}

func (client *Client) handleCreateError(err error, input map[string]interface{}, resourceType string) error {
	parent := input["parent"]
	if errorsHandler.NotFoundError(err) {
		return errors.New(fmt.Sprintf("error creating %s: parent resource not found: %s", resourceType, parent))
	}
	return errors.New(fmt.Sprintf("error creating %s: %s ", resourceType, err.Error()))
}

func (client *Client) handleReadError(err error, resource string, resourceType string) error {
	if errorsHandler.NotFoundError(err) {
		return errors.New(fmt.Sprintf("error reading %s: resource not found: %s", resourceType, resource))
	}
	return errors.New(fmt.Sprintf("error reading %s: %s ", resourceType, err.Error()))
}

func (client *Client) handleUpdateError(err error, input map[string]interface{}, resourceType string) error {
	resource := input["id"]
	if errorsHandler.NotFoundError(err) {
		return errors.New(fmt.Sprintf("error updating %s: resource not found: %s", resourceType, resource))
	}
	return errors.New(fmt.Sprintf("error updating %s: %s ", resourceType, err.Error()))
}
