//the following code shows
//how to access OrderBook fields.
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

	var m polo.OrderBook

	for {
		receive := <-ws.Subs["usdt_btc"]
		updates := receive.([]polo.MarketUpdate)
		for _, v := range updates {
			if v.TypeUpdate == "OrderBookRemove" || v.TypeUpdate == "OrderBookModify" {
				m = v.Data.(polo.OrderBook)

				fmt.Printf("Rate:%f, Type:%s, Amount:%f\n",
					m.Rate, m.TypeOrder, m.Amount)
			}
		}
	}
}
