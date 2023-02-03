package router

import (
	"sync"
	"time"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/logger"
)

// NewRouter makes new inited queries router.
func NewRouter[Ctx ctxType[Ctx], Key comparable, Data any]() *Router[Ctx, Key, Data] {
	return &Router[Ctx, Key, Data]{
		handlers: make(map[Key]Handler[Ctx, Key, Data]),
	}
}

type baseContext interface {
	Ctx() tgot.Context
	Logger() *logger.Logger
}

type ctxType[c baseContext] interface {
	baseContext
	Child(string) c
}

// Handler represents handler with timeout.
type Handler[Ctx ctxType[Ctx], Key comparable, Data any] interface {
	Name() string

	// If unreg is true, handler will be automatically deleted.
	Handle(Ctx, Key, Data) (unreg bool, err error)

	// Specifies when the handler will be automatically unreged.
	Timeout() time.Time

	// Called when the handler times out.
	// The current handler will be automatically unreged so do not
	// call Unreg from this function as this will cause a deadlock.
	Cancel(tgot.Context, Key) error
}

// Router routes queries.
type Router[Ctx ctxType[Ctx], Key comparable, Data any] struct {
	handlers map[Key]Handler[Ctx, Key, Data]
	mut      sync.Mutex
}

func (r *Router[Ctx, Key, Data]) gc(ctx tgot.Context) {
	t := time.Now()
	for key, h := range r.handlers {
		if t.Before(h.Timeout()) {
			continue
		}
		delete(r.handlers, key)
		err := h.Cancel(ctx.Child(h.Name()), key)
		if err != nil {
			ctx.Logger().Error("Cancel '%s' ended with an error: %s", h.Name(), err.Error())
		}
	}
}

// Route routes update.
func (r *Router[Ctx, Key, Data]) Route(ctx Ctx, key Key, data Data) {
	r.mut.Lock()
	r.gc(ctx.Ctx())
	h, ok := r.handlers[key]
	r.mut.Unlock()
	if !ok {
		return
	}

	unreg, err := h.Handle(ctx.Child(h.Name()), key, data)
	if err != nil {
		ctx.Logger().Error("handler '%s' ended with an error: %s", h.Name(), err.Error())
	}
	if unreg {
		r.Unreg(key)
	}
}

// Reg registers handler for key.
func (r *Router[Ctx, Key, Data]) Reg(key Key, h Handler[Ctx, Key, Data]) {
	if h == nil || h.Timeout().Before(time.Now()) {
		return
	}
	r.mut.Lock()
	r.handlers[key] = h
	r.mut.Unlock()
}

// Unreg deletes handler associated with the key.
func (r *Router[Ctx, Key, Data]) Unreg(key Key) {
	r.mut.Lock()
	delete(r.handlers, key)
	r.mut.Unlock()
}
