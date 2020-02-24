package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

func initTelegramBot() (int64, *tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
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

	return userid, bot, updates
}

func initDockerSDK() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return cli
}

func sendMessage(bot *tgbotapi.BotAPI, chatid int64, messageid int, message string) {
	msg := tgbotapi.NewMessage(chatid, message)
	msg.ReplyToMessageID = messageid
	bot.Send(msg)
}

func startMessage(bot *tgbotapi.BotAPI, chatid int64, messageid int) {
	welcome := "Docker Telegram Bot is a bot to manage your docker container in your machine!"
	welcome += "If you want to know what you can do, use /help command."
	sendMessage(bot, chatid, messageid, welcome)
}

func sendContainerList(bot *tgbotapi.BotAPI, chatid int64, messageid int, cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	message := "Container List\n"

	for _, container := range containers {
		fmt.Println("ID: ", container.ID)
		message += "ID: " + container.ID + "\n"
	}

	sendMessage(bot, chatid, messageid, message)
}

func main() {
	//Docker SDK init
	//ctx := context.Background()
	cli := initDockerSDK()

	//Telegram Bot init
	userid, bot, updates := initTelegramBot()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		} else if update.Message.Chat.ID == userid {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Text {
			case "/start":
				startMessage(bot, update.Message.Chat.ID, update.Message.MessageID)
			case "/list":
				sendContainerList(bot, update.Message.Chat.ID, update.Message.MessageID, cli)
			}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad eseguire il comando!")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
}
