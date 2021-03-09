package authentication

import (
	"context"
	"github.com/healthcheck-exporter/cmd/model"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

type AuthClient struct {
	client *http.Client
}

func NewAuthClient(config *model.Config) *AuthClient {
	oauth := &clientcredentials.Config{
		ClientID:     config.Authentication.ClientId,
		ClientSecret: config.Authentication.ClientSecret,
		TokenURL:     config.Authentication.AuthUrl + "/protocol/openid-connect/token",
	}

	ctx := context.Background()
	client := oauth.Client(ctx)

	authClient := AuthClient{
		client: client,
	}

	return &authClient
}

func (ac *AuthClient) GetClient() *http.Client {
	return ac.client
}
