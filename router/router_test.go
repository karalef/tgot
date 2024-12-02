package router

import (
	"testing"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

func TestRouter(t *testing.T) {
	r := NewRouter[tgot.Query[tgot.CallbackAnswer], tgot.MessageID, *tg.CallbackQuery]()
	r.Register(tgot.MessageID{}, nil)
}
