package tg

// Update object represents an incoming update.
type Update struct {
	ID                      int                          `json:"update_id"`
	Message                 *Message                     `json:"message"`
	EditedMessage           *Message                     `json:"edited_message"`
	ChannelPost             *Message                     `json:"channel_post"`
	EditedChannelPost       *Message                     `json:"edited_channel_post"`
	BusinessConnection      *BusinessConnection          `json:"business_connection"`
	BusinessMessage         *Message                     `json:"business_message"`
	EditedBusinessMessage   *Message                     `json:"edited_business_message"`
	DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages"`
	CallbackQuery           *CallbackQuery               `json:"callback_query"`
	MessageReaction         *MessageReactionUpdated      `json:"message_reaction"`
	MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count"`
	InlineQuery             *InlineQuery                 `json:"inline_query"`
	InlineChosen            *InlineChosen                `json:"chosen_inline_result"`
	ShippingQuery           *ShippingQuery               `json:"shipping_query"`
	PreCheckoutQuery        *PreCheckoutQuery            `json:"pre_checkout_query"`
	PurchasedPaidMedia      *PaidMediaPurchased          `json:"purchased_paid_media"`
	Poll                    *Poll                        `json:"poll"`
	PollAnswer              *PollAnswer                  `json:"poll_answer"`
	MyChatMember            *ChatMemberUpdated           `json:"my_chat_member"`
	ChatMember              *ChatMemberUpdated           `json:"chat_member"`
	ChatJoinRequest         *ChatJoinRequest             `json:"chat_join_request"`
	ChatBoost               *ChatBoostUpdated            `json:"chat_boost"`
	RemovedChatBoost        *ChatBoostRemoved            `json:"removed_chat_boost"`
}

// BusinessMessagesDeleted is received when messages are deleted from a connected business account.
type BusinessMessagesDeleted struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Chat                 Chat   `json:"chat"`
	MessageIDs           []int  `json:"message_ids"`
}

// CallbackQuery represents an incoming callback query from a callback
// button in an inline keyboard.
type CallbackQuery struct {
	ID              string                    `json:"id"`
	From            User                      `json:"from"`
	Message         *MaybeInaccessibleMessage `json:"message"`
	InlineMessageID string                    `json:"inline_message_id"`
	ChatInstance    string                    `json:"chat_instance"`
	Data            string                    `json:"data"`
	GameShortName   string                    `json:"game_short_name"`
}

// MessageReactionUpdated represents a change of a reaction on a message performed by a user.
type MessageReactionUpdated struct {
	Chat        Chat           `json:"chat"`
	MessageID   int            `json:"message_id"`
	User        *User          `json:"user"`
	ActorChat   *Chat          `json:"actor_chat"`
	Date        int64          `json:"date"`
	OldReaction []ReactionType `json:"old_reaction"`
	NewReaction []ReactionType `json:"new_reaction"`
}

// MessageReactionCountUpdated represents reaction changes on a message with anonymous reactions.
type MessageReactionCountUpdated struct {
	Chat      Chat            `json:"chat"`
	MessageID int             `json:"message_id"`
	Date      int64           `json:"date"`
	Reactions []ReactionCount `json:"reactions"`
}

// ReactionCount represents a reaction added to a message along with the number of times it was added.
type ReactionCount struct {
	Type       ReactionType `json:"type"`
	TotalCount int          `json:"total_count"`
}

// InlineQuery is an incoming inline query. When the user sends
// an empty query, your bot could return some default or
// trending results.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     User      `json:"from"`
	Query    string    `json:"query"` // up to 256 characters
	Offset   string    `json:"offset"`
	ChatType ChatType  `json:"chat_type"`
	Location *Location `json:"location"`
}

// InlineChosen represents a result of an inline query that was chosen
// by the user and sent to their chat partner.
type InlineChosen struct {
	ResultID        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// ShippingQuery contains information about an incoming shipping query.
type ShippingQuery struct {
	ID              string          `json:"id"`
	From            User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             User       `json:"from"`
	Curency          string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// PaidMediaPurchased contains information about a paid media purchase.
type PaidMediaPurchased struct {
	From             User   `json:"from"`
	PaidMediaPayload string `json:"paid_media_payload"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID    string `json:"poll_id"`
	VoterChat *Chat  `json:"voter_chat"`
	User      *User  `json:"user"`
	Options   []int  `json:"option_ids"`
}

// ChatMemberUpdated represents changes in the status of a chat member.
type ChatMemberUpdated struct {
	Chat                    Chat            `json:"chat"`
	From                    User            `json:"from"`
	Date                    int64           `json:"date"`
	Old                     ChatMember      `json:"old_chat_member"`
	New                     ChatMember      `json:"new_chat_member"`
	InviteLink              *ChatInviteLink `json:"invite_link"`
	ViaJoinRequest          bool            `json:"via_join_request"`
	ViaChatFolderInviteLink bool            `json:"via_chat_folder_invite_link"`
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

// ChatBoostUpdated represents a boost added to a chat or changed.
type ChatBoostUpdated struct {
	Chat  Chat      `json:"chat"`
	Boost ChatBoost `json:"boost"`
}

// ChatBoostRemoved represents a boost removed from a chat.
type ChatBoostRemoved struct {
	Chat       Chat            `json:"chat"`
	BoostID    string          `json:"boost_id"`
	RemoveDate int64           `json:"remove_date"`
	Source     ChatBoostSource `json:"source"`
}
