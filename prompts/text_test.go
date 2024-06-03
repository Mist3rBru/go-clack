package prompts_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Mist3rBru/go-clack/core"
	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/prompts/test"
	"github.com/Mist3rBru/go-clack/prompts/utils"
	"github.com/stretchr/testify/assert"
)

const message = "test message"

func TestTextInitialState(t *testing.T) {
	go prompts.Text(prompts.TextParams{Message: message})
	time.Sleep(1 * time.Millisecond)

	p := test.TextTestingPrompt
	startBar := utils.Color["gray"](utils.S_BAR)
	title := utils.SymbolState(core.StateInitial) + " " + message
	valueWithCursor := utils.Color["cyan"](utils.S_BAR) + " "
	endBar := utils.Color["gray"](utils.S_BAR_END)
	expected := strings.Join([]string{startBar, title, valueWithCursor, endBar}, "\r\n")
	assert.Equal(t, core.StateInitial, p.State)
	assert.Equal(t, expected, p.Frame)
}
