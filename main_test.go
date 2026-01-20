package main

import (
	"os"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{"default range", "100..0", 100, 0, false},
		{"reverse range", "0..100", 0, 100, false},
		{"small range", "10..5", 10, 5, false},
		{"with spaces", "100 .. 0", 100, 0, false},
		{"negative numbers", "-10..10", -10, 10, false},
		{"invalid format no dots", "100-0", 0, 0, true},
		{"invalid start", "abc..0", 0, 0, true},
		{"invalid end", "100..xyz", 0, 0, true},
		{"empty", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := parseRange(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantStart, start)
				assert.Equal(t, tt.wantEnd, end)
			}
		})
	}
}

func TestParseFinalPhase(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		start   int
		end     int
		want    int
		wantErr bool
	}{
		{"absolute number", "5", 100, 0, 5, false},
		{"percentage 10%", "10%", 100, 0, 10, false},
		{"percentage 50%", "50%", 100, 0, 50, false},
		{"percentage with reverse", "10%", 0, 100, 110, false},
		{"invalid", "abc", 100, 0, 0, true},
		{"invalid percent", "abc%", 100, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFinalPhase(tt.val, tt.start, tt.end)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParsePadding(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantV   int
		wantH   int
		wantErr bool
	}{
		{"two values", "1 2", 1, 2, false},
		{"single value", "3", 3, 3, false},
		{"zeros", "0 0", 0, 0, false},
		{"invalid first", "abc 2", 0, 0, true},
		{"invalid second", "1 xyz", 0, 0, true},
		{"too many values", "1 2 3", 0, 0, true},
		{"empty", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, h, err := parsePadding(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantV, v)
				assert.Equal(t, tt.wantH, h)
			}
		})
	}
}

func TestCLIBigFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantBig  bool
		wantErr  bool
	}{
		{"big flag short", []string{"-b"}, true, false},
		{"big flag long", []string{"--big"}, true, false},
		{"big flag with range", []string{"-b", "-r", "10..0"}, true, false},
		{"no big flag", []string{"-r", "10..0"}, false, false},
		{"big flag with other options", []string{"-b", "-s", "dot", "-r", "5..0"}, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original args
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Set up test args
			os.Args = append([]string{"countdown"}, tt.args...)

			var cli CLI
			parser, err := kong.New(&cli,
				kong.Name("countdown"),
				kong.Description("Display spinner while displaying a number which counts downward"),
				kong.UsageOnError(),
			)
			require.NoError(t, err, "kong.New should not error")

			_, err = parser.Parse(tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantBig, cli.Big, "Big flag should match expected value")
			}
		})
	}
}
