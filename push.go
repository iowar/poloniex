package poloniex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

const (
	TICKER     = "1002" /* Ticker Channel Id */
	SUBSBUFFER = 24     /* Subscriptions Buffer */
)

var (
	mutex          = &sync.Mutex{}
	channelsByName = make(map[string]string)
	channelsByID   = make(map[string]string)
	marketChannels []string
)

type subscription struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
}

type WSTicker struct {
	CurrencyPair  string  `json:"currencyPair"`
	Last          float64 `json:"last"`
	LowestAsk     float64 `json:"lowestAsk"`
	HighestBid    float64 `json:"hihgestBid"`
	PercentChange float64 `json:"percentChange"`
	BaseVolume    float64 `json:"baseVolume"`
	QuoteVolume   float64 `json:"quoteVolume"`
	IsFrozen      bool    `json:"isFrozen"`
	High24hr      float64 `json:"high24hr"`
	Low24hr       float64 `json:"low24hr"`
}

type MarketUpdate struct {
	Data       interface{}
	TypeUpdate string `json:"type"`
}

type OrderBook struct {
	Rate      float64 `json:"rate,string"`
	TypeOrder string  `json:"type"`
	Amount    float64 `json:"amount,string"`
}

type OrderBookModify OrderBook

type OrderBookRemove struct {
	Rate      float64 `json:"rate,string"`
	TypeOrder string  `json:"type"`
}

type NewTrade struct {
	TradeId   int64   `json:"tradeID,string"`
	Rate      float64 `json:"rate,string"`
	Amount    float64 `json:"amount,string"`
	Total     float64 `json:"total,string"`
	TypeOrder string  `json:"type"`
}

type WSClient struct {
	Subs       map[string]chan interface{}
	LogBus     chan<- string
	logger     Logger
	wssStopChs map[string]chan bool
	wssLock    *sync.Mutex
	wssClient  *websocket.Conn
}

func setchannelids() (err error) {
	p, err := NewClient("", "")
	if err != nil {
		return err
	}
	resp, err := p.PubReturnTickers()
	if err != nil {
		return err
	}

	for k, v := range resp {
		chid := strconv.Itoa(v.ID)
		channelsByName[k] = chid
		channelsByID[chid] = k
		marketChannels = append(marketChannels, chid)
	}

	channelsByName["TICKER"] = TICKER
	channelsByID[TICKER] = "TICKER"
	return
}

func NewWSClient(args ...bool) (wsclient *WSClient, err error) {
	ws, err := websocket.Dial(pushAPIUrl, "", origin)
	if err != nil {
		return
	}

	wsclient = &WSClient{
		wssClient:  ws,
		Subs:       make(map[string]chan interface{}),
		wssStopChs: make(map[string]chan bool),
		wssLock:    &sync.Mutex{},
	}

	if len(args) > 0 && args[0] {
		logbus := make(chan string)
		wsclient.LogBus = logbus
		wsclient.logger = Logger{isOpen: true, Lock: &sync.Mutex{}}

		go wsclient.logger.LogRoutine(logbus)
	}

	if err = setchannelids(); err != nil {
		return
	}

	if wsclient.logger.isOpen {
		wsclient.LogBus <- "[*] Created New WSClient."
	}

	return
}

func (ws *WSClient) subscribe(chid, chname string) (err error) {
	ws.wssLock.Lock()
	defer ws.wssLock.Unlock()

	if ws.Subs[chname] != nil {
		return Error(SubscribeError)
	}

	ws.Subs[chname] = make(chan interface{}, SUBSBUFFER)
	ws.wssStopChs[chname] = make(chan bool)

	subs := subscription{Command: "subscribe", Channel: chid}
	msg, err := json.Marshal(subs)
	if err != nil {
		return err
	}

	_, err = ws.wssClient.Write(msg)
	if err != nil {
		return err
	}

	if ws.logger.isOpen {
		ws.LogBus <- fmt.Sprintf("[*] Subscribed channel '%s'.\n", chname)
	}

	go func(chid string) {
		var imsg []interface{}
		var wsupdate interface{}
		var rmsg = make([]byte, 128)

		for {
			select {
			case <-ws.wssStopChs[chname]:
				close(ws.wssStopChs[chname])
				delete(ws.wssStopChs, chname)
				if ws.logger.isOpen {
					ws.LogBus <- fmt.Sprintf("[*] Unsubscribed channel '%s'.\n", chname)
				}
				return
			default:
			}

			read_len, err := ws.wssClient.Read(rmsg)
			if err != nil {
				return
			}

			err = json.Unmarshal(rmsg[:read_len], &imsg)
			if err != nil {
				continue
			}

			arg, ok := imsg[0].(float64)
			if !ok {
				continue
			}

			key := strconv.FormatFloat(arg, 'f', 0, 64)

			if key != chid || len(imsg) < 3 {
				continue
			}

			args, ok := imsg[2].([]interface{})
			if !ok {
				continue
			}

			switch chid {
			case TICKER:
				wsupdate, err = convertArgsToTicker(args)
				if err != nil {
					continue
				}
			case stringInSlice(chid, marketChannels):
				wsupdate, err = convertArgsToMarketUpdate(args)
				if err != nil {
					continue
				}
			default:
			}

			select {
			case ws.Subs[chname] <- wsupdate:
			default:

				/* fmt.Println("[BLOCKED] No Message Sent! ")  //DEBUG */
			}
		}
	}(chid)

	return nil
}

