package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// WithSticker returns StickerFile with provided stickerFileID.
func WithSticker(ctx BaseContext, stickerFileID string) StickerFile {
	return StickerFile{
		ctx: ctx.ctx().with(api.NewData().Set("sticker", stickerFileID)),
		id:  stickerFileID,
	}
}

// StickerFile provides methods for working with stickers.
type StickerFile struct {
	ctx *context
	id  string
}

// SetPosition moves a sticker in a set created by the bot to a specific position.
func (s StickerFile) SetPosition(position int) error {
	return s.ctx.method("setStickerPositionInSet", api.NewData().SetInt("position", position))
}

// Delete deletes a sticker from a set created by the bot.
func (s StickerFile) Delete() error {
	return s.ctx.method("deleteStickerFromSet")
}

// SetEmojiList changes the list of emoji assigned to a regular or custom emoji sticker.
func (s StickerFile) SetEmojiList(emojiList []string) error {
	return s.ctx.method("setStickerEmojiList", api.NewData().SetJSON("emoji_list", emojiList))
}

// SetStickerKeywords changes search keywords assigned to a regular or custom emoji sticker.
func (s StickerFile) SetKeywords(keywords []string) error {
	return s.ctx.method("setStickerKeywords", api.NewData().SetJSON("keywords", keywords))
}

// SetMaskPosition changes the mask position of a mask sticker.
func (s StickerFile) SetMaskPosition(maskPos tg.MaskPosition) error {
	return s.ctx.method("setStickerMaskPosition", api.NewData().SetJSON("mask_position", maskPos))
}

// WithStickerSet returns StickerSet with provided setName.
func WithStickerSet(ctx BaseContext, setName string) StickerSet {
	return StickerSet{
		ctx:  ctx.ctx().with(api.NewData().Set("name", setName)),
		name: setName,
	}
}

// StickerSet provides methods for working with sticker sets.
type StickerSet struct {
	ctx  *context
	name string
}

// Get returns a sticker set.
func (s StickerSet) Get() (*tg.StickerSet, error) {
	return method[*tg.StickerSet](s.ctx, "getStickerSet")
}

// SetTitle sets the title of a created sticker set.
func (s StickerSet) SetTitle(title string) error {
	return s.ctx.method("setStickerSetTitle", api.NewData().Set("title", title))
}

// SetCustomEmojiThumbnail sets the thumbnail of a custom emoji sticker set.
func (s StickerSet) SetCustomEmojiThumbnail(customEmojiID string) error {
	d := api.NewData().Set("custom_emoji_id", customEmojiID)
	return s.ctx.method("setCustomEmojiStickerSetThumbnail", d)
}

// Delete deletes a sticker set that was created by the bot.
func (s StickerSet) Delete() error {
	return s.ctx.method("deleteStickerSet")
}

// SetStickerSetThumbnail sets the thumbnail of a sticker set.
func (s StickerSet) SetThumbnail(userID tg.ID, thumb tg.Inputtable, format tg.StickerFormat) error {
	d := api.NewData().SetID("user_id", userID)
	d.SetFile("thumbnail", thumb)
	d.Set("format", string(format))
	return s.ctx.method("setStickerSetThumbnail", d)
}

// Add adds a new sticker to a set created by the bot.
func (s StickerSet) Add(userID tg.ID, sticker tg.InputSticker) error {
	d := api.NewData().SetID("user_id", userID)
	d.SetInput("sticker", sticker)
	return s.ctx.method("addStickerToSet", d)
}

// Replace replace an existing sticker in a sticker set with a new one.
func (s StickerSet) Replace(old string, userID tg.ID, sticker tg.InputSticker) error {
	d := api.NewData().Set("old_sticker", old).SetID("user_id", userID)
	d.SetInput("sticker", sticker)
	return s.ctx.method("replaceStickerInSet", d)
}

// NewStickerSet contains parameters for creating a sticker set.
type NewStickerSet struct {
	UserID          tg.ID             `tg:"user_id"`
	Title           string            `tg:"title"`
	Stickers        []tg.InputSticker `tg:"stickers"`
	Type            tg.StickerType    `tg:"sticker_type"`
	NeedsRepainting bool              `tg:"needs_repainting"`
}

// Create creates a new sticker set owned by a user.
func (s StickerSet) Create(newSet NewStickerSet) error {
	return s.ctx.method("createNewStickerSet", api.NewDataFrom(newSet))
}
