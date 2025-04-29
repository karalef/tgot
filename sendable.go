package tgot

import (
	"github.com/karalef/tgot/api/tg"
)

// Sendable interface for Chat.Send.
type Sendable interface {
	sendMethod() string
}

// CaptionData represents caption with entities and parse mode.
type CaptionData struct {
	Caption   string             `tg:"caption"`
	ParseMode tg.ParseMode       `tg:"parse_mode"`
	Entities  []tg.MessageEntity `tg:"caption_entities"`
}

// SendOptions cointains common send* parameters.
type SendOptions struct {
	MessageEffectID     string             `tg:"message_effect_id"`
	DisableNotification bool               `tg:"disable_notification"`
	AllowPaidBroadcast  bool               `tg:"allow_paid_broadcast"`
	ProtectContent      bool               `tg:"protect_content"`
	ReplyParameters     tg.ReplyParameters `tg:"reply_parameters"`
}

// NewText makes a new text message.
func NewText(text string) Text { return Text{Text: text} }

var _ Sendable = Text{}

// Text contains information about the message to be sent.
type Text struct {
	Text               string                `tg:"text"`
	ParseMode          tg.ParseMode          `tg:"parse_mode"`
	Entities           []tg.MessageEntity    `tg:"entities"`
	LinkPreviewOptions tg.LinkPreviewOptions `tg:"link_preview_options"`
	ReplyMarkup        tg.ReplyMarkup        `tg:"reply_markup"`
}

func (Text) sendMethod() string { return "sendMessage" }

// NewPhoto makes a new photo.
func NewPhoto(photo tg.Inputtable) Photo { return Photo{Photo: photo} }

var _ Sendable = Photo{}

// Photo contains information about the photo to be sent.
type Photo struct {
	Photo tg.Inputtable `tg:"photo"`
	CaptionData
	HasSpoiler            bool           `tg:"has_spoiler"`
	ReplyMarkup           tg.ReplyMarkup `tg:"reply_markup"`
	ShowCaptionAboveMedia bool           `tg:"show_caption_above_media"`
}

func (Photo) sendMethod() string { return "sendPhoto" }

var _ Sendable = Audio{}

// Audio contains information about the audio to be sent.
type Audio struct {
	Audio     tg.Inputtable `tg:"audio"`
	Thumbnail tg.Inputtable `tg:"thumbnail"`
	CaptionData
	Duration    int            `tg:"duration"`
	Performer   string         `tg:"performer"`
	Title       string         `tg:"title"`
	ReplyMarkup tg.ReplyMarkup `tg:"reply_markup"`
}

func (Audio) sendMethod() string { return "sendAudio" }

var _ Sendable = Document{}

// Document contains information about the document to be sent.
type Document struct {
	Document  tg.Inputtable `tg:"document"`
	Thumbnail tg.Inputtable `tg:"thumbnail"`
	CaptionData
	DisableTypeDetection bool           `tg:"disable_content_type_detection"`
	ReplyMarkup          tg.ReplyMarkup `tg:"reply_markup"`
}

func (Document) sendMethod() string { return "sendDocument" }

var _ Sendable = Video{}

// Video contains information about the video to be sent.
type Video struct {
	Video     tg.Inputtable `tg:"video"`
	Thumbnail tg.Inputtable `tg:"thumbnail"`
	CaptionData
	Duration              int             `tg:"duration"`
	Width                 int             `tg:"width"`
	Height                int             `tg:"height"`
	Cover                 []tg.Inputtable `tg:"cover"`
	Start                 int64           `tg:"start_timestamp"`
	HasSpoiler            bool            `tg:"has_spoiler"`
	SupportsStreaming     bool            `tg:"supports_streaming"`
	ReplyMarkup           tg.ReplyMarkup  `tg:"reply_markup"`
	ShowCaptionAboveMedia bool            `tg:"show_caption_above_media"`
}

func (Video) sendMethod() string { return "sendVideo" }

var _ Sendable = Animation{}

// Animation contains information about the animation to be sent.
type Animation struct {
	Animation tg.Inputtable `tg:"animation"`
	Thumbnail tg.Inputtable `tg:"thumbnail"`
	CaptionData
	Duration              int            `tg:"duration"`
	Width                 int            `tg:"width"`
	Height                int            `tg:"height"`
	HasSpoiler            bool           `tg:"has_spoiler"`
	ReplyMarkup           tg.ReplyMarkup `tg:"reply_markup"`
	ShowCaptionAboveMedia bool           `tg:"show_caption_above_media"`
}

func (Animation) sendMethod() string { return "sendAnimation" }

var _ Sendable = Voice{}

// Voice contains information about the voice to be sent.
type Voice struct {
	Voice tg.Inputtable `tg:"voice"`
	CaptionData
	Duration    int            `tg:"duration"`
	ReplyMarkup tg.ReplyMarkup `tg:"reply_markup"`
}

func (Voice) sendMethod() string { return "sendVoice" }

var _ Sendable = VideoNote{}

// VideoNote contains information about the video note to be sent.
type VideoNote struct {
	VideoNote   tg.Inputtable  `tg:"video_note"`
	Thumbnail   tg.Inputtable  `tg:"thumbnail"`
	Duration    int            `tg:"duration"`
	Length      int            `tg:"length"`
	ReplyMarkup tg.ReplyMarkup `tg:"reply_markup"`
}

func (VideoNote) sendMethod() string { return "sendVideoNote" }

var _ Sendable = PaidMedia{}

