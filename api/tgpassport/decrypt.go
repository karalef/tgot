package tgpassport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	sha256pkg "crypto/sha256"
	sha512pkg "crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"sync"

	"github.com/karalef/tgot/api/tg"
)

var (
	sha256pool = sync.Pool{
		New: func() any { return sha256pkg.New() },
	}
	sha512pool = sync.Pool{
		New: func() any { return sha512pkg.New() },
	}
)

func hashSum(pool *sync.Pool, b []byte) []byte {
	h := pool.Get().(hash.Hash)
	defer pool.Put(h)
	h.Reset()
	h.Write(b)
	return h.Sum(nil)
}

func sha256(b []byte) []byte {
	return hashSum(&sha256pool, b)
}

func sha512(b []byte) []byte {
	return hashSum(&sha512pool, b)
}

func decrypt(secret, hash, encrypted []byte) ([]byte, error) {
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}
	if len(hash) == 0 || len(hash) != sha256pkg.Size {
		return nil, errors.New("hash size does not match sha256")
	}
	if len(encrypted) == 0 || len(encrypted)%16 != 0 {
		return nil, errors.New("data length is not divisible by 16")
	}

	// find key and iv
	h := sha512(append(secret, hash...))
	key, iv := h[0:32], h[32:48]

	// decrypt data
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	paddedData := make([]byte, len(encrypted))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(paddedData, encrypted)

	// verify data
	if string(sha256(paddedData)) != string(hash) {
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
func DecryptCredentials(ecreds tg.EncryptedCredentials, priv *rsa.PrivateKey, nonce string) (*Credentials, error) {
	secret, err := fromBase64(ecreds.Secret)
	if err != nil {
		return nil, invalidBase64("secret")
	}
	if len(secret) == 0 {
		return nil, errors.New("empty secret")
	}
	secret, err = rsa.DecryptOAEP(sha256pool.Get().(hash.Hash), rand.Reader, priv, secret, nil)
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

	var creds Credentials
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
func DecryptData[T DataType](creds DataCredentials, data string) (*T, error) {
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
func DecryptFile(creds FileCredentials, fileData []byte) ([]byte, error) {
	return decrypt([]byte(creds.Secret), []byte(creds.FileHash), fileData)
}

type invalidBase64 string

func (e invalidBase64) Error() string {
	return "invalid encoding of " + string(e) + " (must be base64)"
}
