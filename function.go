package main

import (
	"context"
	"log"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

	message := "Container List\n"

	for _, container := range containers {
		message += "ID: " + container.ID[:10] + "\n"
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
