# Poloniex API GO
[![GoDoc](https://godoc.org/github.com/iowar/poloniex?status.svg)](https://godoc.org/github.com/iowar/poloniex)
Poloniex Push, Public and Trading APIs.
# Install
~~~ go
$ go get github.com/iowar/poloniex
~~~ 

# APIs
~~~go
import polo "github.com/iowar/poloniex"
~~~
## Push Api
Create websocket client.
#### NewWSClient()
~~~go
ws, err := polo.NewWSClient()
if err != nil {
    return
}
~~~
* Push Api Methods
    * SubscribeTicker()
    * SubscribeMarket()
    * UnsubscribeTicker()
    * UnsubscribeMarket()

For Enable Logger 
~~~go
ws, err := polo.NewWSClient(true)
~~~
and access Logger
~~~go
ws.LogBus <- "Hello LightSide"
~~~



### Ticker
#### SubscribeTicker()
~~~go
err = ws.SubscribeTicker()
if err != nil {
    return
}
for {
    fmt.Println(<-ws.Subs["ticker"])
}
~~~
#### UnsubscribeTicker()
~~~go
err = ws.SubscribeTicker()
go func() {
    time.Sleep(time.Second * 10)
    ws.UnsubscribeTicker()
}()
for {
    fmt.Println(<-ws.Subs["ticker"])
}
~~~

### OrderBook and Trades
#### SubscribeMarket()
~~~go
err = ws.SubscribeMarket("usdt_btc")
if err != nil {
    return
}
for {
    fmt.Println(<-ws.Subs["usdt_btc"])
}
~~~
#### UnsubscribeMarket()
~~~go
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
~~~~

### Examples
* See [Push Api Examples](https://github.com/iowar/poloniex/tree/master/examples/push)

## Public Api
~~~go
poloniex, err := polo.NewClient(api_key, api_secret, true)
~~~
* Public Api Methods
    * GetTickers()
    * Get24hVolumes()
    * GetOrderBook()
    * GetPublicTradeHistory()
    * GetChartData()
    * GetCurrencies()
    * GetLoanOrders()
    
#### Example
~~~go
resp, err := poloniex.GetTickers()
if err != nil{
    panic(err)
}
fmt.Println(resp)
~~~
* See [Public Api Examples](https://github.com/iowar/poloniex/tree/master/examples/public)

## Trading Api
~~~go
const (
        api_key    = ""
        api_secret = ""
)
~~~
~~~go
poloniex, err := polo.NewClient(api_key, api_secret, true)
~~~ 

* Trading Api Methods
    * GetBalances()
    * GetCompleteBalances()
    * GetAccountBalances()
    * GetDepositAddresses()
    * GenerateNewAddress()
    * GetOpenOrders()
    * GetAllOpenOrders()
    * CancelOrder()
    * GetTradeHistory()
    * GetTradesByOrderID()
    * GetOrderStat()
    * Buy()
    * Sell()


#### Example
~~~go
resp, err := poloniex.Buy("btc_dgb", 0.00000099, 10000)
if err != nil{
    panic(err)
}
fmt.Println(resp)
~~~
* See [Trading Api Examples](https://github.com/iowar/poloniex/tree/master/examples/trading)

License
----
[MIT](https://github.com/iowar/poloniex/blob/master/LICENSE)

