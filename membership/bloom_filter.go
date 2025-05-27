// Package membership provides probabilistic data structures for membership testing
package membership

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/mrtkp9993/probdsgo/utils"
)

const (
	LN2     = math.Ln2
	LN2SQRD = math.Ln2 * math.Ln2
)

// BloomFilter implements a Bloom filter, a space-efficient probabilistic data structure
// used to test whether an element is a member of a set
type BloomFilter struct {
	bitArray      []bool
	bitCount      uint
	hashFuncCount uint
	hashFunctions []*utils.Murmur3
	seeds         []uint32
}

// NewBloomFilter creates a new Bloom filter with specified capacity and error rate
// capacity: maximum number of elements expected to be stored
// errorRate: desired false positive probability
func NewBloomFilter(capacity uint, errorRate float64) (*BloomFilter, error) {
	if capacity < 1 || errorRate <= 0.0 || errorRate >= 1.0 {
		return nil, errors.New("invalid capacity or error rate")
	}

	numberOfBits := uint(math.Ceil(-1.0 * float64(capacity) * math.Log(errorRate) / LN2SQRD))
	numberOfHashes := uint(math.Ceil(-1.0 * math.Log(errorRate) / LN2))

	return NewBloomFilterWithParams(numberOfBits, numberOfHashes)
}

// NewBloomFilterWithParams creates a new Bloom filter with specified bit array size and number of hash functions
// m: size of bit array
// k: number of hash functions
func NewBloomFilterWithParams(m, k uint) (*BloomFilter, error) {
	if m <= 0 || k <= 0 {
		return nil, errors.New("invalid m or k")
	}

	bitArray := make([]bool, m)
	hashFunctions := make([]*utils.Murmur3, k)
	seeds := make([]uint32, k)
	for i := range k {
		seeds[i] = uint32(i + 1)
		hashFunctions[i] = utils.NewMurmur3WithSeed(seeds[i])
	}

	return &BloomFilter{bitArray: bitArray, bitCount: uint(len(bitArray)), hashFuncCount: uint(len(hashFunctions)), hashFunctions: hashFunctions, seeds: seeds}, nil
}

// Add inserts an item into the Bloom filter
// Returns error if insertion fails
func (bf *BloomFilter) Add(item []byte) error {
	if err := validateInput(item); err != nil {
		return err
	}

	for _, hashFunc := range bf.hashFunctions {
		hash := hashFunc.Hash(item)
		position := hash % uint32(len(bf.bitArray))

		bf.bitArray[position] = true
	}

	return nil
}

// Contains checks if an item might be in the Bloom filter
// Returns true if item might be present, false if definitely not present
func (bf *BloomFilter) Contains(item []byte) (bool, error) {
	if err := validateInput(item); err != nil {
		return false, err
	}

	for _, hashFunc := range bf.hashFunctions {
		hash := hashFunc.Hash(item)
		position := hash % uint32(bf.bitCount)

		if !bf.bitArray[position] {
			return false, nil
		}
	}

	return true, nil
}

// Cardinality returns the estimated number of unique items in the Bloom filter
func (bf *BloomFilter) Cardinality() uint {
	X := uint(0)
	for _, bit := range bf.bitArray {
		if bit {
			X++
		}
	}

	if X < bf.hashFuncCount {
		return 0
	}

	if X == bf.hashFuncCount {
		return 1
	}

	if X == bf.bitCount {
		return bf.bitCount / bf.hashFuncCount
	}

	n := -1.0 * (float64(bf.bitCount) / float64(bf.hashFuncCount)) * math.Log(1-(float64(X)/float64(bf.bitCount)))
	log.Println(n, X)
	return uint(math.Round(n))
}

// FalsePositiveRate returns the current false positive rate of the Bloom filter
func (bf *BloomFilter) FalsePositiveRate() float32 {
	X := 0
	for _, bit := range bf.bitArray {
		if bit {
			X++
		}
	}

	if X == 0 {
		return 0
	}

	m := float64(len(bf.bitArray))
	k := float64(len(bf.hashFunctions))

	probability := math.Pow(float64(X)/m, k)

	return float32(probability)
}

func (bf *BloomFilter) Merge(other *BloomFilter) (*BloomFilter, error) {
	if err := checkCompatibility(bf, other); err != nil {
		return nil, fmt.Errorf("cannot merge: %v", err)
	}

	result := &BloomFilter{
		bitArray:      make([]bool, bf.bitCount),
		bitCount:      bf.bitCount,
		hashFuncCount: bf.hashFuncCount,
		hashFunctions: make([]*utils.Murmur3, bf.hashFuncCount),
		seeds:         make([]uint32, bf.hashFuncCount),
	}

	copy(result.seeds, bf.seeds)
	for i, seed := range result.seeds {
		result.hashFunctions[i] = utils.NewMurmur3WithSeed(seed)
	}

	for i := range bf.bitArray {
		result.bitArray[i] = bf.bitArray[i] || other.bitArray[i]
	}

	return result, nil
}

func (bf *BloomFilter) Intersect(other *BloomFilter) (*BloomFilter, error) {
	if err := checkCompatibility(bf, other); err != nil {
		return nil, fmt.Errorf("cannot intersect: %v", err)
	}

	result := &BloomFilter{
		bitArray:      make([]bool, bf.bitCount),
		bitCount:      bf.bitCount,
		hashFuncCount: bf.hashFuncCount,
		hashFunctions: make([]*utils.Murmur3, bf.hashFuncCount),
		seeds:         make([]uint32, bf.hashFuncCount),
	}

	copy(result.seeds, bf.seeds)
	for i, seed := range result.seeds {
		result.hashFunctions[i] = utils.NewMurmur3WithSeed(seed)
	}

	for i := range bf.bitArray {
		result.bitArray[i] = bf.bitArray[i] && other.bitArray[i]
	}

	return result, nil
}

func validateInput(item []byte) error {
	if item == nil {
		return errors.New("input cannot be nil")
	}

	if len(item) == 0 {
		return errors.New("input cannot be empty")
	}

	return nil
}

func checkCompatibility(bf1, bf2 *BloomFilter) error {
	if bf1.bitCount != bf2.bitCount {
		return errors.New("bloom filters have different bit array sizes")
	}
	if bf1.hashFuncCount != bf2.hashFuncCount {
		return errors.New("bloom filters have different numbers of hash functions")
	}

	// Check if they use the same seeds
	for i, seed1 := range bf1.seeds {
		if seed1 != bf2.seeds[i] {
			return errors.New("bloom filters use different hash function seeds")
		}
	}

	return nil
}
