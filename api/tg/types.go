package tg

var (
	// False is a pointer to a boolean value of false.
	False = new(bool)

	// True is a pointer to a boolean value of true.
	True = func() *bool {
		t := new(bool)
		*t = true
		return t
	}()
)

// Contact represents a phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int64  `json:"user_id"`
	Vcard       string `json:"vcard"`
}

// DiceEmoji represenst dice emoji.
type DiceEmoji string

// all available animated emojis.
const (
	DiceCube DiceEmoji = "🎲"
	DiceDart DiceEmoji = "🎯"
	DiceBall DiceEmoji = "🏀"
	DiceGoal DiceEmoji = "⚽"
	DiceSlot DiceEmoji = "🎰"
	DiceBowl DiceEmoji = "🎳"
)

// Dice represents an animated emoji that displays a random value.
type Dice struct {
	Emoji DiceEmoji `json:"emoji"`
	Value int       `json:"value"`
}

// Location represents a point on the map.
type Location struct {
	Long               float32  `json:"longitude"`
	Lat                float32  `json:"latitude"`
	HorizontalAccuracy *float32 `json:"horizontal_accuracy,omitempty"`
	LivePeriod         int      `json:"live_period,omitempty"`
	Heading            int      `json:"heading,omitempty"`
	AlertRadius        int      `json:"proximity_alert_radius,omitempty"`
}

// Venue represents a venue.
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id"`
	FoursquareType  string   `json:"foursquare_type"`
	GooglePlaceID   string   `json:"google_place_id"`
	GooglePlaceType string   `json:"google_place_type"`
}

// StickerSet represents a sticker set.
type StickerSet struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	StickerType StickerType `json:"sticker_type"`
	IsAnimated  bool        `json:"is_animated"`
	IsVideo     bool        `json:"is_video"`
	Stickers    []Sticker   `json:"stickers"`
	Thumbnail   *PhotoSize  `json:"thumb"`
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
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}
