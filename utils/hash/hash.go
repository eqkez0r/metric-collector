package hash

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Hash(b []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(b)
	return h.Sum(nil)
}
