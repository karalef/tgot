package tg

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
