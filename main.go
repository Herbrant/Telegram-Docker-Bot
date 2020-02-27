package main

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	//Docker SDK init
	ctx := context.Background()
	cli := initDockerSDK()

	//Telegram Bot init
	userid, bot, updates := initTelegramBot()

	for update := range updates {
		if update.CallbackQuery != nil {
			fmt.Print(update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message != nil {
			if update.Message.Chat.ID == userid {
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
				case "/stop":
					stopContainer(bot, update.Message.Chat.ID, cli)
				}

			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Non sei autorizzato ad eseguire il comando!")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}
