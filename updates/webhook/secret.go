package webhook

import (
	"crypto/rand"
	"encoding/base64"
)

// Encoding is an interface for encoding.
type Encoding interface {
	EncodeToString([]byte) string
	DecodedLen(int) int
}

// GenerateSecret generates random secret token of specified length.
// It uses the provided encoding (or base64 without padding by default) to get
// the string representation of the secret.
// If length is 0 or bigger than 128, it will be set to 64.
func GenerateSecret(enc Encoding, length uint8) string {
	if enc == nil {
		enc = base64.RawStdEncoding
	}
	if length == 0 || length > 128 {
		length = 64
	}
	buf := make([]byte, enc.DecodedLen(int(length)))
	_, err := rand.Read(buf)
	if err != nil {
		panic("error while generating secret: " + err.Error())
	}
	return enc.EncodeToString(buf)
}
