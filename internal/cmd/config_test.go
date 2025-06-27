package cmd

import (
	"testing"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{
			name:     "true string",
			input:    "true",
			expected: true,
			hasError: false,
		},
		{
			name:     "false string",
			input:    "false",
			expected: false,
			hasError: false,
		},
		{
			name:     "t string",
			input:    "t",
			expected: true,
			hasError: false,
		},
		{
			name:     "f string",
			input:    "f",
			expected: false,
			hasError: false,
		},
		{
			name:     "1 string",
			input:    "1",
			expected: true,
			hasError: false,
		},
		{
			name:     "0 string",
			input:    "0",
			expected: false,
			hasError: false,
		},
		{
			name:     "invalid string",
			input:    "invalid",
			expected: false,
			hasError: true,
		},
		{
			name:     "TRUE uppercase",
			input:    "TRUE",
			expected: true,
			hasError: false,
		},
		{
			name:     "FALSE uppercase",
			input:    "FALSE",
			expected: false,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseBool(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("parseBool(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("parseBool(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result != tt.expected {
				t.Errorf("parseBool(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}