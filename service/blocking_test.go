package gr8service_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	gr8service "github.com/jonasi/gr8go/service"
	"github.com/jonasi/svc"
)

// test that you can call start twice
func TestBlocking_StartTwice(t *testing.T) {
	gr8service := svc.WrapBlocking(nil, nil)
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
	gr8service := svc.WrapBlocking(func(ctx context.Context) error {
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
	gr8service := svc.Multi(
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
