package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

func WithBusiness(ctx BaseContext, id string) *Business {
	return &Business{
		context: ctx.ctx().with(api.NewData().Set("business_connection_id", id)),
		id:      id,
	}
}

// Business provides methods on behalf of a business account.
type Business struct {
	*context
	id string
}

// ReadMessage marks incoming message as read on behalf of a business account.
func (b *Business) ReadMessage(chatID, msgID int64) error {
	return b.method("readBusinessMessage",
		api.NewData().SetInt64("chat_id", chatID).SetInt64("message_id", msgID))
}

// DeleteMessage deletes messages on behalf of a business account.
func (b *Business) DeleteMessages(msgIDs []int64) error {
	return b.method("deleteBusinessMessages",
		api.NewData().SetJSON("message_ids", msgIDs))
}

// SetName changes the first and last name of a managed business account.
func (b *Business) SetName(first, last string) error {
	return b.method("setBusinessAccountName", api.NewData().
		Set("first_name", first).
		Set("last_name", last))
}

// SetUsername changes the username of a managed business account.
func (b *Business) SetUsername(username string) error {
	return b.method("setBusinessAccountUsername",
		api.NewData().Set("username", username))
}

// SetBio changes the bio of a managed business account.
func (b *Business) SetBio(bio string) error {
	return b.method("setBusinessAccountBio", api.NewData().Set("bio", bio))
}

// SetProfilePhoto changes the profile photo of a managed business account.
func (b *Business) SetProfilePhoto(photo tg.InputProfilePhoto, public bool) error {
	d := api.NewData().SetBool("is_public", public)
	d.SetInput("photo", photo)
	return b.method("setBusinessAccountProfilePhoto", d)
}

// RemoveProfilePhoto removes the current profile photo of a managed business
// account.
func (b *Business) RemoveProfilePhoto(public bool) error {
	return b.method("removeBusinessAccountProfilePhoto",
		api.NewData().SetBool("is_public", public))
}

// SetGiftSettings changes the privacy settings pertaining to incoming gifts in
// a managed business account.
func (b *Business) SetGiftSettings(types tg.AcceptedGiftTypes, showButton bool) error {
	return b.method("setBusinessAccountProfilePhoto", api.NewData().
		SetBool("show_gift_button", showButton).
		SetJSON("accepted_gift_types", types))
}

// GetStarBalance returns the amount of Telegram Stars owned by a managed
// business account.
func (b *Business) GetStarBalance() (*tg.StarAmount, error) {
	return method[*tg.StarAmount](b.context, "getBusinessAccountStarBalance")
}

// TransferStars transfers Telegram Stars from the business account balance to
// the bot's balance.
func (b *Business) TransferStars(count uint) error {
	return b.method("transferBusinessAccountStars",
		api.NewData().SetUint("star_count", count))
}

// GetGifts contains parameters for [Business.GetGifts] method.
type GetGifts struct {
	ExcludeUnsaved   bool   `json:"exclude_unsaved"`
	ExcludeSaved     bool   `json:"exclude_saved"`
	ExcludeUnlimited bool   `json:"exclude_unlimited"`
	ExcludeLimited   bool   `json:"exclude_limited"`
	ExcludeUnique    bool   `json:"exclude_unique"`
	SortByPrice      bool   `json:"sort_by_price"`
	Offset           string `json:"offset"`
	Limit            uint   `json:"limit"`
}

// GetGifts returns the gifts received and owned by a managed business account.
func (b *Business) GetGifts(g GetGifts) (*tg.OwnedGifts, error) {
	return method[*tg.OwnedGifts](b.context,
		"getBusinessAccountGifts",
		api.NewDataFrom(g))
}

// ConvertGiftToStars converts a given regular gift to Telegram Stars.
func (b *Business) ConvertGiftToStars(id string) error {
	return b.method("convertGiftToStars",
		api.NewData().Set("owned_gift_id", id))
}

// UpgradeGift upgrades a given regular gift to a unique gift.
func (b *Business) UpgradeGift(id string, stars uint, keepDetails bool) error {
	return b.method("upgradeGift", api.NewData().
		Set("owned_gift_id", id).
		SetBool("keep_original_details", keepDetails).
		SetUint("star_count", stars))
}

// TransferGift transfers an owned unique gift to another user.
func (b *Business) TransferGift(id string, stars uint, newOwnerChatID int64) error {
	return b.method("transferGift", api.NewData().
		Set("owned_gift_id", id).
		SetInt64("new_owner_chat_id", newOwnerChatID).
		SetUint("star_count", stars))
}

// Story is used for story methods.
type Story struct {
	CaptionData
	Areas   []tg.StoryArea       `tg:"areas"`
	Content tg.InputStoryContent `tg:"content"`
}

// PostStory is used for [Business.PostStory] method.
type PostStory struct {
	Story
	ActivePeriod   uint `tg:"active_period"`
	PostToChatPage bool `tg:"post_to_chat_page"`
	ProtectContent bool `tg:"protect_content"`
}

// PostStory posts a story on behalf of a managed business account.
func (b *Business) PostStory(s PostStory) (*tg.Story, error) {
	return method[*tg.Story](b.context, "postStory", api.NewDataFrom(s))
}

// EditStory edits a story previously posted by the bot on behalf of a managed
// business account.
func (b *Business) EditStory(id int, s Story) (*tg.Story, error) {
	return method[*tg.Story](b.context, "editStory",
		api.NewDataFrom(s).SetInt("story_id", id))
}

// DeleteStory deletes a story previously posted by the bot on behalf of a
// managed business account.
func (b *Business) DeleteStory(id int) error {
	return b.method("deleteStory", api.NewData().SetInt("story_id", id))
}
