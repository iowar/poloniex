package poloniex

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func (p *Poloniex) TradeReturnBalances() (balances map[string]string, err error) {
	balances = make(map[string]string)
	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnBalances", nil, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &balances)
	return
}

type Balance struct {
	Available decimal.Decimal `json:"available, string"`
	OnOrders  decimal.Decimal `json:"onOrders, string"`
	BtcValue  decimal.Decimal `json:"btcValue, string"`
}

func (p *Poloniex) TradeReturnCompleteBalances() (completebalances map[string]Balance, err error) {

	completebalances = make(map[string]Balance)
	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnCompleteBalances", nil, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &completebalances)
	return

}

func (p *Poloniex) TradeReturnDepositAdresses() (depositaddresses map[string]string, err error) {

	depositaddresses = make(map[string]string)
	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnDepositAddresses", nil, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &depositaddresses)
	return
}

type NewAddress struct {
	Success  int
	Response string
}

func (p *Poloniex) TradeGenerateNewAddress(currency string) (newaddress NewAddress, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currency": strings.ToUpper(currency)}
	go p.tradingRequest("generateNewAddress", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &newaddress)
	return
}

/*
TODO
returnDepositsWithdrawals
*/

type OpenOrder struct {
	OrderNumber decimal.Decimal `json:"orderNumber, string"`
	Type        string          `json:"type, string"`
	Rate        decimal.Decimal `json:"rate, string"`
	/*StartingAmount decimal.Decimal `json:"startingAmount, string"`*/
	Amount decimal.Decimal `json:"amount, string"`
	Total  decimal.Decimal `json:"total, string"`
	/*Date           string*/
	/*Margin         int*/
}

func (p *Poloniex) TradeReturnOpenOrders(currency string) (openorders []OpenOrder, err error) {

	openorders = make([]OpenOrder, 0)
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currencyPair": strings.ToUpper(currency)}
	go p.tradingRequest("returnOpenOrders", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &openorders)
	return
}

//New Method
//Reason: different data type return
//when currency is 'all'
func (p *Poloniex) TradeReturnAllOpenOrders() (openorders map[string][]OpenOrder, err error) {

	openorders = make(map[string][]OpenOrder, 0)
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currencyPair": "all"}
	go p.tradingRequest("returnOpenOrders", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &openorders)
	return
}

//Self Trade History
type TradeHistory2 struct {
	Date        string          `json:"date"`
	Type        string          `json:"type"`
	Buy         decimal.Decimal `json:"buy, string"`
	Rate        decimal.Decimal `json:"rate, string"`
	Amount      decimal.Decimal `json:"amount, string"`
	Total       decimal.Decimal `json:"total, string"`
	OrderNumber decimal.Decimal `json:"order_number,string"`
	//Category
	//OrderNumber
	//Fee
	//TradeId
	//GlobalTradeId
}

func (p *Poloniex) TradeReturnTradeHistory(currency string, args ...interface{}) (tradehistory []TradeHistory2, err error) {

	tradehistory = make([]TradeHistory2, 0)
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currencyPair": strings.ToUpper(currency)}

	if len(args) >= 2 {
		start, ok := args[0].(time.Time)
		if ok == false {
			return nil, Error(StartTimeError)
		}
		end, ok := args[1].(time.Time)
		if ok == false {
			return nil, Error(EndTimeError)
		}

		parameters["start"] = strconv.FormatInt(start.UnixNano(), 10)
		parameters["end"] = strconv.FormatInt(end.UnixNano(), 10)
	}

	if len(args) == 3 {
		limit, ok := args[2].(int)
		if ok == false {
			return nil, Error(LimitError)
		}

		parameters["limit"] = strconv.Itoa(limit)
	}

	go p.tradingRequest("returnTradeHistory", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &tradehistory)
	return
}

type OrderTrade struct {
	GlobalTradeID decimal.Decimal
	TradeID       decimal.Decimal
	CurrencyPair  string
	Type          string
	Rate          decimal.Decimal
	Amount        decimal.Decimal
	Total         decimal.Decimal
	Fee           decimal.Decimal
	Date          string
}

func (p *Poloniex) TradeReturnOrderTrade(orderNumber int64) (ordertrades []OrderTrade, err error) {

	ordertrades = make([]OrderTrade, 0)
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"orderNumber": strconv.FormatInt(orderNumber, 10)}
	go p.tradingRequest("returnOrderTrades", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &ordertrades)
	return
}

type ResultTrades struct {
	Amount  decimal.Decimal `json:"amount"`
	Date    string          `json:"date"`
	Rate    decimal.Decimal `json:"rate"`
	Total   decimal.Decimal `json:"total"`
	TradeID decimal.Decimal `json:"tradeId"`
	Type    string          `json:"type"`
}

type Buy struct {
	OrderNumber     decimal.Decimal `json:"orderNumber"`
	ResultingTrades []ResultTrades
}

func (p *Poloniex) TradeBuy(currencyPair string, rate, amount float64) (buy Buy, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currencyPair": strings.ToUpper(currencyPair)}
	parameters["rate"] = strconv.FormatFloat(float64(rate), 'f', 8, 64)
	parameters["amount"] = strconv.FormatFloat(float64(amount), 'f', 8, 64)

	go p.tradingRequest("buy", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &buy)
	return
}

type Sell Buy

func (p *Poloniex) TradeSell(currencyPair string, rate, amount float64) (sell Sell, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currencyPair": strings.ToUpper(currencyPair)}
	parameters["rate"] = strconv.FormatFloat(float64(rate), 'f', 8, 64)
	parameters["amount"] = strconv.FormatFloat(float64(amount), 'f', 8, 64)

	go p.tradingRequest("sell", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &sell)
	return
}

type CancelOrder struct {
	Success int `json:"success"`
}

func (p *Poloniex) TradeCancelOrder(orderNumber int64) (cancelorder CancelOrder, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"orderNumber": strconv.FormatInt(orderNumber, 10)}
	go p.tradingRequest("cancelOrder", parameters, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &cancelorder)
	return
}
