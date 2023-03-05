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

// FileID represents the file that has already been uploaded to Telegram.
type FileID string

// FileData returns the file data.
func (f FileID) FileData() (string, io.Reader) {
	return string(f), nil
}

// FileURL represents the file URL that is used without uploading.
type FileURL string

// FileData returns the file data.
func (f FileURL) FileData() (string, io.Reader) {
	return string(f), nil
}

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

// AsMedia sets the link to the multipart field.
func (f *InputFile) AsMedia(field string) *InputFile {
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

// MediaInputter is an interface for InputMedia.
type MediaInputter interface {
	inputMedia()
}

// InputMediaData represents any available input media object.
type InputMediaData interface {
	inputMediaType() string
}

// InputMedia represents the content of a media message to be sent.
type InputMedia[T InputMediaData] struct {
	Media     Inputtable      `json:"media"`
	Caption   string          `json:"caption,omitempty"`
	ParseMode ParseMode       `json:"parse_mode,omitempty"`
	Entities  []MessageEntity `json:"caption_entities,omitempty"`
	Data      T               `json:"-"`
}

func (*InputMedia[T]) inputMedia() {}

func (i InputMedia[T]) MarshalJSON() ([]byte, error) {
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
	HasSpoiler bool `json:"has_spoiler,omitempty"`
}

func (InputMediaPhoto) inputMediaType() string {
	return "photo"
}

// InputMediaVideo represents a video to be sent.
type InputMediaVideo struct {
	Thumbnail         *InputFile `json:"thumb,omitempty"`
	Width             int        `json:"width,omitempty"`
	Height            int        `json:"height,omitempty"`
	Duration          int        `json:"duration,omitempty"`
	SupportsStreaming bool       `json:"supports_streaming,omitempty"`
	HasSpoiler        bool       `json:"has_spoiler,omitempty"`
}

func (InputMediaVideo) inputMediaType() string {
	return "video"
}

// InputMediaAnimation represents an animation file
// (GIF or H.264/MPEG-4 AVC video without sound) to be sent.
type InputMediaAnimation struct {
	Thumbnail  *InputFile `json:"thumb,omitempty"`
	Width      int        `json:"width,omitempty"`
	Height     int        `json:"height,omitempty"`
	Duration   int        `json:"duration,omitempty"`
	HasSpoiler bool       `json:"has_spoiler,omitempty"`
}

func (InputMediaAnimation) inputMediaType() string {
	return "animation"
}

// InputMediaAudio represents an audio file to be treated as music to be sent.
type InputMediaAudio struct {
	Thumbnail *InputFile `json:"thumb,omitempty"`
	Duration  int        `json:"duration,omitempty"`
	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
}

func (InputMediaAudio) inputMediaType() string {
	return "audio"
}

// InputMediaDocument represents a general file to be sent.
type InputMediaDocument struct {
	Thumbnail            *InputFile `json:"thumb,omitempty"`
	DisableTypeDetection bool       `json:"disable_content_type_detection,omitempty"`
}

func (InputMediaDocument) inputMediaType() string {
	return "document"
}
