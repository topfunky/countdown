package countdown

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestGetSpinner(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int // expected number of frames
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
			if len(s.Frames) != tt.wantLen {
				t.Errorf("GetSpinner(%q) has %d frames, want %d", tt.input, len(s.Frames), tt.wantLen)
			}
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

	if m.current != cfg.Start {
		t.Errorf("NewModel current = %d, want %d", m.current, cfg.Start)
	}
	if m.done {
		t.Error("NewModel should not be done initially")
	}
	if m.config.Title != cfg.Title {
		t.Errorf("NewModel title = %q, want %q", m.config.Title, cfg.Title)
	}
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

			if got := m.isInFinalPhase(); got != tt.want {
				t.Errorf("isInFinalPhase() = %v, want %v", got, tt.want)
			}
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
			// For NoColor, check type
			if tt.input == "" {
				if _, ok := got.(lipgloss.NoColor); !ok {
					t.Errorf("parseColor(%q) should return NoColor", tt.input)
				}
				return
			}
			// For Color, compare string representation
			if got != tt.want {
				t.Errorf("parseColor(%q) = %v, want %v", tt.input, got, tt.want)
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

	// View should contain the title and current count
	view := m.View()
	if view == "" {
		t.Error("View() should not return empty string when not done")
	}

	// When done, view should be empty
	m.done = true
	view = m.View()
	if view != "" {
		t.Errorf("View() should return empty string when done, got %q", view)
	}
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

	// View should contain "(killed)" when killed
	view := m.View()
	if view == "" {
		t.Error("View() should not return empty string when not done")
	}
	if !strings.Contains(view, "(killed)") {
		t.Errorf("View() should contain '(killed)' when killed, got %q", view)
	}
}

func TestHighContrastColor(t *testing.T) {
	tests := []struct {
		name    string
		bgColor string
		want    string // "0" for black, "15" for white
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
			if !ok {
				t.Errorf("highContrastColor(%q) did not return lipgloss.Color", tt.bgColor)
				return
			}
			if string(gotColor) != tt.want {
				t.Errorf("highContrastColor(%q) = %q, want %q", tt.bgColor, gotColor, tt.want)
			}
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
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("ansi256ToRGB(%d) = (%d,%d,%d), want (%d,%d,%d)",
					tt.input, r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
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
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("hexToRGB(%q) = (%d,%d,%d), want (%d,%d,%d)",
					tt.input, r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}
