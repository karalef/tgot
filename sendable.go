package bot

import (
	"errors"
	"strconv"
	"tghwbot/bot/tg"
)

// Sendable interface for Chat.Send.
type Sendable interface {
	what() string
	params() params
}

// Fileable interface.
type Fileable interface {
	Sendable
	files() []File
}

// CaptionData represents caption with entities and parse mode.
type CaptionData struct {
	Caption   string
	ParseMode tg.ParseMode
	Entities  []tg.MessageEntity
}

func (c *CaptionData) embed(p params) {
	p.set("caption", c.Caption)
	p.set("parse_mode", string(c.ParseMode))
	p.set("caption_entities", c.Entities)
}

// BaseFile is common structure for single file with thumbnail.
type BaseFile struct {
	File      *tg.InputFile
	Thumbnail *tg.InputFile
}

func (b *BaseFile) files(field string) []File {
	f := make([]File, 0, 2)
	f = append(f, File{field, b.File})
	if b.Thumbnail != nil {
		f = append(f, File{"thumb", b.Thumbnail})
	}
	return f
}

// SendOptions cointains common send* parameters.
type SendOptions struct {
	DisableNotification      bool
	ProtectContent           bool
	ReplyTo                  int
	AllowSendingWithoutReply bool
	ReplyMarkup              tg.ReplyMarkup
}

// NewMessage makes a new message.
func NewMessage(text string) *Message {
	return &Message{
		Text: text,
	}
}

// Message contains information about the message to be sent.
type Message struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
}

func (Message) what() string {
	return "Message"
}

func (m Message) params() params {
	p := params{}
	p.set("text", m.Text)
	p.set("parse_mode", string(m.ParseMode))
	p.set("entities", m.Entities)
	p.set("disable_web_page_preview", m.DisableWebPagePreview)
	return p
}

// NewPhoto makes a new photo.
func NewPhoto(photo *tg.InputFile) *Photo {
	return &Photo{
		Photo: photo,
	}
}

// Photo contains information about the photo to be sent.
type Photo struct {
	Photo *tg.InputFile
	CaptionData
}

func (Photo) what() string {
	return "Photo"
}

func (ph Photo) params() params {
	p := params{}
	ph.CaptionData.embed(p)
	return p
}

func (ph Photo) files() []File {
	return []File{{"photo", ph.Photo}}
}

// NewAudio makes a new audio.
func NewAudio(audio *tg.InputFile) *Audio {
	return &Audio{
		BaseFile: BaseFile{File: audio},
	}
}

// Audio contains information about the audio to be sent.
type Audio struct {
	BaseFile
	CaptionData
	Duration  int
	Performer string
	Title     string
}

func (Audio) what() string {
	return "Audio"
}

func (a Audio) params() params {
	p := params{}
	a.CaptionData.embed(p)
	p.set("duration", a.Duration)
	p.set("performer", a.Performer)
	p.set("title", a.Title)
	return p
}

func (a Audio) files() []File {
	return a.BaseFile.files("audio")
}

// NewDocument makes a new document.
func NewDocument(document *tg.InputFile) *Document {
	return &Document{
		BaseFile: BaseFile{File: document},
	}
}

// Document contains information about the document to be sent.
type Document struct {
	BaseFile
	CaptionData
	DisableTypeDetection bool
}

func (Document) what() string {
	return "Document"
}

func (d Document) params() params {
	p := params{}
	d.CaptionData.embed(p)
	p.set("disable_content_type_detection", d.DisableTypeDetection)
	return p
}

func (d Document) files() []File {
	return d.BaseFile.files("document")
}

// NewVideo makes a new video.
func NewVideo(video *tg.InputFile) *Video {
	return &Video{
		BaseFile: BaseFile{File: video},
	}
}

// Video contains information about the video to be sent.
type Video struct {
	BaseFile
	CaptionData
	Duration          int
	Width             int
	Height            int
	SupportsStreaming bool
}

func (Video) what() string {
	return "Video"
}

