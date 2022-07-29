package main

/*
import (
	"log"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	bot, err := tg.NewBotAPI("5482440199:AAF9F5NyBaLncK_uIW1Yrdjxv0rKBZNb0Cw")
	if err != nil {
		log.Panic(err)
	}

	game := new(Game)
	bot.Debug = true

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {

			switch update.Message.Text {
			case "/start":
				msg := tg.NewMessage(update.Message.Chat.ID, "Сыграем в пятнашки?")
				msg.ReplyMarkup = game.markup.CreateQuestionNewGame()

				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			}

		} else if update.CallbackQuery != nil {

			callback := tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			if update.CallbackQuery.Data == "newgame" {

				game.NewGame(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				game.text = "============== Идет игра =============="
				msg := tg.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, game.text)
				msg.ReplyMarkup = game.markup.CreateBoard(game.board.board)

				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}

			} else {
				v, err := strconv.Atoi(update.CallbackQuery.Data)

				if err == nil {
					if game.Update(v) {
						msg := tg.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, game.text)
						msg.ReplyMarkup = game.markup.CreateBoard(game.board.board)

						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
					}
					if game.isWin() {
						msg := tg.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Победа!\nНачать заново?")
						msg.ReplyMarkup = game.markup.CreateQuestionNewGame()

						if _, err = bot.Send(msg); err != nil {
							panic(err)
						}
					} else {

					}
				}
			}
		}
	}
}
*/
