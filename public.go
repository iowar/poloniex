package poloniex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Ticker struct {
	ID            int             `json:"id, int"`
	Last          decimal.Decimal `json:"last, string"`
	LowestAsk     decimal.Decimal `json:"lowestAsk, string"`
	HighestBid    decimal.Decimal `json:"highestBid, string"`
	PercentChange decimal.Decimal `json:"percentChange, string"`
	BaseVolume    decimal.Decimal `json:"baseVolume, string"`
	QuoteVolume   decimal.Decimal `json:"quoteVolume, string"`
	IsFrozen      int             `json:"isFrozen ,string"`
	High24hr      decimal.Decimal `json:"high24hr, string"`
	Low24hr       decimal.Decimal `json:"low24hr, string"`
}

func (p *Poloniex) PubReturnTickers() (tickers map[string]Ticker, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest("returnTicker", respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &tickers)
	return
}

type Volume struct {
	Volumes   map[string]map[string]decimal.Decimal
	TotalBTC  float64 `json:"totalBTC, string"`
	TotalETH  float64 `json:"totalETH, string"`
	TotalUSDT float64 `json:"totalUSDT, string"`
	TotalXMR  float64 `json:"totalXMR, string"`
	TotalXUSD float64 `json:"totalXUSD, string"`
}

func (v *Volume) UnmarshalJSON(b []byte) error {
	m := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	v.Volumes = make(map[string]map[string]decimal.Decimal)

	for key, value := range m {
		switch key {
		case "totalBTC":
			f, err := parseJSONFloatString(value)
			if err != nil {
				return err
			}

			v.TotalBTC = f

		case "totalETH":
			f, err := parseJSONFloatString(value)
			if err != nil {
				return err
			}

			v.TotalETH = f

		case "totalUSDT":
			f, err := parseJSONFloatString(value)
			if err != nil {
				return err
			}

			v.TotalUSDT = f

		case "totalXMR":
			f, err := parseJSONFloatString(value)
			if err != nil {
				return err
			}

			v.TotalXMR = f

		case "totalXUSD":
			f, err := parseJSONFloatString(value)
			if err != nil {
				return err
			}

			v.TotalXUSD = f

		default:
			g := make(map[string]decimal.Decimal)
			err := json.Unmarshal(value, &g)
			if err != nil {
				return err
			}
			v.Volumes[key] = g
		}
	}

	return err

}

func (p *Poloniex) PubReturn24hVolume() (volumes Volume, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest("return24hVolume", respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &volumes)
	return

}

type Order struct {
	Asks     [][]interface{} `json:"asks"`
	Bids     [][]interface{} `json:"bids"`
	IsFrozen string          `json:"isFrozen"`
}

func (p *Poloniex) PubReturnOrderBook(market string, depth int) (orders Order, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest(fmt.Sprintf("returnOrderBook&currencyPair=%s&depth=%d",
		strings.ToUpper(market), depth), respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &orders)
	return
}

type Trade struct {
	GlobalTradeID uint64          `json:"globalTradeID"`
	TradeID       uint64          `json:"tradeID"`
	Date          string          `json:"date, string"`
	Type          string          `json:"type, string"`
	Rate          decimal.Decimal `json:"rate, string"`
	Amount        decimal.Decimal `json:"amount, string"`
	Total         decimal.Decimal `json:"total, string"`
}

func (p *Poloniex) PubReturnTradeHistory(market string, args ...time.Time) (trades []Trade, err error) {
	trades = make([]Trade, 0)
	respch := make(chan []byte)
	errch := make(chan error)

	action := fmt.Sprintf("returnTradeHistory&currencyPair=%s", strings.ToUpper(market))

	if len(args) == 2 {
		start := strconv.FormatInt(args[0].UnixNano(), 10)
		end := strconv.FormatInt(args[1].UnixNano(), 10)
		action += fmt.Sprintf("&start=%s&end=%s", start, end)
	}

	go p.publicRequest(action, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &trades)
	return
}

type CandleStick struct {
	Date            decimal.Decimal `json:"date"`
	High            decimal.Decimal `json:"high"`
	Low             decimal.Decimal `json:"low"`
	Open            decimal.Decimal `json:"open"`
	Close           decimal.Decimal `json:"close"`
	Volume          decimal.Decimal `json:"volume"`
	QuoteVolume     decimal.Decimal `json:"quoteVolume"`
	WeightedAverage decimal.Decimal `json:"weightedAverage"`
}

func (p *Poloniex) PubReturnChartData(market string, start, end time.Time, period int) (candles []CandleStick, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	var v1, v2 int64
	var periods = []int{300, 900, 1800, 7200, 14400, 86400}

	if intInSlice(period, periods) == false {
		return nil, Error(PeriodError)
	}

	action := fmt.Sprintf("returnChartData&currencyPair=%s",
		strings.ToUpper(market))

	if start.IsZero() == false && end.IsZero() == false {
		v1 = start.Unix()
		v2 = end.Unix()

		if int((v2 - v1)) < period {
			return nil, Error(TimePeriodError)
		}

	} else if start.IsZero() == true && end.IsZero() == true {
		v1 = time.Now().AddDate(0, 0, -1).Unix()
		v2 = time.Now().Unix()
	} else {
		return nil, Error(TimeError)
	}

	action += fmt.Sprintf("&start=%d&end=%d&period=%d",
		v1, v2, period)

	go p.publicRequest(action, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &candles)
	return
}

type Currency struct {
	Id             int             `json:"id"`
	Name           string          `json:"name"`
	TxFee          decimal.Decimal `json:"txFee"`
	MinConf        decimal.Decimal `json:"minConf"`
	DepositAddress string          `json:"depositAddress"`
	Disabled       int             `json:"disabled"`
	Delisted       int             `json:"delisted"`
	Frozen         int             `json:"frozen"`
}

func (p *Poloniex) PubReturnCurrencies() (currencies map[string]Currency, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest("returnCurrencies", respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &currencies)
	return
}

type LoanOrderSc struct {
	Rate     decimal.Decimal `json:"rate, string"`
	Amount   decimal.Decimal `json:"amount, string"`
	RangeMin int             `json:"rangeMin"`
	RangeMax int             `json:"rangeMax"`
}

type LoanOrder struct {
	Offers  []LoanOrderSc `json:"offers"`
	Demands []LoanOrderSc `json:"demands"`
}

func (p *Poloniex) PubReturnLoanOrders(currency string) (loanorders LoanOrder, err error) {

	respch := make(chan []byte)
	errch := make(chan error)

	action := fmt.Sprintf("returnLoanOrders&currency=%s", currency)
	go p.publicRequest(action, respch, errch)

	response := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &loanorders)
	return
}
