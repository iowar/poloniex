# Poloniex API GO
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

### TrollBox
TrollBox is disabled from poloniex. We will give support if it is enable.

### Examples
* See [Push Api Examples](https://github.com/iowar/poloniex/tree/master/examples/push)

## Public Api
~~~go
poloniex, err := polo.NewClient(api_key, api_secret, true)
~~~
* Public Api Methods
    * PubReturnTicker()
    * PubReturn24hVolume()
    * PubReturnOrderBook()
    * PubReturnTradeHistory()
    * PubReturnChartData()
    * PubReturnCurrencies()
    * PubReturnLoanOrders()
    
#### Example
~~~go
resp, err := poloniex.PubReturnTicker()
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
    * TradeReturnBalances()
    * TradeReturnCompleteBalances()
    * TradeReturnDepositAdresses()
    * TradeGenerateNewAddress()
    * TradeReturnOpenOrders()
    * TradeReturnAllOpenOrders()
    * TradeReturnTradeHistory()
    * TradeReturnTradeHistory()
    * TradeReturnOrderTrade()
    * TradeBuy()
    * TradeSell()
    * TradeCancelOrder()

#### Example
~~~go
resp, err := poloniex.TradeBuy("btc_dgb", 0.00000099, 10000)
if err != nil{
    panic(err)
}
fmt.Println(resp)
~~~
* See [Trading Api Examples](https://github.com/iowar/poloniex/tree/master/examples/trading)

License
----
[MIT](https://github.com/iowar/poloniex/blob/master/LICENSE)


## Donations

| Name | Address |
| ------ | ------ |
| BTC | 1JM2rchFeVtLCTMUeinUcwVc2nnd5jtawX |
| LTC | LV46TbxeQxD9GyReQyvd5y366Nvn1MrnuF |
| DGB | DG2R4YxAywenpkVkWS1n3szPWYgBzyxmoZ |
| USDT | 17fGT7stxZjiREJ8ajAsdNegTRPYNW5Ao1 |

