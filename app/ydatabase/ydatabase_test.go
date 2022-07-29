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

	base := NewYandexDatabase(os.Getenv("YDB_CONNECTION_STRING"))

	defer base.cancel()
	defer func() { _ = base.db.Close(base.ctx) }()

	start := time.Now()
	err := base.create_database("HALLO_BASE")
	end := time.Since(start)
	fmt.Println("Time for create database: ", end)
	if err != nil {
		panic(err)
	}

	start = time.Now()
	err = base.write_to_database("HALLO_BASE", m)
	end = time.Since(start)
	fmt.Println("Time for write to database: ", end)
	if err != nil {
		panic(err)
	}

	data := YDBTableMessage{}
	start = time.Now()
	base.read_from_table("HALLO_BASE", 700, &data)
	end = time.Since(start)
	fmt.Println("Time for read from database:", end)

	if data.msg_id != m[0].msg_id {
		t.Fatalf("Read msg_id != write msg_id: %d != %d", data.msg_id, m[0].msg_id)
	}

	for i, val := range data.board {
		if val != m[0].board[i] {
			t.Fatalf("Read value != write value: %d != %d", val, m[0].board[i])
		}
	}

}

//export YDB_ENDPOINT=grpcs://ydb.serverless.yandexcloud.net:2135
//export YDB_CONNECTION_STRING=grpcs://ydb.serverless.yandexcloud.net:2135/ru-central1/b1g7t9amkto3it85f06v/etnl9onhehmsj3fgjkss
//export YDB_ACCESS_TOKEN_CREDENTIALS=t1.9euelZqRlZPNm8yQipHIxp2SyJiXx-3rnpWakouVz8mKjpyKls6OkMiYjo7l8_dcbnBo-e9RNQEr_d3z9xwdbmj571E1ASv9.uZvHWcraQ-XUSG9szzDy7X-vptKr4j2TjBenTmrZP8aXYHTipMVobED9tg8fbe_2iAbyaBlNWxLEpYwIok6PCw
