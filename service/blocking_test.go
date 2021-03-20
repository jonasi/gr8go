package gr8service_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gr8service "github.com/jonasi/gr8go/service"
)

// test that you can call start twice
func TestBlocking_StartTwice(t *testing.T) {
	gr8service := gr8service.FromBlocking(nil, nil)
	ctx := context.Background()
	err := gr8service.Start(ctx)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
	}
	err = gr8service.Start(ctx)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
	}
}

// TestBlocking_StartCancel tests that calling stop
// will send a cancel signal to the start function ctx
func TestBlocking_StartCancel(t *testing.T) {
	var cancelled int32
	gr8service := gr8service.FromBlocking(func(ctx context.Context) error {
		<-ctx.Done()
		atomic.StoreInt32(&cancelled, 1)
		return ctx.Err()
	}, nil)

	ctx := context.Background()
	err := gr8service.Start(ctx)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
	}
	err = gr8service.Stop(ctx)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
	}
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&cancelled) != 1 {
		t.Errorf("Expected cancelled to be true")
	}
}

func TestMulti(t *testing.T) {
	var idx int32
	gr8service := gr8service.Multi(
		gr8service.FromStartStop(func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)

			if idx == 0 {
				idx = 1
			}

			return nil
		}, nil),
		gr8service.FromStartStop(func(ctx context.Context) error {
			if idx == 1 {
				idx = 2
			}
			return nil
		}, nil),
	)

	ctx := context.Background()
	err := gr8service.Start(ctx)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
	}

	if idx != 2 {
		t.Errorf("Expected 2")
	}
}

func TestCancelledContext(t *testing.T) {
	var (
		mu                    sync.Mutex
		stoppedBeforeCanceled bool
		stopped               bool
		cancelled             bool
	)

	gr8service := gr8service.FromBlocking(func(ctx context.Context) error {
		<-ctx.Done()
		mu.Lock()
		defer mu.Unlock()

		cancelled = true
		return ctx.Err()
	}, func(ctx context.Context) error {
		mu.Lock()
		defer mu.Unlock()

		stopped = true
		if !cancelled {
			stoppedBeforeCanceled = true
		}
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	err := gr8service.Start(ctx)
	if err != nil {
		t.Errorf("Received unexpected error: %s", err)
	}

	cancel()
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()

	if !stopped {
		t.Errorf("Expected stopped to be true")
	}

	if !cancelled {
		t.Errorf("Expected cancelled to be true")
	}

	if !stoppedBeforeCanceled {
		t.Errorf("Expected Stop to be called before cancel")
	}
}
