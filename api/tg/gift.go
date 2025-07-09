package tg

import "github.com/karalef/tgot/api/internal/oneof"

// Gift represents a gift that can be sent by the bot.
type Gift struct {
	ID               string  `json:"id"`
	Sticker          Sticker `json:"sticker"`
	StarCount        uint    `json:"star_count"`
	UpgradeStarCount uint    `json:"upgrade_star_count"`
	Total            uint    `json:"total_count"`
	Remaining        uint    `json:"remaining_count"`
}

// Gifts represent a list of gifts.
type Gifts struct {
	Gifts []Gift `json:"gifts"`
}

// OwnedGifts contains the list of gifts received and owned by a user or a chat.
type OwnedGifts struct {
	Count      uint        `json:"total_count"`
	Gifts      []OwnedGift `json:"gifts"`
	NextOffset string      `json:"next_offset"`
}

// AcceptedGiftTypes describes the types of gifts that can be gifted to a user
// or a chat.
type AcceptedGiftTypes struct {
	Unlimited bool `json:"unlimited_gifts"`
	Limited   bool `json:"limited_gifts"`
	Unique    bool `json:"unique_gifts"`
	Premium   bool `json:"premium_subscription"`
}

// StarAmount describes an amount of Telegram Stars.
type StarAmount struct {
	Amount         uint `json:"amount"`
	NanostarAmount uint `json:"nanostar_amount"`
}

// UniqueGiftModel describes the model of a unique gift.
type UniqueGiftModel struct {
	Name          string  `json:"name"`
	Sticker       Sticker `json:"sticker"`
	RarityPerMile uint    `json:"rarity_per_mille"`
}

// UniqueGiftSymbol describes the symbol shown on the pattern of a unique gift.
type UniqueGiftSymbol struct {
	Name          string  `json:"name"`
	Sticker       Sticker `json:"sticker"`
	RarityPerMile uint    `json:"rarity_per_mille"`
}

// UniqueGiftBackdropColors describes the colors of the backdrop of a unique gift.
type UniqueGiftBackdropColors struct {
	Center uint `json:"center_color"`
	Edge   uint `json:"edge_color"`
	Symbol uint `json:"symbol_color"`
	Text   uint `json:"text_color"`
}

// UniqueGiftBackdrop describes the backdrop of a unique gift.
type UniqueGiftBackdrop struct {
	Name          string                   `json:"name"`
	Colors        UniqueGiftBackdropColors `json:"colors"`
	RarityPerMile uint                     `json:"rarity_per_mille"`
}

// UniqueGift describes a unique gift that was upgraded from a regular gift.
type UniqueGift struct {
	BaseName string             `json:"base_name"`
	Name     string             `json:"name"`
	Number   uint               `json:"number"`
	Model    UniqueGiftModel    `json:"model"`
	Symbol   UniqueGiftSymbol   `json:"symbol"`
	Backdrop UniqueGiftBackdrop `json:"backdrop"`
}

// OwnedGiftType represents the type of owned gift.
type OwnedGiftType string

// all available owned gift types..
const (
	OwnedGiftTypeRegular OwnedGiftType = "regular"
	OwnedGiftTypeUnique  OwnedGiftType = "unique"
)

var ownedGiftTypes = oneof.NewMap[OwnedGiftType](
	OwnedGiftRegular{},
	OwnedGiftUnique{},
)

func (OwnedGiftType) TypeFor(t OwnedGiftType) oneof.Type {
	return ownedGiftTypes.TypeFor(t)
}

// OwnedGift describes a gift received and owned by a user or a chat.
type OwnedGift = oneof.Object[OwnedGiftType, oneof.IDTypeType]

// OwnedGiftRegular describes a regular gift owned by a user or a chat.
type OwnedGiftRegular struct {
	Gift                    Gift            `json:"gift"`
	OwnedID                 string          `json:"owned_gift_id"`
	Sender                  *User           `json:"sender_user"`
	Date                    int64           `json:"send_date"`
	Text                    string          `json:"text"`
	Entities                []MessageEntity `json:"entities"`
	IsPrivate               bool            `json:"is_private"`
	IsSaved                 bool            `json:"is_saved"`
	Upgradable              bool            `json:"can_be_upgraded"`
	WasRefunded             bool            `json:"was_refunded"`
	ConvertStarCount        uint            `json:"convert_star_count"`
	PrepaidUpgradeStarCount uint            `json:"prepaid_upgrade_star_count"`
}

func (OwnedGiftRegular) Type() OwnedGiftType { return OwnedGiftTypeRegular }

// OwnedGiftUnique describes a unique gift received and owned by a user or a chat.
type OwnedGiftUnique struct {
	Gift              UniqueGift `json:"gift"`
	OwnedID           string     `json:"owned_gift_id"`
	Sender            *User      `json:"sender_user"`
	Date              int64      `json:"send_date"`
	IsSaved           bool       `json:"is_saved"`
	Transferable      bool       `json:"can_be_transferred"`
	TransferStarCount uint       `json:"transfer_star_count"`
	NextTransferDate  int64      `json:"next_transfer_date"`
}

func (OwnedGiftUnique) Type() OwnedGiftType { return OwnedGiftTypeUnique }

// GiftInfo describes a service message about a regular gift that was sent or received.
type GiftInfo struct {
	Gift                    Gift            `json:"gift"`
	OwnedID                 string          `json:"owned_gift_id"`
	ConvertStarCount        uint            `json:"convert_star_count"`
	PrepaidUpgradeStarCount uint            `json:"prepaid_upgrade_star_count"`
	Upgradable              bool            `json:"can_be_upgraded"`
	Text                    string          `json:"text"`
	Entities                []MessageEntity `json:"entities"`
	IsPrivate               bool            `json:"is_private"`
}

// GiftOrigin represents the origin of a gift.
type GiftOrigin string

// all available gift origins.
const (
	GiftOriginUpgrade  = "upgrade"
	GiftOriginTransfer = "transfer"
	GiftOriginResale   = "resale"
)

// UniqueGiftInfo describes a service message about a unique gift that was sent or received.
type UniqueGiftInfo struct {
	Gift                UniqueGift `json:"gift"`
	Origin              GiftOrigin `json:"origin"`
	LastResaleStarCount uint       `json:"last_resale_star_count,omitempty"`
	OwnedID             string     `json:"owned_gift_id,omitempty"`
	TransferStarCount   uint       `json:"transfer_star_count,omitempty"`
	NextTransferDate    int64      `json:"next_transfer_date,omitempty"`
}
