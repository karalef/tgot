package tgot

import (
	"errors"
	"strconv"

	"github.com/karalef/tgot/tg"
)

// Sendable interface for Chat.Send.
type Sendable interface {
	method() string
	data(params) []file
}

// CaptionData represents caption with entities and parse mode.
type CaptionData struct {
	Caption   string
	ParseMode tg.ParseMode
	Entities  []tg.MessageEntity
}

func (c CaptionData) embed(p params) {
	p.set("caption", c.Caption)
	p.set("parse_mode", string(c.ParseMode))
	p.setJSON("caption_entities", c.Entities)
}

// BaseFile is common structure for single file with thumbnail.
type BaseFile struct {
	File      tg.FileSignature
	Thumbnail tg.FileSignature
}

func (b BaseFile) files(field string) (f []file) {
	if b.Thumbnail != nil {
		f = make([]file, 1, 2)
		f[0] = file{"thumb", b.Thumbnail}
	} else {
		f = make([]file, 0, 1)
	}
	return append(f, file{field, b.File})
}

// BaseSendOptions contains common send* parameters for all send methods.
type BaseSendOptions struct {
	DisableNotification      bool
	ProtectContent           bool
	ReplyTo                  int
	AllowSendingWithoutReply bool
}

func (o BaseSendOptions) embed(p params) {
	p.setBool("disable_notification", o.DisableNotification)
	p.setBool("protect_content", o.ProtectContent)
	p.setInt("reply_to_message_id", o.ReplyTo)
	p.setBool("allow_sending_without_reply", o.AllowSendingWithoutReply)
}

// SendOptions cointains common send* parameters.
type SendOptions[T tg.ReplyMarkup] struct {
	BaseSendOptions
	ReplyMarkup T
}

func (o SendOptions[T]) embed(p params) {
	o.BaseSendOptions.embed(p)
	p.setJSON("reply_markup", o.ReplyMarkup)
}

// MediaGroupSendOptions contains sending options for sendMediaGroup.
type MediaGroupSendOptions = BaseSendOptions

// NewMessage makes a new message.
func NewMessage(text string) Message {
	return Message{
		Text: text,
	}
}

var _ Sendable = Message{}

// Message contains information about the message to be sent.
type Message struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
}

func (Message) method() string {
	return "sendMessage"
}

func (m Message) data(p params) []file {
	p.set("text", m.Text)
	p.set("parse_mode", string(m.ParseMode))
	p.setJSON("entities", m.Entities)
	p.setBool("disable_web_page_preview", m.DisableWebPagePreview)
	return nil
}

// NewPhoto makes a new photo.
func NewPhoto(photo tg.FileSignature) Photo {
	return Photo{
		Photo: photo,
	}
}

var _ Sendable = Photo{}

// Photo contains information about the photo to be sent.
type Photo struct {
	Photo tg.FileSignature
	CaptionData
}

func (Photo) method() string {
	return "sendPhoto"
}

func (ph Photo) data(p params) []file {
	ph.CaptionData.embed(p)
	return []file{{"photo", ph.Photo}}
}

// NewAudio makes a new audio.
func NewAudio(audio tg.FileSignature) Audio {
	return Audio{
		BaseFile: BaseFile{File: audio},
	}
}

var _ Sendable = Audio{}

// Audio contains information about the audio to be sent.
type Audio struct {
	BaseFile
	CaptionData
	Duration  int
	Performer string
	Title     string
}

func (Audio) method() string {
	return "sendAudio"
}

func (a Audio) data(p params) []file {
	a.CaptionData.embed(p)
	p.setInt("duration", a.Duration)
	p.set("performer", a.Performer)
	p.set("title", a.Title)
	return a.BaseFile.files("audio")
}

// NewDocument makes a new document.
func NewDocument(document tg.FileSignature) Document {
	return Document{
		BaseFile: BaseFile{File: document},
	}
}

var _ Sendable = Document{}

// Document contains information about the document to be sent.
type Document struct {
	BaseFile
	CaptionData
	DisableTypeDetection bool
}

func (Document) method() string {
	return "sendDocument"
}

func (d Document) data(p params) []file {
	d.CaptionData.embed(p)
	p.setBool("disable_content_type_detection", d.DisableTypeDetection)
	return d.BaseFile.files("document")
}

// NewVideo makes a new video.
func NewVideo(video tg.FileSignature) Video {
	return Video{
		BaseFile: BaseFile{File: video},
	}
}

var _ Sendable = Video{}

// Video contains information about the video to be sent.
type Video struct {
	BaseFile
	CaptionData
	Duration          int
	Width             int
	Height            int
	SupportsStreaming bool
}

func (Video) method() string {
	return "sendVideo"
}

func (v Video) data(p params) []file {
	p.setInt("duration", v.Duration)
	p.setInt("width", v.Width)
	p.setInt("height", v.Height)
	v.CaptionData.embed(p)
	p.setBool("supports_streaming", v.SupportsStreaming)
	return v.BaseFile.files("video")
}

// NewAnimation makes a new animation.
func NewAnimation(animation tg.FileSignature) Video {
	return Video{
		BaseFile: BaseFile{File: animation},
	}
}

var _ Sendable = Animation{}

// Animation contains information about the animation to be sent.
type Animation struct {
	BaseFile
	CaptionData
	Duration int
	Width    int
	Height   int
}

func (Animation) method() string {
	return "sendAnimation"
}

