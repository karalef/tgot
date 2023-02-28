package tgpassport

import "github.com/karalef/tgot/api/internal"

// PassportElementError represents an error in the Telegram Passport element which was submitted
// that should be resolved by the user.
type PassportElementError struct {
	Type    string
	Message string
	ErrorSource
}

// MarshalJSON implements json.Marshaler.
func (e *PassportElementError) MarshalJSON() ([]byte, error) {
	return internal.MergeJSON(struct {
		Source  string `json:"source"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}{e.ErrorSource.source(), e.Type, e.Message}, e.ErrorSource)
}

// ErrorSource interface.
type ErrorSource interface {
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
