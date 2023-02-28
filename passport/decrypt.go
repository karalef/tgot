package passport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/api/tgpassport"
)

func decrypt(secret, hash, encrypted []byte) ([]byte, error) {
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}
	if len(hash) != sha256.Size {
		return nil, errors.New("hash size does not match sha256")
	}
	if len(encrypted) == 0 || len(encrypted)%16 != 0 {
		return nil, errors.New("data length is not divisible by 16")
	}

	// find key and iv
	h := sha512.Sum512(append(secret, hash...))
	key, iv := h[0:32], h[32:48]

	// decrypt data
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	paddedData := make([]byte, len(encrypted))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(paddedData, encrypted)

	// verify data
	if dataHash := sha256.Sum256(paddedData); string(dataHash[:]) != string(hash) {
		return nil, errors.New("hash does not match")
	}
	padding := paddedData[0]
	if padding < 32 {
		return nil, errors.New("invalid data padding")
	}

	// remove padding
	return paddedData[padding:], err
}

func fromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// DecryptCredentials decrypts telegram encrypted passport credentials.
func DecryptCredentials(ecreds tg.EncryptedCredentials, priv *rsa.PrivateKey, nonce string) (*tgpassport.Credentials, error) {
	secret, err := fromBase64(ecreds.Secret)
	if err != nil {
		return nil, invalidBase64("secret")
	}
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}

	// TODO: maybe SHA-1?
	secret, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, secret, nil)
	if err != nil {
		return nil, err
	}

	dataHash, err := fromBase64(ecreds.Hash)
	if err != nil {
		return nil, invalidBase64("hash")
	}
	data, err := fromBase64(ecreds.Hash)
	if err != nil {
		return nil, invalidBase64("data")
	}
	data, err = decrypt(secret, dataHash, data)
	if err != nil {
		return nil, err
	}

	var creds tgpassport.Credentials
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return nil, err
	}
	if creds.Nonce != nonce {
		return nil, errors.New("nonce does not match")
	}
	return &creds, nil
}

// DecryptData decrypts base64-encoded encrypted data.
func DecryptData[T tgpassport.DataType](creds tgpassport.DataCredentials, data string) (*T, error) {
	encryptedData, err := fromBase64(data)
	if err != nil {
		return nil, invalidBase64("data")
	}
	dataJSON, err := decrypt([]byte(creds.Secret), []byte(creds.DataHash), encryptedData)
	if err != nil {
		return nil, err
	}
	var result T
	err = json.Unmarshal(dataJSON, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DecryptFile decrypts encrypted file data.
func DecryptFile(creds tgpassport.FileCredentials, fileData []byte) ([]byte, error) {
	return decrypt([]byte(creds.Secret), []byte(creds.FileHash), fileData)
}

type invalidBase64 string

func (e invalidBase64) Error() string {
	return "invalid encoding of " + string(e) + " (must be base64)"
}
