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

var paidMediaTypes = oneof.NewMap[PaidMediaType](
	PaidMediaPreview{},
	PaidMediaPhoto{},
	PaidMediaVideo{},
)

func (PaidMediaType) TypeFor(t PaidMediaType) oneof.Type {
	return paidMediaTypes.TypeFor(t)
}

// PaidMedia describes paid media.
type PaidMedia = oneof.Object[PaidMediaType, oneof.IDTypeType]

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

// StoryArea describes a clickable area on a story media.
type StoryArea struct {
	Type     StoryAreaType     `json:"type"`
	Position StoryAreaPosition `json:"position"`
}

// StoryAreaPosition describes the position of a clickable area within a story.
type StoryAreaPosition struct {
	XPercentage            float64 `json:"x_percentage"`
	YPercentage            float64 `json:"y_percentage"`
	WidthPercentage        float64 `json:"width_percentage"`
	HeightPercentage       float64 `json:"height_percentage"`
	RotationAngle          float64 `json:"rotation_angle"`
	CornerRadiusPercentage float64 `json:"corner_radius_percentage"`
}

// LocationAddress describes the physical address of a location.
type LocationAddress struct {
	CountryCode string `json:"country_code"` // two-letter ISO 3166-1 alpha-2 country code
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	Street      string `json:"street,omitempty"`
}

// StoryAreaTypeID represents the ID of story area type.
type StoryAreaTypeID string

// all available story types.
const (
	StoryAreaTypeIDLocation          StoryAreaTypeID = "location"
	StoryAreaTypeIDSuggestedReaction StoryAreaTypeID = "suggested_reaction"
	StoryAreaTypeIDLink              StoryAreaTypeID = "link"
	StoryAreaTypeIDWeather           StoryAreaTypeID = "weather"
	StoryAreaTypeIDUniqueGift        StoryAreaTypeID = "unique_gift"
)

var storyAreaTypes = oneof.NewMap[StoryAreaTypeID](
	StoryAreaTypeLocation{},
	StoryAreaTypeSuggestedReaction{},
	StoryAreaTypeLink{},
	StoryAreaTypeWeather{},
	StoryAreaTypeUniqueGift{},
)

func (StoryAreaTypeID) TypeFor(t StoryAreaTypeID) oneof.Type {
	return storyAreaTypes.TypeFor(t)
}

// StoryAreaType describes the type of a clickable area on a story.
type StoryAreaType = oneof.Object[StoryAreaTypeID, oneof.IDTypeType]

// StoryAreaTypeLocation describes a story area pointing to a location.
type StoryAreaTypeLocation struct {
	Latitude  float32          `json:"latitude"`
	Longitude float32          `json:"longitude"`
	Address   *LocationAddress `json:"address"`
}

func (StoryAreaTypeLocation) Type() StoryAreaTypeID { return StoryAreaTypeIDLocation }

// StoryAreaTypeSuggestedReaction describes a story area pointing to a suggested
// reaction.
type StoryAreaTypeSuggestedReaction struct {
	RectionType ReactionType `json:"reaction_type"`
	IsDark      bool         `json:"is_dark"`
	IsFlipped   bool         `json:"is_flipped"`
}

func (StoryAreaTypeSuggestedReaction) Type() StoryAreaTypeID { return StoryAreaTypeIDLocation }

// StoryAreaTypeLink describes a story area pointing to an HTTP or tg:// link.
type StoryAreaTypeLink struct {
	URL string `json:"url"`
}

func (StoryAreaTypeLink) Type() StoryAreaTypeID { return StoryAreaTypeIDLink }

// StoryAreaTypeWeather describes a story area containing weather information.
type StoryAreaTypeWeather struct {
	Temperature     float64 `json:"temperature"`
	Emoji           string  `json:"emoji"`
	BackgroundColor uint    `json:"background_color"`
}

func (StoryAreaTypeWeather) Type() StoryAreaTypeID { return StoryAreaTypeIDWeather }

// StoryAreaTypeUniqueGift describes a story area pointing to a unique gift.
type StoryAreaTypeUniqueGift struct {
	Name string `json:"name"`
}

func (StoryAreaTypeUniqueGift) Type() StoryAreaTypeID { return StoryAreaTypeIDUniqueGift }
