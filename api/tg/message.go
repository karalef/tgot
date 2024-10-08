package tg

import (
	"time"
)

// Message represents a message.
type Message struct {
	ID                        int                        `json:"message_id"`
	MessageThreadID           int                        `json:"message_thread_id"`
	From                      *User                      `json:"from"`
	SenderChat                *Chat                      `json:"sender_chat"`
	Date                      int64                      `json:"date"`
	Chat                      *Chat                      `json:"chat"`
	ForwardOrigin             *MessageOrigin             `json:"forward_origin"`
	IsTopicMessage            bool                       `json:"is_topic_message"`
	IsAutomaticForward        bool                       `json:"is_automatic_forward"`
	ReplyTo                   *Message                   `json:"reply_to_message"`
	ExternalReply             *ExternalReplyInfo         `json:"external_reply"`
	Quote                     *TextQuote                 `json:"quote"`
	ViaBot                    *User                      `json:"via_bot"`
	EditDate                  int64                      `json:"edit_date"`
	HasProtectedContent       bool                       `json:"has_protected_content"`
	MediaGroupID              string                     `json:"media_group_id"`
	AuthorSignature           string                     `json:"author_signature"`
	Text                      string                     `json:"text"`
	Entities                  []MessageEntity            `json:"entities"`
	LinkPreviewOptions        *LinkPreviewOptions        `json:"link_preview_options"`
	Animation                 *Animation                 `json:"animation"`
	Audio                     *Audio                     `json:"audio"`
	Document                  *Document                  `json:"document"`
	Photo                     []PhotoSize                `json:"photo"`
	Sticker                   *Sticker                   `json:"sticker"`
	Story                     *Story                     `json:"story"`
	Video                     *Video                     `json:"video"`
	VideoNote                 *VideoNote                 `json:"video_note"`
	Voice                     *Voice                     `json:"voice"`
	Caption                   string                     `json:"caption"`
	CaptionEntities           []MessageEntity            `json:"caption_entities"`
	HasMediaSpoiler           bool                       `json:"has_media_spoiler"`
	Contact                   *Contact                   `json:"contact"`
	Dice                      *Dice                      `json:"dice"`
	Game                      *Game                      `json:"game"`
	Poll                      *Poll                      `json:"poll"`
	Venue                     *Venue                     `json:"venue"`
	Location                  *Location                  `json:"location"`
	NewChatMembers            []*User                    `json:"new_chat_members"`
	LeftChatMember            *User                      `json:"left_chat_member"`
	NewChatTitle              string                     `json:"new_chat_title"`
	NewChatPhoto              []PhotoSize                `json:"new_chat_photo"`
	DeleteChatPhoto           bool                       `json:"delete_chat_photo"`
	GroupCreated              bool                       `json:"group_chat_created"`
	SuperGroupCreated         bool                       `json:"supergroup_chat_created"`
	ChannelCreated            bool                       `json:"channel_chat_created"`
	AutoDeleteTimerChanged    *AutoDeleteTimer           `json:"message_auto_delete_timer_changed"`
	MigrateTo                 int64                      `json:"migrate_to_chat_id"`
	MigrateFrom               int64                      `json:"migrate_from_chat_id"`
	PinnedMessage             *MaybeInaccessibleMessage  `json:"pinned_message"`
	Invoice                   *Invoice                   `json:"invoice"`
	SuccessfulPayment         *SuccessfulPayment         `json:"successful_payment"`
	UsersShared               *UsersShared               `json:"users_shared"`
	ChatShared                *ChatShared                `json:"chat_shared"`
	ConnectedWebsite          string                     `json:"connected_website"`
	PassportData              *PassportData              `json:"passport_data"`
	ProximityAlert            *ProximityAlert            `json:"proximity_alert_triggered"`
	ForumTopicCreated         *ForumTopicCreated         `json:"forum_topic_created"`
	ForumTopicEdited          *ForumTopicEdited          `json:"forum_topic_edited"`
	ForumTopicClosed          *ForumTopicClosed          `json:"forum_topic_closed"`
	ForumTopicReopened        *ForumTopicReopened        `json:"forum_topic_reopened"`
	GeneralForumTopicHidden   *GeneralForumTopicHidden   `json:"general_forum_topic_hidden"`
	GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden"`
	GiveawayCreated           *GiveawayCreated           `json:"giveaway_created"`
	Giveaway                  *Giveaway                  `json:"giveaway"`
	GiveawayWinners           *GiveawayWinners           `json:"giveaway_winners"`
	GiveawayCompleted         *GiveawayCompleted         `json:"giveaway_completed"`
	VideoChatScheduled        *VideoChatScheduled        `json:"video_chat_scheduled"`
	VideoChatStarted          *VideoChatStarted          `json:"video_chat_started"`
	VideoChatEnded            *VideoChatEnded            `json:"video_chat_ended"`
	VideoChatInvited          *VideoChatInvited          `json:"video_chat_participants_invited"`
	WebAppData                *WebAppData                `json:"web_app_data"`
	ReplyMarkup               *InlineKeyboardMarkup      `json:"reply_markup"`
}

