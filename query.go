package tgot

import (
	"sync"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Answerable represents any object that is used as an answer to a query.
type Answerable interface {
	answerData(data *api.Data, queryID string) (method string)
}

// WithQuery creates a new Query context.
func WithQuery[T Answerable](ctx BaseContext, queryID string, from tg.User) Query[T] {
	return &queryContext[T]{
		context: ctx.ctx().with(nil),
		once:    new(sync.Once),
		user:    from,
		queryID: queryID,
	}
}

// Query is the context for all queries that require an answer.
// The answer method can be used only once.
type Query[T Answerable] interface {
	Context[Query[T]]
	Answer(T) error
	Sender() *User
}

type queryContext[T Answerable] struct {
	*context
	once    *sync.Once
	user    tg.User
	queryID string
}

func (c *queryContext[T]) WithName(name string) Query[T] {
	return &queryContext[T]{context: c.context.child(name), once: c.once, queryID: c.queryID}
}

func (c *queryContext[T]) Sender() *User { return WithUser(c, c.user.ID) }

func (c *queryContext[T]) Answer(answer T) (err error) {
	c.once.Do(func() {
		data := api.NewData()
		err = c.method(answer.answerData(data, c.queryID), data)
	})
	return
}

// InlineAnswer represents an answer to inline query.
type InlineAnswer struct {
	Results    []tg.InlineQueryResulter     `tg:"results"`
	CacheTime  *int                         `tg:"cache_time,force"`
	IsPersonal bool                         `tg:"is_personal"`
	NextOffset string                       `tg:"next_offset"`
	Button     *tg.InlineQueryResultsButton `tg:"button"`
}

func (a InlineAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("inline_query_id", queryID).AddObject(a)
	return "answerInlineQuery"
}

// CallbackAnswer represents an answer to callback query.
type CallbackAnswer struct {
	Text      string `tg:"text"`
	ShowAlert bool   `tg:"show_alert"`
	URL       string `tg:"url"`
	CacheTime int    `tg:"cache_time"`
}

func (a CallbackAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("callback_query_id", queryID).AddObject(a)
	return "answerCallbackQuery"
}

// ShippingAnswer represents an answer to shipping query.
type ShippingAnswer struct {
	OK              bool                `tg:"ok"`
	ShippingOptions []tg.ShippingOption `tg:"shipping_options"`
	ErrorMessage    string              `tg:"error_message"`
}

func (a ShippingAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("shipping_query_id", queryID).AddObject(a)
	return "answerShippingQuery"
}

// PreCheckoutAnswer represents an answer to pre-checkout query.
type PreCheckoutAnswer struct {
	OK           bool   `tg:"ok"`
	ErrorMessage string `tg:"error_message"`
}

func (a PreCheckoutAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("pre_checkout_query_id", queryID).AddObject(a)
	return "answerPreCheckoutQuery"
}

// WebAppAnswer represents an answer to WebApp query.
type WebAppAnswer struct {
	Result tg.InlineQueryResulter `tg:"result"`
}

func (a WebAppAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("web_app_query_id", queryID).AddObject(a)
	return "answerWebAppQuery"
}
