package tgot

import (
	"errors"
	"strconv"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

// Sendable interface for Chat.Send.
type Sendable interface {
	data() (string, api.Data)
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

// BaseSendOptions contains common send* parameters for all send methods.
type BaseSendOptions struct {
	DisableNotification      bool
	ProtectContent           bool
	ReplyTo                  int
	AllowSendingWithoutReply bool
}

func (o BaseSendOptions) embed(d *api.Data) {
	d.SetBool("disable_notification", o.DisableNotification)
	d.SetBool("protect_content", o.ProtectContent)
	d.SetInt("reply_to_message_id", o.ReplyTo)
	d.SetBool("allow_sending_without_reply", o.AllowSendingWithoutReply)
}

// SendOptions cointains common send* parameters.
type SendOptions[T tg.ReplyMarkup] struct {
	BaseSendOptions
	ReplyMarkup T
}

func (o SendOptions[T]) embed(d *api.Data) {
	o.BaseSendOptions.embed(d)
	d.SetJSON("reply_markup", o.ReplyMarkup)
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
}

func (m Message) data() (string, api.Data) {
	d := api.NewData()
	d.Set("text", m.Text)
	d.Set("parse_mode", string(m.ParseMode))
	d.SetJSON("entities", m.Entities)
	d.SetBool("disable_web_page_preview", m.DisableWebPagePreview)
	return "sendMessage", d
}

// NewPhoto makes a new photo.
func NewPhoto(photo tg.Inputtable) Photo { return Photo{Photo: photo} }

var _ Sendable = Photo{}

// Photo contains information about the photo to be sent.
type Photo struct {
	Photo tg.Inputtable
	CaptionData
}

func (ph Photo) data() (string, api.Data) {
	d := api.NewData()
	d.SetFile("photo", ph.Photo, nil)
	ph.CaptionData.embed(&d)
	return "sendPhoto", d
}

var _ Sendable = Audio{}

// Audio contains information about the audio to be sent.
type Audio struct {
	Audio     tg.Inputtable
	Thumbnail tg.Inputtable
	CaptionData
	Duration  int
	Performer string
	Title     string
}

func (a Audio) data() (string, api.Data) {
	d := api.NewData()
	d.SetFile("audio", a.Audio, a.Thumbnail)
	a.CaptionData.embed(&d)
	d.SetInt("duration", a.Duration)
	d.Set("performer", a.Performer)
	d.Set("title", a.Title)
	return "sendAudio", d
}

var _ Sendable = Document{}

// Document contains information about the document to be sent.
type Document struct {
	Document  tg.Inputtable
	Thumbnail tg.Inputtable
	CaptionData
	DisableTypeDetection bool
}

func (d Document) data() (string, api.Data) {
	data := api.NewData()
	data.SetFile("document", d.Document, d.Thumbnail)
	d.CaptionData.embed(&data)
	data.SetBool("disable_content_type_detection", d.DisableTypeDetection)
	return "sendDocument", data
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
	SupportsStreaming bool
}

func (v Video) data() (string, api.Data) {
	d := api.NewData()
	d.SetFile("video", v.Video, v.Thumbnail)
	v.CaptionData.embed(&d)
	d.SetInt("duration", v.Duration)
	d.SetInt("width", v.Width)
	d.SetInt("height", v.Height)
	d.SetBool("supports_streaming", v.SupportsStreaming)
	return "sendVideo", d
}

var _ Sendable = Animation{}

// Animation contains information about the animation to be sent.
type Animation struct {
	Animation tg.Inputtable
	Thumbnail tg.Inputtable
	CaptionData
	Duration int
	Width    int
	Height   int
}

func (a Animation) data() (string, api.Data) {
	d := api.NewData()
	d.SetFile("animation", a.Animation, a.Thumbnail)
	a.CaptionData.embed(&d)
	d.SetInt("duration", a.Duration)
	d.SetInt("width", a.Width)
	d.SetInt("height", a.Height)
	return "sendAnimation", d
}

var _ Sendable = Voice{}

// Voice contains information about the voice to be sent.
type Voice struct {
	Voice tg.Inputtable
	CaptionData
	Duration int
}

func (v Voice) data() (string, api.Data) {
	d := api.NewData()
	d.SetFile("voice", v.Voice, nil)
	v.CaptionData.embed(&d)
	d.SetInt("duration", v.Duration)
	return "sendVoice", d
}

var _ Sendable = VideoNote{}

// VideoNote contains information about the video note to be sent.
type VideoNote struct {
	VideoNote tg.Inputtable
	Thumbnail tg.Inputtable
	Duration  int
	Length    int
}

func (v VideoNote) data() (string, api.Data) {
	d := api.NewData()
	d.SetFile("video_note", v.VideoNote, v.Thumbnail)
	d.SetInt("duration", v.Duration)
	d.SetInt("length", v.Length)
	return "sendVideoNote", d
}

// MediaGroup contains information about the media group to be sent.
type MediaGroup []tg.MediaInputter

func (g MediaGroup) data() (api.Data, error) {
	d := api.NewData()
	return d, prepareInputMedia(&d, true, g...)
}

func prepareInputMedia(d *api.Data, mediaGroup bool, media ...tg.MediaInputter) error {
	if len(media) == 0 {
		return nil
	}
	for i := range media {
		n := strconv.Itoa(i)
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
			return errors.New("unsupported media group entry " + n + " type")
		}
		if med == nil {
			return errors.New("media group entry " + n + " does not exist")
		}
		if f, ok := med.(*tg.InputFile); ok {
			field := "file-" + n
			d.AddFile(field, f.AsMedia(field))
		}
		if thumb == nil {
			continue
		}
		if f, ok := thumb.(*tg.InputFile); ok {
			field := "file-" + n + "-thumb"
			d.AddFile(field, f.AsMedia(field))
		}
	}

	if !mediaGroup {
		d.SetJSON("media", media[0])
	} else {
		d.SetJSON("media", media)
	}
	return nil
}

var _ Sendable = Location{}

// Location contains information about the location to be sent.
type Location tg.Location

func (l Location) data() (string, api.Data) {
	d := api.NewData()
	d.SetFloat("latitude", l.Lat)
	d.SetFloat("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		d.SetFloat("horizontal_accuracy", *l.HorizontalAccuracy, true)
	}
	d.SetInt("live_period", l.LivePeriod)
	d.SetInt("heading", l.Heading)
	d.SetInt("proximity_alert_radius", l.AlertRadius)
	return "sendLocation", d
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
}

func (v Venue) data() (string, api.Data) {
	d := api.NewData()
	d.SetFloat("latitude", v.Lat)
	d.SetFloat("longitude", v.Long)
	d.Set("title", v.Title)
	d.Set("address", v.Address)
	d.Set("foursquare_id", v.FoursquareID)
	d.Set("foursquare_type", v.FoursquareType)
	d.Set("google_place_id", v.GooglePlaceID)
	d.Set("google_place_type", v.GooglePlaceType)
	return "sendVenue", d
}

var _ Sendable = Contact{}

// Contact contains information about the contact to be sent.
type Contact struct {
	PhoneNumber string
	FirstName   string
	LastName    string
	Vcard       string
}

func (c Contact) data() (string, api.Data) {
	d := api.NewData()
	d.Set("phone_number", c.PhoneNumber)
	d.Set("first_name", c.FirstName)
	d.Set("last_name", c.LastName)
	d.Set("vcard", c.Vcard)
	return "sendContact", d
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
}

func (poll Poll) data() (string, api.Data) {
	d := api.NewData()
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
	return "sendPoll", d
}

var _ Sendable = Dice("")

// Dice contains information about the dice to be sent.
type Dice tg.DiceEmoji

func (d Dice) data() (string, api.Data) {
	return "sendDice", api.NewData().Set("emoji", string(d))
}

// Sticker contains information about the sticker to be sent.
type Sticker struct {
	Sticker tg.Inputtable
}

func (s Sticker) data() (string, api.Data) {
	d := api.NewData()
	d.SetFile("sticker", s.Sticker, nil)
	return "sendSticker", d
}
