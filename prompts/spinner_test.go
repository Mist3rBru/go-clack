package prompts_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

type MockTimer struct {
	mu      sync.Mutex
	waiters []chan struct{}
}

func (t *MockTimer) Sleep(duration time.Duration) {
	waiter := make(chan struct{})
	t.mu.Lock()
	t.waiters = append(t.waiters, waiter)
	t.mu.Unlock()
	<-waiter
}

func (m *MockTimer) ResolveAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, waiter := range m.waiters {
		close(waiter)
	}
	m.waiters = []chan struct{}(nil)
}

type MockWriter struct {
	mu   sync.Mutex
	Data []string
}

func (w *MockWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	w.Data = append(w.Data, string(data))
	w.mu.Unlock()
	return 0, nil
}

func runSpinner() (*prompts.SpinnerController, *MockTimer, *MockWriter) {
	timer := &MockTimer{}
	writer := &MockWriter{}
	s, _ := prompts.Spinner(context.TODO(), prompts.SpinnerOptions{
		Timer:  timer,
		Output: writer,
	})
	return s, timer, writer
}

func TestSpinnerFrameAnimation(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading")
	for i := 0; i < 5; i++ {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}

	assert.Equal(t, "◒ Loading", mw.Data[4])
	assert.Equal(t, "◐ Loading", mw.Data[7])
	assert.Equal(t, "◓ Loading", mw.Data[10])
	assert.Equal(t, "◑ Loading", mw.Data[13])
}

func TestSpinnerDotsAnimation(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading")

	for mw.Data[len(mw.Data)-1] != "◒ Loading" {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))

	for mw.Data[len(mw.Data)-1] != "◒ Loading." {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading.", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))

	for mw.Data[len(mw.Data)-1] != "◒ Loading.." {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading..", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))

	for mw.Data[len(mw.Data)-1] != "◒ Loading..." {
		mt.ResolveAll()
		time.Sleep(time.Microsecond)
	}
	assert.Equal(t, "◒ Loading...", mw.Data[len(mw.Data)-1], fmt.Sprint(len(mw.Data)))
}

func TestSpinnerRemoveDotsFromMessage(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading...")
	time.Sleep(time.Microsecond)
	mt.ResolveAll()
	time.Sleep(time.Microsecond)

	assert.Equal(t, "◒ Loading", mw.Data[4])
}

func TestSpinnerMessageMethod(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading...")
	time.Sleep(time.Millisecond)
	s.Message("Still Loading")
	mt.ResolveAll()
	time.Sleep(time.Millisecond)

	assert.Equal(t, "◐ Still Loading", mw.Data[7])
}

func TestSpinnerStopMessage(t *testing.T) {
	s, mt, mw := runSpinner()

	s.Start("Loading...")
	time.Sleep(time.Millisecond)
	s.Stop("Loaded", 0)
	mt.ResolveAll()
	time.Sleep(time.Millisecond)

	assert.Equal(t, "◇ Loaded\n", mw.Data[8])
}
