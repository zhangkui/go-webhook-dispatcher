package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
)

var ErrSecretVersionNotFound = errors.New("secret version not found")

// SignPayload computes the HMAC-SHA256 signature for payload using the secret
// pinned to the requested version. The version is fixed by the delivery record
// at creation time, so key rotation must never alter how a historical task is
// signed: only the secret at the requested version is used, regardless of
// whether newer versions exist.
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
