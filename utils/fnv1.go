package utils

type FNV1 struct{}

func (f *FNV1) Hash(data []byte) uint64 {
	hash := uint64(0xcbf29ce484222325)
	for _, b := range data {
		hash = hash * 0x100000001b3
		hash = hash ^ uint64(b)
	}
	return hash
}

func NewFNV1() *FNV1 {
	return &FNV1{}
}