func (v Video) params() params {
	p := params{}
	p.set("duration", v.Duration)
	p.set("width", v.Width)
	p.set("height", v.Height)
	v.CaptionData.embed(p)
	p.set("supports_streaming", v.SupportsStreaming)
	return p
}

func (v Video) files() []File {
	return v.BaseFile.files("video")
}

// NewAnimation makes a new animation.
func NewAnimation(animation *tg.InputFile) *Video {
	return &Video{
		BaseFile: BaseFile{File: animation},
	}
}

// Animation contains information about the animation to be sent.
type Animation struct {
	BaseFile
	CaptionData
	Duration int
	Width    int
	Height   int
}

func (Animation) what() string {
	return "Animation"
}

func (a Animation) params() params {
	p := params{}
	p.set("duration", a.Duration)
	p.set("width", a.Width)
	p.set("height", a.Height)
	a.CaptionData.embed(p)
	return p
}

func (a Animation) files() []File {
	return a.BaseFile.files("animation")
}

// NewVoice makes a new voice.
func NewVoice(voice *tg.InputFile) *Voice {
	return &Voice{
		Voice: voice,
	}
}

// Voice contains information about the voice to be sent.
type Voice struct {
	Voice *tg.InputFile
	CaptionData
	Duration int
}

func (Voice) what() string {
	return "Voice"
}

func (v Voice) params() params {
	p := params{}
	v.CaptionData.embed(p)
	p.set("duration", v.Duration)
	return p
}

func (v Voice) files() []File {
	return []File{{"voice", v.Voice}}
}

// NewVideoNote makes a new video note.
func NewVideoNote(videoNote *tg.InputFile) *VideoNote {
	return &VideoNote{
		BaseFile: BaseFile{File: videoNote},
	}
}

// VideoNote contains information about the video note to be sent.
type VideoNote struct {
	BaseFile
	Duration int
	Length   int
}

func (VideoNote) what() string {
	return "VideoNote"
}

func (v VideoNote) params() params {
	p := params{}
	p.set("duration", v.Duration)
	p.set("length", v.Length)
	return p
}

func (v VideoNote) files() []File {
	return v.BaseFile.files("video_note")
}

// MediaGroup contains information about the media group to be sent.
type MediaGroup struct {
	Media                    tg.Media
	DisableNotification      bool
	ProtectContent           bool
	ReplyTo                  int
	AllowSendingWithoutReply bool
}

func (g MediaGroup) data() (params, []File, error) {
	p := params{}
	p.set("disable_notification", g.DisableNotification)
	p.set("protect_content", g.ProtectContent)
	p.set("reply_to_message_id", g.ReplyTo)
	p.set("allow_sending_without_reply", g.AllowSendingWithoutReply)

	var files []File
	for i := range g.Media {
		n := strconv.Itoa(i)
		var media, thumb **tg.InputFile
		switch m := g.Media[i].(type) {
		case *tg.InputMediaPhoto:
			media = &m.Media
		case *tg.InputMediaVideo:
			media, thumb = &m.Media, &m.Thumbnail
		case *tg.InputMediaAudio:
			media, thumb = &m.Media, &m.Thumbnail
		case *tg.InputMediaDocument:
			media, thumb = &m.Media, &m.Thumbnail
		default:
			return nil, nil, errors.New("unsupported media group entry " + n + " type")
		}
		if media == nil {
			return nil, nil, errors.New("media group entry " + n + " does not exist")
		}
		if _, r := (*media).UploadData(); r != nil {
			*media = tg.FileReader("file-"+n, r)
			files = append(files, File{"file-" + n, *media})

		}
		if thumb == nil {
			continue
		}
		if _, r := (*thumb).UploadData(); r != nil {
			*thumb = tg.FileReader("file-"+n+"-thumb", r)
			files = append(files, File{"file-" + n + "-thumb", *thumb})
		}
	}

	return p.set("media", g.Media), files, nil
}

// Location contains information about the location to be sent.
type Location struct {
	tg.Location
}

func (Location) what() string {
	return "Location"
}

