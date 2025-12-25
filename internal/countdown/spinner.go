package countdown

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
)

// SpinnerMap maps spinner names to their configurations.
var SpinnerMap = map[string]spinner.Spinner{
	"dot":       spinner.Dot,
	"line":      spinner.Line,
	"minidot":   spinner.MiniDot,
	"jump":      spinner.Jump,
	"pulse":     spinner.Pulse,
	"points":    spinner.Points,
	"globe":     spinner.Globe,
	"moon":      spinner.Moon,
	"monkey":    spinner.Monkey,
	"meter":     spinner.Meter,
	"hamburger": spinner.Hamburger,
	"bomb":      {Frames: []string{"💣", "💥"}, FPS: time.Second / 2},
	"none":      {Frames: []string{""}, FPS: time.Second},
}

// GetSpinner returns the spinner configuration for the given name.
func GetSpinner(name string) spinner.Spinner {
	if s, ok := SpinnerMap[name]; ok {
		return s
	}
	return spinner.Dot
}
