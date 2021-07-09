package healthcheck

import (
	"context"
	"fmt"
	"github.com/healthcheck-exporter/cmd/authentication"
	"github.com/healthcheck-exporter/cmd/exporter"
	"github.com/healthcheck-exporter/cmd/model"
	"github.com/healthcheck-exporter/cmd/watchdog"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HealthCheck struct {
	config     *model.Config
	authClient *authentication.AuthClient
	status     *model.Status
	wsClient   *WsClient
	exporter   *exporter.Exporter
	watchDog   *watchdog.WatchDog
	httpClient *http.Client
}

func NewHealthCheck(config *model.Config, authClient *authentication.AuthClient, ex *exporter.Exporter, wd *watchdog.WatchDog) *HealthCheck {
	hc := HealthCheck{
		config:     config,
		authClient: authClient,
		status:     &model.Status{
			Tasks: make(map[string]*model.Task),
		},
		wsClient:   NewWsClient(ex),
		exporter:   ex,
		watchDog:   wd,
		httpClient: &http.Client{},
	}

	hc.InitTasks()

	return &hc
}

func (hc *HealthCheck) InitTasks() {
	for i := 0; i < len(hc.config.Jobs); i++ {
		hc.InitTask(&hc.config.Jobs[i])
	}

	for i := 0; i < len(hc.config.Jobs); i++ {
		go hc.StartTask(&hc.config.Jobs[i])
	}
}

func (hc *HealthCheck) getTask(taskId string) *model.Task {
	if hc.status.Tasks[taskId] == nil {
		task := model.Task{
			Id:            taskId,
			Online:        false,
			SuccessChecks: 0,
			FailureChecks: 0,
			RestartTime:   0,
		}
		hc.status.Tasks[taskId] = &task
	}
	return hc.status.Tasks[taskId]
}

func (hc *HealthCheck) StartTask(function *model.Job) {
	counter := 0
	task := hc.getTask(function.Id)
	for true {
		counter++
		active := false
		if function.DependentJob != "" {
			if hc.getTask(function.DependentJob).Online {
				active = true
			}
		} else {
			active = true
		}

		if active {
			if hc.check(function) {
				hc.exporter.SetCounter(function.Id, 0)
				if task.Online {
					log.Info(fmt.Sprintf("%s: Task status updated (is online?): %t",
						function.Id, hc.status.Tasks[function.Id].Online))
				}

				task.Online = true
				task.SuccessChecks++
				task.FailureChecks = 0

				log.Debug(fmt.Sprintf("%s: Task status updated (is online?): %t",
					function.Id, hc.status.Tasks[function.Id].Online))
			} else {
				hc.exporter.AddCounter(function.Id, function.Timeout)

				task.Online = false
				task.FailureChecks++
				log.Info(fmt.Sprintf("%s: Task status updated (is online?): %t, count: %d",
					function.Id, hc.status.Tasks[function.Id].Online, hc.status.Tasks[function.Id].FailureChecks))

				if function.WatchDog.Enabled &&
					task.FailureChecks >= function.WatchDog.FailureThreshold &&
					(time.Now().Unix()-task.RestartTime) > function.WatchDog.AwaitAfterRestart {

					for y := 0; y < len(function.WatchDog.Deployments); y++ {
						err := hc.watchDog.DeletePod(function.WatchDog.Deployments[y], function.WatchDog.Namespace)
						if err != nil {
							log.Error(fmt.Sprintf("Delete pod error: %s", err.Error()))
						}
					}

					task.FailureChecks = 0
					task.RestartTime = time.Now().Unix()
				}
			}
		}

		duration := time.Duration(function.Timeout) * time.Second
		time.Sleep(duration)

	}
}

func (hc *HealthCheck) InitTask(function *model.Job) {
	task := hc.getTask(function.Id)
	log.Info(fmt.Sprintf("%s: Started task", task.Id))

	if function.Location.Type == "kubernetes" {
		podIps, err := hc.watchDog.GetPodIp(function.Location.Deployment, function.Location.Namespace)
		if err != nil {
			log.Error(fmt.Sprintf("%s: error wss last message exceeded timeout", function.Id))
			return
		}

		urls := make([]string, 0)
		for _, u := range function.Urls {
			base, _ := url.Parse(u)

			for _, ip := range podIps {
				base.Host = fmt.Sprintf("%s:%s", ip, function.Location.Port)
				base.Scheme = "http"
				newurl := base.String()
				//newurl := fmt.Sprintf("%s://%s:%s%s", "http", ip, function.Location.Port, base.Path)
				urls = append(urls, newurl)
			}
		}

		fmt.Println(urls)
	}
}

func (hc *HealthCheck) check(function *model.Job) bool {
	switch function.Type {
	case "http_get":
		return hc.checkHttpGet(function)
	case "http_post":
		return hc.checkHttpPost(function)
	case "websocket":
		return hc.checkWs(function)
	}
	return false
}

func (hc *HealthCheck) checkWs(function *model.Job) bool {
	for _, u := range function.Urls {
		difference := hc.wsClient.TimeDifferenceWithLastMessage(function.Id, u)

		if difference > function.Timeout {
			log.Error(fmt.Sprintf("%s: error wss last message exceeded timeout", function.Id))
			return false
		}

		return true
	}

	return true
}

func (hc *HealthCheck) getHttpClient(function *model.Job) *http.Client {
	if function.AuthEnabled {
		return hc.authClient.GetClient()
	} else {
		return hc.httpClient
	}
}

func (hc *HealthCheck) checkHttpGet(function *model.Job) bool {
	start := time.Now()

	for _, u := range function.Urls {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return false
		}

		if function.ResponseTimeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(function.ResponseTimeout)*time.Second)
			req = req.WithContext(ctx)
			defer cancel()
		}

		resp, err := hc.getHttpClient(function).Do(req)
		if err != nil {
			log.Error(fmt.Sprintf("Error http get request: %s", err.Error()))
			return false
		}
		if resp == nil || resp.StatusCode != 200 {
			log.Error(fmt.Sprintf("%s: Empty http get result or invalid response code", function.Id))
			return false
		}
		defer resp.Body.Close()
	}


	hc.exporter.SetGauge(function.Id, float64(time.Since(start).Milliseconds()))
	log.Info(fmt.Sprintf("%s %s", function.Id, time.Since(start)))

	return true
}

func (hc *HealthCheck) checkHttpPost(function *model.Job) bool {
	for _, u := range function.Urls {
		req, err := http.NewRequest("POST", u, strings.NewReader(function.Body))
		if err != nil {
			return false
		}

		req.Header.Add("accept", "*/*")
		req.Header.Add("Content-Type", "application/json")

		if function.ResponseTimeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(function.ResponseTimeout)*time.Second)
			req = req.WithContext(ctx)
			defer cancel()
		}

		resp, err := hc.getHttpClient(function).Do(req)
		if err != nil {
			log.Error(fmt.Sprintf("Error http post request: %s", err.Error()))
			return false
		}
		if resp == nil || resp.StatusCode != 200 {
			log.Error(fmt.Sprintf("Empty http post result or invalid response code"))
			return false
		}
		defer resp.Body.Close()
	}

	return true
}

func (hc *HealthCheck) Status() *model.Status {
	return hc.status
}
