package tg

import (
	"bytes"
	"encoding/json"
	"io"
)

// InputFile represents the contents of a file to be uploaded.
type InputFile struct {
	urlID string
	name  string
	data  io.Reader
}

// MarshalJSON is json.Marshaler implementation
// for sending as media.
func (f InputFile) MarshalJSON() ([]byte, error) {
	if f.urlID != "" {
		return []byte("\"" + f.urlID + "\""), nil
	}
	return []byte("\"" + "attach://" + f.name + "\""), nil
}

// Data returns the file data to send when a file does not need to be uploaded.
func (f InputFile) Data() string {
	return f.urlID
}

// UploadData returns the file name and data reader for the file.
func (f InputFile) UploadData() (string, io.Reader) {
	return f.name, f.data
}

// FileID returns InputFile that has already been uploaded to Telegram.
func FileID(fid string) *InputFile {
	return &InputFile{urlID: fid}
}

// FileURL return InputFile as URL that is used without uploading.
func FileURL(url string) *InputFile {
	return &InputFile{urlID: url}
}

// FileReader returns InputFile that needs to be uploaded.
func FileReader(name string, data io.Reader) *InputFile {
	return &InputFile{
		name: name,
		data: data,
	}
}

// FileBytes returns InputFile as bytes that needs to be uploaded.
func FileBytes(name string, data []byte) *InputFile {
	return &InputFile{
		name: name,
		data: bytes.NewReader(data),
	}
}

// InputMedia represents the content of a media message to be sent.
type InputMedia interface {
	inputMediaType() string
}

var (
	_ InputMedia = InputMediaPhoto{}
	_ InputMedia = InputMediaVideo{}
	_ InputMedia = InputMediaAnimation{}
	_ InputMedia = InputMediaAudio{}
	_ InputMedia = InputMediaDocument{}
)

// Media represents list of InputMedia.
type Media []InputMedia

// MarshalJSON implements json.Marshaler.
func (media Media) MarshalJSON() ([]byte, error) {
	data := make([]json.RawMessage, len(media))
	for i := range media {
		var err error
		data[i], err = mergeJSON(media[i], struct {
			Type string `json:"type"`
		}{media[i].inputMediaType()})
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(data)
}

// BaseInputMedia type.
type BaseInputMedia struct {
	Media     *InputFile      `json:"media"`
	Caption   string          `json:"caption,omitempty"`
	ParseMode ParseMode       `json:"parse_mode,omitempty"`
	Entities  []MessageEntity `json:"caption_entities,omitempty"`
}

// NewInputMediaPhoto creates new InputMediaPhoto object.
func NewInputMediaPhoto(file *InputFile) *InputMediaPhoto {
	return &InputMediaPhoto{
		BaseInputMedia: BaseInputMedia{
			Media: file,
		},
	}
}

// InputMediaPhoto represents a photo to be sent.
type InputMediaPhoto struct {
	BaseInputMedia
}

func (InputMediaPhoto) inputMediaType() string {
	return "photo"
}

// InputMediaVideo represents a video to be sent.
type InputMediaVideo struct {
	BaseInputMedia
	Thumbnail         *InputFile `json:"thumb,omitempty"`
	Width             int        `json:"width,omitempty"`
	Height            int        `json:"height,omitempty"`
	Duration          int        `json:"duration,omitempty"`
	SupportsStreaming bool       `json:"supports_streaming,omitempty"`
}

func (InputMediaVideo) inputMediaType() string {
	return "video"
}

// InputMediaAnimation represents an animation file
// (GIF or H.264/MPEG-4 AVC video without sound) to be sent.
type InputMediaAnimation struct {
	BaseInputMedia
	Thumbnail *InputFile `json:"thumb,omitempty"`
	Width     int        `json:"width,omitempty"`
	Height    int        `json:"height,omitempty"`
	Duration  int        `json:"duration,omitempty"`
}

func (InputMediaAnimation) inputMediaType() string {
	return "animation"
}

// InputMediaAudio represents an audio file to be treated as music to be sent.
type InputMediaAudio struct {
	BaseInputMedia
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
	BaseInputMedia
	Thumbnail            *InputFile `json:"thumb,omitempty"`
	DisableTypeDetection bool       `json:"disable_content_type_detection,omitempty"`
}

func (InputMediaDocument) inputMediaType() string {
	return "document"
}
