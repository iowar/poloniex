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
	resp, err := poloniex.TradeReturnBalances()
	//resp, err := poloniex.TradeReturnCompleteBalances()
	//resp, err := poloniex.TradeReturnDepositAdresses()
	//resp, err := poloniex.TradeGenerateNewAddress("BTC")
	//resp, err := poloniex.TradeReturnOpenOrders("btc_lsk")
	//resp, err := poloniex.TradeReturnTradeHistory("btc_lsk")
	//resp, err := poloniex.TradeReturnTradeHistory("btc_lsk", polo.ZeroTime, time.Now(), 100)
	//resp, err := poloniex.TradeReturnOrderTrade(70899268924)
	//resp, err := poloniex.TradeBuy("btc_bcn", 0.00000064, 100000)
	//resp, err := poloniex.TradeSell("usdt_btc", 15500, 0.01)
	//resp, err := poloniex.TradeCancelOrder(128001639152)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
