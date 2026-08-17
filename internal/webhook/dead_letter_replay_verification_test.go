package webhook

import (
	"testing"
	"time"
)

func TestReplayCreatesNewDeliveryAndPreservesDeadLetter(t *testing.T) {
	service, store, _, clock := newTestService(500)
	original := store.AddDelivery(Delivery{EventID: "evt-dead", SubscriptionID: "sub-dead", Endpoint: "https://example.test/hook", Status: DeliveryDeadLetter, AttemptCount: 3, KeyVersion: 2, CreatedAt: clock.Now(), UpdatedAt: clock.Now()})
	store.AddAttempt(original.ID, Attempt{Number: 3, AttemptedAt: clock.Now(), Error: "failed"})
	clock.now = clock.now.Add(time.Minute)
	replay, err := service.ReplayDeadLetter(original.ID)
	if err != nil { t.Fatal(err) }
	if replay.ID == original.ID || replay.ReplayOf != original.ID || replay.Status != DeliveryPending { t.Fatalf("unexpected replay: %#v", replay) }
	preserved, ok := store.Delivery(original.ID)
	if !ok || preserved.Status != DeliveryDeadLetter || preserved.AttemptCount != 3 { t.Fatalf("original changed: %#v", preserved) }
	if len(store.Attempts(original.ID)) != 1 { t.Fatal("original attempts changed") }
}
