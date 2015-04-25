package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

type Msg map[string]interface{}

func generateRandomMap(n int) Msg {
	msg := make(Msg)
	for i := 0; i < n; i++ {
		msg[string(i)] = i
	}
	return msg
}

func testPerf(msg Msg, n int, method int) {
	start_time := time.Now()
	for i := 0; i < n; i++ {
		if method == 1 {
			ms, _ := json.Marshal(msg)
			json.Unmarshal(ms, make(Msg))
		} else {
			ms, _ := bson.Marshal(msg)
			bson.Unmarshal(ms, make(Msg))
		}
	}
	stop_time := time.Now()
	fmt.Printf("running %d times takes %v\n", n, stop_time.Sub(start_time))
}

func main() {
	msg := generateRandomMap(0)
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	os.Stdout.Write(b)

	b, err = bson.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	os.Stdout.Write(b)

	testPerf(generateRandomMap(1000), 1000, 1)
	testPerf(generateRandomMap(1000), 1000, 2)
}
