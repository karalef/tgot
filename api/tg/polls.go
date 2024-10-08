package tg

// Poll contains information about a poll.
type Poll struct {
	ID                  string          `json:"id"`
	Question            string          `json:"question"`
	QuestionEntities    []MessageEntity `json:"question_entities"`
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
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	VoterCount   int             `json:"voter_count"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll.
type PollAnswer struct {
	PollID    string `json:"poll_id"`
	VoterChat *Chat  `json:"voter_chat"`
	User      *User  `json:"user"`
	Options   []int  `json:"option_ids"`
}

// InputPollOption contains information about one answer option in a poll to be sent.
type InputPollOption struct {
	Text      string          `json:"text"`
	ParseMode ParseMode       `json:"text_parse_mode"`
	Entities  []MessageEntity `json:"text_entities"`
}