// Time converts unixtime to time.Time.
func (m *Message) Time() time.Time {
	return time.Unix(m.Date, 0)
}

// MessageID represents a unique message identifier.
type MessageID struct {
	ID int `json:"message_id"`
}

// MessageEntity represents one special entitty in a text message.
// For example, hashtag, usernames, URLs, etc.
type MessageEntity struct {
	Type          EntityType `json:"type"`
	Offset        int        `json:"offset"` // in UTF-16
	Length        int        `json:"length"`
	URL           string     `json:"url,omitmepty"`
	User          *User      `json:"user,omitmepty"`
	Language      string     `json:"language,omitmepty"`
	CustomEmojiID string     `json:"custom_emoji_id,omitempty"`
}

// EntityType is a MessageEntity type.
type EntityType string

// all available entity types.
const (
	EntityMention       EntityType = "mention"
	EntityHashtag       EntityType = "hashtag"
	EntityCashtag       EntityType = "cashtag"
	EntityCommand       EntityType = "bot_command"
	EntityURL           EntityType = "url"
	EntityEmail         EntityType = "email"
	EntityPhone         EntityType = "phone_number"
	EntityBold          EntityType = "bold"
	EntityItalic        EntityType = "italic"
	EntityUnderline     EntityType = "underline"
	EntityStrikethrough EntityType = "strikethrough"
	EntitySpoiler       EntityType = "spoiler"
	EntityCode          EntityType = "code"
	EntityCodeBlock     EntityType = "pre"
	EntityTextLink      EntityType = "text_link"
	EntityTextMention   EntityType = "text_mention"
	EntityCustomEmoji   EntityType = "custom_emoji"
	EntityBlockquote    EntityType = "blockquote"
)

// ProximityAlert represents the content of a service message,
// sent whenever a user in the chat triggers a proximity alert
// set by another user.
type ProximityAlert struct {
	Traveler *User `json:"traveler"`
	Watcher  *User `json:"watcher"`
	Distance int   `json:"distance"`
}

// AutoDeleteTimer represents a service message about a change
// in auto-delete timer settings.
type AutoDeleteTimer struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

// VideoChatScheduled represents a service message about a video chat scheduled in the chat.
type VideoChatScheduled struct {
	Time int64 `json:"start_time"`
}

// VideoChatStarted represents a service message about a video chat started in the chat.
type VideoChatStarted struct{}

// VideoChatEnded represents a service message about a video chat ended in the chat.
type VideoChatEnded struct {
	Duration int64 `json:"duration"`
}

// VideoChatInvited represents a service message about new members invited to a video chat.
type VideoChatInvited struct {
	Users []User `json:"users"`
}

