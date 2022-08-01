package main

import (
	"app/ydatabase"
	"fmt"
	"os"
)

type Game struct {
	MarkUpGameBoard
	database   ydatabase.YandexDatabase
	chat_id    int64
	msg_id     int
	step_count int
	text       string
	board      Board
	markup     MarkUpGameBoard
}

func (g *Game) Update(val uint8) bool {
	g.step_count = g.step_count + 1
	return g.board.step(val)
}

func (g *Game) NewGame(chat_id int64, msg_id int) bool {

	g.chat_id = chat_id
	g.msg_id = msg_id

	g.board.create()
	g.markup.CreateBoard(g.board.board)

	return true
}

func (g *Game) IsWin() bool {
	if g.board.board == [16]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15} {
		return true
	}
	return false
}

func (g *Game) Save(chat_id uint32, msg_id uint32) error {
	//db := ydatabase.Connect(os.Getenv("YDB_CONNECTION_STRING"))

	msg := [1]ydatabase.YDBTableMessage{}

	//err := db.CreateDatabase(fmt.Sprint(chat_id))
	err := g.database.CreateDatabase(fmt.Sprint(chat_id))
	if err != nil {
		panic(err)
	}

	msg[0].Msg_id = msg_id
	for i, v := range g.board.board {
		msg[0].Board[i] = uint8(v)
	}

	//err = db.WriteToDatabase(fmt.Sprint(chat_id), msg)
	err = g.database.WriteToDatabase(fmt.Sprint(chat_id), msg)
	if err != nil {
		panic(err)
	}
	return nil
}

func int_to_uint8(intval *[]int, uint8val *[16]uint8) {
	for i, v := range *intval {
		uint8val[i] = uint8(v)
	}
}

func (g *Game) Restore(chat_id, msg_id uint32) error {
	//db := ydatabase.Connect(os.Getenv("YDB_CONNECTION_STRING"))

	data := ydatabase.YDBTableMessage{}
	//err := db.ReadFromTable(fmt.Sprint(chat_id), msg_id, &data)
	err := g.database.ReadFromTable(fmt.Sprint(chat_id), msg_id, &data)

	if err != nil {
		panic(err)
	}

	for i, v := range data.Board {
		g.board.board[i] = v
	}

	return nil
}

func (g *Game) Connect() {
	g.database = ydatabase.Connect(os.Getenv("YDB_CONNECTION_STRING"))
}
