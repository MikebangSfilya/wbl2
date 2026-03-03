package unpackstring

import (
	"errors"
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  error
	}{
		{
			name:     "test_case1",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			name:     "test_case2",
			input:    "abcd",
			expected: "abcd",
		},
		{
			name:    "test_case3",
			input:   "45",
			wantErr: errDigit,
		},
		{
			name:     "test_case4",
			input:    "",
			expected: "",
		},
		{
			name:    "test_case5",
			input:   "a45",
			wantErr: errDigit,
		},
		{
			name:     "test_case5",
			input:    "qwe\\4\\5",
			expected: "qwe45",
		},
		{
			name:     "test_case5",
			input:    "qwe\\45",
			expected: "qwe44444",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
			} else {
				if got != tt.expected {
					t.Fatalf("expected %s, got %s", tt.expected, got)
				}
			}
		})
	}
}
