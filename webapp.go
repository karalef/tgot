package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// WebAppContext is the context for WebApp query.
type WebAppContext = QueryContext[WebAppAnswer]

// MakeWebAppContext makes WebApp query context.
func MakeWebAppContext(ctx Context, queryID string) WebAppContext {
	return makeQueryContext[WebAppAnswer](ctx, queryID)
}

// WebAppAnswer represents an answer to WebApp query.
type WebAppAnswer struct {
	Result tg.InlineQueryResulter
}

func (a WebAppAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("web_app_query_id", queryID)
	d.SetJSON("result", a.Result)
	return "answerWebAppQuery"
}
