package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSignPayloadUsesRequestedKeyVersion(t *testing.T) {
	payload := []byte(`{"id":7}`)
	got, err := SignPayload(map[int]string{1: "old-secret", 2: "new-secret"}, 1, 1234, payload)
	if err != nil { t.Fatal(err) }
	mac := hmac.New(sha256.New, []byte("old-secret"))
	mac.Write([]byte("1234."))
	mac.Write(payload)
	want := hex.EncodeToString(mac.Sum(nil))
	if got != want { t.Fatalf("signature = %s, want %s", got, want) }
}
