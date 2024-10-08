package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
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

func prepareInputStickers(d *api.Data, stickers ...tg.InputSticker) {
	for _, s := range stickers {
		d.AddAttach(s.Sticker)
	}
}

// UploadStickerFile uploads a .PNG file with a sticker for later use
// in createNewStickerSet and addStickerToSet methods.
func (c Context) UploadStickerFile(userID int64, sticker *tg.InputFile, format tg.StickerFormat) (*tg.File, error) {
	d := api.NewData()
	d.SetInt64("user_id", userID)
	d.AddFile("sticker", sticker)
	d.Set("sticker_format", string(format))
	return method[*tg.File](c, "uploadStickerFile", d)
}

// NewStickerSet contains parameters for creating a sticker set.
type NewStickerSet struct {
	UserID          int64
	Name            string
	Title           string
	Stickers        []tg.InputSticker
	Type            tg.StickerType
	NeedsRepainting bool
}

// CreateNewStickerSet creates a new sticker set owned by a user.
func (c Context) CreateNewStickerSet(newSet NewStickerSet) error {
	d := api.NewData()
	prepareInputStickers(api.NewData(), newSet.Stickers...)
	d.SetInt64("user_id", newSet.UserID)
	d.Set("name", newSet.Name)
	d.Set("title", newSet.Title)
	d.SetJSON("stickers", newSet.Stickers)
	d.Set("sticker_type", string(newSet.Type))
	d.SetBool("needs_repainting", newSet.NeedsRepainting)
	return c.method("createNewStickerSet", d)
}

// NewSticker contains parameters for adding a sticker to a set.
type NewSticker struct {
	UserID  int64
	SetName string
	Sticker tg.InputSticker
}

func (n NewSticker) data(d *api.Data) {
	prepareInputStickers(d, n.Sticker)
	d.SetInt64("user_id", n.UserID)
	d.Set("name", n.SetName)
	d.SetJSON("sticker", n.Sticker)
}

// AddSticker adds a new sticker to a set created by the bot.
func (c Context) AddSticker(sticker NewSticker) error {
	d := api.NewData()
	sticker.data(d)
	return c.method("addStickerToSet", d)
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

// SetStickerEmojiList changes the list of emoji assigned to a regular or custom emoji sticker.
func (c Context) SetStickerEmojiList(stickerFileID string, emojiList []string) error {
	d := api.NewData()
	d.Set("sticker", stickerFileID)
	d.SetJSON("emoji_list", emojiList)
	return c.method("setStickerEmojiList", d)
}

// SetStickerKeywords changes search keywords assigned to a regular or custom emoji sticker.
func (c Context) SetStickerKeywords(stickerFileID string, keywords []string) error {
	d := api.NewData()
	d.Set("sticker", stickerFileID)
	d.SetJSON("keywords", keywords)
	return c.method("setStickerKeywords", d)
}

// SetStickerMaskPosition changes the mask position of a mask sticker.
func (c Context) SetStickerMaskPosition(stickerFileID string, maskPos tg.MaskPosition) error {
	d := api.NewData()
	d.Set("sticker", stickerFileID)
	d.SetJSON("mask_position", maskPos)
	return c.method("setStickerMaskPosition", d)
}

// SetStickerSetTitle sets the title of a created sticker set.
func (c Context) SetStickerSetTitle(setName string, title string) error {
	d := api.NewData()
	d.Set("name", setName)
	d.Set("title", title)
	return c.method("setStickerSetTitle", d)
}

// SetStickerSetThumbnail sets the thumbnail of a sticker set.
func (c Context) SetStickerSetThumbnail(setName string, userID int64, thumb tg.Inputtable, format tg.StickerFormat) error {
	d := api.NewData()
	d.Set("name", setName)
	d.SetInt64("user_id", userID)
	d.AddFile("thumbnail", thumb)
	d.Set("format", string(format))
	return c.method("setStickerSetThumbnail", d)
}

// SetCustomEmojiStickerSetThumbnail sets the thumbnail of a custom emoji sticker set.
func (c Context) SetCustomEmojiStickerSetThumbnail(setName string, customEmojiID string) error {
	d := api.NewData()
	d.Set("name", setName)
	d.Set("custom_emoji_id", customEmojiID)
	return c.method("setCustomEmojiStickerSetThumbnail", d)
}

// ReplaceStickerInSet replace an existing sticker in a sticker set with a new one.
func (c Context) ReplaceStickerInSet(old string, sticker NewSticker) error {
	d := api.NewData()
	d.Set("old_sticker", old)
	sticker.data(d)
	return c.method("replaceStickerInSet", d)
}

// DeleteStickerSet deletes a sticker set that was created by the bot.
func (c Context) DeleteStickerSet(setName string) error {
	return c.method("deleteStickerSet", api.NewData().Set("name", setName))
}
