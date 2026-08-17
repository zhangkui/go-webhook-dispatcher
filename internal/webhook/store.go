package webhook

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrEventAlreadyExists   = errors.New("event id already exists")
	ErrDeliveryNotFound     = errors.New("delivery not found")
	ErrDeliveryNotDead      = errors.New("delivery is not dead letter")
)

type MemoryStore struct {
	mu            sync.RWMutex
	subscriptions []Subscription
	events        map[string]Event
	deliveries    map[string]Delivery
	attempts      map[string][]Attempt
	sequence      uint64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{events: make(map[string]Event), deliveries: make(map[string]Delivery), attempts: make(map[string][]Attempt)}
}

func (s *MemoryStore) nextID(prefix string) string {
	s.sequence++
	return fmt.Sprintf("%s-%06d", prefix, s.sequence)
}

func cloneSubscription(sub Subscription) Subscription {
	sub.EventTypes = append([]string(nil), sub.EventTypes...)
	sub.Filter = cloneMap(sub.Filter)
	sub.Secrets = cloneIntMap(sub.Secrets)
	return sub
}

func cloneEvent(event Event) Event {
	event.Attributes = cloneMap(event.Attributes)
	event.Payload = append([]byte(nil), event.Payload...)
	return event
}

func cloneMap(source map[string]string) map[string]string {
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func cloneIntMap(source map[int]string) map[int]string {
	result := make(map[int]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func (s *MemoryStore) AddSubscription(sub Subscription) Subscription {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sub.ID == "" {
		sub.ID = s.nextID("sub")
	}
	s.subscriptions = append(s.subscriptions, cloneSubscription(sub))
	return cloneSubscription(sub)
}

func (s *MemoryStore) DisableSubscription(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sub := range s.subscriptions {
		if sub.ID == id {
			sub.Active = false
			return nil
		}
	}
	return ErrSubscriptionNotFound
}

func (s *MemoryStore) Subscriptions() []Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Subscription, 0, len(s.subscriptions))
	for _, sub := range s.subscriptions {
		result = append(result, cloneSubscription(sub))
	}
	return result
}

func (s *MemoryStore) AddEvent(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[event.ID]; exists {
		return ErrEventAlreadyExists
	}
	s.events[event.ID] = cloneEvent(event)
	return nil
}

func (s *MemoryStore) Event(id string) (Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	return cloneEvent(event), ok
}

func (s *MemoryStore) AddDelivery(delivery Delivery) Delivery {
	s.mu.Lock()
	defer s.mu.Unlock()
	if delivery.ID == "" {
		delivery.ID = s.nextID("del")
	}
	s.deliveries[delivery.ID] = delivery
	return delivery
}

func (s *MemoryStore) Delivery(id string) (Delivery, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	delivery, ok := s.deliveries[id]
	return delivery, ok
}

func (s *MemoryStore) UpdateDelivery(delivery Delivery) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.deliveries[delivery.ID]; !ok {
		return ErrDeliveryNotFound
	}
	s.deliveries[delivery.ID] = delivery
	return nil
}

func (s *MemoryStore) DueDeliveries(now time.Time) []Delivery {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Delivery, 0)
	for id, delivery := range s.deliveries {
		if (delivery.Status == DeliveryPending || delivery.Status == DeliveryRetrying) && !delivery.NextAttemptAt.After(now) {
			delivery.Status = DeliveryProcessing
			delivery.UpdatedAt = now
			s.deliveries[id] = delivery
			result = append(result, delivery)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
	return result
}

func (s *MemoryStore) AddAttempt(deliveryID string, attempt Attempt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attempts[deliveryID] = append(s.attempts[deliveryID], attempt)
}

func (s *MemoryStore) Attempts(deliveryID string) []Attempt {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Attempt, len(s.attempts[deliveryID]))
	copy(result, s.attempts[deliveryID])
	return result
}

func (s *MemoryStore) ListDeliveries(filter DeliveryFilter) []Delivery {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Delivery, 0)
	for _, delivery := range s.deliveries {
		if filter.Status != "" && delivery.Status != filter.Status {
			continue
		}
		if filter.Endpoint != "" && delivery.Endpoint != filter.Endpoint {
			continue
		}
		if filter.EventID != "" && delivery.EventID != filter.EventID {
			continue
		}
		if filter.SubscriptionID != "" && delivery.SubscriptionID != filter.SubscriptionID {
			continue
		}
		result = append(result, delivery)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
	return result
}
