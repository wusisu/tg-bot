package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	os.MkdirAll(viper.GetString("DownloadDir"), os.ModePerm)

	bot, err := tgbotapi.NewBotAPI(viper.GetString("BotToken"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	app := NewApp()
	app.bot = bot

	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	app.db = db

	log.Debugf("Authorized on account %s", bot.Self.UserName)

	app.readUpdates()
}
