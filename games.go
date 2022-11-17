package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// Game contains information about the game to be sent.
type Game struct {
	ShortName string
}

func (g Game) data() api.Data {
	return api.NewData().Set("game_short_name", g.ShortName)
}

// SetGameScore contains parameters for setting the game score.
type SetGameScore struct {
	Score       int
	Force       bool
	DisableEdit bool
}

// SetGameScore sets the score of the specified user in a game message.
func (c Context) SetGameScore(sig MessageSignature, userID int64, s SetGameScore) (*tg.Message, error) {
	d := api.NewData()
	d.SetInt64("user_id", userID)
	d.SetInt("score", s.Score, true)
	d.SetBool("force", s.Force)
	d.SetBool("disable_edit_message", s.DisableEdit)
	return c.sig(sig, "setGameScore", d)
}

// GetGameHighScores returns data for high score tables.
func (c Context) GetGameHighScores(sig MessageSignature, userID int64) ([]tg.GameHighScore, error) {
	d := api.NewData().SetInt64("user_id", userID)
	sig.signature(&d)
	return method[[]tg.GameHighScore](c, "getGameHighScores", d)
}
