package ydatabase

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestYDataBase(t *testing.T) {
	ts := time.Now().Add(time.Duration(60) * time.Minute)

	m := [1]YDBTableMessage{}
	m[0] = YDBTableMessage{700, [16]uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, ts}

	base := Connect(os.Getenv("YDB_CONNECTION_STRING"))

	defer base.cancel()
	defer func() { _ = base.db.Close(base.ctx) }()

	start := time.Now()
	err := base.CreateDatabase("TEST_BASE")
	end := time.Since(start)
	fmt.Println("Time for create database: ", end)
	if err != nil {
		panic(err)
	}

	start = time.Now()
	err = base.WriteToDatabase("TEST_BASE", m)
	end = time.Since(start)
	fmt.Println("Time for write to database: ", end)
	if err != nil {
		panic(err)
	}

	data := YDBTableMessage{}
	start = time.Now()
	base.ReadFromTable("TEST_BASE", 700, &data)
	end = time.Since(start)
	fmt.Println("Time for read from database:", end)

	if data.Msg_id != m[0].Msg_id {
		t.Fatalf("Read msg_id != write msg_id: %d != %d", data.Msg_id, m[0].Msg_id)
	}

	for i, val := range data.Board {
		if val != m[0].Board[i] {
			t.Fatalf("Read value != write value: %d != %d", val, m[0].Board[i])
		}
	}

}
