package router

import (
	"testing"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

func TestRouter(t *testing.T) {
	r := NewRouter[tgot.CallbackContext, tgot.MsgSignature, *tg.CallbackQuery]()
	r.Reg(tgot.MsgSignature{}, nil)
}
