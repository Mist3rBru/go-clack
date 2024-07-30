package prompts_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/stretchr/testify/assert"
)

func runSpinner() (*prompts.SpinnerController, *MockTimer, *MockWriter) {
	timer := &MockTimer{}
	writer := &MockWriter{}
	s := prompts.Spinner(prompts.SpinnerOptions{
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
