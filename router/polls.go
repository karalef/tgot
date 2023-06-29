package router

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// NewPolls makes new poll answers router.
func NewPolls() *Polls {
	return &Polls{NewRouter[tgot.Context, string, *tg.PollAnswer]()}
}

// Polls routes poll answers.
type Polls struct {
	*Router[tgot.Context, string, *tg.PollAnswer]
}

// Route handles poll answers.
//
// It can be used as [Handler.OnPollAnswer].
func (p *Polls) Route(ctx tgot.Context, q *tg.PollAnswer) {
	p.Router.Route(ctx, q.PollID, q)
}

// PollHandler represents poll answers handler.
type PollHandler = Handler[tgot.Context, string, *tg.PollAnswer]

// Register registers poll answers handler.
func (p *Polls) Register(pollID string, h PollHandler) {
	if pollID != "" {
		p.Router.Register(pollID, h)
	}
}

// RegisterOneTime registers poll answers handler which will be unregistered after first call.
func (p *Polls) RegisterOneTime(pollID string, h PollHandler) {
	if pollID != "" {
		p.Router.RegisterOneTime(pollID, h)
	}
}
