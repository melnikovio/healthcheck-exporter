package main

import (
	"encoding/json"
	"fmt"
	"github.com/healthcheck-exporter/cmd/api"
	"github.com/healthcheck-exporter/cmd/authentication"
	"github.com/healthcheck-exporter/cmd/bot"
	"github.com/healthcheck-exporter/cmd/healthcheck"
	"github.com/healthcheck-exporter/cmd/model"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func main() {
	log.Info("Starting service")

	var config *model.Config
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil || config == nil {
		panic(err)
	}

	log.Info(fmt.Sprintf("Config file is: %s", config))

	authClient := authentication.NewAuthClient(config)

	botClient := bot.NewBot(config)

	hcClient := healthcheck.NewHealthCheck(config, authClient, botClient)

	// initialize api
	router := api.NewRouter(hcClient)

	// enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	log.Info(fmt.Sprintf(http.ListenAndServe(":8080",
		corsHandler.Handler(router)).Error()))
}
