package main

import (
	"math/rand"
	"time"
)

type Board struct {
	board [16]uint8
}

func (b *Board) create() bool {
	rand.Seed(time.Now().UnixNano())
	temp := rand.Perm(16)
	for i, v := range temp {
		b.board[i] = uint8(v)
	}
	return true
}

func (b *Board) swap(i, j uint8) bool {
	temp := b.board[i]
	b.board[i] = b.board[j]
	b.board[j] = temp
	return true
}

func (b *Board) step(val uint8) bool {
	for i, v := range b.board {
		if v == val {
			for j, z := range b.board {
				if z == 0 {
					if (i % 4) != 0 {
						if (i - 1) == j {
							b.swap(uint8(i), uint8(j))
							return true
						}
					}
					if ((i + 1) % 4) != 0 {
						if (i + 1) == j {
							b.swap(uint8(i), uint8(j))
							return true
						}
					}
					if (i - 4) >= 0 {
						if (i - 4) == j {
							b.swap(uint8(i), uint8(j))
							return true
						}
					}
					if (i + 4) < 16 {
						if (i + 4) == j {
							b.swap(uint8(i), uint8(j))
							return true
						}
					}
				}
			}
		}
	}
	return false
}
