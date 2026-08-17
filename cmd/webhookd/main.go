package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhangkui/go-webhook-dispatcher/internal/httpapi"
	"github.com/zhangkui/go-webhook-dispatcher/internal/webhook"
)

func main() {
	addr := os.Getenv("WEBHOOK_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	clock := webhook.RealClock{}
	store := webhook.NewMemoryStore()
	sender := webhook.NewRecordingSender(http.StatusAccepted)
	service := webhook.NewService(store, sender, clock, webhook.RetryPolicy{BaseDelay: time.Second, MaxDelay: time.Minute, MaxAttempts: 5})
	scheduler := webhook.NewScheduler(service, clock, 100*time.Millisecond)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	scheduler.Start(ctx)
	server := &http.Server{Addr: addr, Handler: httpapi.New(service), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	log.Printf("webhook dispatcher listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	scheduler.Wait()
}
