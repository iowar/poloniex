package main

import (
	"fmt"

	polo "github.com/iowar/poloniex"
)

func main() {

	ws, err := polo.NewWSClient()
	if err != nil {
		return
	}

	err = ws.SubscribeMarket("usdt_btc")
	if err != nil {
		return
	}

	for {
		fmt.Println(<-ws.Subs["usdt_btc"])
	}
}
