// Пакет hash определяет функции для хеширования
package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// Функция хеширования, которая получает data - хешируемые данные,
// и key - ключ для хеширования
func Hash(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return h.Sum(nil)
}

// Функция подписи данных, которая получает data - подписываемые данные,
// и key - ключ для подписи
func Sign(data []byte, key string) string {
	return base64.StdEncoding.EncodeToString(Hash(data, key))
}
