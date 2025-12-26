package countdown

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Config holds the countdown configuration.
type Config struct {
	SpinnerType       string
	Title             string
	Start             int
	End               int
	TimeInterval      int
	Decrement         int
	FinalPhase        int
	SpinnerForeground string
	SpinnerBackground string
	TitleForeground   string
	TitleBackground   string
	PaddingVertical   int
	PaddingHorizontal int
}

// Model represents the Bubbletea model for the countdown.
type Model struct {
	config         Config
	spinner        spinner.Model
	current        int
	done           bool
	spinnerStyle   lipgloss.Style
	titleStyle     lipgloss.Style
	countStyle     lipgloss.Style
	containerStyle lipgloss.Style
}

// tickMsg is sent when the countdown should decrement.
type tickMsg struct{}

// NewModel creates a new countdown model.
func NewModel(cfg Config) Model {
	s := spinner.New()
	s.Spinner = GetSpinner(cfg.SpinnerType)

	// Build spinner style
	spinnerStyle := lipgloss.NewStyle()
	if cfg.SpinnerForeground != "" {
		spinnerStyle = spinnerStyle.Foreground(parseColor(cfg.SpinnerForeground))
	}
	if cfg.SpinnerBackground != "" {
		spinnerStyle = spinnerStyle.Background(parseColor(cfg.SpinnerBackground))
	}
	s.Style = spinnerStyle

	// Build title style
	titleStyle := lipgloss.NewStyle()
	if cfg.TitleForeground != "" {
		titleStyle = titleStyle.Foreground(parseColor(cfg.TitleForeground))
	}
	if cfg.TitleBackground != "" {
		titleStyle = titleStyle.Background(parseColor(cfg.TitleBackground))
	}

	// Count style (same as title by default)
	countStyle := titleStyle

	// Container style with padding
	containerStyle := lipgloss.NewStyle().
		PaddingTop(cfg.PaddingVertical).
		PaddingBottom(cfg.PaddingVertical).
		PaddingLeft(cfg.PaddingHorizontal).
		PaddingRight(cfg.PaddingHorizontal)

	return Model{
		config:         cfg,
		spinner:        s,
		current:        cfg.Start,
		spinnerStyle:   spinnerStyle,
		titleStyle:     titleStyle,
		countStyle:     countStyle,
		containerStyle: containerStyle,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tick(m.config.TimeInterval))
}

// tick returns a command that sends a tickMsg after the specified interval.
func tick(seconds int) tea.Cmd {
	return tea.Tick(time.Duration(seconds)*time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.done = true
			return m, tea.Quit
		}

	case tickMsg:
		// Determine direction
		if m.config.Start > m.config.End {
			m.current -= m.config.Decrement
			if m.current <= m.config.End {
				m.current = m.config.End
				m.done = true
				return m, tea.Quit
			}
		} else {
			m.current += m.config.Decrement
			if m.current >= m.config.End {
				m.current = m.config.End
				m.done = true
				return m, tea.Quit
			}
		}
		return m, tick(m.config.TimeInterval)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the model.
func (m Model) View() string {
	if m.done {
		return ""
	}

	// Check if we're in final phase
	inFinalPhase := m.isInFinalPhase()

	// Build the spinner view
	spinnerView := m.spinner.View()

	// Build the title and count with potential style swap in final phase.
	//
	// Add space to title for unbroken display when inverted.
	titleStr := fmt.Sprintf("%s ", m.config.Title)
	countStr := strconv.Itoa(m.current)
	var titleView string
	var countView string
	if inFinalPhase && m.current%2 == 1 {
		// Final phase: foreground becomes background, text is high-contrast
		finalStyle := lipgloss.NewStyle()

		// Determine the foreground color to use as background
		fgColor := m.config.SpinnerForeground // Default foreground
		if m.config.TitleForeground != "" {
			fgColor = m.config.TitleForeground
		}

		// Set the original foreground as the new background
		if fgColor != "" {
			finalStyle = finalStyle.Background(parseColor(fgColor))
		} else {
			// Use default spinner color (212) as background
			fgColor = "212"
			finalStyle = finalStyle.Background(parseColor(fgColor))
		}

		// Calculate high-contrast foreground for readability
		finalStyle = finalStyle.Foreground(highContrastColor(fgColor))
		finalStyle = finalStyle.Bold(true)

		titleView = finalStyle.Render(titleStr)
		countView = finalStyle.Render(countStr)
	} else {
		titleView = m.titleStyle.Render(titleStr)
		countView = m.countStyle.Render(countStr)
	}

	// Combine all parts
	content := fmt.Sprintf("%s %s%s", spinnerView, titleView, countView)

	return m.containerStyle.Render(content)
}

// isInFinalPhase checks if the current count is in the final phase.
func (m Model) isInFinalPhase() bool {
	if m.config.Start > m.config.End {
		// Counting down
		return m.current <= m.config.FinalPhase
	}
	// Counting up
	return m.current >= m.config.FinalPhase
}

// parseColor parses a color string and returns a lipgloss.TerminalColor.
func parseColor(s string) lipgloss.TerminalColor {
	s = strings.TrimSpace(s)
	if s == "" {
		return lipgloss.NoColor{}
	}

	// Check if it's a number (ANSI color)
	if _, err := strconv.Atoi(s); err == nil {
		return lipgloss.Color(s)
	}

	// Otherwise treat as hex or named color
	return lipgloss.Color(s)
}

// highContrastColor returns a high-contrast foreground color (black or white)
// for the given background color string.
func highContrastColor(bgColor string) lipgloss.TerminalColor {
	bgColor = strings.TrimSpace(bgColor)
	if bgColor == "" {
		return lipgloss.Color("15") // White for default/empty background
	}

	r, g, b := colorToRGB(bgColor)
	luminance := calcLuminance(r, g, b)

	// Hard code a threshold which works in practice
	if luminance > 0.4 {
		return lipgloss.Color("0") // Black
	}
	return lipgloss.Color("15") // White
}

// colorToRGB converts a color string to RGB values (0-255).
func colorToRGB(s string) (r, g, b uint8) {
	s = strings.TrimSpace(s)

	// Handle hex colors
	if strings.HasPrefix(s, "#") {
		return hexToRGB(s)
	}

	// Handle ANSI 256 colors
	if num, err := strconv.Atoi(s); err == nil {
		return ansi256ToRGB(num)
	}

	// Default to white if unknown
	return 255, 255, 255
}

// hexToRGB converts a hex color string to RGB.
func hexToRGB(hex string) (r, g, b uint8) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		// Short form #RGB -> #RRGGBB
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}

	if len(hex) != 6 {
		return 255, 255, 255
	}

	var rgb uint64
	rgb, _ = strconv.ParseUint(hex, 16, 32)
	return uint8(rgb >> 16), uint8((rgb >> 8) & 0xFF), uint8(rgb & 0xFF)
}

