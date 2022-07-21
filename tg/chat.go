package tg

// Chat represents a chat.
type Chat struct {
	ID        int64    `json:"id"`
	Type      ChatType `json:"type"`
	Title     string   `json:"title"`
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`

	// Returned only in getChat.
	Photo               *ChatPhoto       `json:"photo"`
	Bio                 string           `json:"bio"`
	HasPrivateForwards  bool             `json:"has_private_forwards"`
	JoinToSend          bool             `json:"join_to_send_messages"`
	JoinByRequest       bool             `json:"join_by_request"`
	Description         string           `json:"description"`
	InviteLink          string           `json:"invite_link"`
	PinnedMessage       *Message         `json:"pinned_message"`
	Permissions         *ChatPermissions `json:"permissions"`
	SlowModeDelay       int              `json:"slow_mode_delay"`
	AutoDeleteTime      int              `json:"message_auto_delete_time"`
	HasProtectedContent bool             `json:"has_protected_content"`
	StickerSetNme       string           `json:"sticker_set_name"`
	CanSetStickerSet    bool             `json:"can_set_sticker_set"`
	LinkedChatID        int64            `json:"linked_chat_id"`
	Location            *ChatLocation    `json:"location"`
}

// IsPrivate returns true if chat is private.
func (c *Chat) IsPrivate() bool {
	return c.Type == ChatPrivate
}

// IsGroup returns true if chat is group.
func (c *Chat) IsGroup() bool {
	return c.Type == ChatGroup
}

// IsSuperGroup returns true if chat is supergroup.
func (c *Chat) IsSuperGroup() bool {
	return c.Type == ChatSuperGroup
}

// IsAnyGroup returns true if chat is group or supergroup.
func (c *Chat) IsAnyGroup() bool {
	return c.Type == ChatGroup || c.Type == ChatSuperGroup
}

// IsChannel returns true if chat is channel.
func (c *Chat) IsChannel() bool {
	return c.Type == ChatChannel
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
	CanSendMessages bool `json:"can_send_messages"`
	CanSendMedia    bool `json:"can_send_media_messages"`
	CanSendPolls    bool `json:"can_send_polls"`
	CanSendOther    bool `json:"can_send_other_messages"`
	CanAddPreviews  bool `json:"can_add_web_page_previews"`
	CanChangeInfo   bool `json:"can_change_info"`
	CanInviteUsers  bool `json:"can_invite_users"`
	CanPinMessages  bool `json:"can_pin_messages"`
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
	CanPostMessages     bool `json:"can_post_messages"`
	CanEditMessages     bool `json:"can_edit_messages"`
	CanPinMessages      bool `json:"can_pin_messages"`
}

// ChatLocation represents a location to which a chat is connected.
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}
