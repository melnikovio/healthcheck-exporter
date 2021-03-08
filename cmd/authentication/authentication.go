package authentication

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v8"
	"github.com/healthcheck-exporter/cmd/model"
	log "github.com/sirupsen/logrus"
	"time"
)

type AuthClient struct {
	config *model.Config
	client gocloak.GoCloak
	token  *gocloak.JWT
}

func NewAuthClient(config *model.Config) *AuthClient {
	authClient := AuthClient{
		config: config,
	}
	authClient.GetToken()

	return &authClient
}

func (authClient *AuthClient) GetToken() *gocloak.JWT {
	var err error
	if authClient.token == nil {
		authClient.client = gocloak.NewClient(authClient.config.Authentication.AuthUrl)
		ctx := context.Background()
		authClient.token, err = authClient.client.LoginClient(ctx,
			authClient.config.Authentication.ClientId, authClient.config.Authentication.ClientSecret, authClient.config.Authentication.Realm)
		if err != nil {
			log.Error(fmt.Sprintf("Error receiving access token: %s", err.Error()))
		}
		log.Info("Token is received")

		go authClient.RefreshToken()
	}

	return authClient.token
}

func (authClient *AuthClient) RefreshToken() {
	var err error
	for true {
		duration := time.Duration(authClient.token.ExpiresIn-10) * time.Second
		log.Info(fmt.Sprintf("Token will be refreshed in %d seconds", authClient.token.ExpiresIn))
		time.Sleep(duration)

		ctx := context.Background()
		authClient.token, err = authClient.client.RefreshToken(ctx, authClient.token.RefreshToken,
			authClient.config.Authentication.ClientId, authClient.config.Authentication.ClientSecret, authClient.config.Authentication.Realm)
		if err != nil {
			log.Error(fmt.Sprintf("Error refreshing access token: %s", err.Error()))
		}

		log.Info(fmt.Sprintf("Token is refreshed"))
	}
}
