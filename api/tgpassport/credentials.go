package tgpassport

// Credentials is a JSON-serialized object.
type Credentials struct {
	SecureData SecureData `json:"secure_data"`
	Nonce      string     `json:"nonce"`
}

// SecureData represents the credentials required to decrypt encrypted data.
// All fields are optional and depend on fields that were requested.
type SecureData struct {
	PersonalDetails       *SecureValue `json:"personal_details"`
	Passport              *SecureValue `json:"passport"`
	InternalPassport      *SecureValue `json:"internal_passport"`
	DriverLicense         *SecureValue `json:"driver_license"`
	IdentityCard          *SecureValue `json:"identity_card"`
	Address               *SecureValue `json:"address"`
	UtilityBill           *SecureValue `json:"utility_bill"`
	BankStatement         *SecureValue `json:"bank_statement"`
	RentalAgreement       *SecureValue `json:"rental_agreement"`
	PassportRegistration  *SecureValue `json:"passport_registration"`
	TemporaryRegistration *SecureValue `json:"temporary_registration"`
}

// SecureValue represents the credentials required to decrypt encrypted values.
// All fields are optional and depend on the type of fields that were requested.
type SecureValue struct {
	Data        *DataCredentials  `json:"data"`
	FrontSide   *FileCredentials  `json:"front_side"`
	ReverseSide *FileCredentials  `json:"reverse_side"`
	Selfie      *FileCredentials  `json:"selfie"`
	Translation []FileCredentials `json:"translation"`
	Files       []FileCredentials `json:"files"`
}

// DataCredentials can be used to decrypt encrypted data from the data field in EncryptedPassportElement.
type DataCredentials struct {
	DataHash string `json:"data_hash"`
	Secret   string `json:"secret"`
}

// FileCredentials can be used to decrypt encrypted files from the front_side, reverse_side, selfie, files
// and translation fields in EncryptedPassportElement.
type FileCredentials struct {
	FileHash string `json:"file_hash"`
	Secret   string `json:"secret"`
}
