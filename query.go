package tgot

import (
	"sync"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

type answerable interface {
	answerData(data *api.Data, queryID string) (method string)
}

func makeQueryContext[T answerable](ctx Context, queryID string) QueryContext[T] {
	return QueryContext[T]{
		Context: ctx,
		queryID: queryID,
		once:    new(sync.Once),
	}
}

// QueryContext is the common context for all queries that require an answer.
type QueryContext[T answerable] struct {
	Context
	queryID string

	once *sync.Once
}

func (c QueryContext[T]) Ctx() Context {
	return c.Context
}

func (c QueryContext[T]) Child(name string) QueryContext[T] {
	c.Context = c.Context.Child(name)
	return c
}

func (c QueryContext[T]) Answer(answer T) (err error) {
	c.once.Do(func() {
		data := api.NewData()
		err = c.method(answer.answerData(data, c.queryID), data)
	})
	return
}

// InlineContext type.
type InlineContext = QueryContext[InlineAnswer]

// InlineAnswer represents an answer to inline query.
type InlineAnswer struct {
	Results    []tg.InlineQueryResulter
	CacheTime  *int
	IsPersonal bool
	NextOffset string
	Button     *tg.InlineQueryResultsButton
}

func (a InlineAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("inline_query_id", queryID)
	d.SetJSON("results", a.Results)
	if a.CacheTime != nil {
		d.SetInt("cache_time", *a.CacheTime, true)
	}
	d.SetBool("is_personal", a.IsPersonal)
	d.Set("next_offset", a.NextOffset)
	d.SetJSON("button", a.Button)
	return "answerInlineQuery"
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

func (a CallbackAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("callback_query_id", queryID)
	d.Set("text", a.Text)
	d.SetBool("show_alert", a.ShowAlert)
	d.Set("url", a.URL)
	d.SetInt("cache_time", a.CacheTime)
	return "answerCallbackQuery"
}
