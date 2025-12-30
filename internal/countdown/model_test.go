package countdown

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSpinner(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
	}{
		{"dot spinner", "dot", 8},
		{"line spinner", "line", 4},
		{"moon spinner", "moon", 8},
		{"bomb spinner", "bomb", 2},
		{"none spinner", "none", 1},
		{"unknown defaults to dot", "unknown", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GetSpinner(tt.input)
			assert.Equal(t, tt.wantLen, len(s.Frames), fmt.Sprintf("Spinner %s", tt.name))
		})
	}
}

func TestNewModel(t *testing.T) {
	cfg := Config{
		SpinnerType:       "dot",
		Title:             "Test",
		Start:             10,
		End:               0,
		TimeInterval:      1,
		Decrement:         1,
		FinalPhase:        2,
		SpinnerForeground: "212",
		PaddingVertical:   1,
		PaddingHorizontal: 2,
	}

	m := NewModel(cfg)

	assert.Equal(t, cfg.Start, m.current)
	assert.False(t, m.done)
	assert.Equal(t, cfg.Title, m.config.Title)
}

func TestModelIsInFinalPhase(t *testing.T) {
	tests := []struct {
		name       string
		start      int
		end        int
		current    int
		finalPhase int
		want       bool
	}{
		{"countdown not in final", 100, 0, 50, 5, false},
		{"countdown in final", 100, 0, 3, 5, true},
		{"countdown at final", 100, 0, 5, 5, true},
		{"countup not in final", 0, 100, 50, 95, false},
		{"countup in final", 0, 100, 97, 95, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				config: Config{
					Start:      tt.start,
					End:        tt.end,
					FinalPhase: tt.finalPhase,
				},
				current: tt.current,
			}

			assert.Equal(t, tt.want, m.isInFinalPhase())
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  lipgloss.TerminalColor
	}{
		{"empty string", "", lipgloss.NoColor{}},
		{"ansi number", "212", lipgloss.Color("212")},
		{"hex color", "#ff0000", lipgloss.Color("#ff0000")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseColor(tt.input)
			if tt.input == "" {
				_, ok := got.(lipgloss.NoColor)
				assert.True(t, ok, "parseColor should return NoColor for empty string")
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestModelView(t *testing.T) {
	cfg := Config{
		SpinnerType:  "none",
		Title:        "Test",
		Start:        10,
		End:          0,
		TimeInterval: 1,
		Decrement:    1,
		FinalPhase:   2,
	}

	m := NewModel(cfg)

	view := m.View()
	assert.NotEmpty(t, view, "View() should not return empty string when not done")

	m.done = true
	view = m.View()
	assert.Empty(t, view, "View() should return empty string when done")
}

func TestModelViewWithKilled(t *testing.T) {
	cfg := Config{
		SpinnerType:  "none",
		Title:        "Test",
		Start:        10,
		End:          0,
		TimeInterval: 1,
		Decrement:    1,
		FinalPhase:   2,
	}

	m := NewModel(cfg)
	m.killed = true

	view := m.View()
	assert.NotEmpty(t, view)
	assert.Contains(t, view, "(killed)")
}

func TestHighContrastColor(t *testing.T) {
	tests := []struct {
		name    string
		bgColor string
		want    string
	}{
		{"empty defaults to white text", "", "15"},
		{"black bg gets white text", "0", "15"},
		{"white bg gets black text", "15", "0"},
		{"bright yellow gets black text", "11", "0"},
		{"dark blue gets white text", "4", "15"},
		{"pink 212 gets black text", "212", "0"},
		{"hex white gets black text", "#ffffff", "0"},
		{"hex black gets white text", "#000000", "15"},
		{"hex red gets white text", "#ff0000", "15"},
		{"grayscale light gets black", "255", "0"},
		{"grayscale dark gets white", "232", "15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := highContrastColor(tt.bgColor)
			gotColor, ok := got.(lipgloss.Color)
			require.True(t, ok, "highContrastColor should return lipgloss.Color")
			assert.Equal(t, tt.want, string(gotColor))
		})
	}
}

func TestAnsi256ToRGB(t *testing.T) {
	tests := []struct {
		name  string
		input int
		wantR uint8
		wantG uint8
		wantB uint8
	}{
		{"black", 0, 0, 0, 0},
		{"white", 15, 255, 255, 255},
		{"red", 1, 128, 0, 0},
		{"bright red", 9, 255, 0, 0},
		{"color cube start", 16, 0, 0, 0},
		{"grayscale mid", 244, 128, 128, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b := ansi256ToRGB(tt.input)
			assert.Equal(t, tt.wantR, r)
			assert.Equal(t, tt.wantG, g)
			assert.Equal(t, tt.wantB, b)
		})
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name  string
		input string
		wantR uint8
		wantG uint8
		wantB uint8
	}{
		{"white", "#ffffff", 255, 255, 255},
		{"black", "#000000", 0, 0, 0},
		{"red", "#ff0000", 255, 0, 0},
		{"green", "#00ff00", 0, 255, 0},
		{"blue", "#0000ff", 0, 0, 255},
		{"short white", "#fff", 255, 255, 255},
		{"short black", "#000", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b := hexToRGB(tt.input)
			assert.Equal(t, tt.wantR, r)
			assert.Equal(t, tt.wantG, g)
			assert.Equal(t, tt.wantB, b)
		})
	}
}
