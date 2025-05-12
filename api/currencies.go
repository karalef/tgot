package api

import (
	"context"
	"math"
	"net/http"
)

// CurrenciesDataURL is the url that provides data about currencies.
const CurrenciesDataURL = "https://core.telegram.org/bots/payments/currencies.json"

// GetCurrenciesData returns information about currencies supported by the telegram api.
func (a *API) GetCurrenciesData(ctx context.Context) (map[string]Currency, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, CurrenciesDataURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, raw, err := DecodeJSON[map[string]Currency](resp.Body)
	if err != nil {
		return nil, &JSONError{
			baseError: makeError("/bots/payments/currencies.json", nil, err),
			Status:    resp.StatusCode,
			Response:  raw,
		}
	}
	return *r, nil
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
