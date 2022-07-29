package main

type Game struct {
	YandexDatabase
	chat_id    int64
	msg_id     int
	step_count int
	text       string
	board      Board
	markup     MarkUpGameBoard
}

func (g *Game) Update(val int) bool {
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

func (g *Game) isWin() bool {
	for i, v := range g.board.board {
		if i < 16 {
			if (i + 1) != v {
				return false
			}
		}
	}
	return true
}
