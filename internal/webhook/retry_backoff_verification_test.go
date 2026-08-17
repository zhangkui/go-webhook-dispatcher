package webhook

import (
	"testing"
	"time"
)

func TestRetryPolicyStartsAtBaseDelay(t *testing.T) {
	policy := RetryPolicy{BaseDelay: 5 * time.Second, MaxDelay: 30 * time.Second, MaxAttempts: 5}
	want := []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second, 30 * time.Second}
	for index, expected := range want {
		if got := policy.NextDelay(index + 1); got != expected { t.Fatalf("attempt %d delay = %s, want %s", index+1, got, expected) }
	}
}
