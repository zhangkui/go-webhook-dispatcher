package webhook

import "time"

type DeliveryStatus string

const (
	DeliveryPending    DeliveryStatus = "pending"
	DeliveryProcessing DeliveryStatus = "processing"
	DeliverySucceeded  DeliveryStatus = "succeeded"
	DeliveryRetrying   DeliveryStatus = "retrying"
	DeliveryDeadLetter DeliveryStatus = "dead_letter"
)

type Subscription struct {
	ID         string            `json:"id"`
	Endpoint   string            `json:"endpoint"`
	EventTypes []string          `json:"event_types"`
	Filter     map[string]string `json:"filter,omitempty"`
	Secrets    map[int]string    `json:"-"`
	KeyVersion int               `json:"key_version"`
	Active     bool              `json:"active"`
	CreatedAt  time.Time         `json:"created_at"`
}

type Event struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	Payload     []byte            `json:"payload"`
	PublishedAt time.Time         `json:"published_at"`
}

type Delivery struct {
	ID             string         `json:"id"`
	EventID        string         `json:"event_id"`
	SubscriptionID string         `json:"subscription_id"`
	Endpoint       string         `json:"endpoint"`
	Status         DeliveryStatus `json:"status"`
	AttemptCount   int            `json:"attempt_count"`
	NextAttemptAt  time.Time      `json:"next_attempt_at"`
	KeyVersion     int            `json:"key_version"`
	ReplayOf       string         `json:"replay_of,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Attempt struct {
	Number      int       `json:"number"`
	AttemptedAt time.Time `json:"attempted_at"`
	StatusCode  int       `json:"status_code"`
	Error       string    `json:"error,omitempty"`
	Signature   string    `json:"signature"`
	Timestamp   int64     `json:"timestamp"`
	KeyVersion  int       `json:"key_version"`
}

type DeliveryFilter struct {
	Status         DeliveryStatus
	Endpoint       string
	EventID        string
	SubscriptionID string
}
