package poloniex

import (
	"encoding/json"
	"strconv"
	"time"
)

var ZeroTime = time.Time{}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringInSlice(a string, list []string) string {
	for _, b := range list {
		if b == a {
			return b
		}
	}
	return ""
}

func parseJSONFloatString(data json.RawMessage) (float64, error) {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(s, 64)
}
