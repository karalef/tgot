package tgot

import (
	"sync"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// MakeWebAppContext makes WebApp query context.
func MakeWebAppContext(ctx Context, queryID string) *WebAppContext {
	return &WebAppContext{
		Context: ctx,
		queryID: queryID,
	}
}

// WebAppContext is the context for WebApp query.
type WebAppContext struct {
	Context
	queryID string

	once *sync.Once
}

// Child creates sub context.
func (c WebAppContext) Child(name string) WebAppContext {
	c.Context = c.Context.Child(name)
	return c
}

// Answer answers to the WebApp query.
func (c WebAppContext) Answer(result tg.InlineQueryResulter) (sent *tg.SentWebAppMessage, err error) {
	c.once.Do(func() {
		d := api.NewData()
		d.Set("web_app_query_id", c.queryID)
		d.SetJSON("result", result)
		sent, err = method[*tg.SentWebAppMessage](c.Context, "answerWebAppQuery", d)
	})
	return
}
