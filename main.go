package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Config rappresents configuration parameters for telegram bot
type Config struct {
	TOKEN, USERID string
}

func initConfig() Config {
	data, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println("error:", err)
	}

	var config Config

	err = json.Unmarshal(data, &config)

	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("TOKEN: ", config.TOKEN)
	fmt.Println("USERID: ", config.USERID)
	return config
}

func startMessage(bot *tgbotapi.BotAPI, chatid int64, messageid int) {
	welcome := "Docker Telegram Bot is a bot to manage your docker container in your machine!"
	welcome += "If you want to know what you can do, use /help command."

	msg := tgbotapi.NewMessage(chatid, welcome)
	msg.ReplyToMessageID = messageid
	bot.Send(msg)
}

func main() {
	config := initConfig()
	userid, err := strconv.ParseInt(config.USERID, 10, 64)

	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		} else if update.Message.Chat.ID == userid {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "/start" {
				startMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
			} else if update.Message.Text == "/help" {

			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad eseguire il comando!")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
}
