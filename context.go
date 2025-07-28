package tgot

import (
	stdcontext "context"

	"github.com/karalef/tgot/api"
)

// Context represents any context.
type Context[T BaseContext] interface {
	BaseContext

	// Base copies the context without data.
	Base() Empty

	// WithName copies the context with nested name.
	WithName(name string) T
}

// BaseContext represents base context.
type BaseContext interface {
	stdcontext.Context

	// Path returns context path.
	Path() string

	// Bot returns bot instance.
	Bot() *Bot

	ctx() *context
}

// Empty represents empty context.
type Empty Context[Empty]

// NewContext creates new context.
func (b *Bot) NewContext(ctx stdcontext.Context, name string) Empty {
	return newContext(ctx, name, b, nil)
}

func newContext(c stdcontext.Context, path string, bot *Bot, data *api.Data) *context {
	if data != nil {
		data = data.Copy()
	}
	return &context{c, bot, path, data}
}

var _ Empty = &context{}

type context struct {
	stdcontext.Context
	bot  *Bot
	path string
	data *api.Data
}

func (c *context) ctx() *context              { return c }
func (c *context) Bot() *Bot                  { return c.bot }
func (c *context) Path() string               { return c.path }
func (c *context) Base() Empty                { return c.with(nil) }
func (c *context) WithName(name string) Empty { return c.child(name) }

func (c *context) with(d *api.Data) *context {
	if d == nil && c.data == nil {
		return c
	}
	return newContext(c.Context, c.path, c.bot, d)
}

func (c *context) add(d *api.Data) *context {
	if d == nil {
		return c
	}
	return newContext(c.Context, c.path, c.bot, c.data.WriteTo(d))
}

func (c *context) child(name string) *context {
	if name == "" {
		return c
	}
	if c.path != "" {
		name = c.path + "::" + name
	}
	return newContext(c.Context, name, c.bot, c.data)
}

func (c *context) method(meth string, d ...*api.Data) error {
	_, err := method[api.Empty](c, meth, d...)
	return err
}

func method[T any](c BaseContext, method string, d ...*api.Data) (T, error) {
	data := c.ctx().data
	if len(d) > 0 {
		data = data.WriteTo(d[0])
		defer data.Put()
	}

	return api.Request[T](c, c.Bot().API(), method, data)
}
