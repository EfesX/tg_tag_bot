package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func send_msg(bot tg.BotAPI, msg tg.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

type ReqBody struct {
	UpdateID string `json:"update_id"`
	Message  string `json:"message"`
}

func HandlerUpdate(bot *tg.BotAPI, update *tg.Update) {

	game := Game{}
	game.Connect()

	if update.Message != nil {
		switch update.Message.Text {
		case "/start":
			msg := tg.NewMessage(update.Message.Chat.ID, "Сыграем в пятнашки?")
			msg.ReplyMarkup = game.CreateQuestionNewGame()
			send_msg(*bot, msg)

			msg = tg.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyMarkup = tg.NewReplyKeyboard(
				tg.NewKeyboardButtonRow(
					tg.NewKeyboardButton("Новая игра"),
				),
				tg.NewKeyboardButtonRow(
					tg.NewKeyboardButton("О боте"),
					tg.NewKeyboardButton("Правила"),
				),
			)
			send_msg(*bot, msg)

		case "Новая игра":
			msg := tg.NewMessage(update.Message.Chat.ID, "Сыграем в пятнашки?")
			msg.ReplyMarkup = game.CreateQuestionNewGame()
			send_msg(*bot, msg)

		case "О боте":
			msg := tg.NewMessage(update.Message.Chat.ID, "Бот для игры в пятнашки\nРазработчик: @efesxzc")
			send_msg(*bot, msg)

		case "Правила":
			text := "Для победы необходимо расположить по порядку в четыре строки числа от 1 до 15"
			msg := tg.NewMessage(update.Message.Chat.ID, text)
			send_msg(*bot, msg)
		}

	} else if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "newgame":
			game.NewGame(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			game.Save(uint32(update.CallbackQuery.Message.Chat.ID), uint32(update.CallbackQuery.Message.MessageID))

			msg := tg.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "================= Играем =================")
			msg.ReplyMarkup = game.CreateBoard(game.board.board)

			send_msg(*bot, msg)

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15":
			callback := tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			game.Restore(uint32(update.CallbackQuery.Message.Chat.ID), uint32(update.CallbackQuery.Message.MessageID))

			num, err := strconv.ParseInt(update.CallbackQuery.Data, 10, 8)
			if err != nil {
				panic(err)
			}
			if game.Update(uint8(num)) {
				if game.IsWin() {
					msg := tg.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "================= Победа! =================\nЕще разок?")
					msg.ReplyMarkup = game.CreateQuestionNewGame()
					send_msg(*bot, msg)
				} else {
					game.Save(uint32(update.CallbackQuery.Message.Chat.ID), uint32(update.CallbackQuery.Message.MessageID))
					msg := tg.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "================= Играем =================")
					msg.ReplyMarkup = game.CreateBoard(game.board.board)
					send_msg(*bot, msg)
				}
			}
		}
	}
}

func Handler(rw http.ResponseWriter, req *http.Request) {
	bot, err := tg.NewBotAPI(os.Getenv("TG_TAG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	update, uerr := bot.HandleUpdate(req)
	if uerr != nil {
		log.Panic(uerr)
	}

	HandlerUpdate(bot, update)

	rw.Header().Set("X-Custom-Header", "Test")
	rw.WriteHeader(200)
	io.WriteString(rw, fmt.Sprintf("OK"))
}

func main() {

	bot, err := tg.NewBotAPI(os.Getenv("TG_TAG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		HandlerUpdate(bot, &update)
	}
}
