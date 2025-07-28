package tg

import (
	"bytes"
	"io"

	"github.com/karalef/tgot/api/internal"
	"github.com/karalef/tgot/api/internal/oneof"
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

// Inputter is an interface for any type that contains Inputtable.
type Inputter interface {
	GetInput() []Inputtable
}

var _ Inputter = InputSticker{}

// InputSticker describes a sticker to be added to a sticker set.
type InputSticker struct {
	Sticker      Inputtable    `json:"sticker"`
	Format       StickerFormat `json:"format"`
	EmojiList    []string      `json:"emoji_list"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	Keywords     []string      `json:"keywords,omitempty"`
}

func (i InputSticker) GetInput() []Inputtable { return []Inputtable{i.Sticker} }

// StickerFormat is a Sticker format.
type StickerFormat string

// all available sticker formats.
const (
	StickerStatic   StickerFormat = "static"
	StickerAnimated StickerFormat = "animated"
	StickerVideo    StickerFormat = "video"
)

// InputMediaData represents any available input media object.
type InputMediaData interface {
	Inputter
	inputMediaType() string
}

var (
	_ Inputter       = (*InputMedia)(nil)
	_ InputMediaData = (*InputMediaPhoto)(nil)
	_ InputMediaData = (*InputMediaVideo)(nil)
	_ InputMediaData = (*InputMediaAnimation)(nil)
	_ InputMediaData = (*InputMediaAudio)(nil)
	_ InputMediaData = (*InputMediaDocument)(nil)
)

// InputMedia represents the content of a media message to be sent.
type InputMedia struct {
	Media     Inputtable      `json:"media"`
	Caption   string          `json:"caption,omitempty"`
	ParseMode ParseMode       `json:"parse_mode,omitempty"`
	Entities  []MessageEntity `json:"caption_entities,omitempty"`
	Data      InputMediaData  `json:"-"`
}

func (i InputMedia) GetInput() []Inputtable {
	return append(i.Data.GetInput(), i.Media)
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
func (*InputMediaPhoto) GetInput() []Inputtable { return nil }

// InputMediaVideo represents a video to be sent.
type InputMediaVideo struct {
	Thumbnail             *InputFile `json:"thumbnail,omitempty"`
	Cover                 *InputFile `json:"cover,omitempty"`
	Start                 Duration   `json:"start_timestamp"`
	Width                 int        `json:"width,omitempty"`
	Height                int        `json:"height,omitempty"`
	Duration              Duration   `json:"duration,omitempty"`
	SupportsStreaming     bool       `json:"supports_streaming,omitempty"`
	HasSpoiler            bool       `json:"has_spoiler,omitempty"`
	ShowCaptionAboveMedia bool       `json:"show_caption_above_media,omitempty"`
}

func (*InputMediaVideo) inputMediaType() string   { return "video" }
func (v *InputMediaVideo) GetInput() []Inputtable { return []Inputtable{v.Thumbnail} }

// InputMediaAnimation represents an animation file
// (GIF or H.264/MPEG-4 AVC video without sound) to be sent.
type InputMediaAnimation struct {
	Thumbnail             *InputFile `json:"thumbnail,omitempty"`
	Width                 int        `json:"width,omitempty"`
	Height                int        `json:"height,omitempty"`
	Duration              Duration   `json:"duration,omitempty"`
	HasSpoiler            bool       `json:"has_spoiler,omitempty"`
	ShowCaptionAboveMedia bool       `json:"show_caption_above_media,omitempty"`
}

func (*InputMediaAnimation) inputMediaType() string   { return "animation" }
func (a *InputMediaAnimation) GetInput() []Inputtable { return []Inputtable{a.Thumbnail} }

// InputMediaAudio represents an audio file to be treated as music to be sent.
type InputMediaAudio struct {
	Thumbnail *InputFile `json:"thumbnail,omitempty"`
	Duration  Duration   `json:"duration,omitempty"`
	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
}

func (*InputMediaAudio) inputMediaType() string   { return "audio" }
func (a *InputMediaAudio) GetInput() []Inputtable { return []Inputtable{a.Thumbnail} }

// InputMediaDocument represents a general file to be sent.
type InputMediaDocument struct {
	Thumbnail            *InputFile `json:"thumbnail,omitempty"`
	DisableTypeDetection bool       `json:"disable_content_type_detection,omitempty"`
}

func (*InputMediaDocument) inputMediaType() string   { return "document" }
func (d *InputMediaDocument) GetInput() []Inputtable { return []Inputtable{d.Thumbnail} }

// InputPaidMediaData represents any available input paid media object.
type InputPaidMediaData interface {
	Inputter
	inputPaidMediaType() string
}

var (
	_ Inputter           = InputPaidMedia{}
	_ InputPaidMediaData = (*InputPaidMediaPhoto)(nil)
	_ InputPaidMediaData = (*InputPaidMediaVideo)(nil)
)

// InputPaidMedia describes the paid media to be sent.
type InputPaidMedia struct {
	Media Inputtable         `json:"media"`
	Data  InputPaidMediaData `json:"-"`
}

func (i InputPaidMedia) GetInput() []Inputtable {
	return append(i.Data.GetInput(), i.Media)
}

func (i InputPaidMedia) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(i.Data, struct {
		Type  string     `json:"type"`
		Media Inputtable `json:"media"`
	}{i.Data.inputPaidMediaType(), i.Media})
}

// InputPaidMediaPhoto is the paid media to send is a photo.
type InputPaidMediaPhoto struct{}

func (*InputPaidMediaPhoto) inputPaidMediaType() string { return "photo" }
func (*InputPaidMediaPhoto) GetInput() []Inputtable     { return nil }

// InputPaidMediaVideo is the paid media to send is a video.
type InputPaidMediaVideo struct {
	Thumbnail         *InputFile `json:"thumbnail,omitempty"`
	Cover             *InputFile `json:"cover,omitempty"`
	Start             Duration   `json:"start_timestamp"`
	Width             int        `json:"width,omitempty"`
	Height            int        `json:"height,omitempty"`
	Duration          Duration   `json:"duration,omitempty"`
	SupportsStreaming bool       `json:"supports_streaming,omitempty"`
}

func (*InputPaidMediaVideo) inputPaidMediaType() string { return "video" }
func (v *InputPaidMediaVideo) GetInput() []Inputtable   { return []Inputtable{v.Thumbnail, v.Cover} }

// InputProfilePhotoData represents any available input profile photo object.
type InputProfilePhotoData interface {
	Media() *InputFile
	inputProfilePhotoType() string
}

var (
	_ Inputter              = InputProfilePhoto{}
	_ InputProfilePhotoData = InputProfilePhotoStatic{}
	_ InputProfilePhotoData = InputProfilePhotoAnimated{}
)

// InputProfilePhoto describes a profile photo to set.
type InputProfilePhoto struct {
	Data InputProfilePhotoData
}

func (i InputProfilePhoto) GetInput() []Inputtable { return []Inputtable{i.Data.Media()} }
func (i InputProfilePhoto) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(i.Data, oneof.IDTypeType{Type: i.Data.inputProfilePhotoType()})
}

// InputProfilePhotoStatic is a static profile photo in the .JPG format.
type InputProfilePhotoStatic struct {
	Photo *InputFile `json:"photo"`
}

func (InputProfilePhotoStatic) inputProfilePhotoType() string { return "static" }
func (i InputProfilePhotoStatic) Media() *InputFile           { return i.Photo }

// InputProfilePhotoAnimated is an animated profile photo in the MPEG4 format.
type InputProfilePhotoAnimated struct {
	Animation *InputFile `json:"animation"`
	MainFrame float64    `json:"main_frame_timestamp,omitempty"`
}

func (InputProfilePhotoAnimated) inputProfilePhotoType() string { return "animated" }
func (i InputProfilePhotoAnimated) Media() *InputFile           { return i.Animation }

// InputStoryContentData represents any available input story content object.
type InputStoryContentData interface {
	Media() *InputFile
	inputStoryContentType() string
}

var (
	_ Inputter              = InputStoryContent{}
	_ InputStoryContentData = InputStoryContentPhoto{}
	_ InputStoryContentData = InputStoryContentVideo{}
)

// InputStoryContent describes the content of a story to post.
type InputStoryContent struct {
	Data InputStoryContentData
}

func (i InputStoryContent) GetInput() []Inputtable { return []Inputtable{i.Data.Media()} }
func (i InputStoryContent) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(i.Data, oneof.IDTypeType{Type: i.Data.inputStoryContentType()})
}

// InputStoryContentPhoto describes a photo to post as a story.
type InputStoryContentPhoto struct {
	Photo *InputFile `json:"photo"`
}

func (InputStoryContentPhoto) inputStoryContentType() string { return "photo" }
func (i InputStoryContentPhoto) Media() *InputFile           { return i.Photo }

// InputStoryContentVideo describes a video to post as a story.
type InputStoryContentVideo struct {
	Video       *InputFile `json:"video"`
	Duration    float64    `json:"duration,omitempty"`
	CoverFrame  float64    `json:"cover_frame_timestamp,omitempty"`
	IsAnimation bool       `json:"is_animation"`
}

func (InputStoryContentVideo) inputStoryContentType() string { return "video" }
func (i InputStoryContentVideo) Media() *InputFile           { return i.Video }
