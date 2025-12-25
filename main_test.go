package main

import (
	"testing"
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
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if start != tt.wantStart {
					t.Errorf("parseRange(%q) start = %d, want %d", tt.input, start, tt.wantStart)
				}
				if end != tt.wantEnd {
					t.Errorf("parseRange(%q) end = %d, want %d", tt.input, end, tt.wantEnd)
				}
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
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFinalPhase(%q) error = %v, wantErr %v", tt.val, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseFinalPhase(%q) = %d, want %d", tt.val, got, tt.want)
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
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePadding(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if v != tt.wantV {
					t.Errorf("parsePadding(%q) vertical = %d, want %d", tt.input, v, tt.wantV)
				}
				if h != tt.wantH {
					t.Errorf("parsePadding(%q) horizontal = %d, want %d", tt.input, h, tt.wantH)
				}
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
	}

	for _, tt := range tests {
		if got := abs(tt.input); got != tt.want {
			t.Errorf("abs(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
