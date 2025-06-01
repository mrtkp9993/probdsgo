package membership

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/mrtkp9993/probdsgo/utils"
)

var LOAD_FACTOR_MAP = map[uint]float64{
	2: 0.84,
	4: 0.95,
	8: 0.98,
}

const (
	// MaxKicks is the maximum number of relocations before declaring failure
	MaxKicks = 500
)

// Bucket represents a bucket in the Cuckoo filter
type Bucket struct {
	fingerprints []byte
	size         uint
}

// CuckooFilter represents a probabilistic data structure for membership testing
type CuckooFilter struct {
	buckets             []Bucket
	hashFunc            *utils.Murmur3
	fingerprintHashFunc *utils.Murmur3
	rng                 *rand.Rand

	bucketSize      uint
	fingerprintSize uint
	count           uint
}

func NewCuckooFilter(capacity uint, bucketSize uint, desiredFpr float64) (*CuckooFilter, error) {
	if capacity < 1 || desiredFpr <= 0 || desiredFpr >= 1 {
		return nil, errors.New("invalid capacity or desired false positive rate")
	}

	if bucketSize == 0 {
		return nil, errors.New("bucket size cannot be zero")
	}

	if _, ok := LOAD_FACTOR_MAP[bucketSize]; !ok {
		return nil, errors.New("invalid bucket size, must be 2, 4, or 8")
	}

	fingerprintSize := uint(math.Ceil(math.Log2(2 * float64(bucketSize) / desiredFpr)))
	if fingerprintSize > 8 {
		fingerprintSize = 8
	} else if fingerprintSize < 1 {
		fingerprintSize = 1
	}

	loadFactor := LOAD_FACTOR_MAP[bucketSize]
	numBuckets := uint(math.Ceil(float64(capacity) / (loadFactor * float64(bucketSize))))
	numBuckets = getNextPowerOf2(numBuckets)

	buckets := make([]Bucket, numBuckets)
	for i := range buckets {
		buckets[i] = Bucket{
			fingerprints: make([]byte, bucketSize),
			size:         0,
		}
	}

	hashFunc := utils.NewMurmur3WithSeed(0)
	fingerprintHashFunc := utils.NewMurmur3WithSeed(1)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &CuckooFilter{
		buckets:             buckets,
		hashFunc:            hashFunc,
		fingerprintHashFunc: fingerprintHashFunc,
		count:               0,
		rng:                 rng,
		bucketSize:          bucketSize,
		fingerprintSize:     fingerprintSize,
	}, nil
}

func getNextPowerOf2(n uint) uint {
	if n == 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}

// generateFingerprint creates a fingerprint for the given item
func (cf *CuckooFilter) generateFingerprint(item []byte) byte {
	hash := cf.fingerprintHashFunc.Hash(item)
	mask := byte((1 << cf.fingerprintSize) - 1)
	fp := byte(hash) & mask
	if fp == 0 {
		fp = 1
	}
	return fp
}

// getIndices returns the two possible bucket indices for an item
func (cf *CuckooFilter) getIndices(item []byte, fingerprint byte) (uint, uint) {
	hash1 := cf.hashFunc.Hash(item)
	index1 := uint(hash1) & (uint(len(cf.buckets)) - 1)

	hash2 := cf.fingerprintHashFunc.Hash([]byte{fingerprint})
	index2 := index1 ^ (uint(hash2) & (uint(len(cf.buckets)) - 1))

	return index1, index2
}

// Insert adds an item to the filter
func (cf *CuckooFilter) Insert(item []byte) error {
	loadFactor := float64(cf.count) / (float64(len(cf.buckets)) * float64(cf.bucketSize))
	maxLoadFactor := LOAD_FACTOR_MAP[cf.bucketSize]
	if loadFactor >= maxLoadFactor {
		return errors.New("filter is full")
	}

	fingerprint := cf.generateFingerprint(item)
	i1, i2 := cf.getIndices(item, fingerprint)

	// Try to insert into either bucket
	if cf.buckets[i1].size < cf.bucketSize {
		cf.buckets[i1].fingerprints[cf.buckets[i1].size] = fingerprint
		cf.buckets[i1].size++
		cf.count++
		return nil
	}

	if cf.buckets[i2].size < cf.bucketSize {
		cf.buckets[i2].fingerprints[cf.buckets[i2].size] = fingerprint
		cf.buckets[i2].size++
		cf.count++
		return nil
	}

	currentFp := fingerprint
	currentIndex := i1
	if cf.rng.Intn(2) == 0 {
		currentIndex = i2
	}

	for range MaxKicks {
		randPos := uint(cf.rng.Intn(int(cf.bucketSize)))
		oldFp := cf.buckets[currentIndex].fingerprints[randPos]
		cf.buckets[currentIndex].fingerprints[randPos] = currentFp

		currentFp = oldFp

		hash := cf.fingerprintHashFunc.Hash([]byte{currentFp})
		alternateIndex := currentIndex ^ (uint(hash) & (uint(len(cf.buckets)) - 1))

		if cf.buckets[alternateIndex].size < cf.bucketSize {
			cf.buckets[alternateIndex].fingerprints[cf.buckets[alternateIndex].size] = currentFp
			cf.buckets[alternateIndex].size++
			cf.count++
			return nil
		}

		currentIndex = alternateIndex
	}

	return errors.New("maximum number of relocations reached")
}

// Lookup checks if an item might be in the filter
func (cf *CuckooFilter) Lookup(item []byte) bool {
	fingerprint := cf.generateFingerprint(item)
	i1, i2 := cf.getIndices(item, fingerprint)

	for i := uint(0); i < cf.buckets[i1].size; i++ {
		if cf.buckets[i1].fingerprints[i] == fingerprint {
			return true
		}
	}

	for i := uint(0); i < cf.buckets[i2].size; i++ {
		if cf.buckets[i2].fingerprints[i] == fingerprint {
			return true
		}
	}

	return false
}

func (cf *CuckooFilter) Delete(item []byte) bool {
	fingerprint := cf.generateFingerprint(item)
	i1, i2 := cf.getIndices(item, fingerprint)

	for i := uint(0); i < cf.buckets[i1].size; i++ {
		if cf.buckets[i1].fingerprints[i] == fingerprint {
			if i < cf.buckets[i1].size-1 {
				cf.buckets[i1].fingerprints[i] = cf.buckets[i1].fingerprints[cf.buckets[i1].size-1]
			}
			cf.buckets[i1].size--
			cf.count--
			return true
		}
	}

	for i := uint(0); i < cf.buckets[i2].size; i++ {
		if cf.buckets[i2].fingerprints[i] == fingerprint {
			if i < cf.buckets[i2].size-1 {
				cf.buckets[i2].fingerprints[i] = cf.buckets[i2].fingerprints[cf.buckets[i2].size-1]
			}
			cf.buckets[i2].size--
			cf.count--
			return true
		}
	}

	return false
}

func (cf *CuckooFilter) Count() uint {
	return cf.count
}

func (cf *CuckooFilter) LoadFactor() float64 {
	return float64(cf.count) / (float64(len(cf.buckets)) * float64(cf.bucketSize))
}

func (cf *CuckooFilter) Size() uint {
	return uint(len(cf.buckets)) * cf.bucketSize
}
