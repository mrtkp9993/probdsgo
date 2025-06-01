package membership

import (
	"fmt"
	"testing"
)

func TestNewCuckooFilter(t *testing.T) {
	tests := []struct {
		name       string
		capacity   uint
		bucketSize uint
		fpr        float64
		wantErr    bool
	}{
		{
			name:       "valid parameters",
			capacity:   1000,
			bucketSize: 4,
			fpr:        0.01,
			wantErr:    false,
		},
		{
			name:       "invalid capacity",
			capacity:   0,
			bucketSize: 4,
			fpr:        0.01,
			wantErr:    true,
		},
		{
			name:       "invalid bucket size",
			capacity:   1000,
			bucketSize: 3,
			fpr:        0.01,
			wantErr:    true,
		},
		{
			name:       "invalid fpr",
			capacity:   1000,
			bucketSize: 4,
			fpr:        1.5,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCuckooFilter(tt.capacity, tt.bucketSize, tt.fpr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCuckooFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCuckooFilter_BasicOperations(t *testing.T) {
	cf, err := NewCuckooFilter(1000, 4, 0.01)
	if err != nil {
		t.Fatalf("Failed to create CuckooFilter: %v", err)
	}

	// Test Insert and Lookup
	item1 := []byte("test1")
	item2 := []byte("test2")

	if err := cf.Insert(item1); err != nil {
		t.Errorf("Insert() error = %v", err)
	}

	if !cf.Lookup(item1) {
		t.Error("Lookup() failed to find inserted item")
	}

	if cf.Lookup(item2) {
		t.Error("Lookup() found non-existent item")
	}

	// Test Delete
	if !cf.Delete(item1) {
		t.Error("Delete() failed to remove existing item")
	}

	if cf.Delete(item2) {
		t.Error("Delete() removed non-existent item")
	}

	if cf.Lookup(item1) {
		t.Error("Lookup() found deleted item")
	}
}

func TestCuckooFilter_FullCapacity(t *testing.T) {
	capacity := uint(1000)
	bucketSize := uint(4)
	cf, err := NewCuckooFilter(capacity, bucketSize, 0.01)
	if err != nil {
		t.Fatalf("Failed to create CuckooFilter: %v", err)
	}

	// Calculate maximum items based on load factor
	maxItems := uint(float64(capacity) * LOAD_FACTOR_MAP[bucketSize])
	insertedCount := uint(0)

	// Try to insert items
	for i := uint(0); i < capacity*2; i++ {
		err := cf.Insert([]byte(fmt.Sprintf("item%d", i)))
		if err != nil {
			if insertedCount >= maxItems {
				// Expected to fail after reaching load factor
				return
			}
			t.Errorf("Insert failed before reaching load factor: %v", err)
			return
		}
		insertedCount++
	}

	if insertedCount > maxItems {
		t.Errorf("Inserted more items than allowed by load factor: got %d, want <= %d", insertedCount, maxItems)
	}
}
