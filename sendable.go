package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Sendable interface for Chat.Send.
type Sendable interface {
	sendData(*api.Data) string
}

func embed[T interface{ embed(*api.Data) }](d *api.Data, e []T) {
	if len(e) > 0 {
		e[0].embed(d)
	}
}

// CaptionData represents caption with entities and parse mode.
type CaptionData struct {
	Caption   string
	ParseMode tg.ParseMode
	Entities  []tg.MessageEntity
}

func (c CaptionData) embed(d *api.Data) {
	d.Set("caption", c.Caption)
	d.Set("parse_mode", string(c.ParseMode))
	d.SetJSON("caption_entities", c.Entities)
}

// SendOptions cointains common send* parameters.
type SendOptions struct {
	DisableNotification      bool
	ProtectContent           bool
	ReplyTo                  int
	AllowSendingWithoutReply bool
}

func (o SendOptions) embed(d *api.Data) {
	d.SetBool("disable_notification", o.DisableNotification)
	d.SetBool("protect_content", o.ProtectContent)
	d.SetInt("reply_to_message_id", o.ReplyTo)
	d.SetBool("allow_sending_without_reply", o.AllowSendingWithoutReply)
}

// NewMessage makes a new message.
func NewMessage(text string) Message { return Message{Text: text} }

var _ Sendable = Message{}

// Message contains information about the message to be sent.
type Message struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
	ReplyMarkup           tg.ReplyMarkup
}

func (m Message) sendData(d *api.Data) string {
	d.Set("text", m.Text)
	d.Set("parse_mode", string(m.ParseMode))
	d.SetJSON("entities", m.Entities)
	d.SetBool("disable_web_page_preview", m.DisableWebPagePreview)
	d.SetJSON("reply_markup", m.ReplyMarkup)
	return "sendMessage"
}

// NewPhoto makes a new photo.
func NewPhoto(photo tg.Inputtable) Photo { return Photo{Photo: photo} }

var _ Sendable = Photo{}

// Photo contains information about the photo to be sent.
type Photo struct {
	Photo tg.Inputtable
	CaptionData
	HasSpoiler  bool
	ReplyMarkup tg.ReplyMarkup
}

func (ph Photo) sendData(d *api.Data) string {
	d.SetFile("photo", ph.Photo, nil)
	ph.CaptionData.embed(d)
	d.SetBool("has_spoiler", ph.HasSpoiler)
	d.SetJSON("reply_markup", ph.ReplyMarkup)
	return "sendPhoto"
}

var _ Sendable = Audio{}

// Audio contains information about the audio to be sent.
type Audio struct {
	Audio     tg.Inputtable
	Thumbnail tg.Inputtable
	CaptionData
	Duration    int
	Performer   string
	Title       string
	ReplyMarkup tg.ReplyMarkup
}

func (a Audio) sendData(d *api.Data) string {
	d.SetFile("audio", a.Audio, a.Thumbnail)
	a.CaptionData.embed(d)
	d.SetInt("duration", a.Duration)
	d.Set("performer", a.Performer)
	d.Set("title", a.Title)
	d.SetJSON("reply_markup", a.ReplyMarkup)
	return "sendAudio"
}

var _ Sendable = Document{}

// Document contains information about the document to be sent.
type Document struct {
	Document  tg.Inputtable
	Thumbnail tg.Inputtable
	CaptionData
	DisableTypeDetection bool
	ReplyMarkup          tg.ReplyMarkup
}

func (d Document) sendData(data *api.Data) string {
	data.SetFile("document", d.Document, d.Thumbnail)
	d.CaptionData.embed(data)
	data.SetBool("disable_content_type_detection", d.DisableTypeDetection)
	data.SetJSON("reply_markup", d.ReplyMarkup)
	return "sendDocument"
}

var _ Sendable = Video{}

// Video contains information about the video to be sent.
type Video struct {
	Video     tg.Inputtable
	Thumbnail tg.Inputtable
	CaptionData
	Duration          int
	Width             int
	Height            int
	HasSpoiler        bool
	SupportsStreaming bool
	ReplyMarkup       tg.ReplyMarkup
}

func (v Video) sendData(d *api.Data) string {
	d.SetFile("video", v.Video, v.Thumbnail)
	v.CaptionData.embed(d)
	d.SetInt("duration", v.Duration)
	d.SetInt("width", v.Width)
	d.SetInt("height", v.Height)
	d.SetBool("has_spoiler", v.HasSpoiler)
	d.SetBool("supports_streaming", v.SupportsStreaming)
	d.SetJSON("reply_markup", v.ReplyMarkup)
	return "sendVideo"
}

var _ Sendable = Animation{}

// Animation contains information about the animation to be sent.
type Animation struct {
	Animation tg.Inputtable
	Thumbnail tg.Inputtable
	CaptionData
	Duration    int
	Width       int
	Height      int
	HasSpoiler  bool
	ReplyMarkup tg.ReplyMarkup
}