// ForumTopicCreated represents a service message about a new forum topic created in the chat.
type ForumTopicCreated struct {
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

// ForumTopicClosed represents a service message about a forum topic closed in the chat.
// Currently holds no information.
type ForumTopicClosed struct{}

// ForumTopicEdited represents a service message about an edited forum topic.
type ForumTopicEdited struct {
	Name              string `json:"name"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

// ForumTopicReopened represents a service message about a forum topic reopened in the chat.
// Currently holds no information.
type ForumTopicReopened struct{}

// GeneralForumTopicHidden represents a service message about General forum topic hidden in the chat.
// Currently holds no information.
type GeneralForumTopicHidden struct{}

// GeneralForumTopicUnhidden represents a service message about General forum topic unhidden in the chat.
// Currently holds no information.
type GeneralForumTopicUnhidden struct{}

// UsersShared contains information about the users whose identifiers were shared with the bot using a KeyboardButtonRequestUsers button.
type UsersShared struct {
	RequestID int          `json:"request_id"`
	Users     []SharedUser `json:"users"`
}

// SharedUser contains information about a user that was shared with the bot using a KeyboardButtonRequestUsers button.
type SharedUser struct {
	UserID    int64       `json:"user_id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Username  string      `json:"username"`
	Photo     []PhotoSize `json:"photo"`
}

// ChatShared contains information about the chat whose identifier was shared with the bot using a KeyboardButtonRequestChat button.
type ChatShared struct {
	RequestID int   `json:"request_id"`
	ChatID    int64 `json:"chat_id"`
}

// WriteAccessAllowed represents a service message about a user allowing a bot to write messages after adding the bot to the attachment menu or launching a Web App from a link.
type WriteAccessAllowed struct {
	FromRequest        bool   `json:"from_request"`
	WebAppName         string `json:"web_app_name"`
	FromAttachmentMenu bool   `json:"from_attachment_menu"`
}

// ReactionType describes the type of a reaction.
type ReactionType struct {
	Type string `json:"type"`

	// ReactionTypeEmoji is the reaction is based on an emoji.
	Emoji string `json:"emoji"`

	// ReactionTypeCustomEmoji is the reaction is based on a custom emoji.
	CustomEmojiID string `json:"custom_emoji_id"`
}

// ExternalReplyInfo contains information about a message that is being replied to,
// which may come from another chat or forum topic.
type ExternalReplyInfo struct {
	Origin             MessageOrigin       `json:"origin"`
	Chat               *Chat               `json:"chat"`
	MessageID          int                 `json:"message_id"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options"`
	Animation          *Animation          `json:"animation"`
	Audio              *Audio              `json:"audio"`
	Document           *Document           `json:"document"`
	Photo              []PhotoSize         `json:"photo"`
	Sticker            *Sticker            `json:"sticker"`
	Story              *Story              `json:"story"`
	Video              *Video              `json:"video"`
	VideoNote          *VideoNote          `json:"video_note"`
	Voice              *Voice              `json:"voice"`
	HasMediaSpoiler    bool                `json:"has_media_spoiler"`
	Contact            *Contact            `json:"contact"`
	Dice               *Dice               `json:"dice"`
	Game               *Game               `json:"game"`
	Giveaway           *Giveaway           `json:"giveaway"`
	GiveawayWinners    *GiveawayWinners    `json:"giveaway_winners"`
	Invoice            *Invoice            `json:"invoice"`
	Location           *Location           `json:"location"`
	Poll               *Poll               `json:"poll"`
	Venue              *Venue              `json:"venue"`
}

// TextQuote contains information about the quoted part of a message that
// is replied to by the given message.
type TextQuote struct {
	Text     string          `json:"text"`
	Entities []MessageEntity `json:"entities"`
	Position int             `json:"position"`
	IsManual bool            `json:"is_manual"`
}

// GiveawayCreated represents a service message about the creation of a scheduled giveaway.
type GiveawayCreated struct {
	PrizeStarCount int `json:"prize_star_count"`
}

// Giveaway represents a message about a scheduled giveaway.
type Giveaway struct {
	Chats                         []Chat   `json:"chats"`
	WinnersSelectionDate          int64    `json:"winners_selection_date"`
	WinnerCount                   int      `json:"winner_count"`
	OnlyNewMembers                bool     `json:"only_new_members"`
	HasPublicWinners              bool     `json:"has_public_winners"`
	PrizeDescription              string   `json:"prize_description"`
	CountryCodes                  []string `json:"country_codes"`
	PrizeStarCount                int      `json:"prize_star_count"`
	PremiumSubscriptionMonthCount int      `json:"premium_subscription_month_count"`
}

// GiveawayWinners represents a message about the completion of a giveaway with public winners.
type GiveawayWinners struct {
	Chat                          Chat   `json:"chat"`
	GiveawayMessageID             int    `json:"giveaway_message_id"`
	WinnersSelectionDate          int64  `json:"winners_selection_date"`
	WinnerCount                   int    `json:"winner_count"`
	Winners                       []User `json:"winners"`
	AdditionalChatCount           int    `json:"additional_chat_count"`
	PrizeStarCount                int    `json:"prize_star_count"`
	PremiumSubscriptionMonthCount int    `json:"premium_subscription_month_count"`
	UnclaimedPrizeCount           int    `json:"unclaimed_prize_count"`
	OnlyNewMembers                bool   `json:"only_new_members"`
	WasRefunded                   bool   `json:"was_refunded"`
	PrizeDescription              string `json:"prize_description"`
}

// GiveawayCompleted represents a service message about the completion of a giveaway without public winners.
type GiveawayCompleted struct {
	WinnerCount         int     `json:"winner_count"`
	UnclaimedPrizeCount int     `json:"unclaimed_prize_count"`
	GiveawayMessage     Message `json:"giveaway_message"`
	IsStarGiveaway      bool    `json:"is_star_giveaway"`
}

// MessageOriginType represents the type of message origin.
type MessageOriginType string

const (
	MessageOriginUser       MessageOriginType = "user"
	MessageOriginHiddenUser MessageOriginType = "hidden_user"
	MessageOriginChat       MessageOriginType = "chat"
	MessageOriginChannel    MessageOriginType = "channel"
)

// MessageOrigin describes the origin of a message.
type MessageOrigin struct {
	Type MessageOriginType `json:"type"`
	Date int64             `json:"date"`

	// user
	SenderUser *User `json:"sender_user"`

	// hidden_user
	SenderUserName string `json:"sender_user_name"`

	// chat
	SenderChat *Chat `json:"sender_chat"`

	// channel
	Chat      *Chat `json:"chat"`
	MessageID int   `json:"message_id"`

	// chat & channel
	AuthorSignature string `json:"author_signature"`
}

// InaccessibleMessage describes a message that was deleted or is otherwise inaccessible to the bot.
type InaccessibleMessage struct {
	Chat Chat  `json:"chat"`
	ID   int   `json:"message_id"`
	Date int64 `json:"date"`
}

// MaybeInaccessibleMessage describes a message that can be inaccessible to the bot.
type MaybeInaccessibleMessage struct {
	*InaccessibleMessage
	*Message
}

func (m MaybeInaccessibleMessage) ID() int {
	if m.InaccessibleMessage != nil {
		return m.InaccessibleMessage.ID
	}
	return m.Message.ID
}

func (m MaybeInaccessibleMessage) Date() int64 {
	if m.InaccessibleMessage != nil {
		return m.InaccessibleMessage.Date
	}
	return m.Message.Date
}

func (m MaybeInaccessibleMessage) Chat() *Chat {
	if m.InaccessibleMessage != nil {
		return &m.InaccessibleMessage.Chat
	}
	return m.Message.Chat
}

func (m MaybeInaccessibleMessage) IsInaccessible() bool {
	return m.Message == nil
}

// ParseMode type.
type ParseMode string

// all available parse modes.
const (
	Markdown   ParseMode = "Markdown"
	MarkdownV2 ParseMode = "MarkdownV2"
	HTML       ParseMode = "HTML"
)
