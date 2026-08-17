package webhook

import "time"

type Clock interface {
	Now() time.Time
	NewTicker(time.Duration) Ticker
}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type RealClock struct{}

func (RealClock) Now() time.Time                   { return time.Now().UTC() }
func (RealClock) NewTicker(d time.Duration) Ticker { return realTicker{Ticker: time.NewTicker(d)} }

type realTicker struct{ *time.Ticker }

func (t realTicker) C() <-chan time.Time { return t.Ticker.C }
