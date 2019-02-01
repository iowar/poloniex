package poloniex

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func (p *Poloniex) GetBalances() (balances map[string]string, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnBalances", nil, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &balances)
	return
}

type Balance struct {
	Available decimal.Decimal `json:"available, string"`
	OnOrders  decimal.Decimal `json:"onOrders, string"`
	BtcValue  decimal.Decimal `json:"btcValue, string"`
}

func (p *Poloniex) GetCompleteBalances() (completebalances map[string]Balance, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnCompleteBalances", nil, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &completebalances)
	return
}

type Accounts struct {
	Margin   map[string]decimal.Decimal `json:"margin"`
	Lending  map[string]decimal.Decimal `json:"lending"`
	Exchange map[string]decimal.Decimal `json:"exchange"`
}

func (p *Poloniex) GetAccountBalances() (accounts Accounts, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnAvailableAccountBalances", nil, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &accounts)
	return
}

func (p *Poloniex) GetDepositAddresses() (depositaddresses map[string]string, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnDepositAddresses", nil, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &depositaddresses)
	return
}

type NewAddress struct {
	Success  int    `json:"success"`
	Response string `json:"response"`
}

func (p *Poloniex) GenerateNewAddress(currency string) (newaddress NewAddress, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currency": strings.ToUpper(currency)}
	go p.tradingRequest("generateNewAddress", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &newaddress)
	return
}

type OpenOrder struct {
	OrderNumber    string          `json:"orderNumber"`
	Type           string          `json:"type"`
	Price          decimal.Decimal `json:"rate, string"`
	StartingAmount decimal.Decimal `json:"startingAmount, string"`
	Amount         decimal.Decimal `json:"amount, string"`
	Total          decimal.Decimal `json:"total, string"`
	Date           string          `json:"date"`
	Margin         int             `json:"margin"`
}

// Send market to get open orders.
func (p *Poloniex) GetOpenOrders(market string) (openorders []OpenOrder, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currencyPair": strings.ToUpper(market)}
	go p.tradingRequest("returnOpenOrders", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &openorders)
	return
}

// This method returns all open orders.
func (p *Poloniex) GetAllOpenOrders() (openorders map[string][]OpenOrder, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"currencyPair": "all"}
	go p.tradingRequest("returnOpenOrders", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &openorders)
	if err != nil {
		return
	}

	for k, v := range openorders {
		if len(v) == 0 {
			delete(openorders, k)
		}
	}
	return
}

type CancelOrder struct {
	Success int `json:"success"`
}

func (p *Poloniex) CancelOrder(orderNumber string) (cancelorder CancelOrder, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"orderNumber": orderNumber}
	go p.tradingRequest("cancelOrder", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &cancelorder)
	return
}

type TradeHistory struct {
	GlobalTradeID int             `json:"globalTradeId"`
	TradeID       string          `json:"tradeId"`
	Date          string          `json:"date"`
	Price         decimal.Decimal `json:"rate, string"`
	Amount        decimal.Decimal `json:"amount, string"`
	Total         decimal.Decimal `json:"total, string"`
	Fee           decimal.Decimal `json:"fee,string"`
	OrderNumber   decimal.Decimal `json:"orderNumber,string"`
	Type          string          `json:"type"`
	Category      string          `json:"category"`
}

func (p *Poloniex) GetTradeHistory(market string, start, end time.Time, limit int) (tradehistory []TradeHistory, err error) {
	parameters := map[string]string{
		"currencyPair": strings.ToUpper(market),
		"start":        strconv.FormatInt(start.Unix(), 10),
		"end":          strconv.FormatInt(end.Unix(), 10),
		"limit":        strconv.Itoa(limit),
	}

	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("returnTradeHistory", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &tradehistory)
	return
}

type OrderTrade struct {
	GlobalTradeID decimal.Decimal `json:"globalTradeId"`
	TradeID       decimal.Decimal `json:"tradeId"`
	Market        string          `json:"currencyPair"`
	Type          string          `json:"type"`
	Price         decimal.Decimal `json:"rate"`
	Amount        decimal.Decimal `json:"amount"`
	Total         decimal.Decimal `json:"total"`
	Fee           decimal.Decimal `json:"fee"`
	Date          string          `json:"date"`
}

func (p *Poloniex) GetTradesByOrderID(orderNumber string) (ordertrades []OrderTrade, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"orderNumber": orderNumber}
	go p.tradingRequest("returnOrderTrades", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &ordertrades)
	return
}

type OrderStat struct {
	Status         string          `json:"status"`
	Rate           decimal.Decimal `json:"rate"`
	Amount         decimal.Decimal `json:"amount"`
	CurrencyPair   string          `json:"currencyPair"`
	Date           string          `json:"date"`
	Total          decimal.Decimal `json:"total"`
	Type           string          `json:"type"`
	StartingAmount decimal.Decimal `json:"startingAmount"`
}

// error result
type OrderStat1 struct {
	Success int `json:"success"`
	Result  struct {
		Error string `json:"error"`
	} `json:"result"`
}

// success result
type OrderStat2 struct {
	Success int                  `json:"success"`
	Result  map[string]OrderStat `json:"result"`
}

func (p *Poloniex) GetOrderStat(orderNumber string) (orderstat OrderStat, err error) {
	var check1 OrderStat1
	var check2 OrderStat2

	respch := make(chan []byte)
	errch := make(chan error)

	parameters := map[string]string{"orderNumber": orderNumber}
	go p.tradingRequest("returnOrderStatus", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	// check error
	err = json.Unmarshal(resp, &check1)
	if err != nil {
		return
	}
	if check1.Success == 0 && len(check1.Result.Error) > 0 {
		err = errors.New(check1.Result.Error)
		return

	}

	// check success
	err = json.Unmarshal(resp, &check2)
	if err != nil {
		return
	}
	if check2.Success == 1 {
		orderstat = check2.Result[orderNumber]
		return
	}

	err = errors.New("Unexpected Result!")
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
	OrderNumber     string `json:"orderNumber"`
	ResultingTrades []ResultTrades
}

func (p *Poloniex) Buy(market string, price, amount float64) (buy Buy, err error) {
	parameters := map[string]string{
		"currencyPair": strings.ToUpper(market),
		"rate":         strconv.FormatFloat(float64(price), 'f', 8, 64),
		"amount":       strconv.FormatFloat(float64(amount), 'f', 8, 64),
	}

	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("buy", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &buy)
	return
}

type Sell Buy

func (p *Poloniex) Sell(market string, price, amount float64) (sell Sell, err error) {
	parameters := map[string]string{
		"currencyPair": strings.ToUpper(market),
		"rate":         strconv.FormatFloat(float64(price), 'f', 8, 64),
		"amount":       strconv.FormatFloat(float64(amount), 'f', 8, 64),
	}

	respch := make(chan []byte)
	errch := make(chan error)

	go p.tradingRequest("sell", parameters, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &sell)
	return
}
