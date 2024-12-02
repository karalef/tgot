package tg

import "github.com/karalef/tgot/api/internal/oneof"

// PaidMediaType represents the type of a paid media.
type PaidMediaType string

// all available paid media types.
const (
	PaidMediaTypePreview PaidMediaType = "preview"
	PaidMediaTypePhoto   PaidMediaType = "photo"
	PaidMediaTypeVideo   PaidMediaType = "video"
)

var paidMediaTypes = oneof.Map[PaidMediaType]{
	PaidMediaTypePreview: PaidMediaPreview{},
	PaidMediaTypePhoto:   PaidMediaPhoto{},
	PaidMediaTypeVideo:   PaidMediaVideo{},
}

type oneOfPaidMedia struct{}

func (oneOfPaidMedia) New(t PaidMediaType) (oneof.Value[PaidMediaType], bool) {
	return paidMediaTypes.New(t)
}

// PaidMedia describes paid media.
type PaidMedia = oneof.Object[PaidMediaType, oneOfPaidMedia]

// PaidMediaPreview represents the paid media that isn't available before the payment.
type PaidMediaPreview struct {
	Width    uint `json:"width"`
	Height   uint `json:"height"`
	Duration uint `json:"duration"`
}

func (PaidMediaPreview) Type() PaidMediaType { return PaidMediaTypePreview }

// PaidMediaPhoto represents photo paid media.
type PaidMediaPhoto struct {
	Photo []PhotoSize `json:"photo"`
}

func (PaidMediaPhoto) Type() PaidMediaType { return PaidMediaTypePhoto }

// PaidMediaVideo represents video paid media.
type PaidMediaVideo struct {
	Video Video `json:"video"`
}

func (PaidMediaVideo) Type() PaidMediaType { return PaidMediaTypeVideo }

// PaidMediaInfo describes the paid media added to a message.
type PaidMediaInfo struct {
	StarCount uint        `json:"star_count"`
	PaidMedia []PaidMedia `json:"paid_media"`
}

// StickerSet represents a sticker set.
type StickerSet struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	StickerType StickerType `json:"sticker_type"`
	Stickers    []Sticker   `json:"stickers"`
	Thumbnail   *PhotoSize  `json:"thumbnail"`
}

// StickerType is a Sticker type.
type StickerType string

// all available sticker types.
const (
	StickerRegular     StickerType = "regular"
	StickerMask        StickerType = "mask"
	StickerCustomEmoji StickerType = "custom_emoji"
)

// MaskPosition describes the position on faces where a mask should be placed by default.
type MaskPosition struct {
	Point  MaskPositionPoint `json:"point"`
	XShift float32           `json:"x_shift"`
	YShift float32           `json:"y_shift"`
	Scale  float32           `json:"scale"`
}

// MaskPositionPoint is a MaskPosition point.
type MaskPositionPoint string

// all available mask position points.
const (
	MaskPositionPointForehead MaskPositionPoint = "forehead"
	MaskPositionPointEyes     MaskPositionPoint = "eyes"
	MaskPositionPointMouth    MaskPositionPoint = "mouth"
	MaskPositionPointChin     MaskPositionPoint = "chin"
)
