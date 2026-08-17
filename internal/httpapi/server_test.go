package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zhangkui/go-webhook-dispatcher/internal/webhook"
)

type apiClock struct{ now time.Time }

func (c apiClock) Now() time.Time                         { return c.now }
func (c apiClock) NewTicker(time.Duration) webhook.Ticker { return apiTicker{make(chan time.Time)} }

type apiTicker struct{ ch chan time.Time }

func (t apiTicker) C() <-chan time.Time { return t.ch }
func (t apiTicker) Stop()               {}

func TestRegisterSubscription(t *testing.T) {
	service := webhook.NewService(webhook.NewMemoryStore(), webhook.NewRecordingSender(202), apiClock{time.Unix(1, 0)}, webhook.RetryPolicy{BaseDelay: time.Second, MaxAttempts: 3})
	request := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(`{"endpoint":"https://example.test/hook","event_types":["invoice.paid"],"secrets":{"1":"secret"},"key_version":1}`))
	response := httptest.NewRecorder()
	New(service).ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
