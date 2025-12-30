package main

import (
	"testing"

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
