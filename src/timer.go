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
	return &Timer{}
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
	d := t.Elapsed().Round(time.Second)
	return fmt.Sprintf("%02d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
}
