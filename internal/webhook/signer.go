package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
)

var ErrSecretVersionNotFound = errors.New("secret version not found")

func SignPayload(secrets map[int]string, version int, timestamp int64, payload []byte) (string, error) {
	secret, ok := secrets[version]
	if !ok {
		return "", ErrSecretVersionNotFound
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
	mac.Write([]byte("."))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil)), nil
}
