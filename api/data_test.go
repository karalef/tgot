package api_test

import (
	"testing"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

type Embed struct {
	EmbedVal int
}

type Struct struct {
	StructVal int `json:"struct_val"`
}

type TestType struct {
	Str         string
	Int         int
	Media       []tg.InputMedia
	ReplyMarkup *tg.InlineKeyboardMarkup `tg:"reply_markup"`
	Embed
	Force  *int
	Struct Struct
}

func TestMarshal(t *testing.T) {
	testt := TestType{
		Str: "test",
		Int: 1,
		Media: []tg.InputMedia{
			{
				Media: tg.FileBytes("test", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}),
				Data: &tg.InputMediaVideo{
					Thumbnail:  tg.FileBytes("test-thumb", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}),
					HasSpoiler: true,
				},
			},
		},
		ReplyMarkup: &tg.InlineKeyboardMarkup{
			Keyboard: [][]tg.InlineKeyboardButton{{{Text: "test"}}},
		},
		Embed: Embed{1},
		Force: new(int),
		Struct: Struct{
			StructVal: 1,
		},
	}

	d := api.NewDataFrom(testt)

	t.Log("params")
	for k, v := range d.Params {
		t.Logf("%s: %s\n", k, v)
	}
	t.Log()
	t.Log("files")
	for k, v := range d.Upload {
		t.Logf("%s: %s\n", k, v.Name)
	}
}
