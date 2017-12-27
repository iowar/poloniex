package main

import (
	"fmt"
	"time"

	polo "github.com/iowar/poloniex"
)

func main() {

	ws, err := polo.NewWSClient(true)
	if err != nil {
		return
	}

	err = ws.SubscribeTicker()
	if err != nil {
		return
	}

	go func() {
		time.Sleep(time.Second * 10)
		ws.UnsubscribeTicker()
	}()

	for {
		fmt.Println(<-ws.Subs["ticker"])
	}
}
