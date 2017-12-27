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

	err = ws.SubscribeMarket("usdt_btc")
	if err != nil {
		return
	}

	go func() {
		/* If the logger is enabled, LogBus can be used */
		ws.LogBus <- "[*] Starting Unsubscribe goroutine"
		time.Sleep(time.Second * 10)
		ws.UnsubscribeMarket("usdt_btc")
	}()

	for {
		fmt.Println(<-ws.Subs["usdt_btc"])
	}
}
