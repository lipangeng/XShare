package controller

import (
	"errors"
	"sync"
	"testing"
)

func TestController_StartStop_StateTransitions(t *testing.T) {
	ctrl := NewController()

	if err := ctrl.StartForward(); err != nil {
		t.Fatalf("expected first StartForward to succeed, got error: %v", err)
	}

	err := ctrl.StartForward()
	if err == nil {
		t.Fatalf("expected second StartForward to fail when already running")
	}
	if !errors.Is(err, ErrForwardAlreadyRunning) {
		t.Fatalf("expected ErrForwardAlreadyRunning, got: %v", err)
	}

	if err := ctrl.StopForward(); err != nil {
		t.Fatalf("expected StopForward to succeed from running state, got error: %v", err)
	}
}

func TestController_StopForward_WhenIdle_ReturnsError(t *testing.T) {
	ctrl := NewController()

	err := ctrl.StopForward()
	if err == nil {
		t.Fatalf("expected StopForward to fail when controller is idle")
	}
	if !errors.Is(err, ErrForwardAlreadyIdle) {
		t.Fatalf("expected ErrForwardAlreadyIdle, got: %v", err)
	}
}

func TestController_StartForward_ConcurrentFromIdle_ExactlyOneSuccess(t *testing.T) {
	ctrl := NewController()
	const n = 24

	start := make(chan struct{})
	results := make(chan error, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			<-start
			results <- ctrl.StartForward()
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	successes := 0
	alreadyRunning := 0
	for err := range results {
		if err == nil {
			successes++
			continue
		}
		if errors.Is(err, ErrForwardAlreadyRunning) {
			alreadyRunning++
			continue
		}
		t.Fatalf("unexpected error from StartForward: %v", err)
	}

	if successes != 1 {
		t.Fatalf("expected exactly one successful StartForward, got %d", successes)
	}
	if alreadyRunning != n-1 {
		t.Fatalf("expected %d ErrForwardAlreadyRunning results, got %d", n-1, alreadyRunning)
	}
}

func TestController_StopForward_ConcurrentFromRunning_ExactlyOneSuccess(t *testing.T) {
	ctrl := NewController()
	if err := ctrl.StartForward(); err != nil {
		t.Fatalf("expected setup StartForward to succeed, got: %v", err)
	}

	const n = 24
	start := make(chan struct{})
	results := make(chan error, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			<-start
			results <- ctrl.StopForward()
		}()
	}

	close(start)
	wg.Wait()
	close(results)

	successes := 0
	alreadyIdle := 0
	for err := range results {
		if err == nil {
			successes++
			continue
		}
		if errors.Is(err, ErrForwardAlreadyIdle) {
			alreadyIdle++
			continue
		}
		t.Fatalf("unexpected error from StopForward: %v", err)
	}

	if successes != 1 {
		t.Fatalf("expected exactly one successful StopForward, got %d", successes)
	}
	if alreadyIdle != n-1 {
		t.Fatalf("expected %d ErrForwardAlreadyIdle results, got %d", n-1, alreadyIdle)
	}
}
