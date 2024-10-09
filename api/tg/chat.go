package tg

// Chat represents a chat.
type Chat struct {
	ID        int64    `json:"id"`
	Type      ChatType `json:"type"`
	Title     string   `json:"title"`
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	IsForum   bool     `json:"is_forum"`
}

// ChatFullInfo contains full information about a chat.
type ChatFullInfo struct {
	Chat
	AccentColorID                  uint8                 `json:"accent_color_id"`
	MaxReactionCount               uint                  `json:"max_reaction_count"`
	Photo                          *ChatPhoto            `json:"photo"`
	ActiveUsernames                []string              `json:"active_usernames"`
	Birthdate                      *Birthdate            `json:"birthdate"`
	BusinessIntro                  *BusinessIntro        `json:"business_intro"`
	BusinessLocation               *BusinessLocation     `json:"business_location"`
	BusinessOpeningHours           *BusinessOpeningHours `json:"business_opening_hours"`
	PersonalChat                   *Chat                 `json:"personal_chat"`
	AvailableReactions             []ReactionType        `json:"available_reactions"`
	BackgroundCustomEmojiID        string                `json:"background_custom_emoji_id"`
	ProfileAccentColorID           uint8                 `json:"profile_accent_color_id"`
	ProfileBackgroundCustomEmojiID string                `json:"profile_background_custom_emoji_id"`
	EmojiStatusCustomEmoji         string                `json:"emoji_status_custom_emoji_id"`
	EmojiStatusExpirationDate      int64                 `json:"emoji_status_expiration_date"`
	Bio                            string                `json:"bio"`
	HasPrivateForwards             bool                  `json:"has_private_forwards"`
	HasRestrictedVoiceAndVideo     bool                  `json:"has_restricted_voice_and_video_messages"`
	JoinToSend                     bool                  `json:"join_to_send_messages"`
	JoinByRequest                  bool                  `json:"join_by_request"`
	Description                    string                `json:"description"`
	InviteLink                     string                `json:"invite_link"`
	PinnedMessage                  *Message              `json:"pinned_message"`
	Permissions                    *ChatPermissions      `json:"permissions"`
	CanSendPaidMedia               bool                  `json:"can_send_paid_media"`
	SlowModeDelay                  int                   `json:"slow_mode_delay"`
	UnrestrictBoostCount           int                   `json:"unrestrict_boost_count"`
	AutoDeleteTime                 int                   `json:"message_auto_delete_time"`
	HasAgressiveAntiSpam           bool                  `json:"has_aggressive_anti_spam_enabled"`
	HasHiddenMembers               bool                  `json:"has_hidden_members"`
	HasProtectedContent            bool                  `json:"has_protected_content"`
	HasVisibleHistory              bool                  `json:"has_visible_history"`
	StickerSetName                 string                `json:"sticker_set_name"`
	CanSetStickerSet               bool                  `json:"can_set_sticker_set"`
	CustomEmojiStickerSetName      string                `json:"custom_emoji_sticker_set_name"`
	LinkedChatID                   int64                 `json:"linked_chat_id"`
	Location                       *ChatLocation         `json:"location"`
}

// Is returns true if the chat type matches the one specified.
func (c *Chat) Is(t ChatType) bool {
	return c.Type == t
}

// ChatType represents one of the possible chat types.
type ChatType string

// all available chat types.
const (
	ChatSender     ChatType = "sender"
	ChatPrivate    ChatType = "private"
	ChatGroup      ChatType = "group"
	ChatSuperGroup ChatType = "supergroup"
	ChatChannel    ChatType = "channel"
)

// ChatAction represents one the possible chat actions.
type ChatAction string

// all available chat actions.
const (
	ActionTyping          ChatAction = "typing"
	ActionUploadPhoto     ChatAction = "upload_photo"
	ActionRecordVideo     ChatAction = "record_video"
	ActionUploadVideo     ChatAction = "upload_video"
	ActionRecordVoice     ChatAction = "record_voice"
	ActionUploadVoice     ChatAction = "upload_voice"
	ActionUploadDocument  ChatAction = "upload_document"
	ActionChooseStocker   ChatAction = "choose_sticker"
	ActionFindLocation    ChatAction = "find_location"
	ActionRecordVideoNote ChatAction = "record_video_note"
	ActionUploadVideoNote ChatAction = "upload_video_note"
)

// ChatPhoto represents a chat photo.
type ChatPhoto struct {
	// 160x160 chat photo
	SmallFileID   string `json:"small_file_id"`
	SmallUniqueID string `json:"small_file_unique_id"`

	// 640x640 chat photo
	BigFileID   string `json:"big_file_id"`
	BigUniqueID string `json:"big_file_unique_id"`
}

// ChatPermissions describes actions that a non-administrator user is allowed to take in a chat.
type ChatPermissions struct {
	CanSendMessages   bool `json:"can_send_messages,omitempty"`
	CanSendAudios     bool `json:"can_send_audios,omitempty"`
	CanSendDocuments  bool `json:"can_send_documents,omitempty"`
	CanSendPhotos     bool `json:"can_send_photos,omitempty"`
	CanSendVideos     bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls      bool `json:"can_send_polls,omitempty"`
	CanSendOther      bool `json:"can_send_other_messages,omitempty"`
	CanAddPreviews    bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo     bool `json:"can_change_info,omitempty"`
	CanInviteUsers    bool `json:"can_invite_users,omitempty"`
	CanPinMessages    bool `json:"can_pin_messages,omitempty"`
	CanManageTopics   bool `json:"can_manage_topics,omitempty"`
}

