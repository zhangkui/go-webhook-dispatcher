package webhook

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

type Service struct {
	store  *MemoryStore
	sender Sender
	clock  Clock
	retry  RetryPolicy
}

func NewService(store *MemoryStore, sender Sender, clock Clock, retry RetryPolicy) *Service {
	return &Service{store: store, sender: sender, clock: clock, retry: retry}
}

func (s *Service) RegisterSubscription(sub Subscription) (Subscription, error) {
	parsed, err := url.ParseRequestURI(sub.Endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return Subscription{}, errors.New("invalid endpoint")
	}
	if len(sub.EventTypes) == 0 {
		return Subscription{}, errors.New("at least one event type is required")
	}
	if sub.KeyVersion <= 0 || sub.Secrets[sub.KeyVersion] == "" {
		return Subscription{}, errors.New("active secret version is required")
	}
	sub.Active = true
	sub.CreatedAt = s.clock.Now()
	return s.store.AddSubscription(sub), nil
}

func (s *Service) DisableSubscription(id string) error { return s.store.DisableSubscription(id) }

func (s *Service) Publish(event Event) ([]Delivery, error) {
	if strings.TrimSpace(event.ID) == "" || strings.TrimSpace(event.Type) == "" || len(event.Payload) == 0 {
		return nil, errors.New("event id, type and payload are required")
	}
	now := s.clock.Now()
	event.PublishedAt = now
	deliveries := make([]Delivery, 0)
	for _, sub := range s.store.Subscriptions() {
		if !matches(sub, event) {
			continue
		}
		delivery := s.store.AddDelivery(Delivery{
			EventID: event.ID, SubscriptionID: sub.ID, Endpoint: sub.Endpoint,
			Status: DeliveryPending, NextAttemptAt: now, KeyVersion: sub.KeyVersion,
			CreatedAt: now, UpdatedAt: now,
		})
		deliveries = append(deliveries, delivery)
	}
	if err := s.store.AddEvent(event); err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (s *Service) DispatchDue(ctx context.Context) int {
	deliveries := s.store.DueDeliveries(s.clock.Now())
	for _, delivery := range deliveries {
		s.dispatchOne(ctx, delivery)
	}
	return len(deliveries)
}

func (s *Service) dispatchOne(ctx context.Context, delivery Delivery) {
	event, eventOK := s.store.Event(delivery.EventID)
	sub, subOK := s.subscription(delivery.SubscriptionID)
	now := s.clock.Now()
	if !eventOK || !subOK {
		delivery.Status = DeliveryDeadLetter
		delivery.UpdatedAt = now
		_ = s.store.UpdateDelivery(delivery)
		return
	}
	timestamp := now.Unix()
	signature, sendErr := SignPayload(sub.Secrets, delivery.KeyVersion, timestamp, event.Payload)
	result := SendResult{}
	if sendErr == nil {
		result, sendErr = s.sender.Send(ctx, DeliveryRequest{
			DeliveryID: delivery.ID, EventID: event.ID, Endpoint: delivery.Endpoint,
			Payload: event.Payload, Timestamp: timestamp, Signature: signature, KeyVersion: delivery.KeyVersion,
		})
	}
	delivery.AttemptCount++
	attempt := Attempt{Number: delivery.AttemptCount, AttemptedAt: now, StatusCode: result.StatusCode, Signature: signature, Timestamp: timestamp, KeyVersion: delivery.KeyVersion}
	if sendErr != nil {
		attempt.Error = sendErr.Error()
	}
	s.store.AddAttempt(delivery.ID, attempt)
	if sendErr == nil && result.StatusCode >= 200 && result.StatusCode < 300 {
		delivery.Status = DeliverySucceeded
	} else if delivery.AttemptCount >= s.retry.MaxAttempts {
		delivery.Status = DeliveryDeadLetter
	} else {
		delivery.Status = DeliveryRetrying
		delivery.NextAttemptAt = now.Add(s.retry.NextDelay(delivery.AttemptCount))
	}
	delivery.UpdatedAt = now
	_ = s.store.UpdateDelivery(delivery)
}

func (s *Service) subscription(id string) (Subscription, bool) {
	for _, sub := range s.store.Subscriptions() {
		if sub.ID == id {
			return sub, true
		}
	}
	return Subscription{}, false
}

func (s *Service) Deliveries(filter DeliveryFilter) []Delivery { return s.store.ListDeliveries(filter) }
func (s *Service) Attempts(deliveryID string) []Attempt        { return s.store.Attempts(deliveryID) }

func (s *Service) ReplayDeadLetter(id string) (Delivery, error) {
	delivery, ok := s.store.Delivery(id)
	if !ok {
		return Delivery{}, ErrDeliveryNotFound
	}
	if delivery.Status != DeliveryDeadLetter {
		return Delivery{}, ErrDeliveryNotDead
	}
	now := s.clock.Now()
	replay := s.store.AddDelivery(Delivery{
		EventID:        delivery.EventID,
		SubscriptionID: delivery.SubscriptionID,
		Endpoint:       delivery.Endpoint,
		Status:         DeliveryPending,
		NextAttemptAt:  now,
		KeyVersion:     delivery.KeyVersion,
		ReplayOf:       delivery.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	return replay, nil
}
