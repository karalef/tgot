package passport

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/karalef/tgot/api/tg"
)

const (
	// ScopeIDDocument is an alias for one of "passport", "driver_license", or "identity_card".
	ScopeIDDocument tg.PassportElementType = "id_document"

	// ScopeAddressDocument is an alias for one of "utility_bill", "bank_statement", or "rental_agreement".
	ScopeAddressDocument tg.PassportElementType = "address_document"
)

// Scope represents the data to be requested.
type Scope struct {
	Data []ScopeElement `json:"data"`
	V    int            `json:"v"` // must be 1
}

// ScopeElement represents a requested element.
type ScopeElement interface {
	passportScopeElement()
}

// ScopeElementOneOfSeveral represents several elements one of which must be provided.
type ScopeElementOneOfSeveral struct {
	// List of elements one of which must be provided;
	// must contain either several of “passport”, “driver_license”, “identity_card”,
	// “internal_passport” or several of “utility_bill”, “bank_statement”, “rental_agreement”,
	// “passport_registration”, “temporary_registration”
	OneOf []ScopeElementOne `json:"one_of"`

	// Use this parameter if you want to request a selfie with the document from this list that the user chooses to upload.
	Selfie bool `json:"selfie,omitempty"`

	// Use this parameter if you want to request a translation of the document from this list that the user chooses to upload.
	Translation bool `json:"translation,omitempty"`
}

func (ScopeElementOneOfSeveral) passportScopeElement() {}

// ScopeElementOne represents one particular element that must be provided.
type ScopeElementOne struct {
	// Element type. One of "personal_details", "passport", "driver_license", "identity_card",
	// "internal_passport", "address", "utility_bill", "bank_statement", "rental_agreement",
	// "passport_registration", "temporary_registration", "phone_number", "email"
	Type tg.PassportElementType `json:"type"`

	// Use this parameter if you want to request a selfie with the document as well.
	// Available for "passport", "driver_license", "identity_card" and "internal_passport"
	Selfie bool `json:"selfie,omitempty"`

	// Use this parameter if you want to request a translation of the document as well.
	// Available for "passport", "driver_license", "identity_card", "internal_passport",
	// "utility_bill", "bank_statement", "rental_agreement", "passport_registration" and "temporary_registration".
	Translation bool `json:"translation,omitempty"`

	// Use this parameter to request the first, last and middle name of the user in the language of the user's country of residence.
	// Available for "personal_details"
	NativeNames bool `json:"native_names,omitempty"`
}

func (ScopeElementOne) passportScopeElement() {}

// RequestParams contains parameters to request information.
type RequestParams struct {
	BotID int
	Scope Scope

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

// URI generates request URI.
func (p *RequestParams) URI() (string, error) {
	q, err := p.Query()
	if err != nil {
		return "", err
	}
	return "tg://passport?" + q.Encode(), nil
}
