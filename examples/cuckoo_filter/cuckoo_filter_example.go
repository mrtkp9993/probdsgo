package main

import (
	"fmt"

	"github.com/mrtkp9993/probdsgo/membership"
)

func main() {
	// 8-bit fingerprint example
	cf8, err := membership.NewCuckooFilter(10000, 4, membership.FINGERPRINT_SIZE_8)
	if err != nil {
		panic(err)
	}

	cf8.Insert([]byte("hello"))

	fmt.Println("hello exists?:", cf8.Lookup([]byte("hello")))
	fmt.Println("world exists?:", cf8.Lookup([]byte("world")))

	cf8.Delete([]byte("hello"))
	fmt.Println("hello exists after delete?:", cf8.Lookup([]byte("hello")))

	cf8.Insert([]byte("hello"))
	cf8.Insert([]byte("world"))
	cf8.Insert([]byte("test"))
	cf8.Insert([]byte("test2"))
	cf8.Insert([]byte("test3"))

	fmt.Println("Cardinality:", cf8.Count())
	fmt.Println("Load Factor:", cf8.LoadFactor())
}
