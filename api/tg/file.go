package tg

// FileData contains common files info.
type FileData struct {
	FileID   string `json:"file_id"`
	UniqueID string `json:"file_unique_id"`
	FileSize int64  `json:"file_size"`
}

// PhotoSize represents one size of a photo or a
// file / sticker thumbnail.
type PhotoSize struct {
	FileData

	Width  int `json:"width"`
	Height int `json:"height"`
}

// Animation object represents a animation file.
type Animation struct {
	FileData

	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Duration  int        `json:"duration"`
	Thumbnail *PhotoSize `json:"thumbnail"`
	FileName  string     `json:"file_name"`
	MimeType  string     `json:"mime_type"`
}

// Audio represents an audio file to be treated as music
// by the Telegram clients.
type Audio struct {
	FileData

	Duration  int        `json:"duration"`
	Performer string     `json:"performer"`
	Title     string     `json:"title"`
	FileName  string     `json:"file_name"`
	MimeType  string     `json:"mime_type"`
	Thumbnail *PhotoSize `json:"thumbnail"`
}

// Document represents a general file
// (as opposed to photos, voice messages and audio files).
type Document struct {
	FileData

	Thumbnail *PhotoSize `json:"thumbnail"`
	FileName  string     `json:"file_name"`
	MimeType  string     `json:"mime_type"`
}

// Video represents a video file.
type Video struct {
	FileData

	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Duration  int        `json:"duration"`
	Thumbnail *PhotoSize `json:"thumbnail"`
	FileName  string     `json:"file_name"`
	MimeType  string     `json:"mime_type"`
}

// VideoNote represents a video message.
type VideoNote struct {
	FileData

	Length    int        `json:"length"`
	Duration  int        `json:"duration"`
	Thumbnail *PhotoSize `json:"thumbnail"`
}

// Voice represents a voice note.
type Voice struct {
	FileData

	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
}

// Sticker represents a sticker.
type Sticker struct {
	FileData

	Type            StickerType   `json:"type"`
	Width           int           `json:"width"`
	Height          int           `json:"height"`
	IsAnimated      bool          `json:"is_animated"`
	IsVideo         bool          `json:"is_video"`
	Thumbnail       *PhotoSize    `json:"thumbnail"`
	Emoji           string        `json:"emoji"`
	SetName         string        `json:"set_name"`
	Premium         *File         `json:"premium_animation"`
	MaskPosition    *MaskPosition `json:"mask_position"`
	CustomEmojiID   string        `json:"custom_emoji_id"`
	NeedsRepainting bool          `json:"needs_repainting"`
}

// Story represents a story.
type Story struct{}

// File object represents a file ready to be downloaded.
type File struct {
	FileData
	FilePath string `json:"file_path"`
}
