package hash

import (
	"crypto/hmac"
	"crypto/sha256"
)

func GetSHA256Hash(key string, data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write(data)
	if err != nil {
		return make([]byte, 0), err
	}

	return h.Sum(nil), nil
}
