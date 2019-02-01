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

func (p *Poloniex) GetTickers() (tickers map[string]Ticker, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest("returnTicker", respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &tickers)
	return
}

type Volume struct {
	Volumes   map[string]map[string]decimal.Decimal
	TotalBTC  float64 `json:"totalBTC, string"`
	TotalETH  float64 `json:"totalETH, string"`
	TotalUSDC float64 `json:"totalUSDC, string"`
	TotalUSDT float64 `json:"totalUSDT, string"`
	TotalXMR  float64 `json:"totalXMR, string"`
	TotalXUSD float64 `json:"totalXUSD, string"`
}

func (v *Volume) UnmarshalJSON(b []byte) error {
	rmsg := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &rmsg)
	if err != nil {
		return err
	}

	v.Volumes = make(map[string]map[string]decimal.Decimal)

	for key, value := range rmsg {
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

		case "totalUSDC":
			f, err := parseJSONFloatString(value)
			if err != nil {
				return err
			}

			v.TotalUSDC = f

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
			cf := make(map[string]decimal.Decimal)
			err := json.Unmarshal(value, &cf)
			if err != nil {
				return err
			}
			v.Volumes[key] = cf
		}
	}

	return err
}

func (p *Poloniex) Get24hVolumes() (volumes Volume, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest("return24hVolume", respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &volumes)
	return

}

type Book struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

func (bk *Book) UnmarshalJSON(b []byte) error {
	var msg []interface{}

	err := json.Unmarshal(b, &msg)
	if err != nil {
		return err
	}

	price, err := strconv.ParseFloat(msg[0].(string), 64)
	if err != nil {
		return err
	}

	bk.Price = price
	bk.Quantity = msg[1].(float64)
	return nil
}

type OrderBook struct {
	Asks     []Book `json:"asks"`
	Bids     []Book `json:"bids"`
	IsFrozen string `json:"isFrozen"`
	Seq      int    `json:"seq"`
}

func (p *Poloniex) GetOrderBook(market string, depth int) (orderbook OrderBook, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest(fmt.Sprintf("returnOrderBook&currencyPair=%s&depth=%d",
		strings.ToUpper(market), depth), respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &orderbook)
	return
}

type PublicTrade struct {
	GlobalTradeID uint64          `json:"globalTradeID"`
	TradeID       uint64          `json:"tradeID"`
	Date          string          `json:"date, string"`
	Type          string          `json:"type, string"`
	Rate          decimal.Decimal `json:"rate, string"`
	Amount        decimal.Decimal `json:"amount, string"`
	Total         decimal.Decimal `json:"total, string"`
}

func (p *Poloniex) GetPublicTradeHistory(market string, args ...time.Time) (trades []PublicTrade, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	action := fmt.Sprintf("returnTradeHistory&currencyPair=%s", strings.ToUpper(market))

	if len(args) == 2 {
		action += fmt.Sprintf("&start=%d&end=%d", args[0].Unix(), args[1].Unix())
	}

	go p.publicRequest(action, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &trades)
	return
}

type CandleStick struct {
	Date            int64   `json:"date"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Open            float64 `json:"open"`
	Close           float64 `json:"close"`
	Volume          float64 `json:"volume"`
	QuoteVolume     float64 `json:"quoteVolume"`
	WeightedAverage float64 `json:"weightedAverage"`
}

func (p *Poloniex) GetChartData(market string, start, end time.Time, period string) (candles []CandleStick, err error) {
	var periodSec int
	var v1, v2 int64

	switch period {
	case "5m":
		periodSec = 300
	case "15m":
		periodSec = 900
	case "30m":
		periodSec = 1800
	case "2h":
		periodSec = 7200
	case "4h":
		periodSec = 14400
	case "1d":
		periodSec = 86400
	default:
		return nil, Error(PeriodError)
	}

	action := fmt.Sprintf("returnChartData&currencyPair=%s",
		strings.ToUpper(market))

	if start.IsZero() == false && end.IsZero() == false {
		v1 = start.Unix()
		v2 = end.Unix()

		if int((v2 - v1)) < periodSec {
			return nil, Error(TimePeriodError)
		}

	} else if start.IsZero() == true && end.IsZero() == true {
		v1 = time.Now().AddDate(0, 0, -1).Unix()
		v2 = time.Now().Unix()
	} else {
		return nil, Error(TimeError)
	}

	respch := make(chan []byte)
	errch := make(chan error)

	action += fmt.Sprintf("&start=%d&end=%d&period=%d",
		v1, v2, periodSec)

	go p.publicRequest(action, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &candles)
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

func (p *Poloniex) GetCurrencies() (currencies map[string]Currency, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	go p.publicRequest("returnCurrencies", respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &currencies)
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

func (p *Poloniex) GetLoanOrders(currency string) (loanorders LoanOrder, err error) {
	respch := make(chan []byte)
	errch := make(chan error)

	action := fmt.Sprintf("returnLoanOrders&currency=%s", currency)
	go p.publicRequest(action, respch, errch)

	resp := <-respch
	err = <-errch

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &loanorders)
	return
}