// ansi256ToRGB converts an ANSI 256 color number to RGB.
func ansi256ToRGB(n int) (r, g, b uint8) {
	if n < 0 || n > 255 {
		return 255, 255, 255
	}

	// Standard colors (0-15)
	if n < 16 {
		// Basic ANSI colors approximate RGB values
		standard := [][3]uint8{
			{0, 0, 0},       // 0: Black
			{128, 0, 0},     // 1: Red
			{0, 128, 0},     // 2: Green
			{128, 128, 0},   // 3: Yellow
			{0, 0, 128},     // 4: Blue
			{128, 0, 128},   // 5: Magenta
			{0, 128, 128},   // 6: Cyan
			{192, 192, 192}, // 7: White
			{128, 128, 128}, // 8: Bright Black
			{255, 0, 0},     // 9: Bright Red
			{0, 255, 0},     // 10: Bright Green
			{255, 255, 0},   // 11: Bright Yellow
			{0, 0, 255},     // 12: Bright Blue
			{255, 0, 255},   // 13: Bright Magenta
			{0, 255, 255},   // 14: Bright Cyan
			{255, 255, 255}, // 15: Bright White
		}
		return standard[n][0], standard[n][1], standard[n][2]
	}

	// Color cube (16-231): 6x6x6 cube
	if n < 232 {
		n -= 16
		ri := n / 36
		gi := (n % 36) / 6
		bi := n % 6

		// Convert 0-5 to 0-255 (0, 95, 135, 175, 215, 255)
		toVal := func(i int) uint8 {
			if i == 0 {
				return 0
			}
			return uint8(55 + i*40)
		}
		return toVal(ri), toVal(gi), toVal(bi)
	}

	// Grayscale (232-255): 24 shades
	gray := uint8(8 + (n-232)*10)
	return gray, gray, gray
}

// calcLuminance calculates relative luminance using sRGB.
func calcLuminance(r, g, b uint8) float64 {
	// Convert to linear RGB
	toLinear := func(v uint8) float64 {
		f := float64(v) / 255.0
		if f <= 0.03928 {
			return f / 12.92
		}
		return math.Pow((f+0.055)/1.055, 2.4)
	}

	rLin := toLinear(r)
	gLin := toLinear(g)
	bLin := toLinear(b)

	// Calculate luminance (ITU-R BT.709)
	return 0.2126*rLin + 0.7152*gLin + 0.0722*bLin
}

// Run starts the countdown application.
func Run(cfg Config) error {
	p := tea.NewProgram(NewModel(cfg))
	_, err := p.Run()
	return err
}
