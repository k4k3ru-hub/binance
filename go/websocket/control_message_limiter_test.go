package websocket

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestControlMessageLimiterAllowsFirstSlotImmediately(t *testing.T) {
	limiter := newControlMessageLimiter(defaultControlMessageInterval)
	started := time.Now()
	if err := limiter.Wait(context.Background()); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}
	if elapsed := time.Since(started); elapsed >= 50*time.Millisecond {
		t.Fatalf("Wait() elapsed = %v, want immediate", elapsed)
	}
}

func TestControlMessageLimiterSpacesSlotsAndLimitsAnyOneSecondWindow(t *testing.T) {
	limiter := newControlMessageLimiter(defaultControlMessageInterval)
	times := make([]time.Time, 4)
	for index := range times {
		if err := limiter.Wait(context.Background()); err != nil {
			t.Fatalf("Wait(%d) error = %v", index, err)
		}
		times[index] = time.Now()
	}
	for index := 1; index < len(times); index++ {
		if interval := times[index].Sub(times[index-1]); interval < defaultControlMessageInterval-5*time.Millisecond {
			t.Fatalf("interval[%d] = %v, want approximately >= %v", index, interval, defaultControlMessageInterval)
		}
	}
	if window := times[3].Sub(times[0]); window < time.Second {
		t.Fatalf("four slots fit in %v, want >= 1s", window)
	}
}

func TestControlMessageLimiterWaitCanBeCanceled(t *testing.T) {
	limiter := newControlMessageLimiter(defaultControlMessageInterval)
	if err := limiter.Wait(context.Background()); err != nil {
		t.Fatalf("first Wait() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- limiter.Wait(ctx) }()
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Wait() error = %v, want context.Canceled", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait() did not return promptly after cancellation")
	}
}

func TestControlMessageLimiterPreservesDeadlineError(t *testing.T) {
	limiter := newControlMessageLimiter(defaultControlMessageInterval)
	if err := limiter.Wait(context.Background()); err != nil {
		t.Fatalf("first Wait() error = %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := limiter.Wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Wait() error = %v, want context.DeadlineExceeded", err)
	}
}

func TestControlMessageLimiterSerializesConcurrentWaiters(t *testing.T) {
	const interval = 20 * time.Millisecond
	limiter := newControlMessageLimiter(interval)
	const count = 6
	times := make(chan time.Time, count)
	start := make(chan struct{})
	var waitGroup sync.WaitGroup
	for index := 0; index < count; index++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			<-start
			if err := limiter.Wait(context.Background()); err != nil {
				t.Errorf("Wait() error = %v", err)
				return
			}
			times <- time.Now()
		}()
	}
	close(start)
	waitGroup.Wait()
	close(times)

	var previous time.Time
	for acquired := range times {
		if !previous.IsZero() && acquired.Sub(previous) < interval-5*time.Millisecond {
			t.Fatalf("concurrent interval = %v, want approximately >= %v", acquired.Sub(previous), interval)
		}
		previous = acquired
	}
}
