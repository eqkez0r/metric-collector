package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func Hash(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return h.Sum(nil)
}

func Sign(data []byte, key string) string {
	return base64.StdEncoding.EncodeToString(Hash(data, key))
}
