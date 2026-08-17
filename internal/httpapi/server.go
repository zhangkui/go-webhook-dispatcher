package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/zhangkui/go-webhook-dispatcher/internal/webhook"
)

type Server struct{ service *webhook.Service }

func New(service *webhook.Service) http.Handler {
	server := &Server{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscriptions", server.registerSubscription)
	mux.HandleFunc("POST /subscriptions/{id}/disable", server.disableSubscription)
	mux.HandleFunc("POST /events", server.publishEvent)
	mux.HandleFunc("GET /deliveries", server.listDeliveries)
	mux.HandleFunc("POST /deliveries/{id}/replay", server.replayDelivery)
	return mux
}

func (s *Server) registerSubscription(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Endpoint   string            `json:"endpoint"`
		EventTypes []string          `json:"event_types"`
		Filter     map[string]string `json:"filter"`
		Secrets    map[string]string `json:"secrets"`
		KeyVersion int               `json:"key_version"`
	}
	if err := decode(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	secrets := make(map[int]string, len(request.Secrets))
	for rawVersion, secret := range request.Secrets {
		version, err := strconv.Atoi(rawVersion)
		if err != nil {
			writeError(w, http.StatusBadRequest, errors.New("secret versions must be integers"))
			return
		}
		secrets[version] = secret
	}
	sub, err := s.service.RegisterSubscription(webhook.Subscription{Endpoint: request.Endpoint, EventTypes: request.EventTypes, Filter: request.Filter, Secrets: secrets, KeyVersion: request.KeyVersion})
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, sub)
}

func (s *Server) disableSubscription(w http.ResponseWriter, r *http.Request) {
	if err := s.service.DisableSubscription(r.PathValue("id")); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) publishEvent(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID         string            `json:"id"`
		Type       string            `json:"type"`
		Attributes map[string]string `json:"attributes"`
		Payload    json.RawMessage   `json:"payload"`
	}
	if err := decode(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	deliveries, err := s.service.Publish(webhook.Event{ID: request.ID, Type: request.Type, Attributes: request.Attributes, Payload: request.Payload})
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, webhook.ErrEventAlreadyExists) {
			status = http.StatusConflict
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"deliveries": deliveries})
}

func (s *Server) listDeliveries(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	writeJSON(w, http.StatusOK, s.service.Deliveries(webhook.DeliveryFilter{Status: webhook.DeliveryStatus(query.Get("status")), Endpoint: query.Get("endpoint"), EventID: query.Get("event_id"), SubscriptionID: query.Get("subscription_id")}))
}

func (s *Server) replayDelivery(w http.ResponseWriter, r *http.Request) {
	delivery, err := s.service.ReplayDeadLetter(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusAccepted, delivery)
}

func decode(w http.ResponseWriter, r *http.Request, destination any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": strings.TrimSpace(err.Error())})
}
