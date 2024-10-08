package router

import (
	"testing"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

func TestRouter(t *testing.T) {
	r := NewRouter[tgot.CallbackContext, tgot.MsgID, *tg.CallbackQuery]()
	r.Register(tgot.MsgID{}, nil)
}
