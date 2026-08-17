package webhook

import "testing"

func TestDisabledSubscriptionStopsFutureDeliveries(t *testing.T) {
	service, _, _, _ := newTestService(202)
	disabled, err := service.RegisterSubscription(Subscription{Endpoint: "https://disabled.test/hook", EventTypes: []string{"order.created"}, Secrets: map[int]string{1: "a"}, KeyVersion: 1})
	if err != nil { t.Fatal(err) }
	_, err = service.RegisterSubscription(Subscription{Endpoint: "https://active.test/hook", EventTypes: []string{"order.created"}, Secrets: map[int]string{1: "b"}, KeyVersion: 1})
	if err != nil { t.Fatal(err) }
	if err := service.DisableSubscription(disabled.ID); err != nil { t.Fatal(err) }
	deliveries, err := service.Publish(Event{ID: "evt-disabled", Type: "order.created", Payload: []byte("{}")})
	if err != nil { t.Fatal(err) }
	if len(deliveries) != 1 || deliveries[0].Endpoint != "https://active.test/hook" { t.Fatalf("unexpected deliveries: %#v", deliveries) }
}
