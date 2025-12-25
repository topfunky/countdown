package countdown

import (
	"fmt"
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
	countStyle := titleStyle.Copy()

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

	// Build the title
	titleView := m.titleStyle.Render(m.config.Title)

	// Build the count with potential style swap in final phase
	countStr := strconv.Itoa(m.current)
	var countView string
	if inFinalPhase {
		// Swap foreground and background for final phase
		finalStyle := m.countStyle.Copy()
		if m.config.TitleForeground != "" || m.config.SpinnerForeground != "" {
			fg := m.config.TitleForeground
			if fg == "" {
				fg = m.config.SpinnerForeground
			}
			finalStyle = finalStyle.Background(parseColor(fg))
		}
		if m.config.TitleBackground != "" || m.config.SpinnerBackground != "" {
			bg := m.config.TitleBackground
			if bg == "" {
				bg = m.config.SpinnerBackground
			}
			finalStyle = finalStyle.Foreground(parseColor(bg))
		}
		// If no colors were set, use default inverse
		if m.config.TitleForeground == "" && m.config.SpinnerForeground == "" {
			finalStyle = finalStyle.Reverse(true)
		}
		countView = finalStyle.Render(countStr)
	} else {
		countView = m.countStyle.Render(countStr)
	}

	// Combine all parts
	content := fmt.Sprintf("%s %s %s", spinnerView, titleView, countView)

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

// Run starts the countdown application.
func Run(cfg Config) error {
	p := tea.NewProgram(NewModel(cfg))
	_, err := p.Run()
	return err
}
