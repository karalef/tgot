package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

// GetStickerSet returns a sticker set.
func (c Context) GetStickerSet(name string) (*tg.StickerSet, error) {
	d := api.NewData().Set("name", name)
	return method[*tg.StickerSet](c, "getStickerSet", d)
}

// GetCustomEmojiStickers returns information about custom emoji stickers by their identifiers.
func (c Context) GetCustomEmojiStickers(ids ...string) ([]tg.Sticker, error) {
	d := api.NewData().SetJSON("custom_emoji_ids", ids)
	return method[[]tg.Sticker](c, "getCustomEmojiStickers", d)
}

// UploadStickerFile uploads a .PNG file with a sticker for later use
// in createNewStickerSet and addStickerToSet methods.
func (c Context) UploadStickerFile(userID int64, pngSticker *tg.InputFile) (*tg.File, error) {
	d := api.NewData()
	d.SetInt64("user_id", userID)
	d.SetFile("png_sticker", pngSticker, nil)
	return method[*tg.File](c, "uploadStickerFile", d)
}

// NewSticker contains common parameters for creating a sticker set
// and adding a sticker to a set.
type NewSticker struct {
	UserID       int64
	SetName      string
	PngSticker   tg.Inputtable
	TgsSticker   *tg.InputFile
	WebmSticker  *tg.InputFile
	Emojis       string
	MaskPosition *tg.MaskPosition
}

func (n NewSticker) data() api.Data {
	d := api.NewData()
	d.SetInt64("user_id", n.UserID)
	d.Set("name", n.SetName)
	d.Set("emojis", n.Emojis)
	d.SetJSON("mask_position", n.MaskPosition)
	switch {
	case n.PngSticker != nil:
		d.SetFile("png_sticker", n.PngSticker, nil)
	case n.TgsSticker != nil:
		d.SetFile("tgs_sticker", n.TgsSticker, nil)
	case n.WebmSticker != nil:
		d.SetFile("webm_sticker", n.WebmSticker, nil)
	}
	return d
}

// NewStickerSet contains parameters for creating a sticker set.
type NewStickerSet struct {
	NewSticker
	Title string
	Type  tg.StickerType
}

// CreateNewStickerSet creates a new sticker set owned by a user.
func (c Context) CreateNewStickerSet(newSet NewStickerSet) error {
	d := newSet.NewSticker.data()
	d.Set("title", newSet.Title)
	d.Set("sticker_type", string(newSet.Type))
	return c.method("createNewStickerSet", d)
}

// AddSticker adds a new sticker to a set created by the bot.
func (c Context) AddSticker(sticker NewSticker) error {
	return c.method("addStickerToSet", sticker.data())
}

// SetStickerPosition moves a sticker in a set created by the bot to a specific position.
func (c Context) SetStickerPosition(stickerFileID string, position int) error {
	d := api.NewData()
	d.Set("sticker", stickerFileID)
	d.SetInt("position", position)
	return c.method("setStickerPositionInSet", d)
}

// DeleteSticker deletes a sticker from a set created by the bot.
func (c Context) DeleteSticker(stickerFileID string) error {
	d := api.NewData().Set("sticker", stickerFileID)
	return c.method("deleteStickerFromSet", d)
}

// SetStickerSetThumb sets the thumbnail of a sticker set.
func (c Context) SetStickerSetThumb(setName string, userID int64, thumb tg.Inputtable) error {
	d := api.NewData()
	d.Set("name", setName)
	d.SetInt64("user_id", userID)
	d.AddFile("thumb", thumb)
	return c.method("setStickerSetThumb", d)
}
