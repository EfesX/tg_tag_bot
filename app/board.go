package main

import (
	"math/rand"
	"time"
)

type Board struct {
	board []int
}

func (b *Board) create() bool {
	rand.Seed(time.Now().UnixNano())
	b.board = rand.Perm(16)
	return true
}

func (b *Board) swap(i, j int) bool {
	temp := b.board[i]
	b.board[i] = b.board[j]
	b.board[j] = temp
	return true
}

func (b *Board) step(val int) bool {
	for i, v := range b.board {
		if v == val {
			for j, z := range b.board {
				if z == 0 {
					if (i % 4) != 0 {
						if (i - 1) == j {
							b.swap(i, j)
							return true
						}
					}
					if ((i + 1) % 4) != 0 {
						if (i + 1) == j {
							b.swap(i, j)
							return true
						}
					}
					if (i - 4) >= 0 {
						if (i - 4) == j {
							b.swap(i, j)
							return true
						}
					}
					if (i + 4) < 16 {
						if (i + 4) == j {
							b.swap(i, j)
							return true
						}
					}
				}
			}
		}
	}
	return false
}