// ChatAdministratorRights represents the rights of an administrator in a chat.
type ChatAdministratorRights struct {
	IsAnonymous         bool `json:"is_anonymous"`
	CanManageChat       bool `json:"can_manage_chat"`
	CanDeleteMessages   bool `json:"can_delete_messages"`
	CanManageVideoChats bool `json:"can_manage_video_chats"`
	CanRestrictMembers  bool `json:"can_restrict_members"`
	CanPromoteMembers   bool `json:"can_promote_members"`
	CanChangeInfo       bool `json:"can_change_info"`
	CanInviteUsers      bool `json:"can_invite_users"`
	CanPostStories      bool `json:"can_post_stories"`
	CanEditStories      bool `json:"can_edit_stories"`
	CanDeleteStories    bool `json:"can_delete_stories"`
	CanPostMessages     bool `json:"can_post_messages,omitempty"`
	CanEditMessages     bool `json:"can_edit_messages,omitempty"`
	CanPinMessages      bool `json:"can_pin_messages,omitempty"`
	CanManageTopics     bool `json:"can_manage_topics,omitempty"`
}

// ChatLocation represents a location to which a chat is connected.
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

// ForumTopic represents a forum topic.
type ForumTopic struct {
	MessageThreadID   int    `json:"message_thread_id"`
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id"`
}

// ChatInviteLink represents an invite link for a chat.
type ChatInviteLink struct {
	InviteLink         string `json:"invite_link"`
	Creator            *User  `json:"creator"`
	CreatesJoinRequest bool   `json:"creates_join_request"`
	IsPrimary          bool   `json:"is_primary"`
	IsRevoked          bool   `json:"is_revoked"`
	Name               string `json:"name,omitempty"`
	ExpireDate         int64  `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	PendingCount       int    `json:"pending_join_request_count,omitempty"`
}

// ChatJoinRequest represents a join request sent to a chat.
type ChatJoinRequest struct {
	Chat       *Chat           `json:"chat"`
	From       *User           `json:"from"`
	UserChatID int64           `json:"user_chat_id"`
	Date       int64           `json:"date"`
	Bio        string          `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

// ChatMemberUpdated represents changes in the status of a chat member.
type ChatMemberUpdated struct {
	Chat                    *Chat           `json:"chat"`
	From                    *User           `json:"from"`
	Date                    int64           `json:"date"`
	Old                     *ChatMember     `json:"old_chat_member"`
	New                     *ChatMember     `json:"new_chat_member"`
	InviteLink              *ChatInviteLink `json:"invite_link"`
	ViaJoinRequest          bool            `json:"via_join_request"`
	ViaChatFolderInviteLink bool            `json:"via_chat_folder_invite_link"`
}

// ChatMember contains information about one member of a chat.
type ChatMember struct {
	Status      MemberStatus `json:"status"`
	User        *User        `json:"user"`
	CustomTitle string       `json:"custom_title"`
	CanBeEdited bool         `json:"can_be_edited"`
	IsMember    bool         `json:"is_member"`
	UntilDate   int64        `json:"until_date"`

	CanSendMessages   bool `json:"can_send_messages"`
	CanSendAudios     bool `json:"can_send_audios,omitempty"`
	CanSendDocuments  bool `json:"can_send_documents,omitempty"`
	CanSendPhotos     bool `json:"can_send_photos,omitempty"`
	CanSendVideos     bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls      bool `json:"can_send_polls"`
	CanSendOther      bool `json:"can_send_other_messages"`
	CanAddPreviews    bool `json:"can_add_web_page_previews"`
	ChatAdministratorRights
}

// Is returns true if the user status matches the one specified.
func (m *ChatMember) Is(s MemberStatus) bool {
	return m.Status == s
}

// MemberStatus represents one of the possible member status.
type MemberStatus string

// all available member statuses.
const (
	MemberCreator    MemberStatus = "creator"
	MemberAdmin      MemberStatus = "administrator"
	MemberDefault    MemberStatus = "member"
	MemberRestricted MemberStatus = "restricted"
	MemberLeft       MemberStatus = "left"
	MemberKicked     MemberStatus = "kicked"
)

// BusinessIntro contains information about the start page settings of a Telegram Business account.
type BusinessIntro struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Sticker *Sticker `json:"sticker"`
}

// BusinessLocation contains information about the location of a Telegram Business account.
type BusinessLocation struct {
	Address  string    `json:"address"`
	Location *Location `json:"location"`
}

// BusinessOpeningHoursInterval describes an interval of time during which a business is open.
type BusinessOpeningHoursInterval struct {
	OpeningMinute int `json:"opening_minute"`
	ClosingMinute int `json:"closing_minute"`
}

// BusinessOpeningHours describes the opening hours of a business.
type BusinessOpeningHours struct {
	TimeZoneName string                         `json:"time_zone_name"`
	OpeningHours []BusinessOpeningHoursInterval `json:"opening_hours"`
}

// Birthdate describes the birthdate of a user.
type Birthdate struct {
	Day   uint8  `json:"day"`
	Month uint8  `json:"month"`
	Year  uint16 `json:"year"`
}
