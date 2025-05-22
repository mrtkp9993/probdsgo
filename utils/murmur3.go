// This code is a direct translation of the murmur3_32 function from C to Go.
// The C code is taken from Wikipedia and Github
// https://en.wikipedia.org/wiki/MurmurHash
// https://github.com/jwerle/murmurhash.c/blob/master/murmurhash.c
package utils

func murmur32Scramble(k uint32) uint32 {
	k *= 0xcc9e2d51
	k = (k << 15) | (k >> (17))
	k *= 0x1b873593
	return k
}

func Murmur3_32(key []byte, seed uint32) uint32 {
	h := seed
	var k uint32

	nblocks := len(key) / 4
	for i := 0; i < nblocks; i++ {
		k = uint32(key[i*4]) |
			uint32(key[i*4+1])<<8 |
			uint32(key[i*4+2])<<16 |
			uint32(key[i*4+3])<<24
		h ^= murmur32Scramble(k)
		h = (h << 13) | (h >> 19)
		h = h*5 + 0xe6546b64
	}

	k = 0
	tail := key[nblocks*4:]
	switch len(tail) {
	case 3:
		k ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k ^= uint32(tail[0])
		k = murmur32Scramble(k)
		h ^= k
	}

	h ^= uint32(len(key))
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}
