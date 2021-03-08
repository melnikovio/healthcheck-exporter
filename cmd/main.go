package main

import (
	"encoding/json"
	"fmt"
	"github.com/healthcheck-exporter/cmd/authentication"
	"github.com/healthcheck-exporter/cmd/exporter"
	"github.com/healthcheck-exporter/cmd/healthcheck"
	"github.com/healthcheck-exporter/cmd/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	ex := exporter.NewExporter(config)

	healthcheck.NewHealthCheck(config, authClient, ex)

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":2112", nil)
	if err != nil {
		panic(err)
	}
}
