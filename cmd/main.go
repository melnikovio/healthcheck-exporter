package main

import (
	"fmt"
	"github.com/healthcheck-exporter/cmd/api"
	"github.com/healthcheck-exporter/cmd/authentication"
	"github.com/healthcheck-exporter/cmd/common"
	"github.com/healthcheck-exporter/cmd/configuration"
	"github.com/healthcheck-exporter/cmd/exporter"
	"github.com/healthcheck-exporter/cmd/healthcheck"
	"github.com/healthcheck-exporter/cmd/watchdog"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	fmt.Println(common.Logo)

	config := configuration.GetConfiguration()
	authClient := authentication.NewAuthClient(config)

	ex := exporter.NewExporter(config)

	wd := watchdog.NewWatchDog()

	hcClient := healthcheck.NewHealthCheck(config, authClient, ex, wd)
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
