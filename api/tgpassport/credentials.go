package tgpassport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/karalef/tgot/api/tg"
)

// DecryptCredentials decrypts telegram encrypted passport credentials.
func DecryptCredentials(ecreds tg.EncryptedCredentials, priv *rsa.PrivateKey, nonce string) (*Credentials, error) {
	secret, err := decodeString(ecreds.Secret, "EncryptedCredentials.Secret")
	if err != nil {
		return nil, err
	}
	secret, err = rsa.DecryptOAEP(sha1.New(), nil, priv, secret, nil)
	if err != nil {
		return nil, err
	}

	dataHash, err := decodeString(ecreds.Hash, "EncryptedCredentials.Hash")
	if err != nil {
		return nil, err
	}
	data, err := decodeString(ecreds.Data, "EncryptedCredentials.Data")
	if err != nil {
		return nil, err
	}
	data, err = decrypt(secret, dataHash, data)
	if err != nil {
		return nil, err
	}

	var creds Credentials
	if err = json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	if creds.Nonce != nonce {
		return nil, errors.New("nonce does not match")
	}
	return &creds, nil
}

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

// DecryptData decrypts base64-encoded encrypted data.
// dst must satisfy [DataType] interface.
func (c DataCredentials) DecryptData(dst any, data string) error {
	encryptedData, err := decodeString(data, "data")
	if err != nil {
		return err
	}
	dataJSON, err := decrypt([]byte(c.Secret), []byte(c.DataHash), encryptedData)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataJSON, dst)
}

// FileCredentials can be used to decrypt encrypted files from the front_side, reverse_side, selfie, files
// and translation fields in EncryptedPassportElement.
type FileCredentials struct {
	FileHash string `json:"file_hash"`
	Secret   string `json:"secret"`
}

// Decrypt decrypts encrypted file data.
func (c FileCredentials) Decrypt(fileData []byte) ([]byte, error) {
	return decrypt([]byte(c.Secret), []byte(c.FileHash), fileData)
}

func decrypt(secret, hash, data []byte) ([]byte, error) {
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}
	if len(hash) != sha256.Size {
		return nil, errors.New("hash size does not match sha256")
	}
	if len(data) == 0 || len(data)%16 != 0 {
		return nil, errors.New("data length is not divisible by 16")
	}

	h := sha512.New()
	h.Write(secret)
	h.Write(hash)
	dig := h.Sum(nil)

	// decrypt data
	block, err := aes.NewCipher(dig[0:32])
	if err != nil {
		panic(err)
	}
	cipher.NewCBCDecrypter(block, dig[32:48]).CryptBlocks(data, data)

	// verify data
	if dig := sha256.Sum256(data); string(dig[:]) != string(hash) {
		return nil, errors.New("hash does not match")
	}
	padding := data[0]
	if padding < 32 {
		return nil, errors.New("invalid data padding")
	}

	return data[padding:], err
}

func decodeString(s, name string) ([]byte, error) {
	if len(s) == 0 {
		return nil, errors.New("empty " + name)
	}
	dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.New(
			"invalid encoding of " + string(name) + " (must be base64)")
	}
	return dec, nil
}
