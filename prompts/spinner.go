package prompts

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Mist3rBru/go-clack/prompts/symbols"
	isunicodesupported "github.com/Mist3rBru/go-clack/third_party/is-unicode-supported"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
	"github.com/Mist3rBru/go-clack/third_party/sisteransi"
)

type Timer interface {
	Sleep(duration time.Duration)
}

type DefaultTimer struct{}

func (t *DefaultTimer) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

type SpinnerOptions struct {
	Timer  Timer
	Output io.Writer
}

type SpinnerController struct {
	Start   func(msg string)
	Stop    func(msg string, code int)
	Message func(msg string)
}

func Spinner(ctx context.Context, options SpinnerOptions) *SpinnerController {
	_ctx, cancel := context.WithCancel(ctx)

	var timer Timer
	if options.Timer == nil {
		timer = &DefaultTimer{}
	} else {
		timer = options.Timer
	}

	var output io.Writer
	if options.Output == nil {
		output = os.Stdout
	} else {
		output = options.Output
	}

	var message, prevMessage string

	var frames []string
	var frameIndex, frameInterval int

	const dotsInterval float32 = 0.125
	var dotsTimer float32

	if isunicodesupported.IsUnicodeSupported() {
		frames = []string{"◒", "◐", "◓", "◑"}
		frameInterval = 80
	} else {
		frames = []string{"•", "o", "O", "0"}
		frameInterval = 120
	}

	write := func(str string) {
		output.Write([]byte(str))
	}

	clearPrevMessage := func() {
		write(sisteransi.MoveCursor(-len(strings.Split(prevMessage, "\n"))+1, -999))
		write(sisteransi.EraseDown())
	}

	return &SpinnerController{
		Start: func(msg string) {
			write(sisteransi.HideCursor())
			write(picocolors.Gray(symbols.BAR) + "\n")

			frameIndex = 0
			dotsTimer = 0
			message = parseMessage(msg)

			go func() {
				for {
					select {
					case <-_ctx.Done():
						return
					default:
						clearPrevMessage()
						prevMessage = message
						frame := picocolors.Magenta(frames[frameIndex])
						loadingDots := strings.Repeat(".", min(int(math.Floor(float64(dotsTimer))), 3))
						write(fmt.Sprintf("%s %s%s", frame, message, loadingDots))
						if frameIndex+1 < len(frames) {
							frameIndex++
						} else {
							frameIndex = 0
						}
						if int(dotsTimer) < 4 {
							dotsTimer += dotsInterval
						} else {
							dotsTimer = 0
						}
						timer.Sleep(time.Duration(frameInterval) * time.Millisecond)
					}
				}
			}()
		},
		Message: func(msg string) {
			message = parseMessage(msg)
		},
		Stop: func(msg string, code int) {
			cancel()
			clearPrevMessage()
			var step string
			switch code {
			case 0:
				step = picocolors.Green(symbols.STEP_SUBMIT)
			case 1:
				step = picocolors.Red(symbols.STEP_CANCEL)
			default:
				step = picocolors.Red(symbols.STEP_ERROR)
			}
			if msg != "" {
				message = parseMessage(msg)
			}
			write(sisteransi.ShowCursor())
			write(fmt.Sprintf("%s %s\n", step, message))
		},
	}
}

func parseMessage(msg string) string {
	dotsRegex := regexp.MustCompile(`\.+$`)
	return dotsRegex.ReplaceAllString(msg, "")
}
