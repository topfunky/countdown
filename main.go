package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/countdown/countdown/internal/countdown"
)

var version = "dev"

// CLI defines the command-line interface.
type CLI struct {
	Version      bool   `short:"v" help:"Print the version number"`
	Spinner      string `short:"s" default:"dot" help:"Spinner type" env:"COUNTDOWN_SPINNER" enum:"dot,line,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger,bomb,none"`
	Title        string `default:"Liftoff in" help:"Text to display to user while counting" env:"COUNTDOWN_TITLE"`
	Range        string `short:"r" default:"100..0" help:"Numbers to count from and to"`
	TimeInterval int    `short:"t" default:"1" help:"Number of seconds between each iteration"`
	Decrement    int    `short:"d" default:"1" help:"Number subtracted from current count at each iteration"`
	FinalPhase   string `short:"f" default:"5" help:"Number at which the final phase starts. At this number, the foreground and background colors are swapped. Can be a number such as '5' or a percentage such as '10%'"`

	SpinnerStyle SpinnerStyle `embed:"" prefix:"spinner."`
	TitleStyle   TitleStyle   `embed:"" prefix:"title."`
	Padding      string       `default:"0 0" help:"Padding" env:"COUNTDOWN_PADDING"`
}

// SpinnerStyle defines styling for the spinner.
type SpinnerStyle struct {
	Foreground string `default:"212" help:"Foreground Color" env:"COUNTDOWN_SPINNER_FOREGROUND"`
	Background string `default:"" help:"Background Color" env:"COUNTDOWN_SPINNER_BACKGROUND"`
}

// TitleStyle defines styling for the title.
type TitleStyle struct {
	Foreground string `default:"" help:"Foreground Color" env:"COUNTDOWN_TITLE_FOREGROUND"`
	Background string `default:"" help:"Background Color" env:"COUNTDOWN_TITLE_BACKGROUND"`
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("countdown"),
		kong.Description("Display spinner while displaying a number which counts downward"),
		kong.UsageOnError(),
	)

	if cli.Version {
		fmt.Printf("countdown %s\n", version)
		os.Exit(0)
	}

	// Parse range
	start, end, err := parseRange(cli.Range)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	// Parse final phase
	finalPhase, err := parseFinalPhase(cli.FinalPhase, start, end)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	// Parse padding
	padV, padH, err := parsePadding(cli.Padding)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	config := countdown.Config{
		SpinnerType:       cli.Spinner,
		Title:             cli.Title,
		Start:             start,
		End:               end,
		TimeInterval:      cli.TimeInterval,
		Decrement:         cli.Decrement,
		FinalPhase:        finalPhase,
		SpinnerForeground: cli.SpinnerStyle.Foreground,
		SpinnerBackground: cli.SpinnerStyle.Background,
		TitleForeground:   cli.TitleStyle.Foreground,
		TitleBackground:   cli.TitleStyle.Background,
		PaddingVertical:   padV,
		PaddingHorizontal: padH,
	}

	if err := countdown.Run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// parseRange parses a range string like "100..0" into start and end values.
func parseRange(r string) (int, int, error) {
	parts := strings.Split(r, "..")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format: %s (expected format: start..end)", r)
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start value in range: %s", parts[0])
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end value in range: %s", parts[1])
	}

	return start, end, nil
}

// parseFinalPhase parses the final phase value which can be a number or percentage.
func parseFinalPhase(val string, start, end int) (int, error) {
	val = strings.TrimSpace(val)

	if strings.HasSuffix(val, "%") {
		percentStr := strings.TrimSuffix(val, "%")
		percent, err := strconv.Atoi(percentStr)
		if err != nil {
			return 0, fmt.Errorf("invalid percentage in final-phase: %s", val)
		}

		total := abs(start - end)
		return end + (total * percent / 100), nil
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid final-phase value: %s", val)
	}

	return num, nil
}

// parsePadding parses padding string "vertical horizontal" into two values.
func parsePadding(p string) (int, int, error) {
	parts := strings.Fields(p)
	if len(parts) == 1 {
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid padding value: %s", parts[0])
		}
		return v, v, nil
	}
	if len(parts) == 2 {
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid vertical padding: %s", parts[0])
		}
		h, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid horizontal padding: %s", parts[1])
		}
		return v, h, nil
	}
	return 0, 0, fmt.Errorf("invalid padding format: %s (expected 'v h' or 'v')", p)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
