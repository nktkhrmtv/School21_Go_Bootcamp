package main

import (
	"math"
	"os"
	"testing"
)

func TestReadNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
		err      string
	}{
		{
			name:     "valid input",
			input:    "1\n2\n3\nq\n",
			expected: []int{1, 2, 3},
			err:      "",
		},
		{
			name:     "invalid input (not a number)",
			input:    "1\nabc\nq\n",
			expected: nil,
			err:      "invalid input: abc is not a number",
		},
		{
			name:     "out of range input",
			input:    "100001\nq\n",
			expected: nil,
			err:      "number 100001 is out of range (-100000 to 100000)",
		},
		{
			name:     "empty input",
			input:    "q\n",
			expected: nil,
			err:      "no numbers provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "testinput")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name()) 
			
			if _, err := tmpFile.WriteString(tt.input); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			tmpFile.Close() 

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			file, err := os.Open(tmpFile.Name())
			if err != nil {
				t.Fatalf("failed to open temp file: %v", err)
			}
			os.Stdin = file

			var numbers []int
			err = readNumbers(&numbers)

			if err != nil && err.Error() != tt.err {
				t.Errorf("readNumbers() error = %v, wantErr %v", err, tt.err)
				return
			}
			if err == nil && !equalSlices(numbers, tt.expected) {
				t.Errorf("readNumbers() = %v, want %v", numbers, tt.expected)
			}
		})
	}
}

func TestStatistics(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		mean     float64
		median   float64
		mode     int
		sd       float64
	}{
		{
			name:     "simple case",
			numbers:  []int{1, 2, 3, 4, 5},
			mean:     3.0,
			median:   3.0,
			mode:     1,
			sd:       1.41,
		},
		{
			name:     "even number of elements",
			numbers:  []int{1, 2, 3, 4},
			mean:     2.5,
			median:   2.5,
			mode:     1,
			sd:       1.12,
		},
		{
			name:     "repeated numbers",
			numbers:  []int{2, 2, 3, 3, 3},
			mean:     2.6,
			median:   3.0,
			mode:     3,
			sd:       0.49,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := Init(tt.numbers)
			if mean := stats.Mean(); math.Abs(mean-tt.mean) > 0.01 {
				t.Errorf("Mean() = %v, want %v", mean, tt.mean)
			}
			if median := stats.Median(); math.Abs(median-tt.median) > 0.01 {
				t.Errorf("Median() = %v, want %v", median, tt.median)
			}
			if mode := stats.Mode(); mode != tt.mode {
				t.Errorf("Mode() = %v, want %v", mode, tt.mode)
			}
			if sd := stats.SD(); math.Abs(sd-tt.sd) > 0.01 {
				t.Errorf("SD() = %v, want %v", sd, tt.sd)
			}
		})
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

