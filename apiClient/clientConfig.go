package apiClient

type ClientConfig struct {
	Credentials ClientCredentials
}

type ClientCredentials struct {
	ClientID     string
	ClientSecret string
	Endpoint     string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}
