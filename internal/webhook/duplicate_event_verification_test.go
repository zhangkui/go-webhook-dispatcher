package webhook

import (
	"errors"
	"testing"
)

func TestDuplicateEventDoesNotCreateAdditionalDeliveries(t *testing.T) {
	service, _, _, _ := newTestService(202)
	_, err := service.RegisterSubscription(Subscription{Endpoint: "https://example.test/hook", EventTypes: []string{"invoice.paid"}, Secrets: map[int]string{1: "secret"}, KeyVersion: 1})
	if err != nil { t.Fatal(err) }
	event := Event{ID: "evt-duplicate", Type: "invoice.paid", Payload: []byte("{}")}
	if _, err := service.Publish(event); err != nil { t.Fatal(err) }
	if _, err := service.Publish(event); !errors.Is(err, ErrEventAlreadyExists) { t.Fatalf("error = %v", err) }
	deliveries := service.Deliveries(DeliveryFilter{EventID: event.ID})
	if len(deliveries) != 1 { t.Fatalf("deliveries = %d, want 1", len(deliveries)) }
}
