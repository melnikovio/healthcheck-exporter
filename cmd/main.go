package main

import (
	"encoding/json"
	"fmt"
	"github.com/healthcheck-exporter/cmd/api"
	"github.com/healthcheck-exporter/cmd/authentication"
	"github.com/healthcheck-exporter/cmd/exporter"
	"github.com/healthcheck-exporter/cmd/healthcheck"
	"github.com/healthcheck-exporter/cmd/model"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func main() {
	var config *model.Config
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Error("Couldn't load configuration")
		panic(err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil || config == nil {
		log.Error("Couldn't parse configuration")
		panic(err)
	}

	authClient := authentication.NewAuthClient(config)

	ex := exporter.NewExporter(config)

	hcClient := healthcheck.NewHealthCheck(config, authClient, ex)
	//
	//http.Handle("/metrics", promhttp.Handler())
	//
	//http.Handle("/probe", promhttp.Handler())


	// initialize api
	router := api.NewRouter(hcClient)

	// enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET"},
	})

	log.Info(fmt.Sprintf(http.ListenAndServe(":2112",
		corsHandler.Handler(router)).Error()))



	//log.Info(fmt.Sprintf(http.ListenAndServe(":2112", nil).Error()))



	//
	//
	//err = http.ListenAndServe(":2112", nil)
	//if err != nil {
	//	panic(err)
	//}
}