// PaidMedia contains information about the paid media to be sent.
type PaidMedia struct {
	StarCount uint                `tg:"star_count"`
	Media     []tg.InputPaidMedia `tg:"media"`
	Payload   string              `tg:"payload"`
	CaptionData
	ShowCaptionAboveMedia bool `tg:"show_caption_above_media"`
}

func (PaidMedia) sendMethod() string { return "sendPaidMedia" }

var _ Sendable = Location{}

// Location contains information about the location to be sent.
type Location struct {
	tg.Location `tg:",json"`
	ReplyMarkup tg.ReplyMarkup `tg:"reply_markup"`
}

func (Location) sendMethod() string { return "sendLocation" }

var _ Sendable = Venue{}

// Venue contains information about the venue to be sent.
type Venue struct {
	Lat             float32        `tg:"latitude"`
	Long            float32        `tg:"longitude"`
	Title           string         `tg:"title"`
	Address         string         `tg:"address"`
	FoursquareID    string         `tg:"foursquare_id"`
	FoursquareType  string         `tg:"foursquare_type"`
	GooglePlaceID   string         `tg:"google_place_id"`
	GooglePlaceType string         `tg:"google_place_type"`
	ReplyMarkup     tg.ReplyMarkup `tg:"reply_markup"`
}

func (Venue) sendMethod() string { return "sendVenue" }

var _ Sendable = Contact{}

// Contact contains information about the contact to be sent.
type Contact struct {
	PhoneNumber string         `tg:"phone_number"`
	FirstName   string         `tg:"first_name"`
	LastName    string         `tg:"last_name"`
	Vcard       string         `tg:"vcard"`
	ReplyMarkup tg.ReplyMarkup `tg:"reply_markup"`
}

func (Contact) sendMethod() string { return "sendContact" }

var _ Sendable = Poll{}

// Poll contains information about the poll to be sent.
type Poll struct {
	Question             string               `tg:"question"`
	ParseMode            tg.ParseMode         `tg:"parse_mode"`
	Entities             []tg.MessageEntity   `tg:"entities"`
	Options              []tg.InputPollOption `tg:"options"`
	IsAnonymous          bool                 `tg:"is_anonymous"`
	Type                 tg.PollType          `tg:"type"`
	MultipleAnswers      bool                 `tg:"allows_multiple_answers"`
	CorrectOption        int                  `tg:"correct_option_id"`
	Explanation          string               `tg:"explanation"`
	ExplanationParseMode tg.ParseMode         `tg:"explanation_parse_mode"`
	ExplanationEntities  []tg.MessageEntity   `tg:"explanation_entities"`
	OpenPeriod           int                  `tg:"open_period"`
	CloseDate            int64                `tg:"close_date"`
	IsClosed             bool                 `tg:"is_closed"`
	ReplyMarkup          tg.ReplyMarkup       `tg:"reply_markup"`
}

func (Poll) sendMethod() string { return "sendPoll" }

var _ Sendable = Dice{}

// Dice contains information about the dice to be sent.
type Dice struct {
	Emoji       tg.DiceEmoji   `tg:"emoji"`
	ReplyMarkup tg.ReplyMarkup `tg:"reply_markup"`
}

func (Dice) sendMethod() string { return "sendDice" }

var _ Sendable = Sticker{}

// Sticker contains information about the sticker to be sent.
type Sticker struct {
	Sticker     tg.Inputtable  `tg:"sticker"`
	Emoji       string         `tg:"emoji"`
	ReplyMarkup tg.ReplyMarkup `tg:"reply_markup"`
}

func (Sticker) sendMethod() string { return "sendSticker" }

var _ Sendable = Game{}

// Game contains information about the game to be sent.
type Game struct {
	ShortName   string                   `tg:"game_short_name"`
	ReplyMarkup *tg.InlineKeyboardMarkup `tg:"reply_markup"`
}

func (Game) sendMethod() string { return "sendGame" }

var _ Sendable = Invoice{}

// Invoice contains information about the invoice to be sent.
type Invoice struct {
	Title                     string                   `tg:"title"`
	Description               string                   `tg:"description"`
	Payload                   string                   `tg:"payload"`
	ProviderToken             string                   `tg:"provider_token"`
	Currency                  string                   `tg:"currency"`
	Prices                    []tg.LabeledPrice        `tg:"prices"`
	MaxTipAmount              int                      `tg:"max_tip_amount"`
	SuggestedTipAmounts       []int                    `tg:"suggested_tip_amounts"`
	ProviderData              string                   `tg:"provider_data"`
	PhotoURL                  string                   `tg:"photo_url"`
	PhotoSize                 int                      `tg:"photo_size"`
	PhotoWidth                int                      `tg:"photo_width"`
	PhotoHeight               int                      `tg:"photo_height"`
	NeedName                  bool                     `tg:"need_name"`
	NeedPhoneNumber           bool                     `tg:"need_phone_number"`
	NeedEmail                 bool                     `tg:"need_email"`
	NeedShippingAddress       bool                     `tg:"need_shipping_address"`
	SendPhoneNumberToProvider bool                     `tg:"send_phone_number_to_provider"`
	SendEmailToProvider       bool                     `tg:"send_email_to_provider"`
	IsFlexible                bool                     `tg:"is_flexible"`
	StartParameter            string                   `tg:"start_parameter"`
	ReplyMarkup               *tg.InlineKeyboardMarkup `tg:"reply_markup"`
}

func (Invoice) sendMethod() string { return "sendInvoice" }
