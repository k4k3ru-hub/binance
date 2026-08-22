package websocket

import (
	"context"
	"time"
)

const defaultControlMessageInterval = 350 * time.Millisecond

type controlMessageLimiter interface {
	Wait(context.Context) error
}

type intervalControlMessageLimiter struct {
	interval time.Duration
	last     time.Time
	turn     chan struct{}
}

func newControlMessageLimiter(interval time.Duration) controlMessageLimiter {
	limiter := &intervalControlMessageLimiter{interval: interval, turn: make(chan struct{}, 1)}
	limiter.turn <- struct{}{}
	return limiter
}

func (l *intervalControlMessageLimiter) Wait(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.turn:
	}
	defer func() { l.turn <- struct{}{} }()

	wait := time.Until(l.last.Add(l.interval))
	if wait <= 0 {
		l.last = time.Now()
		return nil
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		l.last = time.Now()
		return nil
	}
}