func (l Location) params() params {
	p := params{}
	p.set("latitude", l.Lat)
	p.set("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		p.set("horizontal_accuracy", *l.HorizontalAccuracy)
	}
	p.set("live_period", l.LivePeriod)
	p.set("heading", l.Heading)
	p.set("proximity_alert_radius", l.AlertRadius)
	return p
}

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

func (Venue) what() string {
	return "Venue"
}

func (v Venue) params() params {
	p := params{}
	p.set("latitude", v.Lat)
	p.set("longitude", v.Long)
	p.set("title", v.Title)
	p.set("address", v.Address)
	p.set("foursquare_id", v.FoursquareID)
	p.set("foursquare_type", v.FoursquareType)
	p.set("google_place_id", v.GooglePlaceID)
	p.set("google_place_type", v.GooglePlaceType)
	return p
}

// Contact contains information about the contact to be sent.
type Contact struct {
	PhoneNumber string
	FirstName   string
	LastName    string
	Vcard       string
}

func (Contact) what() string {
	return "Contact"
}

func (c Contact) params() params {
	p := params{}
	p.set("phone_number", c.PhoneNumber)
	p.set("first_name", c.FirstName)
	p.set("last_name", c.LastName)
	p.set("vcard", c.Vcard)
	return p
}

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

func (Poll) what() string {
	return "Poll"
}

func (poll Poll) params() params {
	p := params{}
	p.set("question", poll.Question)
	p.set("options", poll.Options)
	p.set("is_anonymous", poll.IsAnonymous)
	p.set("type", string(poll.Type))
	p.set("allows_multiple_answers", poll.MultipleAnswers)
	p.set("correct_option_id", poll.CorrectOption)
	p.set("explanation", poll.Explanation)
	p.set("explanation_parse_mode", string(poll.ExplanationParseMode))
	p.set("explanation_entities", poll.ExplanationEntities)
	p.set("open_period", poll.OpenPeriod)
	p.set("close_date", poll.CloseDate)
	p.set("is_closed", poll.IsClosed)
	return p
}

// Dice contains information about the dice to be sent.
type Dice tg.DiceEmoji

func (Dice) what() string {
	return "Dice"
}

func (d Dice) params() params {
	return params{}.set("emoji", string(d))
}

// Invoice contains information about the invoice to be sent.
type Invoice struct {
	tg.InputInvoiceMessageContent
	StartParameter           string
	DisableNotification      bool
	ProtectContent           bool
	ReplyTo                  int
	AllowSendingWithoutReply bool
	ReplyMarkup              *tg.InlineKeyboardMarkup
}

func (i Invoice) params() params {
	p := params{}
	p.set("title", i.Title)
	p.set("description", i.Description)
	p.set("payload", i.Payload)
	p.set("provider_token", i.ProviderToken)
	p.set("currency", i.Currency)
	p.set("prices", i.Prices)
	p.set("max_tip_amount", i.MaxTipAmount)
	p.set("suggested_tip_amounts", i.SuggestedTipAmounts)
	p.set("start_parameter", i.StartParameter)
	p.set("provider_data", i.ProviderData)
	p.set("photo_url", i.PhotoURL)
	p.set("photo_size", i.PhotoSize)
	p.set("photo_width", i.PhotoWidth)
	p.set("photo_height", i.PhotoHeight)
	p.set("need_name", i.NeedName)
	p.set("need_phone_number", i.NeedPhoneNumber)
	p.set("need_email", i.NeedEmail)
	p.set("need_shipping_address", i.NeedShippingAddress)
	p.set("send_phone_number_to_provider", i.SendPhoneNumberToProvider)
	p.set("send_email_to_provider", i.SendEmailToProvider)
	p.set("is_flexible", i.IsFlexible)
	p.set("disable_notification", i.DisableNotification)
	p.set("protect_content", i.ProtectContent)
	p.set("reply_to_message_id", i.ReplyTo)
	p.set("allow_sending_without_reply", i.AllowSendingWithoutReply)
	if i.ReplyMarkup != nil {
		p.set("reply_markup", i.ReplyMarkup)
	}
	return p
}
