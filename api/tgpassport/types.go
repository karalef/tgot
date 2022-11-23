package tgpassport

// DataType represents any type of data.
type DataType interface {
	PersonalDetails | ResidentialAddress | IDDocumentData
}

// Gender type.
type Gender string

// genders.
const (
	Male   Gender = "male"
	Female Gender = "female"
)

// PersonalDetails represents personal details.
type PersonalDetails struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	MiddleName           string `json:"middle_name"`
	BirthDate            string `json:"birth_date"` // in DD.MM.YYYY format
	Gender               Gender `json:"gender"`
	CountryCode          string `json:"country_code"`           // ISO 3166-1 alpha-2 country code
	ResidenceCountryCode string `json:"residence_country_code"` // ISO 3166-1 alpha-2 country code
	FirstNameNative      string `json:"first_name_native"`
	LastNameNative       string `json:"last_name_native"`
	MiddleNameNative     string `json:"middle_name_native"`
}

// ResidentialAddress represents a residential address.
type ResidentialAddress struct {
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	CountryCode string `json:"country_code"` // ISO 3166-1 alpha-2 country code
	PostCode    string `json:"post_code"`
}

// IDDocumentData represents the data of an identity document.
type IDDocumentData struct {
	DocumentNo string `json:"document_no"`
	ExpiryDate string `json:"expiry_date"` // in DD.MM.YYYY format
}
