package utils

import (
	"testing"
)

func TestMurmur3_32(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		seed     uint32
		expected uint32
	}{
		{
			name:     "Empty string with seed 0",
			input:    "",
			seed:     0,
			expected: 0x00000000,
		},
		{
			name:     "Single digit '0' with seed 0",
			input:    "0",
			seed:     0,
			expected: 0xd271c07f,
		},
		{
			name:     "Two digits '01' with seed 0",
			input:    "01",
			seed:     0,
			expected: 0x61ec6600,
		},
		{
			name:     "Three digits '012' with seed 0",
			input:    "012",
			seed:     0,
			expected: 0xec6cff8c,
		},
		{
			name:     "Four digits '0123' with seed 0",
			input:    "0123",
			seed:     0,
			expected: 0xd41994a0,
		},
		{
			name:     "Five digits '01234' with seed 0",
			input:    "01234",
			seed:     0,
			expected: 0x19d02170,
		},
		{
			name:     "Single digit '2' with seed 0",
			input:    "2",
			seed:     0,
			expected: 0x0129e217,
		},
		{
			name:     "Two digits '88' with seed 0",
			input:    "88",
			seed:     0,
			expected: 0x7a0040a5,
		},
		{
			name:     "String 'asdfqwer' with seed 0",
			input:    "asdfqwer",
			seed:     0,
			expected: 0xa46b5209,
		},
		{
			name:     "String 'asdfqwerty' with seed 0",
			input:    "asdfqwerty",
			seed:     0,
			expected: 0xa3cfe04b,
		},
		{
			name:     "String 'asd' with seed 0",
			input:    "asd",
			seed:     0,
			expected: 0x14570c6f,
		},
		{
			name:     "String 'Hello' with seed 0",
			input:    "Hello",
			seed:     0,
			expected: 0x12da77c8,
		},
		{
			name:     "String 'Hello1' with seed 0",
			input:    "Hello1",
			seed:     0,
			expected: 0x6357e0a6,
		},
		{
			name:     "String 'Hello2' with seed 0",
			input:    "Hello2",
			seed:     0,
			expected: 0xe5ce223e,
		},
		{
			name:     "String 'hey' with seed 0",
			input:    "hey",
			seed:     0,
			expected: 0x12f94418,
		},
		{
			name:     "String 'dude' with seed 0",
			input:    "dude",
			seed:     0,
			expected: 0xef0487f3,
		},
		{
			name:     "String 'test' with seed 0",
			input:    "test",
			seed:     0,
			expected: 0xba6bd213,
		},
		{
			name:     "String 'kinkajou' with seed 0",
			input:    "kinkajou",
			seed:     0,
			expected: 0xb6d99cf8,
		},
		{
			name:     "Empty string with seed 1",
			input:    "",
			seed:     1,
			expected: 0x514e28b7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Murmur3_32([]byte(tt.input), tt.seed)
			if got != tt.expected {
				t.Errorf("Murmur3_32() = %x, want %x", got, tt.expected)
			}
		})
	}
}