func (a Animation) sendData(d *api.Data) string {
	d.SetFile("animation", a.Animation, a.Thumbnail)
	a.CaptionData.embed(d)
	d.SetInt("duration", a.Duration)
	d.SetInt("width", a.Width)
	d.SetInt("height", a.Height)
	d.SetBool("has_spoiler", a.HasSpoiler)
	d.SetJSON("reply_markup", a.ReplyMarkup)
	return "sendAnimation"
}

var _ Sendable = Voice{}

// Voice contains information about the voice to be sent.
type Voice struct {
	Voice tg.Inputtable
	CaptionData
	Duration    int
	ReplyMarkup tg.ReplyMarkup
}

func (v Voice) sendData(d *api.Data) string {
	d.SetFile("voice", v.Voice, nil)
	v.CaptionData.embed(d)
	d.SetInt("duration", v.Duration)
	d.SetJSON("reply_markup", v.ReplyMarkup)
	return "sendVoice"
}

var _ Sendable = VideoNote{}

// VideoNote contains information about the video note to be sent.
type VideoNote struct {
	VideoNote   tg.Inputtable
	Thumbnail   tg.Inputtable
	Duration    int
	Length      int
	ReplyMarkup tg.ReplyMarkup
}

func (v VideoNote) sendData(d *api.Data) string {
	d.SetFile("video_note", v.VideoNote, v.Thumbnail)
	d.SetInt("duration", v.Duration)
	d.SetInt("length", v.Length)
	d.SetJSON("reply_markup", v.ReplyMarkup)
	return "sendVideoNote"
}

// MediaGroup contains information about the media group to be sent.
type MediaGroup struct {
	Media       []tg.MediaInputter
	ReplyMarkup *tg.InlineKeyboardMarkup
}

func (g MediaGroup) data(d *api.Data) {
	prepareInputMedia(d, g.Media...)
	d.SetJSON("media", g.Media)
	d.SetJSON("reply_markup", g.ReplyMarkup)
}

func prepareInputMedia(d *api.Data, media ...tg.MediaInputter) {
	if len(media) == 0 {
		return
	}
	for i := range media {
		var med, thumb tg.Inputtable
		switch m := media[i].(type) {
		case *tg.InputMedia[tg.InputMediaPhoto]:
			med = m.Media
		case *tg.InputMedia[tg.InputMediaVideo]:
			med, thumb = m.Media, m.Data.Thumbnail
		case *tg.InputMedia[tg.InputMediaAudio]:
			med, thumb = m.Media, m.Data.Thumbnail
		case *tg.InputMedia[tg.InputMediaDocument]:
			med, thumb = m.Media, m.Data.Thumbnail
		case *tg.InputMedia[tg.InputMediaAnimation]:
			med, thumb = m.Media, m.Data.Thumbnail
		default:
			panic("unsupported media type")
		}
		if med == nil {
			continue
		}
		d.AddAttach(med)
		d.AddAttach(thumb)
	}
}

var _ Sendable = Location{}

// Location contains information about the location to be sent.
type Location struct {
	tg.Location
	ReplyMarkup tg.ReplyMarkup
}

func (l Location) sendData(d *api.Data) string {
	d.SetFloat("latitude", l.Lat)
	d.SetFloat("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		d.SetFloat("horizontal_accuracy", *l.HorizontalAccuracy, true)
	}
	d.SetInt("live_period", l.LivePeriod)
	d.SetInt("heading", l.Heading)
	d.SetInt("proximity_alert_radius", l.AlertRadius)
	d.SetJSON("reply_markup", l.ReplyMarkup)
	return "sendLocation"
}

var _ Sendable = Venue{}

// Venue contains information about the venue to be sent.
type Venue struct {
	Lat             float32
	Long            float32
	Title           string
	Address         string
	FoursquareID    string
	FoursquareType  string
	GooglePlaceID   string
	GooglePlaceType string
	ReplyMarkup     tg.ReplyMarkup
}

func (v Venue) sendData(d *api.Data) string {
	d.SetFloat("latitude", v.Lat)
	d.SetFloat("longitude", v.Long)
	d.Set("title", v.Title)
	d.Set("address", v.Address)
	d.Set("foursquare_id", v.FoursquareID)
	d.Set("foursquare_type", v.FoursquareType)
	d.Set("google_place_id", v.GooglePlaceID)
	d.Set("google_place_type", v.GooglePlaceType)
	d.SetJSON("reply_markup", v.ReplyMarkup)
	return "sendVenue"
}

var _ Sendable = Contact{}

// Contact contains information about the contact to be sent.
type Contact struct {
	PhoneNumber string
	FirstName   string
	LastName    string
	Vcard       string
	ReplyMarkup tg.ReplyMarkup
}

func (c Contact) sendData(d *api.Data) string {
	d.Set("phone_number", c.PhoneNumber)
	d.Set("first_name", c.FirstName)
	d.Set("last_name", c.LastName)
	d.Set("vcard", c.Vcard)
	d.SetJSON("reply_markup", c.ReplyMarkup)
	return "sendContact"
}

