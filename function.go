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

//TelegramConfig rappresents configuration parameters for telegram bot
type TelegramConfig struct {
	TOKEN, USERID string
}

func initConfig() TelegramConfig {
	data, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println("error:", err)
	}

	var config TelegramConfig

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

func replyMessage(bot *tgbotapi.BotAPI, chatid int64, messageid int, message string) {
	msg := tgbotapi.NewMessage(chatid, message)
	msg.ReplyToMessageID = messageid
	bot.Send(msg)
}

func sendMessage(bot *tgbotapi.BotAPI, chatid int64, message string) {
	msg := tgbotapi.NewMessage(chatid, message)
	msg.ParseMode = "markdown"
	bot.Send(msg)
}

func startMessage(bot *tgbotapi.BotAPI, chatid int64) {
	welcome := "Docker Telegram Bot is a bot to manage your docker container in your machine!"
	welcome += "If you want to know what you can do, use /help command."
	sendMessage(bot, chatid, welcome)
}

func sendContainersList(bot *tgbotapi.BotAPI, chatid int64, cli *client.Client) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	message := "*Container List*\n"

	for _, container := range containers {
		message += "*Name*: " + container.Names[0][1:]
		message += " *Image*: " + container.Image
		message += " *Status*: " + container.State
		message += "\n"
	}

	sendMessage(bot, chatid, message)
}

func stopAllContainer(ctx context.Context, bot *tgbotapi.BotAPI, chatid int64, cli *client.Client) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		message := "Container " + container.ID[:10] + " stopped."
		sendMessage(bot, chatid, message)
	}
}

func sendImagesList(bot *tgbotapi.BotAPI, chatid int64, cli *client.Client) {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	message := "List Images\n"
	for _, image := range images {
		message += image.ID + "\n"
	}

	sendMessage(bot, chatid, message)
}

func stopContainer(bot *tgbotapi.BotAPI, chatid int64, cli *client.Client) {
	var row []tgbotapi.KeyboardButton
	var msg tgbotapi.MessageConfig

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		row = append(row, tgbotapi.NewKeyboardButton(container.ID[:10]))
	}

	if row == nil {
		msg = tgbotapi.NewMessage(chatid, "No container available.")
	} else {
		var Keyboard = tgbotapi.NewReplyKeyboard(
			row,
		)
		msg = tgbotapi.NewMessage(chatid, "Select a container")
		msg.ReplyMarkup = Keyboard
	}

	bot.Send(msg)
}
