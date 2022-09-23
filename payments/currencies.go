package payments

import (
	"encoding/json"
	"math"
	"net/http"
)

// DataURL is the url that provides data about currencies.
const DataURL = "https://core.telegram.org/bots/payments/currencies.json"

// GetCurrenciesData returns information about currencies supported by the telegram api.
func GetCurrenciesData(client *http.Client) (map[string]Currency, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(DataURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var r map[string]Currency
	return r, json.NewDecoder(resp.Body).Decode(&r)
}

// Currency represents single currency data.
type Currency struct {
	Code         string `json:"code"`
	Title        string `json:"title"`
	Symbol       string `json:"symbol"`
	Native       string `json:"native"`
	ThousandsSep string `json:"thousands_sep"`
	DecimalSep   string `json:"decimal_sep"`
	SymbolLeft   bool   `json:"symbol_left"`
	SpaceBetween bool   `json:"space_between"`
	Exp          int    `json:"exp"`
	MinAmount    int    `json:"min_amount,string"`
	MaxAmount    int    `json:"max_amount,string"`
}

// Amount converts the floating point price into a value for the payment api.
func (c Currency) Amount(price float32) int {
	if c.Exp == 0 {
		return int(price)
	}
	return int(price * float32(math.Pow10(c.Exp)))
}
