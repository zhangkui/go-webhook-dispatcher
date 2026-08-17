package webhook

import (
	"context"
	"testing"
	"time"
)

type fixedClock struct{ now time.Time }

func (c *fixedClock) Now() time.Time                 { return c.now }
func (c *fixedClock) NewTicker(time.Duration) Ticker { return stoppedTicker{ch: make(chan time.Time)} }

type stoppedTicker struct{ ch chan time.Time }

func (t stoppedTicker) C() <-chan time.Time { return t.ch }
func (t stoppedTicker) Stop()               {}

func newTestService(status int) (*Service, *MemoryStore, *RecordingSender, *fixedClock) {
	clock := &fixedClock{now: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}
	store := NewMemoryStore()
	sender := NewRecordingSender(status)
	service := NewService(store, sender, clock, RetryPolicy{BaseDelay: time.Second, MaxDelay: time.Minute, MaxAttempts: 3})
	return service, store, sender, clock
}

func TestPublishFiltersSubscriptions(t *testing.T) {
	service, _, _, _ := newTestService(202)
	_, err := service.RegisterSubscription(Subscription{Endpoint: "https://example.test/orders", EventTypes: []string{"order.created"}, Filter: map[string]string{"region": "us"}, Secrets: map[int]string{1: "secret"}, KeyVersion: 1})
	if err != nil {
		t.Fatal(err)
	}
	deliveries, err := service.Publish(Event{ID: "evt-1", Type: "order.created", Attributes: map[string]string{"region": "eu"}, Payload: []byte(`{"id":1}`)})
	if err != nil {
		t.Fatal(err)
	}
	if len(deliveries) != 0 {
		t.Fatalf("deliveries = %d, want 0", len(deliveries))
	}
}

func TestDispatchRecordsSuccessfulAttempt(t *testing.T) {
	service, _, sender, _ := newTestService(204)
	sub, err := service.RegisterSubscription(Subscription{Endpoint: "https://example.test/hooks", EventTypes: []string{"*"}, Secrets: map[int]string{1: "secret"}, KeyVersion: 1})
	if err != nil {
		t.Fatal(err)
	}
	deliveries, err := service.Publish(Event{ID: "evt-2", Type: "invoice.paid", Payload: []byte(`{"paid":true}`)})
	if err != nil {
		t.Fatal(err)
	}
	if len(deliveries) != 1 || deliveries[0].SubscriptionID != sub.ID {
		t.Fatalf("unexpected deliveries: %#v", deliveries)
	}
	if got := service.DispatchDue(context.Background()); got != 1 {
		t.Fatalf("dispatched = %d, want 1", got)
	}
	current := service.Deliveries(DeliveryFilter{EventID: "evt-2"})
	if len(current) != 1 || current[0].Status != DeliverySucceeded {
		t.Fatalf("unexpected status: %#v", current)
	}
	if len(sender.Requests()) != 1 || len(service.Attempts(current[0].ID)) != 1 {
		t.Fatal("delivery evidence not recorded")
	}
}

func TestListDeliveriesFiltersByEndpoint(t *testing.T) {
	service, _, _, _ := newTestService(202)
	_, _ = service.RegisterSubscription(Subscription{Endpoint: "https://a.test/hook", EventTypes: []string{"x"}, Secrets: map[int]string{1: "a"}, KeyVersion: 1})
	_, _ = service.RegisterSubscription(Subscription{Endpoint: "https://b.test/hook", EventTypes: []string{"x"}, Secrets: map[int]string{1: "b"}, KeyVersion: 1})
	_, _ = service.Publish(Event{ID: "evt-3", Type: "x", Payload: []byte("x")})
	got := service.Deliveries(DeliveryFilter{Endpoint: "https://b.test/hook"})
	if len(got) != 1 || got[0].Endpoint != "https://b.test/hook" {
		t.Fatalf("unexpected deliveries: %#v", got)
	}
}