var _ Sendable = Poll{}

// Poll contains information about the poll to be sent.
type Poll struct {
	Question             string
	Options              []string
	IsAnonymous          bool
	Type                 tg.PollType
	MultipleAnswers      bool
	CorrectOption        int
	Explanation          string
	ExplanationParseMode tg.ParseMode
	ExplanationEntities  []tg.MessageEntity
	OpenPeriod           int
	CloseDate            int64
	IsClosed             bool
	ReplyMarkup          tg.ReplyMarkup
}

func (poll Poll) sendData(d *api.Data) string {
	d.Set("question", poll.Question)
	d.SetJSON("options", poll.Options)
	d.SetBool("is_anonymous", poll.IsAnonymous)
	d.Set("type", string(poll.Type))
	d.SetBool("allows_multiple_answers", poll.MultipleAnswers)
	d.SetInt("correct_option_id", poll.CorrectOption)
	d.Set("explanation", poll.Explanation)
	d.Set("explanation_parse_mode", string(poll.ExplanationParseMode))
	d.SetJSON("explanation_entities", poll.ExplanationEntities)
	d.SetInt("open_period", poll.OpenPeriod)
	d.SetInt64("close_date", poll.CloseDate)
	d.SetBool("is_closed", poll.IsClosed)
	d.SetJSON("reply_markup", poll.ReplyMarkup)
	return "sendPoll"
}

var _ Sendable = Dice{}

// Dice contains information about the dice to be sent.
type Dice struct {
	Emoji       tg.DiceEmoji
	ReplyMarkup tg.ReplyMarkup
}

func (d Dice) sendData(data *api.Data) string {
	data.Set("emoji", string(d.Emoji))
	data.SetJSON("reply_markup", d.ReplyMarkup)
	return "sendDice"
}

var _ Sendable = Sticker{}

// Sticker contains information about the sticker to be sent.
type Sticker struct {
	Sticker     tg.Inputtable
	Emoji       string
	ReplyMarkup tg.ReplyMarkup
}

func (s Sticker) sendData(d *api.Data) string {
	d.SetFile("sticker", s.Sticker, nil)
	d.Set("emoji", s.Emoji)
	d.SetJSON("reply_markup", s.ReplyMarkup)
	return "sendSticker"
}

var _ Sendable = Game{}

// Game contains information about the game to be sent.
type Game struct {
	ShortName   string
	ReplyMarkup *tg.InlineKeyboardMarkup
}

func (g Game) sendData(d *api.Data) string {
	d.Set("game_short_name", g.ShortName)
	d.SetJSON("reply_markup", g.ReplyMarkup)
	return "sendGame"
}

var _ Sendable = Invoice{}

// Invoice contains information about the invoice to be sent.
type Invoice struct {
	Title                     string
	Description               string
	Payload                   string
	ProviderToken             string
	Currency                  string
	Prices                    []tg.LabeledPrice
	MaxTipAmount              int
	SuggestedTipAmounts       []int
	ProviderData              string
	PhotoURL                  string
	PhotoSize                 int
	PhotoWidth                int
	PhotoHeight               int
	NeedName                  bool
	NeedPhoneNumber           bool
	NeedEmail                 bool
	NeedShippingAddress       bool
	SendPhoneNumberToProvider bool
	SendEmailToProvider       bool
	IsFlexible                bool
	StartParameter            string
	ReplyMarkup               *tg.InlineKeyboardMarkup
}

func (i Invoice) sendData(d *api.Data) string {
	d.Set("title", i.Title)
	d.Set("description", i.Description)
	d.Set("payload", i.Payload)
	d.Set("provider_token", i.ProviderToken)
	d.Set("currency", i.Currency)
	d.SetJSON("prices", i.Prices)
	d.SetInt("max_tip_amount", i.MaxTipAmount)
	d.SetJSON("suggested_tip_amounts", i.SuggestedTipAmounts)
	d.Set("start_parameter", i.StartParameter)
	d.Set("provider_data", i.ProviderData)
	d.Set("photo_url", i.PhotoURL)
	d.SetInt("photo_size", i.PhotoSize)
	d.SetInt("photo_width", i.PhotoWidth)
	d.SetInt("photo_height", i.PhotoHeight)
	d.SetBool("need_name", i.NeedName)
	d.SetBool("need_phone_number", i.NeedPhoneNumber)
	d.SetBool("need_email", i.NeedEmail)
	d.SetBool("need_shipping_address", i.NeedShippingAddress)
	d.SetBool("send_phone_number_to_provider", i.SendPhoneNumberToProvider)
	d.SetBool("send_email_to_provider", i.SendEmailToProvider)
	d.SetBool("is_flexible", i.IsFlexible)
	d.SetJSON("reply_markup", i.ReplyMarkup)
	return "sendInvoice"
}
