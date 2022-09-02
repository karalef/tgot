package bot

import (
	"sync"
	"tghwbot/bot/tg"
)

// QueryContext is the common context for all queries that require an answer.
type QueryContext[T interface {
	InlineAnswer | CallbackAnswer |
		ShippingAnswer | PreCheckoutAnswer
	answerData(p params, queryID string) (method string)
}] struct {
	Context
	queryID string

	once sync.Once
}

func (c *QueryContext[T]) Answer(answer T) {
	c.once.Do(func() {
		p := params{}
		c.api(answer.answerData(p, c.queryID), p)
	})
}

var _ QueryContext[InlineAnswer]

// InlineAnswer represents an answer to inline query.
type InlineAnswer struct {
	Results           []tg.InlineQueryResulter
	CacheTime         *int
	IsPersonal        bool
	NextOffset        string
	SwitchPMText      string
	SwitchPMParameter string
}

func (a InlineAnswer) answerData(p params, queryID string) string {
	p.set("inline_query_id", queryID)
	p.setJSON("results", a.Results)
	if a.CacheTime != nil {
		p.setInt("cache_time", *a.CacheTime)
	}
	p.setBool("is_personal", a.IsPersonal)
	p.set("next_offset", a.NextOffset)
	p.set("switch_pm_text", a.SwitchPMText)
	p.set("switch_pm_parameter", a.SwitchPMParameter)
	return "answerInlineQuery"
}

var _ QueryContext[CallbackAnswer]

// CallbackAnswer represents an answer to callback query.
type CallbackAnswer struct {
	Text      string
	ShowAlert bool
	URL       string
	CacheTime *int
}

func (a CallbackAnswer) answerData(p params, queryID string) string {
	p.set("callback_query_id", queryID)
	p.set("text", a.Text)
	p.setBool("show_alert", a.ShowAlert)
	p.set("url", a.URL)
	if a.CacheTime != nil {
		p.setInt("cache_time", *a.CacheTime)
	}
	return "answerCallbackQuery"
}
