package tg

// PaidMediaType represents the type of a paid media.
type PaidMediaType string

// all available paid media types.
const (
	PaidMediaTypePreview PaidMediaType = "preview"
	PaidMediaTypePhoto   PaidMediaType = "photo"
	PaidMediaTypeVideo   PaidMediaType = "video"
)

// PaidMedia describes paid media.
type PaidMedia struct {
	Type PaidMediaType `json:"type"`

	*PaidMediaPreview
	*PaidMediaPhoto
	*PaidMediaVideo
}

// PaidMediaPreview represents the paid media that isn't available before the payment.
type PaidMediaPreview struct {
	Width    uint `json:"width"`
	Height   uint `json:"height"`
	Duration uint `json:"duration"`
}

// PaidMediaPhoto represents photo paid media.
type PaidMediaPhoto struct {
	Photo []PhotoSize `json:"photo"`
}

// PaidMediaVideo represents video paid media.
type PaidMediaVideo struct {
	Video Video `json:"video"`
}

// PaidMediaInfo describes the paid media added to a message.
type PaidMediaInfo struct {
	StarCount uint        `json:"star_count"`
	PaidMedia []PaidMedia `json:"paid_media"`
}
