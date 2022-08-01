package main

import (
	"fmt"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	game := Game{}
	game.Connect()

	start := time.Now()
	err := game.Restore(590579622, uint32(210))
	end := time.Since(start)
	fmt.Println("Time for Restore Game: ", end)
	if err != nil {
		panic(err)
	}

	start = time.Now()
	game.Update(10)
	end = time.Since(start)
	fmt.Println("Time for Game Step: ", end)

	start = time.Now()
	game.Save(590579622, 0)
	end = time.Since(start)
	fmt.Println("Time for Game Save: ", end)

}
