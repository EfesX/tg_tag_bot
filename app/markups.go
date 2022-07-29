package main

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type GameMarkUp interface {
}

type MarkUpGameBoard struct {
	markup tg.InlineKeyboardMarkup
	btns   [16]tg.InlineKeyboardButton
}

func (m *MarkUpGameBoard) CreateBoard(b []int) *tg.InlineKeyboardMarkup {
	for i, v := range b {
		if v == 0 {
			m.btns[i] = tg.NewInlineKeyboardButtonData(" ", "0")
		} else {
			m.btns[i] = tg.NewInlineKeyboardButtonData(fmt.Sprint(v), fmt.Sprint(v))
		}
	}

	m.markup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			m.btns[0], m.btns[1], m.btns[2], m.btns[3],
		),
		tg.NewInlineKeyboardRow(
			m.btns[4], m.btns[5], m.btns[6], m.btns[7],
		),
		tg.NewInlineKeyboardRow(
			m.btns[8], m.btns[9], m.btns[10], m.btns[11],
		),
		tg.NewInlineKeyboardRow(
			m.btns[12], m.btns[13], m.btns[14], m.btns[15],
		),
	)

	return &m.markup
}

func (m *MarkUpGameBoard) CreateQuestionNewGame() *tg.InlineKeyboardMarkup {
	m.markup = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Погнали!", "newgame"),
		),
	)
	return &m.markup
}
