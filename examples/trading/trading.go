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

	resp, err := poloniex.GetBalances()
	//resp, err := poloniex.GetCompleteBalances()
	//resp, err := poloniex.GetAccountBalances()
	//resp, err := poloniex.GetDepositAdresses()
	//resp, err := poloniex.GenerateNewAddress("USDT")
	//resp, err := poloniex.GetOpenOrders("btc_dgb")
	//resp, err := poloniex.GetAllOpenOrders()
	//resp, err := poloniex.CancelOrder("36121803064")
	//resp, err := poloniex.GetTradeHistory("btc_eth", time.Now().AddDate(0, 0, -600), time.Now(), 1)
	//resp, err := poloniex.GetTradesByOrderID("414366201166")
	//resp, err := poloniex.GetOrderStat("36121689178")
	//resp, err := poloniex.Buy("btc_dgb", 0.00000001, 23000)
	//resp, err := poloniex.Sell("btc_dgb", 1, 23.1)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
