package main

import (
	"fmt"

	polo "github.com/iowar/poloniex"
)

const (
	api_key    = ""
	api_secret = ""
)

func main() {
	poloniex, err := polo.NewClient(api_key, api_secret, true)
	resp, err := poloniex.PubReturnTicker()
	//resp, err := poloniex.PubReturn24hVolume()
	//resp, err := poloniex.PubReturnOrderBook("btc_dgb", 5)
	//resp, err := poloniex.PubReturnTradeHistory("btc_dgb")
	//resp, err := poloniex.PubReturnChartData("btc_dgb", polo.ZeroTime, polo.ZeroTime, 300)
	//resp, err := poloniex.PubReturnChartData("btc_dgb", time.Now().AddDate(0, 0, -5), time.Now(), 86400)
	//resp, err := poloniex.PubReturnCurrencies()
	//resp, err := poloniex.PubReturnLoanOrders("BTC")

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
