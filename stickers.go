package bot

import "tghwbot/bot/tg"

// GetStickerSet returns a sticker set.
func (c Context) GetStickerSet(name string) (*tg.StickerSet, error) {
	p := params{}.set("name", name)
	return api[*tg.StickerSet](c, "getStickerSet", p)
}

// GetCustomEmojiStickers returns information about custom emoji stickers by their identifiers.
func (c Context) GetCustomEmojiStickers(ids ...string) ([]tg.Sticker, error) {
	p := params{}
	p.setJSON("custom_emoji_ids", ids)
	return api[[]tg.Sticker](c, "getCustomEmojiStickers", p)
}

// UploadStickerFile uploads a .PNG file with a sticker for later use
// in createNewStickerSet and addStickerToSet methods.
func (c Context) UploadStickerFile(userID int64, pngSticker *tg.InputFile) (*tg.File, error) {
	p := params{}
	p.setInt64("user_id", userID)
	return api[*tg.File](c, "uploadStickerFile", p, file{
		field: "png_sticker", FileSignature: pngSticker,
	})
}

// NewSticker contains common parameters for creating a sticker set
// and adding a sticker to a set.
type NewSticker struct {
	UserID       int64
	SetName      string
	PngSticker   tg.FileSignature
	TgsSticker   *tg.InputFile
	WebmSticker  *tg.InputFile
	Emojis       string
	MaskPosition *tg.MaskPosition
}

func (n NewSticker) data(p params) (f file) {
	p.setInt64("user_id", n.UserID)
	p.set("name", n.SetName)
	p.set("emojis", n.Emojis)
	p.setJSON("mask_position", n.MaskPosition)
	switch {
	case n.PngSticker != nil:
		f.field = "png_sticker"
		f.FileSignature = n.PngSticker
	case n.TgsSticker != nil:
		f.field = "tgs_sticker"
		f.FileSignature = n.TgsSticker
	case n.WebmSticker != nil:
		f.field = "webm_sticker"
		f.FileSignature = n.WebmSticker
	}
	return f
}

// NewStickerSet contains parameters for creating a sticker set.
type NewStickerSet struct {
	NewSticker
	Title string
	Type  tg.StickerType
}

// CreateNewStickerSet creates a new sticker set owned by a user.
func (c Context) CreateNewStickerSet(newSet NewStickerSet) error {
	p := params{}
	p.set("title", newSet.Title)
	p.set("sticker_type", string(newSet.Type))
	f := newSet.data(p)
	return c.api("createNewStickerSet", p, f)
}

// AddSticker adds a new sticker to a set created by the bot.
func (c Context) AddSticker(sticker NewSticker) error {
	p := params{}
	f := sticker.data(p)
	return c.api("addStickerToSet", p, f)
}

// SetStickerPosition moves a sticker in a set created by the bot to a specific position.
func (c Context) SetStickerPosition(stickerFileID string, position int) error {
	p := params{}
	p.set("sticker", stickerFileID)
	p.setInt("position", position)
	return c.api("setStickerPositionInSet", p)
}

// DeleteSticker deletes a sticker from a set created by the bot.
func (c Context) DeleteSticker(stickerFileID string) error {
	p := params{}
	p.set("sticker", stickerFileID)
	return c.api("deleteStickerFromSet", p)
}

// SetStickerSetThumb sets the thumbnail of a sticker set.
func (c Context) SetStickerSetThumb(setName string, userID int64, thumb tg.FileSignature) error {
	p := params{}
	p.set("name", setName)
	p.setInt64("user_id", userID)
	return c.api("setStickerSetThumb", p, file{
		field: "thumb", FileSignature: thumb,
	})
}
