package webhook

import (
	"context"
	"sync"
)

type DeliveryRequest struct {
	DeliveryID string
	EventID    string
	Endpoint   string
	Payload    []byte
	Timestamp  int64
	Signature  string
	KeyVersion int
}

type SendResult struct{ StatusCode int }

type Sender interface {
	Send(context.Context, DeliveryRequest) (SendResult, error)
}

type RecordingSender struct {
	mu         sync.Mutex
	statusCode int
	requests   []DeliveryRequest
}

func NewRecordingSender(statusCode int) *RecordingSender {
	return &RecordingSender{statusCode: statusCode}
}

func (s *RecordingSender) Send(ctx context.Context, request DeliveryRequest) (SendResult, error) {
	if err := ctx.Err(); err != nil {
		return SendResult{}, err
	}
	s.mu.Lock()
	request.Payload = append([]byte(nil), request.Payload...)
	s.requests = append(s.requests, request)
	s.mu.Unlock()
	return SendResult{StatusCode: s.statusCode}, nil
}

func (s *RecordingSender) Requests() []DeliveryRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]DeliveryRequest, len(s.requests))
	copy(result, s.requests)
	return result
}
