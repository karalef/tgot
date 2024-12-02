package passport

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/karalef/tgot/api/tgpassport"
	"github.com/karalef/tgot/deeplinks"
)

// RequestParams contains parameters to request information.
type RequestParams struct {
	BotID int
	Scope tgpassport.Scope

	// Public key of the bot.
	PublicKey string

	// Bot-specified nonce.
	// For security purposes it should be a cryptographically secure unique identifier of the request.
	Nonce string
}

// Query generates query part of the request URI.
func (p *RequestParams) Query() (url.Values, error) {
	if p.BotID == 0 || p.PublicKey == "" || p.Nonce == "" {
		return nil, errors.New("all fields are required")
	}
	if p.Scope.V != 1 {
		return nil, errors.New("invalid version")
	}
	scope, err := json.Marshal(p.Scope)
	if err != nil {
		return nil, err
	}
	return url.Values{
		"bot_id":     {strconv.Itoa(p.BotID)},
		"scope":      {string(scope)},
		"public_key": {p.PublicKey},
		"nonce":      {p.Nonce},
	}, nil
}

// Deeplink generates request deeplink.
func (p *RequestParams) Deeplink() (deeplinks.Deeplink, error) {
	q, err := p.Query()
	if err != nil {
		return deeplinks.Deeplink{}, err
	}
	return deeplinks.New("passport", q), nil
}

// URI generates request URI.
func (p *RequestParams) URI() (string, error) {
	dl, err := p.Deeplink()
	if err != nil {
		return "", err
	}
	return dl.TG(true), nil
}