func (a Animation) data(p params) []file {
	p.setInt("duration", a.Duration)
	p.setInt("width", a.Width)
	p.setInt("height", a.Height)
	a.CaptionData.embed(p)
	return a.BaseFile.files("animation")
}

// NewVoice makes a new voice.
func NewVoice(voice tg.FileSignature) Voice {
	return Voice{
		Voice: voice,
	}
}

var _ Sendable = Voice{}

// Voice contains information about the voice to be sent.
type Voice struct {
	Voice tg.FileSignature
	CaptionData
	Duration int
}

func (Voice) method() string {
	return "sendVoice"
}

func (v Voice) data(p params) []file {
	v.CaptionData.embed(p)
	p.setInt("duration", v.Duration)
	return []file{{"voice", v.Voice}}
}

// NewVideoNote makes a new video note.
func NewVideoNote(videoNote tg.FileSignature) VideoNote {
	return VideoNote{
		BaseFile: BaseFile{File: videoNote},
	}
}

var _ Sendable = VideoNote{}

// VideoNote contains information about the video note to be sent.
type VideoNote struct {
	BaseFile
	Duration int
	Length   int
}

func (VideoNote) method() string {
	return "sendVideoNote"
}

func (v VideoNote) data(p params) []file {
	p.setInt("duration", v.Duration)
	p.setInt("length", v.Length)
	return v.BaseFile.files("video_note")
}

// MediaGroup contains information about the media group to be sent.
type MediaGroup []tg.MediaInputter

func (g MediaGroup) data(p params) ([]file, error) {
	return prepareInputMedia(p, true, g...)
}

func prepareInputMedia(p params, mediaGroup bool, media ...tg.MediaInputter) ([]file, error) {
	if len(media) == 0 {
		return nil, nil
	}
	var files []file
	for i := range media {
		n := strconv.Itoa(i)
		var med, thumb tg.FileSignature
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
			return nil, errors.New("unsupported media group entry " + n + " type")
		}
		if med == nil {
			return nil, errors.New("media group entry " + n + " does not exist")
		}
		if f, ok := med.(*tg.InputFile); ok {
			field := "file-" + n
			files = append(files, file{field, f.AsMedia(field)})
		}
		if thumb == nil {
			continue
		}
		if f, ok := thumb.(*tg.InputFile); ok {
			field := "file-" + n + "-thumb"
			files = append(files, file{field, f.AsMedia(field)})
		}
	}

	if !mediaGroup {
		p.setJSON("media", media[0])
	} else {
		p.setJSON("media", media)
	}
	return files, nil
}

var _ Sendable = Location{}

// Location contains information about the location to be sent.
type Location tg.Location

func (Location) method() string {
	return "sendLocation"
}

func (l Location) data(p params) []file {
	p.setFloat("latitude", l.Lat)
	p.setFloat("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		p.setFloat("horizontal_accuracy", *l.HorizontalAccuracy, true)
	}
	p.setInt("live_period", l.LivePeriod)
	p.setInt("heading", l.Heading)
	p.setInt("proximity_alert_radius", l.AlertRadius)
	return nil
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

func (Venue) method() string {
	return "sendVenue"
}

func (v Venue) data(p params) []file {
	p.setFloat("latitude", v.Lat)
	p.setFloat("longitude", v.Long)
	p.set("title", v.Title)
	p.set("address", v.Address)
	p.set("foursquare_id", v.FoursquareID)
	p.set("foursquare_type", v.FoursquareType)
	p.set("google_place_id", v.GooglePlaceID)
	p.set("google_place_type", v.GooglePlaceType)
	return nil
}

var _ Sendable = Contact{}

// Contact contains information about the contact to be sent.
type Contact struct {
	PhoneNumber string
	FirstName   string
	LastName    string
	Vcard       string
}

func (Contact) method() string {
	return "sendContact"
}

func (c Contact) data(p params) []file {
	p.set("phone_number", c.PhoneNumber)
	p.set("first_name", c.FirstName)
	p.set("last_name", c.LastName)
	p.set("vcard", c.Vcard)
	return nil
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

func (Poll) method() string {
	return "sendPoll"
}

func (poll Poll) data(p params) []file {
	p.set("question", poll.Question)
	p.setJSON("options", poll.Options)
	p.setBool("is_anonymous", poll.IsAnonymous)
	p.set("type", string(poll.Type))
	p.setBool("allows_multiple_answers", poll.MultipleAnswers)
	p.setInt("correct_option_id", poll.CorrectOption)
	p.set("explanation", poll.Explanation)
	p.set("explanation_parse_mode", string(poll.ExplanationParseMode))
	p.setJSON("explanation_entities", poll.ExplanationEntities)
	p.setInt("open_period", poll.OpenPeriod)
	p.setInt64("close_date", poll.CloseDate)
	p.setBool("is_closed", poll.IsClosed)
	return nil
}

var _ Sendable = Dice("")

// Dice contains information about the dice to be sent.
type Dice tg.DiceEmoji

func (Dice) method() string {
	return "sendDice"
}

func (d Dice) data(p params) []file {
	p.set("emoji", string(d))
	return nil
}

// Sticker contains information about the sticker to be sent.
type Sticker struct {
	Sticker tg.FileSignature
}

func (Sticker) method() string {
	return "sendSticker"
}

func (s Sticker) data(params) []file {
	return []file{{field: "sticker", FileSignature: s.Sticker}}
}
