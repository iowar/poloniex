//the following code shows
//how to access NewTrade fields.
package main

import (
	"fmt"

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

	var n polo.NewTrade

	for {
		receive := <-ws.Subs["usdt_btc"]
		updates := receive.([]polo.MarketUpdate)
		for _, v := range updates {
			if v.TypeUpdate == "NewTrade" {
				n = v.Data.(polo.NewTrade)
				fmt.Printf("TradeId:%d, Rate:%f, Amount:%f, Total:%f, Type:%s\n",
					n.TradeId, n.Rate, n.Amount, n.Total, n.TypeOrder)
			}
		}
	}
}
