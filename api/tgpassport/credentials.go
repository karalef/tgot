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
func DecryptCredentials(ecreds tg.EncryptedCredentials, priv *rsa.PrivateKey, nonce string) (*tg.Credentials, error) {
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

	var creds tg.Credentials
	if err = json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	if creds.Nonce != nonce {
		return nil, errors.New("nonce does not match")
	}
	return &creds, nil
}

// DecryptData decrypts base64-encoded encrypted data.
func DecryptData[T tg.PassportDataType](dst T, data string, creds tg.DataCredentials) error {
	encryptedData, err := decodeString(data, "data")
	if err != nil {
		return err
	}
	dataJSON, err := decrypt([]byte(creds.Secret), []byte(creds.DataHash), encryptedData)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataJSON, dst)
}

// Decrypt decrypts encrypted file data.
func DecryptFile(fileData []byte, creds tg.FileCredentials) ([]byte, error) {
	return decrypt([]byte(creds.Secret), []byte(creds.FileHash), fileData)
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
