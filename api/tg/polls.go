package tg

// Poll contains information about a poll.
type Poll struct {
	ID                  string          `json:"id"`
	Question            string          `json:"question"`
	Options             []PollOption    `json:"options"`
	VoterCount          int             `json:"total_voter_count"`
	IsClosed            bool            `json:"is_closed"`
	IsAnonymous         bool            `json:"is_anonymous"`
	Type                PollType        `json:"type"`
	MultipleAnswers     bool            `json:"allows_multiple_answers"`
	CorrectOption       int             `json:"correct_option_id"`
	Explanation         string          `json:"explanation"`
	ExplanationEntities []MessageEntity `json:"explanation_entities"`
	OpenPeriod          int             `json:"open_period"`
	CloseDate           int64           `json:"close_date"`
}

// PollType represents poll type.
type PollType string

// all available poll types.
const (
	PollQuiz    PollType = "quiz"
	PollRegular PollType = "regular"
)

// PollOption contains information about one answer option in a poll.
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID  string `json:"poll_id"`
	User    *User  `json:"user"`
	Options []int  `json:"option_ids"`
}