func convertArgsToTicker(args []interface{}) (wsticker WSTicker, err error) {
	wsticker.CurrencyPair = channelsByID[strconv.FormatFloat(args[0].(float64), 'f', 0, 64)]
	wsticker.Last, err = strconv.ParseFloat(args[1].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "Last")
		return
	}
	wsticker.LowestAsk, err = strconv.ParseFloat(args[2].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "LowestAsk")
		return
	}
	wsticker.HighestBid, err = strconv.ParseFloat(args[3].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "HighestBid")
		return
	}

	wsticker.PercentChange, err = strconv.ParseFloat(args[4].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "PercentChange")
		return
	}

	wsticker.BaseVolume, err = strconv.ParseFloat(args[5].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "BaseVolume")
		return
	}

	wsticker.QuoteVolume, err = strconv.ParseFloat(args[6].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "QuoteVolume")
		return
	}

	if v, ok := args[7].(float64); ok {
		if v == 0 {
			wsticker.IsFrozen = false
		} else {
			wsticker.IsFrozen = true
		}
	} else {
		err = Error(WSTickerError, "IsFrozen")
		return
	}

	wsticker.High24hr, err = strconv.ParseFloat(args[8].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "High24hr")
		return
	}

	wsticker.Low24hr, err = strconv.ParseFloat(args[9].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "Low24hr")
		return
	}

	return
}

func convertArgsToMarketUpdate(args []interface{}) (res []MarketUpdate, err error) {
	res = make([]MarketUpdate, len(args))
	for i, val := range args {
		vals := val.([]interface{})
		marketupdate := MarketUpdate{}
		orderdatafield := OrderBook{}
		tradedatafield := NewTrade{}

		switch vals[0].(string) {
		case "o":
			if vals[3].(string) == "0.00000000" {
				marketupdate.TypeUpdate = "OrderBookRemove"
			} else {
				marketupdate.TypeUpdate = "OrderBookModify"
			}

			if vals[1].(float64) == 1 {
				orderdatafield.TypeOrder = "bid"
			} else {
				orderdatafield.TypeOrder = "ask"
			}

			orderdatafield.Rate, err = strconv.ParseFloat(vals[2].(string), 64)
			if err != nil {
				err = Error(OrderBookError, "Rate")
				return
			}

			orderdatafield.Amount, err = strconv.ParseFloat(vals[3].(string), 64)
			if err != nil {
				err = Error(OrderBookError, "Amount")
				return
			}

			marketupdate.Data = orderdatafield

		case "t":
			marketupdate.TypeUpdate = "NewTrade"
			tradedatafield.TradeId, err = strconv.ParseInt(vals[1].(string), 10, 64)
			if err != nil {
				err = Error(NewTradeError, "TradeId")
				return
			}

			if vals[2].(float64) == 1 {
				tradedatafield.TypeOrder = "buy"
			} else {
				tradedatafield.TypeOrder = "sell"
			}

			tradedatafield.Rate, err = strconv.ParseFloat(vals[3].(string), 64)
			if err != nil {
				err = Error(NewTradeError, "Rate")
				return
			}

			tradedatafield.Amount, err = strconv.ParseFloat(vals[4].(string), 64)
			if err != nil {
				err = Error(NewTradeError, "Amount")
				return
			}

			tradedatafield.Total = vals[5].(float64)

			marketupdate.Data = tradedatafield
		}

		res[i] = marketupdate
	}
	return res, nil

}

func (ws *WSClient) unsubscribe(chid string) (err error) {
	ws.wssLock.Lock()
	defer ws.wssLock.Unlock()

	if ws.Subs[chid] == nil {
		return
	}

	subs := subscription{Command: "unsubscribe", Channel: chid}
	msg, err := json.Marshal(subs)
	if err != nil {
		return err
	}

	_, err = ws.wssClient.Write(msg)
	if err != nil {
		return err
	}

	ws.wssStopChs[chid] <- true
	close(ws.Subs[chid])
	delete(ws.Subs, chid)
	return
}

func (ws *WSClient) SubscribeTicker() error {
	return (ws.subscribe(TICKER, "ticker"))
}

func (ws *WSClient) UnsubscribeTicker() error {
	return (ws.unsubscribe("ticker"))
}

func (ws *WSClient) SubscribeMarket(chname string) error {
	chid := channelsByName[strings.ToUpper(chname)]
	if chid == "" {
		return Error(ChannelError, chname)
	}
	return (ws.subscribe(chid, chname))
}

func (ws *WSClient) UnsubscribeMarket(chname string) error {
	chid := channelsByName[strings.ToUpper(chname)]
	if chid == "" {
		return Error(ChannelError, chname)
	}
	return (ws.unsubscribe(chname))
}
