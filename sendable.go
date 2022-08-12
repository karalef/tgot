package bot

import (
	"errors"
	"strconv"
	"tghwbot/bot/tg"
)

// Sendable interface for Chat.Send.
type Sendable interface {
	what() string
	params(params)
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
	p.setJSON("caption_entities", c.Entities)
}

// BaseFile is common structure for single file with thumbnail.
type BaseFile struct {
	File      *tg.InputFile
	Thumbnail *tg.InputFile
}

func (b *BaseFile) files(field string) []File {
	var f []File
	if b.Thumbnail != nil {
		f = make([]File, 1, 2)
		f[0] = File{"thumb", b.Thumbnail}
	} else {
		f = make([]File, 0, 1)
	}
	return append(f, File{field, b.File})
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
type SendOptions struct {
	BaseSendOptions
	ReplyMarkup tg.ReplyMarkup
}

func (o SendOptions) embed(p params) {
	o.BaseSendOptions.embed(p)
	p.setJSON("reply_markup", o.ReplyMarkup)
}

// MediaGroupSendOptions contains sending options for sendMediaGroup.
type MediaGroupSendOptions struct {
	BaseSendOptions
}

// InvoiceSendOptions contains sending options for sendInvoice.
type InvoiceSendOptions struct {
	BaseSendOptions
	ReplyMarkup *tg.InlineKeyboardMarkup
}

func (o InvoiceSendOptions) embed(p params) {
	o.BaseSendOptions.embed(p)
	p.setJSON("reply_markup", o.ReplyMarkup)
}

// NewMessage makes a new message.
func NewMessage(text string) *Message {
	return &Message{
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

func (Message) what() string {
	return "Message"
}

func (m Message) params(p params) {
	p.set("text", m.Text)
	p.set("parse_mode", string(m.ParseMode))
	p.setJSON("entities", m.Entities)
	p.setBool("disable_web_page_preview", m.DisableWebPagePreview)
}

// NewPhoto makes a new photo.
func NewPhoto(photo *tg.InputFile) *Photo {
	return &Photo{
		Photo: photo,
	}
}

var _ Sendable = Photo{}

// Photo contains information about the photo to be sent.
type Photo struct {
	Photo *tg.InputFile
	CaptionData
}

func (Photo) what() string {
	return "Photo"
}

func (ph Photo) params(p params) {
	ph.CaptionData.embed(p)
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

var _ Sendable = Audio{}

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

func (a Audio) params(p params) {
	a.CaptionData.embed(p)
	p.setInt("duration", a.Duration)
	p.set("performer", a.Performer)
	p.set("title", a.Title)
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

var _ Sendable = Document{}

// Document contains information about the document to be sent.
type Document struct {
	BaseFile
	CaptionData
	DisableTypeDetection bool
}

func (Document) what() string {
	return "Document"
}

func (d Document) params(p params) {
	d.CaptionData.embed(p)
	p.setBool("disable_content_type_detection", d.DisableTypeDetection)
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

func (Video) what() string {
	return "Video"
}

func (v Video) params(p params) {
	p.setInt("duration", v.Duration)
	p.setInt("width", v.Width)
	p.setInt("height", v.Height)
	v.CaptionData.embed(p)
	p.setBool("supports_streaming", v.SupportsStreaming)
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

var _ Sendable = Animation{}

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

func (a Animation) params(p params) {
	p.setInt("duration", a.Duration)
	p.setInt("width", a.Width)
	p.setInt("height", a.Height)
	a.CaptionData.embed(p)
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

var _ Sendable = Voice{}

// Voice contains information about the voice to be sent.
type Voice struct {
	Voice *tg.InputFile
	CaptionData
	Duration int
}

func (Voice) what() string {
	return "Voice"
}

func (v Voice) params(p params) {
	v.CaptionData.embed(p)
	p.setInt("duration", v.Duration)
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

var _ Sendable = VideoNote{}

// VideoNote contains information about the video note to be sent.
type VideoNote struct {
	BaseFile
	Duration int
	Length   int
}

func (VideoNote) what() string {
	return "VideoNote"
}

func (v VideoNote) params(p params) {
	p.setInt("duration", v.Duration)
	p.setInt("length", v.Length)
}

func (v VideoNote) files() []File {
	return v.BaseFile.files("video_note")
}

// MediaGroup contains information about the media group to be sent.
type MediaGroup []tg.MediaInputter

func (g MediaGroup) data(p params) ([]File, error) {
	return prepareInputMedia(p, true, g...)
}

func prepareInputMedia(p params, mediaGroup bool, media ...tg.MediaInputter) ([]File, error) {
	if len(media) == 0 {
		return nil, nil
	}
	var files []File
	for i := range media {
		n := strconv.Itoa(i)
		var med, thumb **tg.InputFile
		switch m := media[i].(type) {
		case tg.InputMedia[tg.InputMediaPhoto]:
			med = &m.Media
		case tg.InputMedia[tg.InputMediaVideo]:
			med, thumb = &m.Media, &m.Data.Thumbnail
		case tg.InputMedia[tg.InputMediaAudio]:
			med, thumb = &m.Media, &m.Data.Thumbnail
		case tg.InputMedia[tg.InputMediaDocument]:
			med, thumb = &m.Media, &m.Data.Thumbnail
		case tg.InputMedia[tg.InputMediaAnimation]:
			med, thumb = &m.Media, &m.Data.Thumbnail
		default:
			return nil, errors.New("unsupported media group entry " + n + " type")
		}
		if *med == nil {
			return nil, errors.New("media group entry " + n + " does not exist")
		}
		if _, r := (*med).UploadData(); r != nil {
			*med = tg.FileReader("file-"+n, r)
			files = append(files, File{"file-" + n, *med})
		}
		if thumb == nil {
			continue
		}
		if _, r := (*thumb).UploadData(); r != nil {
			*thumb = tg.FileReader("file-"+n+"-thumb", r)
			files = append(files, File{"file-" + n + "-thumb", *thumb})
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
type Location struct {
	tg.Location
}

func (Location) what() string {
	return "Location"
}

func (l Location) params(p params) {
	p.setFloat("latitude", l.Lat)
	p.setFloat("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		p.setFloat("horizontal_accuracy", *l.HorizontalAccuracy)
	}
	p.setInt("live_period", l.LivePeriod)
	p.setInt("heading", l.Heading)
	p.setInt("proximity_alert_radius", l.AlertRadius)
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

func (Venue) what() string {
	return "Venue"
}

func (v Venue) params(p params) {
	p.setFloat("latitude", v.Lat)
	p.setFloat("longitude", v.Long)
	p.set("title", v.Title)
	p.set("address", v.Address)
	p.set("foursquare_id", v.FoursquareID)
	p.set("foursquare_type", v.FoursquareType)
	p.set("google_place_id", v.GooglePlaceID)
	p.set("google_place_type", v.GooglePlaceType)
}

var _ Sendable = Contact{}

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

func (c Contact) params(p params) {
	p.set("phone_number", c.PhoneNumber)
	p.set("first_name", c.FirstName)
	p.set("last_name", c.LastName)
	p.set("vcard", c.Vcard)
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

func (Poll) what() string {
	return "Poll"
}

func (poll Poll) params(p params) {
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
}

var _ Sendable = Dice("")

// Dice contains information about the dice to be sent.
type Dice tg.DiceEmoji

func (Dice) what() string {
	return "Dice"
}

func (d Dice) params(p params) {
	p.set("emoji", string(d))
}

// Invoice contains information about the invoice to be sent.
type Invoice struct {
	tg.InputInvoiceMessageContent
	StartParameter string
}

func (i Invoice) params(p params) {
	p.set("title", i.Title)
	p.set("description", i.Description)
	p.set("payload", i.Payload)
	p.set("provider_token", i.ProviderToken)
	p.set("currency", i.Currency)
	p.setJSON("prices", i.Prices)
	p.setInt("max_tip_amount", i.MaxTipAmount)
	p.setJSON("suggested_tip_amounts", i.SuggestedTipAmounts)
	p.set("start_parameter", i.StartParameter)
	p.set("provider_data", i.ProviderData)
	p.set("photo_url", i.PhotoURL)
	p.setInt("photo_size", i.PhotoSize)
	p.setInt("photo_width", i.PhotoWidth)
	p.setInt("photo_height", i.PhotoHeight)
	p.setBool("need_name", i.NeedName)
	p.setBool("need_phone_number", i.NeedPhoneNumber)
	p.setBool("need_email", i.NeedEmail)
	p.setBool("need_shipping_address", i.NeedShippingAddress)
	p.setBool("send_phone_number_to_provider", i.SendPhoneNumberToProvider)
	p.setBool("send_email_to_provider", i.SendEmailToProvider)
	p.setBool("is_flexible", i.IsFlexible)
}
