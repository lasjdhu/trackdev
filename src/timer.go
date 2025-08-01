package src

import (
	"fmt"
	"time"
)

type Timer struct {
	start   time.Time
	elapsed time.Duration
	running bool
}

func NewTimer() *Timer {
	return &Timer{
		elapsed: 0,
		running: false,
	}
}

func (t *Timer) Start() {
	if !t.running {
		t.start = time.Now().Add(-t.elapsed)
		t.running = true
	}
}

func (t *Timer) Stop() {
	if t.running {
		t.elapsed = time.Since(t.start)
		t.running = false
	}
}

func (t *Timer) Toggle() {
	if t.running {
		t.Stop()
	} else {
		t.Start()
	}
}

func (t *Timer) Reset() {
	t.elapsed = 0
	t.running = false
}

func (t *Timer) IsRunning() bool {
	return t.running
}

func (t *Timer) Elapsed() time.Duration {
	if t.running {
		return time.Since(t.start)
	}
	return t.elapsed
}

func (t *Timer) String() string {
	duration := t.Elapsed().Round(time.Second)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

