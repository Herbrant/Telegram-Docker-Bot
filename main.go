package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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

func main() {
	//Docker SDK init
	ctx := context.Background()
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
				startMessage(bot, update.Message.Chat.ID)
			case "/listcontainer":
				sendContainersList(bot, update.Message.Chat.ID, cli)
			case "/stopall":
				stopAllContainer(ctx, bot, update.Message.Chat.ID, cli)
			case "/listimage":
				sendImagesList(bot, update.Message.Chat.ID, cli)
			}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad eseguire il comando!")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
}
