package utils

import (
	"testing"
)

func TestFNV1_Hash(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected uint64
	}{
		{
			name:     "Empty string",
			input:    []byte(""),
			expected: 0xcbf29ce484222325,
		},
		{
			name:     "Single byte '0'",
			input:    []byte("0"),
			expected: 0xaf63bd4c8601b7ef,
		},
		{
			name:     "Two bytes '01'",
			input:    []byte("01"),
			expected: 0x08329807b4eb8b2c,
		},
		{
			name:     "String 'hello'",
			input:    []byte("hello"),
			expected: 0x7b495389bdbdd4c7,
		},
		{
			name:     "String 'test'",
			input:    []byte("test"),
			expected: 0x8c093f7e9fccbf69,
		},
		{
			name:     "Long string",
			input:    []byte("The quick brown fox jumps over the lazy dog"),
			expected: 0xa8b2f3117de37ace,
		},
	}

	fnv := NewFNV1()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fnv.Hash(tt.input)
			if result != tt.expected {
				t.Errorf("Hash(%s) = %x; want %x", tt.input, result, tt.expected)
			}
		})
	}
}
