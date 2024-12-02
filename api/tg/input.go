package tg

import (
	"bytes"
	"io"

	"github.com/karalef/tgot/api/internal"
)

// Inputtable is an interface for FileID, FileURL and InputFile.
type Inputtable interface {
	FileData() (string, io.Reader)
}

type fileString string

// FileData returns the file data.
func (f fileString) FileData() (string, io.Reader) {
	return string(f), nil
}

// FileID represents the file that has already been uploaded to Telegram.
type FileID = fileString

// FileURL represents the file URL that is used without uploading.
type FileURL = fileString

// InputFile represents the contents of a file to be uploaded.
type InputFile struct {
	Name string
	Data io.Reader

	field string
}

// FileData returns the file data.
func (f *InputFile) FileData() (string, io.Reader) {
	return f.Name, f.Data
}

// AsAttachment sets the link to the multipart field.
func (f *InputFile) AsAttachment(field string) *InputFile {
	f.field = field
	return f
}

// MarshalJSON is json.Marshaler implementation
// for sending as media.
func (f InputFile) MarshalJSON() ([]byte, error) {
	return []byte("\"" + "attach://" + f.field + "\""), nil
}

// FileReader creates the InputFile from reader.
func FileReader(name string, r io.Reader) *InputFile {
	return &InputFile{
		Name: name,
		Data: r,
	}
}

// FileBytes creates the InputFile from bytes.
func FileBytes(name string, data []byte) *InputFile {
	return FileReader(name, bytes.NewReader(data))
}

// InputSticker describes a sticker to be added to a sticker set.
type InputSticker struct {
	Sticker      Inputtable    `json:"sticker"`
	Format       StickerFormat `json:"format"`
	EmojiList    []string      `json:"emoji_list"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	Keywords     []string      `json:"keywords,omitempty"`
}

// StickerFormat is a Sticker format.
type StickerFormat string

// all available sticker formats.
const (
	StickerStatic   StickerFormat = "static"
	StickerAnimated StickerFormat = "animated"
	StickerVideo    StickerFormat = "video"
)

// Inputter is an interface for any input.
type Inputter interface {
	GetMedia() (main, thumb Inputtable)
}

// Thumbnailer is an interface for any media data.
type Thumbnailer interface {
	Thumb() Inputtable
}

// InputMediaData represents any available input media object.
type InputMediaData interface {
	Thumbnailer
	inputMediaType() string
}

var _ Inputter = (*InputMedia)(nil)

// InputMedia represents the content of a media message to be sent.
type InputMedia struct {
	Media     Inputtable      `json:"media"`
	Caption   string          `json:"caption,omitempty"`
	ParseMode ParseMode       `json:"parse_mode,omitempty"`
	Entities  []MessageEntity `json:"caption_entities,omitempty"`
	Data      InputMediaData  `json:"-"`
}

func (i InputMedia) GetMedia() (Inputtable, Inputtable) {
	return i.Media, i.Data.Thumb()
}

func (i InputMedia) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(i.Data, struct {
		Type      string          `json:"type"`
		Media     Inputtable      `json:"media"`
		Caption   string          `json:"caption,omitempty"`
		ParseMode ParseMode       `json:"parse_mode,omitempty"`
		Entities  []MessageEntity `json:"caption_entities,omitempty"`
	}{i.Data.inputMediaType(), i.Media, i.Caption, i.ParseMode, i.Entities})
}

// InputMediaPhoto represents a photo to be sent.
type InputMediaPhoto struct {
	HasSpoiler            bool `json:"has_spoiler,omitempty"`
	ShowCaptionAboveMedia bool `json:"show_caption_above_media,omitempty"`
}

func (*InputMediaPhoto) inputMediaType() string { return "photo" }
func (*InputMediaPhoto) Thumb() Inputtable      { return nil }

// InputMediaVideo represents a video to be sent.
type InputMediaVideo struct {
	Thumbnail             *InputFile `json:"thumbnail,omitempty"`
	Width                 int        `json:"width,omitempty"`
	Height                int        `json:"height,omitempty"`
	Duration              int        `json:"duration,omitempty"`
	SupportsStreaming     bool       `json:"supports_streaming,omitempty"`
	HasSpoiler            bool       `json:"has_spoiler,omitempty"`
	ShowCaptionAboveMedia bool       `json:"show_caption_above_media,omitempty"`
}

func (*InputMediaVideo) inputMediaType() string { return "video" }
func (v *InputMediaVideo) Thumb() Inputtable    { return v.Thumbnail }

// InputMediaAnimation represents an animation file
// (GIF or H.264/MPEG-4 AVC video without sound) to be sent.
type InputMediaAnimation struct {
	Thumbnail             *InputFile `json:"thumbnail,omitempty"`
	Width                 int        `json:"width,omitempty"`
	Height                int        `json:"height,omitempty"`
	Duration              int        `json:"duration,omitempty"`
	HasSpoiler            bool       `json:"has_spoiler,omitempty"`
	ShowCaptionAboveMedia bool       `json:"show_caption_above_media,omitempty"`
}

func (*InputMediaAnimation) inputMediaType() string { return "animation" }
func (a *InputMediaAnimation) Thumb() Inputtable    { return a.Thumbnail }

// InputMediaAudio represents an audio file to be treated as music to be sent.
type InputMediaAudio struct {
	Thumbnail *InputFile `json:"thumbnail,omitempty"`
	Duration  int        `json:"duration,omitempty"`
	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
}

func (*InputMediaAudio) inputMediaType() string { return "audio" }
func (a *InputMediaAudio) Thumb() Inputtable    { return a.Thumbnail }

// InputMediaDocument represents a general file to be sent.
type InputMediaDocument struct {
	Thumbnail            *InputFile `json:"thumbnail,omitempty"`
	DisableTypeDetection bool       `json:"disable_content_type_detection,omitempty"`
}

func (*InputMediaDocument) inputMediaType() string { return "document" }
func (d *InputMediaDocument) Thumb() Inputtable    { return d.Thumbnail }

// InputPaidMediaData represents any available input paid media object.
type InputPaidMediaData interface {
	Thumbnailer
	inputPaidMediaType() string
}

var _ Inputter = InputPaidMedia{}

// InputPaidMedia describes the paid media to be sent.
type InputPaidMedia struct {
	Media Inputtable         `json:"media"`
	Data  InputPaidMediaData `json:"-"`
}

func (i InputPaidMedia) GetMedia() (Inputtable, Inputtable) { return i.Media, i.Data.Thumb() }
func (i InputPaidMedia) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(i.Data, struct {
		Type  string     `json:"type"`
		Media Inputtable `json:"media"`
	}{i.Data.inputPaidMediaType(), i.Media})
}

// InputPaidMediaPhoto is the paid media to send is a photo.
type InputPaidMediaPhoto struct{}

func (*InputPaidMediaPhoto) inputPaidMediaType() string { return "photo" }
func (*InputPaidMediaPhoto) Thumb() Inputtable          { return nil }

// InputPaidMediaVideo is the paid media to send is a video.
type InputPaidMediaVideo struct {
	Thumbnail         *InputFile `json:"thumbnail,omitempty"`
	Width             uint       `json:"width,omitempty"`
	Height            uint       `json:"height,omitempty"`
	Duration          uint       `json:"duration,omitempty"`
	SupportsStreaming bool       `json:"supports_streaming,omitempty"`
}

func (*InputPaidMediaVideo) inputPaidMediaType() string { return "video" }
func (v *InputPaidMediaVideo) Thumb() Inputtable        { return v.Thumbnail }
