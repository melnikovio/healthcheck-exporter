package main

import (
	"flag"
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
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")


func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
		}(f) // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

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

	//if *memprofile != "" {
	//	f, err := os.Create(*memprofile)
	//	if err != nil {
	//		log.Fatal("could not create memory profile: ", err)
	//	}
	//	defer f.Close() // error handling omitted for example
	//	runtime.GC() // get up-to-date statistics
	//	if err := pprof.WriteHeapProfile(f); err != nil {
	//		log.Fatal("could not write memory profile: ", err)
	//	}
	//}

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
