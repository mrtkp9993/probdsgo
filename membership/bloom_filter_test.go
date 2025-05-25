package membership

import (
	"testing"
)

func TestNewBloomFilter(t *testing.T) {
	tests := []struct {
		name        string
		capacity    uint
		errorRate   float64
		shouldError bool
	}{
		{
			name:        "Valid parameters",
			capacity:    216553,
			errorRate:   0.01,
			shouldError: false,
		},
		{
			name:        "Zero capacity",
			capacity:    0,
			errorRate:   0.01,
			shouldError: true,
		},
		{
			name:        "Error rate too high",
			capacity:    1000,
			errorRate:   1.5,
			shouldError: true,
		},
		{
			name:        "Error rate too low",
			capacity:    1000,
			errorRate:   0.0,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf, err := NewBloomFilter(tt.capacity, tt.errorRate)
			if tt.shouldError {
				if err == nil {
					t.Errorf("NewBloomFilter() error = nil, expected an error")
				}
				return
			}
			if err != nil {
				t.Errorf("NewBloomFilter() error = %v, expected no error", err)
				return
			}
			if bf == nil {
				t.Error("NewBloomFilter() returned nil BloomFilter")
			}
		})
	}
}

func TestNewBloomFilterWithParams(t *testing.T) {
	tests := []struct {
		name        string
		m           uint
		k           uint
		shouldError bool
	}{
		{
			name:        "Valid parameters",
			m:           1000,
			k:           3,
			shouldError: false,
		},
		{
			name:        "Zero bit array size",
			m:           0,
			k:           3,
			shouldError: true,
		},
		{
			name:        "Zero hash functions",
			m:           1000,
			k:           0,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf, err := NewBloomFilterWithParams(tt.m, tt.k)
			if tt.shouldError {
				if err == nil {
					t.Errorf("NewBloomFilterWithParams() error = nil, expected an error")
				}
				return
			}
			if err != nil {
				t.Errorf("NewBloomFilterWithParams() error = %v, expected no error", err)
				return
			}
			if bf == nil {
				t.Error("NewBloomFilterWithParams() returned nil BloomFilter")
			}
		})
	}
}

func TestBloomFilter_Add(t *testing.T) {
	bf, err := NewBloomFilter(1000, 0.01)
	if err != nil {
		t.Fatalf("Failed to create BloomFilter: %v", err)
	}

	tests := []struct {
		name        string
		input       []byte
		shouldError bool
	}{
		{
			name:        "Add valid item",
			input:       []byte("test item"),
			shouldError: false,
		},
		{
			name:        "Add empty item",
			input:       []byte{},
			shouldError: true,
		},
		{
			name:        "Add nil item",
			input:       nil,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bf.Add(tt.input)
			if tt.shouldError {
				if err == nil {
					t.Error("Add() error = nil, expected an error")
				}
				return
			}
			if err != nil {
				t.Errorf("Add() error = %v, expected no error", err)
			}
		})
	}
}

func TestBloomFilter_Contains(t *testing.T) {
	bf, err := NewBloomFilter(1000, 0.01)
	if err != nil {
		t.Fatalf("Failed to create BloomFilter: %v", err)
	}

	// Add some items
	items := [][]byte{
		[]byte("item1"),
		[]byte("item2"),
		[]byte("item3"),
	}
	for _, item := range items {
		if err := bf.Add(item); err != nil {
			t.Fatalf("Failed to add item: %v", err)
		}
	}

	tests := []struct {
		name        string
		input       []byte
		shouldExist bool
		shouldError bool
	}{
		{
			name:        "Check existing item",
			input:       []byte("item1"),
			shouldExist: true,
			shouldError: false,
		},
		{
			name:        "Check non-existing item",
			input:       []byte("non-existing"),
			shouldExist: false,
			shouldError: false,
		},
		{
			name:        "Check empty item",
			input:       []byte{},
			shouldExist: false,
			shouldError: true,
		},
		{
			name:        "Check nil item",
			input:       nil,
			shouldExist: false,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := bf.Contains(tt.input)
			if tt.shouldError {
				if err == nil {
					t.Error("Contains() error = nil, expected an error")
				}
				return
			}
			if err != nil {
				t.Errorf("Contains() error = %v, expected no error", err)
				return
			}
			if exists != tt.shouldExist {
				t.Errorf("Contains() = %v, want %v", exists, tt.shouldExist)
			}
		})
	}
}
