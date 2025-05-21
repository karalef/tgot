package api

import (
	"context"
	"math"
)

// CurrenciesDataURL is the url that provides data about currencies.
const CurrenciesDataURL = "https://core.telegram.org/bots/payments/currencies.json"

// GetCurrenciesData returns information about currencies supported by the telegram api.
func (a *API) GetCurrenciesData(ctx context.Context) (map[string]Currency, error) {
	body, herr := a.get(ctx, CurrenciesDataURL)
	if herr != nil {
		return nil, herr
	}
	defer body.Close()

	r, raw, err := DecodeJSON[map[string]Currency](body)
	if err != nil {
		return nil, &JSONError{
			baseError: makeError("/bots/payments/currencies.json", nil, err),
			Status:    200,
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
