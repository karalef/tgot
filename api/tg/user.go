package tg

// User represents a Telegram user or bot.
type User struct {
	ID            int64  `json:"id"`
	IsBot         bool   `json:"is_bot"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Username      string `json:"username"`
	LanguageCode  string `json:"language_code"`
	IsPremium     bool   `json:"is_premium"`
	AddedToAttach bool   `json:"added_to_attachment_menu"`

	// Returned only in getMe

	CanJoinGroups        bool `json:"can_join_groups"`
	CanReadAllMessages   bool `json:"can_read_all_group_messages"`
	SupportsInline       bool `json:"supports_inline_queries"`
	CanConnectToBusiness bool `json:"can_connect_to_business"`
	HasMainWebApp        bool `json:"has_main_web_app"`
}

// UserProfilePhotos represent a user's profile pictures.
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}
