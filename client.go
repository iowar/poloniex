package poloniex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	origin        = "https://api2.poloniex.com/"
	pushAPIUrl    = "wss://api2.poloniex.com/realm1"
	publicAPIUrl  = "https://poloniex.com/public?command="
	tradingAPIUrl = "https://poloniex.com/tradingApi"
)

var (
	throttle = time.Tick(time.Second / 5)
)

type Poloniex struct {
	key        string
	secret     string
	httpClient *http.Client
}

func NewClient(key, secret string) (client *Poloniex, err error) {
	client = &Poloniex{
		key:        key,
		secret:     secret,
		httpClient: &http.Client{Timeout: time.Second * 10},
	}

	return
}

// Create public api request.
func (p *Poloniex) publicRequest(action string, respch chan<- []byte, errch chan<- error) {
	defer close(respch)
	defer close(errch)

	rawurl := publicAPIUrl + action

	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		respch <- nil
		errch <- Error(RequestError)
		return
	}

	req.Header.Add("Accept", "application/json")

	<-throttle
	resp, err := p.httpClient.Do(req)
	if err != nil {
		respch <- nil
		errch <- Error(ConnectError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respch <- body
		errch <- err
		return
	}

	respch <- body
	errch <- nil
}

func (p *Poloniex) sign(formData string) (Sign string, err error) {
	if len(p.key) == 0 || len(p.secret) == 0 {
		err = Error(SetApiError)
		return
	}

	mac := hmac.New(sha512.New, []byte(p.secret))
	_, err = mac.Write([]byte(formData))
	if err != nil {
		return
	}

	Sign = hex.EncodeToString(mac.Sum(nil))
	return
}

type checkErr struct {
	Error string `json:"error"`
}

func checkServerError(response []byte) error {
	var check checkErr

	err := json.Unmarshal(response, &check)
	if err != nil {
		return nil
	}
	if check.Error != "" {
		return Error(ServerError, check.Error)
	} else {
		return nil
	}
}

// Create trading api request.
func (p *Poloniex) tradingRequest(action string, parameters map[string]string,
	respch chan<- []byte, errch chan<- error) {

	defer close(respch)
	defer close(errch)

	if parameters == nil {
		parameters = make(map[string]string)
	}
	parameters["command"] = action
	parameters["nonce"] = strconv.FormatInt(time.Now().UnixNano(), 10)

	formValues := url.Values{}

	for k, v := range parameters {
		formValues.Set(k, v)
	}

	formData := formValues.Encode()

	sign, err := p.sign(formData)
	if err != nil {
		respch <- nil
		errch <- err
		return
	}

	req, err := http.NewRequest("POST", tradingAPIUrl,
		strings.NewReader(formData))
	if err != nil {
		respch <- nil
		errch <- Error(RequestError)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Key", p.key)
	req.Header.Add("Sign", sign)

	<-throttle
	resp, err := p.httpClient.Do(req)
	if err != nil {
		respch <- nil
		errch <- Error(ConnectError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respch <- body
		errch <- err
		return
	}

	err = checkServerError(body)
	if err != nil {
		respch <- nil
		errch <- err
	}

	respch <- body
	errch <- nil
}
