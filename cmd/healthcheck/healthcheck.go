package healthcheck

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/healthcheck-exporter/cmd/authentication"
	"github.com/healthcheck-exporter/cmd/bot"
	"github.com/healthcheck-exporter/cmd/model"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type HealthCheck struct {
	config     *model.Config
	authClient *authentication.AuthClient
	status     *model.Status
	wsClient   *WsClient
	bot        *bot.Bot
}

func NewHealthCheck(config *model.Config, authClient *authentication.AuthClient, bot *bot.Bot) *HealthCheck {
	hc := HealthCheck{
		config:     config,
		authClient: authClient,
		status:     &model.Status{},
		wsClient:   NewWsClient(),
		bot:        bot,
	}

	hc.InitTasks()

	return &hc
}

func (hc *HealthCheck) InitTasks() {
	for _, function := range hc.config.Functions {
		go hc.InitTask(function)
	}
}

func (hc *HealthCheck) InitTask(function model.Function) {
	log.Info(fmt.Sprintf("Started task with Id: %s", function.Id))
	hc.status.Task = append(hc.status.Task, model.Task{
		Id:            function.Id,
		Status:        "Init",
		SuccessChecks: 0,
		FailureChecks: 0,
	})

	for true {
		if hc.check(&function) {
			for i := 0; i < len(hc.status.Task); i++ {
				if hc.status.Task[i].Id == function.Id {
					hc.status.Task[i].Status = "Online"
					hc.status.Task[i].SuccessChecks++
				}
			}
		} else {
			for i := 0; i < len(hc.status.Task); i++ {
				if hc.status.Task[i].Id == function.Id {
					hc.status.Task[i].Status = "Failure"
					hc.status.Task[i].FailureChecks++
				}
			}
			msg := tgbotapi.NewMessage(47352032, fmt.Sprintf("Function on NLMK failed: %s", function.Id))
			_, err := hc.bot.Bot.Send(msg)
			if err != nil {
				log.Error(fmt.Sprintf("Error while send message: ", err.Error()))
			}
		}

		duration := time.Duration(function.Timeout) * time.Second
		time.Sleep(duration)

		log.Info(fmt.Sprintf("Updated task with Id: %s", function.Id))
	}
}

func (hc *HealthCheck) check(function *model.Function) bool {
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

func (hc *HealthCheck) checkWs(function *model.Function) bool {
	for _, url := range function.Urls {
		difference := hc.wsClient.DifferenceLastMessageTime(url)

		if difference > function.Timeout {
			return false
		}

		return true
	}

	return true
}

func (hc *HealthCheck) checkHttpGet(function *model.Function) bool {
	for _, url := range function.Urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return false
		}
		client := &http.Client{}
		if function.AuthEnabled {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", hc.authClient.GetToken().AccessToken))
		}

		resp, err := client.Do(req)
		if err != nil || resp == nil || resp.StatusCode != 200 {
			return false
		}
	}

	return true
}

func (hc *HealthCheck) checkHttpPost(function *model.Function) bool {
	for _, url := range function.Urls {
		req, err := http.NewRequest("POST", url, strings.NewReader(function.Body))
		if err != nil {
			return false
		}
		client := &http.Client{}
		if function.AuthEnabled {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", hc.authClient.GetToken().AccessToken))
		}
		req.Header.Add("accept", "*/*")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil || resp == nil || resp.StatusCode != 200 {
			return false
		}
	}

	return true
}

func (hc *HealthCheck) Status() *model.Status {
	return hc.status
}
