package tg

import "github.com/karalef/tgot/api/internal"

// InlineQueryResulter is an interface for InlineQueryResult.
type InlineQueryResulter interface {
	inlineQueryResult()
}

// InlineQueryResultData represents any available inline query result object.
type InlineQueryResultData interface {
	inlineQueryResultType() string
}

// InlineQueryResult represents one result of an inline query.
type InlineQueryResult[T InlineQueryResultData] struct {
	ID     string
	Result T
}

func (InlineQueryResult[T]) inlineQueryResult() {}

// MarshalJSON implements json.Marshaler.
func (r *InlineQueryResult[T]) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{r.Result.inlineQueryResultType(), r.ID}, r.Result)
}

// InputMessageContent represents the content of a message to be sent
// as a result of an inline query.
type InputMessageContent interface {
	inputMessageContent()
}

func (InputTextMessageContent) inputMessageContent()     {}
func (InputLocationMessageContent) inputMessageContent() {}
func (InputVenueMessageContent) inputMessageContent()    {}
func (InputContactMessageContent) inputMessageContent()  {}
func (InputInvoiceMessageContent) inputMessageContent()  {}

// InputTextMessageContent represents the content of a text message to be sent
// as the result of an inline query.
type InputTextMessageContent struct {
	Text               string             `json:"message_text"`
	ParseMode          ParseMode          `json:"parse_mode,omitempty"`
	Entities           []MessageEntity    `json:"entities,omitempty"`
	LinkPreviewOptions LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

// InputLocationMessageContent represents the content of a location message to be sent
// as the result of an inline query.
type InputLocationMessageContent struct {
	Location
}

// InputVenueMessageContent represents the content of a venue message to be sent
// as the result of an inline query.
type InputVenueMessageContent struct {
	Long            float32 `json:"longitude"`
	Lat             float32 `json:"latitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FoursquareID    string  `json:"foursquare_id,omitempty"`
	FoursquareType  string  `json:"foursquare_type,omitempty"`
	GooglePlaceID   string  `json:"google_place_id,omitempty"`
	GooglePlaceType string  `json:"google_place_type,omitempty"`
}

// InputContactMessageContent represents the content of a contact message to be sent
// as the result of an inline query.
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	Vcard       string `json:"vcard,omitempty"`
}

// InputInvoiceMessageContent represents the content of an invoice message to be sent
// as the result of an inline query.
type InputInvoiceMessageContent struct {
	Title                     string         `json:"title"`
	Description               string         `json:"description"`
	Payload                   string         `json:"payload"`
	ProviderToken             string         `json:"provider_token"`
	Currency                  string         `json:"currency"`
	Prices                    []LabeledPrice `json:"prices"`
	MaxTipAmount              int            `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts       []int          `json:"suggested_tip_amounts,omitempty"`
	ProviderData              string         `json:"provider_data,omitempty"`
	PhotoURL                  string         `json:"photo_url,omitempty"`
	PhotoSize                 int            `json:"photo_size,omitempty"`
	PhotoWidth                int            `json:"photo_width,omitempty"`
	PhotoHeight               int            `json:"photo_height,omitempty"`
	NeedName                  bool           `json:"need_name,omitempty"`
	NeedPhoneNumber           bool           `json:"need_phone_number,omitempty"`
	NeedEmail                 bool           `json:"need_email,omitempty"`
	NeedShippingAddress       bool           `json:"need_shipping_address,omitempty"`
	SendPhoneNumberToProvider bool           `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider       bool           `json:"send_email_to_provider,omitempty"`
	IsFlexible                bool           `json:"is_flexible,omitempty"`
}

// InlineQueryResultsButton represents a button to be shown above inline query results.
type InlineQueryResultsButton struct {
	Text           string      `json:"text"`
	WebApp         *WebAppInfo `json:"web_app,omitempty"`
	StartParameter string      `json:"start_parameter,omitempty"`
}

// InlineQueryResultCachedAudio is an inline query response with cached audio.
type InlineQueryResultCachedAudio struct {
	AudioID             string                `json:"audio_file_id"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           ParseMode             `json:"parse_mode,omitempty"`
	Entities            []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedAudio) inlineQueryResultType() string {
	return "audio"
}

// InlineQueryResultCachedDocument is an inline query response with cached document.
type InlineQueryResultCachedDocument struct {
	Title               string                `json:"title"`
	DocumentID          string                `json:"document_file_id"`
	Description         string                `json:"description,omitempty"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           ParseMode             `json:"parse_mode,omitempty"`
	Entities            []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedDocument) inlineQueryResultType() string {
	return "document"
}

// InlineQueryResultCachedGIF is an inline query response with cached gif.
type InlineQueryResultCachedGIF struct {
	GIFID                 string                `json:"gif_file_id"`
	Title                 string                `json:"title,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultCachedGIF) inlineQueryResultType() string {
	return "gif"
}

// InlineQueryResultCachedMPEG4GIF is an inline query response with cached
// H.264/MPEG-4 AVC video without sound gif.
type InlineQueryResultCachedMPEG4GIF struct {
	MPEG4FileID           string                `json:"mpeg4_file_id"`
	Title                 string                `json:"title,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultCachedMPEG4GIF) inlineQueryResultType() string {
	return "mpeg4_gif"
}

// InlineQueryResultCachedPhoto is an inline query response with cached photo.
type InlineQueryResultCachedPhoto struct {
	PhotoID               string                `json:"photo_file_id"`
	Title                 string                `json:"title,omitempty"`
	Description           string                `json:"description,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultCachedPhoto) inlineQueryResultType() string {
	return "photo"
}

// InlineQueryResultCachedSticker is an inline query response with cached sticker.
type InlineQueryResultCachedSticker struct {
	StickerID           string                `json:"sticker_file_id"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedSticker) inlineQueryResultType() string {
	return "sticker"
}

// InlineQueryResultCachedVideo is an inline query response with cached video.
type InlineQueryResultCachedVideo struct {
	VideoID               string                `json:"video_file_id"`
	Title                 string                `json:"title"`
	Description           string                `json:"description,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultCachedVideo) inlineQueryResultType() string {
	return "video"
}

// InlineQueryResultCachedVoice is an inline query response with cached voice.
type InlineQueryResultCachedVoice struct {
	VoiceID             string                `json:"voice_file_id"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           ParseMode             `json:"parse_mode,omitempty"`
	Entities            []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (InlineQueryResultCachedVoice) inlineQueryResultType() string {
	return "voice"
}

// InlineQueryResultArticle represents a link to an article or web page.
type InlineQueryResultArticle struct {
	Title               string                `json:"title"`
	InputMessageContent InputMessageContent   `json:"input_message_content"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	URL                 string                `json:"url,omitempty"`
	Description         string                `json:"description,omitempty"`
	ThumbnailURL        string                `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                   `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                   `json:"thumbnail_height,omitempty"`
}

func (InlineQueryResultArticle) inlineQueryResultType() string {
	return "article"
}

// InlineQueryResultAudio is an inline query response audio.
type InlineQueryResultAudio struct {
	URL                 string                `json:"audio_url"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           ParseMode             `json:"parse_mode,omitempty"`
	Entities            []MessageEntity       `json:"caption_entities,omitempty"`
	Performer           string                `json:"performer,omitempty"`
	Duration            int                   `json:"audio_duration,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (InlineQueryResultAudio) inlineQueryResultType() string {
	return "audio"
}

// InlineQueryResultContact is an inline query response contact.
type InlineQueryResultContact struct {
	PhoneNumber         string                `json:"phone_number"`
	FirstName           string                `json:"first_name"`
	LastName            string                `json:"last_name,omitempty"`
	VCard               string                `json:"vcard,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailURL        string                `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                   `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                   `json:"thumbnail_height,omitempty"`
}

func (InlineQueryResultContact) inlineQueryResultType() string {
	return "contact"
}

// InlineQueryResultGame is an inline query response game.
type InlineQueryResultGame struct {
	GameShortName string                `json:"game_short_name"`
	ReplyMarkup   *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (InlineQueryResultGame) inlineQueryResultType() string {
	return "game"
}

// InlineQueryResultDocument is an inline query response document.
type InlineQueryResultDocument struct {
	Title               string                `json:"title"`
	URL                 string                `json:"document_url"`
	MimeType            string                `json:"mime_type"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           ParseMode             `json:"parse_mode,omitempty"`
	Entities            []MessageEntity       `json:"caption_entities,omitempty"`
	Description         string                `json:"description,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailURL        string                `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                   `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                   `json:"thumbnail_height,omitempty"`
}

func (InlineQueryResultDocument) inlineQueryResultType() string {
	return "document"
}

// InlineQueryResultGIF is an inline query response GIF.
type InlineQueryResultGIF struct {
	URL                   string                `json:"gif_url"`
	ThumbnailURL          string                `json:"thumbnail_url"`
	Width                 int                   `json:"gif_width,omitempty"`
	Height                int                   `json:"gif_height,omitempty"`
	Duration              int                   `json:"gif_duration,omitempty"`
	ThumbnailMIMEType     string                `json:"thumbnail_mime_type,omitempty"`
	Title                 string                `json:"title,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultGIF) inlineQueryResultType() string {
	return "gif"
}

// InlineQueryResultLocation is an inline query response location.
type InlineQueryResultLocation struct {
	Location
	Title               string                `json:"title"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailURL        string                `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                   `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                   `json:"thumbnail_height,omitempty"`
}

func (InlineQueryResultLocation) inlineQueryResultType() string {
	return "location"
}

// InlineQueryResultMPEG4GIF is an inline query response MPEG4 GIF.
type InlineQueryResultMPEG4GIF struct {
	URL                   string                `json:"mpeg4_url"`
	ThumbnailURL          string                `json:"thumbnail_url"`
	Width                 int                   `json:"mpeg4_width,omitempty"`
	Height                int                   `json:"mpeg4_height,omitempty"`
	Duration              int                   `json:"mpeg4_duration,omitempty"`
	Title                 string                `json:"title,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultMPEG4GIF) inlineQueryResultType() string {
	return "mpeg4_gif"
}

// InlineQueryResultPhoto is an inline query response photo.
type InlineQueryResultPhoto struct {
	URL                   string                `json:"photo_url"`
	ThumbnailURL          string                `json:"thumbnail_url"`
	Width                 int                   `json:"photo_width,omitempty"`
	Height                int                   `json:"photo_height,omitempty"`
	Title                 string                `json:"title,omitempty"`
	Description           string                `json:"description,omitempty"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultPhoto) inlineQueryResultType() string {
	return "photo"
}

// InlineQueryResultVenue is an inline query response venue.
type InlineQueryResultVenue struct {
	Lat                 float32               `json:"latitude"`
	Long                float32               `json:"longitude"`
	Title               string                `json:"title"`
	Address             string                `json:"address"`
	FoursquareID        string                `json:"foursquare_id,omitempty"`
	FoursquareType      string                `json:"foursquare_type,omitempty"`
	GooglePlaceID       string                `json:"google_place_id,omitempty"`
	GooglePlaceType     string                `json:"google_place_type,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailURL        string                `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                   `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                   `json:"thumbnail_height,omitempty"`
}

func (InlineQueryResultVenue) inlineQueryResultType() string {
	return "venue"
}

// InlineQueryResultVideo is an inline query response video.
type InlineQueryResultVideo struct {
	URL                   string                `json:"video_url"`
	MimeType              string                `json:"mime_type"`
	ThumbnailURL          string                `json:"thumbnail_url"`
	Title                 string                `json:"title"`
	Caption               string                `json:"caption,omitempty"`
	ParseMode             ParseMode             `json:"parse_mode,omitempty"`
	Entities              []MessageEntity       `json:"caption_entities,omitempty"`
	Width                 int                   `json:"video_width,omitempty"`
	Height                int                   `json:"video_height,omitempty"`
	Duration              int                   `json:"video_duration,omitempty"`
	Description           string                `json:"description,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ShowCaptionAboveMedia bool                  `json:"show_caption_above_media,omitempty"`
}

func (InlineQueryResultVideo) inlineQueryResultType() string {
	return "video"
}

// InlineQueryResultVoice is an inline query response voice.
type InlineQueryResultVoice struct {
	URL                 string                `json:"voice_url"`
	Title               string                `json:"title"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           ParseMode             `json:"parse_mode,omitempty"`
	Entities            []MessageEntity       `json:"caption_entities,omitempty"`
	Duration            int                   `json:"voice_duration,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (InlineQueryResultVoice) inlineQueryResultType() string {
	return "voice"
}

// PreparedInlineMessage describes an inline message to be sent by a user of a
// Mini App.
type PreparedInlineMessage struct {
	ID      string `json:"id"`
	Expires int64  `json:"expiration_date"`
}
