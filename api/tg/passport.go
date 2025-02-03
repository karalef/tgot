package tg

import "github.com/karalef/tgot/api/internal"

// PassportData describes Telegram Passport data shared with the bot by the user.
type PassportData struct {
	// Array with information about documents and other Telegram Passport elements that was shared with the bot.
	Data []EncryptedPassportElement `json:"data"`

	// Encrypted credentials required to decrypt the data.
	Credentials EncryptedCredentials `json:"credentials"`
}

// PassportFile represents a file uploaded to Telegram Passport.
type PassportFile struct {
	FileData
	FileDate int64 `json:"file_date"`
}

// EncryptedPassportElement describes documents or other Telegram Passport elements
// shared with the bot by the user.
type EncryptedPassportElement struct {
	// Element type. One of “personal_details”, “passport”, “driver_license”,
	// “identity_card”, “internal_passport”, “address”, “utility_bill”, “bank_statement”,
	// “rental_agreement”, “passport_registration”, “temporary_registration”, “phone_number”, “email”.
	Type PassportElementType `json:"type"`

	// Base64-encoded encrypted Telegram Passport element data provided by the user,
	// available for “personal_details”, “passport”, “driver_license”, “identity_card”,
	// “internal_passport” and “address” types.
	// Can be decrypted and verified using the accompanying EncryptedCredentials.
	Data string `json:"data"`

	// User's verified phone number, available only for “phone_number” type.
	PhoneNumber string `json:"phone_number"`

	// User's verified email address, available only for “email” type.
	Email string `json:"email"`

	// Array of encrypted files with documents provided by the user, available for
	// “utility_bill”, “bank_statement”, “rental_agreement”, “passport_registration” and
	// “temporary_registration” types.
	// Files can be decrypted and verified using the accompanying EncryptedCredentials.
	Files []PassportFile `json:"files"`

	// Encrypted file with the front side of the document, provided by the user.
	// Available for “passport”, “driver_license”, “identity_card” and “internal_passport”.
	// The file can be decrypted and verified using the accompanying EncryptedCredentials.
	FrontSide *PassportFile `json:"front_side"`

	// Encrypted file with the reverse side of the document, provided by the user.
	// Available for “driver_license” and “identity_card”.
	// The file can be decrypted and verified using the accompanying EncryptedCredentials.
	ReverseSide *PassportFile `json:"reverse_side"`

	// Encrypted file with the selfie of the user holding a document, provided by the user;
	// available for “passport”, “driver_license”, “identity_card” and “internal_passport”.
	// The file can be decrypted and verified using the accompanying EncryptedCredentials.
	Selfie *PassportFile `json:"selfie"`

	// Array of encrypted files with translated versions of documents provided by the user.
	// Available if requested for “passport”, “driver_license”, “identity_card”,
	// “internal_passport”, “utility_bill”, “bank_statement”, “rental_agreement”,
	// “passport_registration” and “temporary_registration” types.
	// Files can be decrypted and verified using the accompanying EncryptedCredentials.
	Translation []PassportFile `json:"translation"`

	// Base64-encoded element hash for using in PassportElementErrorUnspecified.
	Hash string `json:"hash"`
}

// PassportElementType represents passport element type.
type PassportElementType string

// passport element types.
const (
	ElementPersonalDetails       PassportElementType = "personal_details"
	ElementPassport              PassportElementType = "passport"
	ElementDriverLicense         PassportElementType = "driver_license"
	ElementIdentityCard          PassportElementType = "identity_card"
	ElementInternalPassport      PassportElementType = "internal_passport"
	ElementAddress               PassportElementType = "address"
	ElementUtilityBill           PassportElementType = "utility_bill"
	ElementBankStatement         PassportElementType = "bank_statement"
	ElementRentalAgreement       PassportElementType = "rental_agreement"
	ElementPassportRegistration  PassportElementType = "passport_registration"
	ElementTemporaryRegistration PassportElementType = "temporary_registration"
	ElementPhoneNumber           PassportElementType = "phone_number"
	ElementEmail                 PassportElementType = "email"
)

// EncryptedCredentials describes data required for decrypting and authenticating EncryptedPassportElement.
type EncryptedCredentials struct {
	// Base64-encoded encrypted JSON-serialized data with unique user's payload, data hashes
	// and secrets required for EncryptedPassportElement decryption and authentication.
	Data string `json:"data"`

	// Base64-encoded data hash for data authentication.
	Hash string `json:"hash"`

	// Base64-encoded secret, encrypted with the bot's public RSA key, required for data decryption.
	Secret string `json:"secret"`
}

// PassportElementError represents an error in the Telegram Passport element which was submitted
// that should be resolved by the user.
type PassportElementError struct {
	Type    string
	Message string
	PassportElementErrorSource
}

// MarshalJSON implements json.Marshaler.
func (e *PassportElementError) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(struct {
		Source  string `json:"source"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}{e.PassportElementErrorSource.source(), e.Type, e.Message}, e.PassportElementErrorSource)
}

// PassportElementErrorSource interface.
type PassportElementErrorSource interface {
	source() string
}

// PassportElementErrorDataField represents an issue in one of the data fields that was provided by the user.
type PassportElementErrorDataField struct {
	FieldName string `json:"field_name"`
	DataHash  string `json:"data_hash"` // base64-encoded
}

func (e PassportElementErrorDataField) source() string { return "data" }

// PassportElementErrorFrontSide represents an issue with the front side of a document.
type PassportElementErrorFrontSide struct {
	FileHash string `json:"file_hash"` // base64-encoded
}

func (e PassportElementErrorFrontSide) source() string { return "front_side" }

// PassportElementErrorReverseSide represents an issue with the reverse side of a document.
type PassportElementErrorReverseSide struct {
	FileHash string `json:"file_hash"` // base64-encoded
}

func (e PassportElementErrorReverseSide) source() string { return "reverse_side" }

// PassportElementErrorSelfie represents an issue with the selfie with a document.
type PassportElementErrorSelfie struct {
	FileHash string `json:"file_hash"` // base64-encoded
}

func (e PassportElementErrorSelfie) source() string { return "selfie" }

// PassportElementErrorFile represents an issue with a document scan.
type PassportElementErrorFile struct {
	FileHash string `json:"file_hash"` // base64-encoded
}

func (e PassportElementErrorFile) source() string { return "file" }

// PassportElementErrorFiles represents an issue with a list of scans.
type PassportElementErrorFiles struct {
	FileHashes []string `json:"file_hashes"` // base64-encoded
}

func (e PassportElementErrorFiles) source() string { return "files" }

// PassportElementErrorTranslationFile represents an issue with one of the files that constitute the translation of a document.
type PassportElementErrorTranslationFile struct {
	FileHash string `json:"file_hash"` // base64-encoded
}

func (e PassportElementErrorTranslationFile) source() string { return "translation_file" }

// PassportElementErrorTranslationFiles represents an issue with the translated version of a document.
type PassportElementErrorTranslationFiles struct {
	FileHashes []string `json:"file_hashes"` // base64-encoded
}

func (e PassportElementErrorTranslationFiles) source() string { return "translation_files" }

// PassportElementErrorUnspecified represents an issue in an unspecified place.
type PassportElementErrorUnspecified struct {
	ElementHash string `json:"element_hash"` // base64-encoded
}

func (e PassportElementErrorUnspecified) source() string { return "unspecified" }
