package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/healthcheck-exporter/cmd/model"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	config *model.Config
	Bot    *tgbotapi.BotAPI
}

func NewBot(config *model.Config) *Bot {

	////
	bot, err := tgbotapi.NewBotAPI("1527699463:AAGUIIrQA_d99AGx43c84sKHBWqo9wJY4mU")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	botClient := Bot{
		config: config,
		Bot:    bot,
	}

	go botClient.Start()

	return &botClient
}

func (bot *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//
		//bot.Bot.Send(msg)
	}
	/////
}
