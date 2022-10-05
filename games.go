package tgot

import "github.com/karalef/tgot/tg"

// Game contains information about the game to be sent.
type Game struct {
	ShortName string
}

func (g Game) params(p params) {
	p.set("game_short_name", g.ShortName)
}

// SetGameScore contains parameters for setting the game score.
type SetGameScore struct {
	Score       int
	Force       bool
	DisableEdit bool
}

// SetGameScore sets the score of the specified user in a game message.
func (c Context) SetGameScore(sig MessageSignature, userID int64, s SetGameScore) (*tg.Message, error) {
	p := params{}
	p.setInt64("user_id", userID)
	p.setInt("score", s.Score, true)
	p.setBool("force", s.Force)
	p.setBool("disable_edit_message", s.DisableEdit)
	return c.sig(sig, "setGameScore", p)
}

// GetGameHighScores returns data for high score tables.
func (c Context) GetGameHighScores(sig MessageSignature, userID int64) ([]tg.GameHighScore, error) {
	p := params{}.setInt64("user_id", userID)
	sig.signature(p)
	return api[[]tg.GameHighScore](c, "getGameHighScores", p)
}
