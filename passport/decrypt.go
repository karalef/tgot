package passport

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
	"github.com/karalef/tgot/api/tgpassport"
)

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

	// find key and iv
	dig := sha512.Sum512(append(secret, hash...))
	key, iv := dig[0:32], dig[32:48]

	// decrypt data
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(data, data)

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

func fromBase64(s, name string) ([]byte, error) {
	if len(s) == 0 {
		return nil, errors.New("empty " + name)
	}
	dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, invalidBase64(name)
	}
	return dec, nil
}

type invalidBase64 string

func (e invalidBase64) Error() string {
	return "invalid encoding of " + string(e) + " (must be base64)"
}

// DecryptCredentials decrypts telegram encrypted passport credentials.
func DecryptCredentials(ecreds tg.EncryptedCredentials, priv *rsa.PrivateKey, nonce string) (*tgpassport.Credentials, error) {
	secret, err := fromBase64(ecreds.Secret, "secret")
	if err != nil {
		return nil, err
	}

	secret, err = rsa.DecryptOAEP(sha1.New(), nil, priv, secret, nil)
	if err != nil {
		return nil, err
	}

	dataHash, err := fromBase64(ecreds.Hash, "hash")
	if err != nil {
		return nil, err
	}
	data, err := fromBase64(ecreds.Data, "data")
	if err != nil {
		return nil, err
	}
	data, err = decrypt(secret, dataHash, data)
	if err != nil {
		return nil, err
	}

	var creds tgpassport.Credentials
	if err = json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	if creds.Nonce != nonce {
		return nil, errors.New("nonce does not match")
	}
	return &creds, nil
}

// DecryptData decrypts base64-encoded encrypted data.
func DecryptData[T tgpassport.DataType](creds tgpassport.DataCredentials, data string) (*T, error) {
	encryptedData, err := fromBase64(data, "data")
	if err != nil {
		return nil, err
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
