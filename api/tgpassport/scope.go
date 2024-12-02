package tgpassport

import "github.com/karalef/tgot/api/tg"

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
