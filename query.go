package tgot

import (
	"sync"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

type answerable interface {
	InlineAnswer | CallbackAnswer | ShippingAnswer | PreCheckoutAnswer
	answerData(queryID string) (method string, data api.Data)
}

// QueryContext is the common context for all queries that require an answer.
type QueryContext[T answerable] struct {
	Context
	queryID string

	once sync.Once
}

func (c *QueryContext[T]) Answer(answer T) (err error) {
	c.once.Do(func() {
		err = c.method(answer.answerData(c.queryID))
	})
	return
}

// InlineContext type.
type InlineContext = QueryContext[InlineAnswer]

// InlineAnswer represents an answer to inline query.
type InlineAnswer struct {
	Results           []tg.InlineQueryResulter
	CacheTime         *int
	IsPersonal        bool
	NextOffset        string
	SwitchPMText      string
	SwitchPMParameter string
}

func (a InlineAnswer) answerData(queryID string) (string, api.Data) {
	d := api.NewData()
	d.Set("inline_query_id", queryID)
	d.SetJSON("results", a.Results)
	if a.CacheTime != nil {
		d.SetInt("cache_time", *a.CacheTime, true)
	}
	d.SetBool("is_personal", a.IsPersonal)
	d.Set("next_offset", a.NextOffset)
	d.Set("switch_pm_text", a.SwitchPMText)
	d.Set("switch_pm_parameter", a.SwitchPMParameter)
	return "answerInlineQuery", d
}

// CallbackContext type.
type CallbackContext = QueryContext[CallbackAnswer]

// CallbackAnswer represents an answer to callback query.
type CallbackAnswer struct {
	Text      string
	ShowAlert bool
	URL       string
	CacheTime int
}

func (a CallbackAnswer) answerData(queryID string) (string, api.Data) {
	d := api.NewData()
	d.Set("callback_query_id", queryID)
	d.Set("text", a.Text)
	d.SetBool("show_alert", a.ShowAlert)
	d.Set("url", a.URL)
	d.SetInt("cache_time", a.CacheTime)
	return "answerCallbackQuery", d
}
