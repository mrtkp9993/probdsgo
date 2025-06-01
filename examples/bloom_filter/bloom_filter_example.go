package main

import (
	"fmt"

	"github.com/mrtkp9993/probdsgo/membership"
)

func main() {
	bf, err := membership.NewBloomFilter(10000, 0.01)
	if err != nil {
		panic(err)
	}

	bf.Add([]byte("hello"))

	exists, _ := bf.Contains([]byte("hello"))
	fmt.Println("hello exists?:", exists)

	exists, _ = bf.Contains([]byte("world"))
	fmt.Println("world exists?:", exists)

	bf.Add([]byte("world"))

	exists, _ = bf.Contains([]byte("world"))
	fmt.Println("world exists?:", exists)

	fmt.Println("Cardinality:", bf.Cardinality())
	fmt.Println("False Positive Rate:", bf.FalsePositiveRate())

	bf2, err := membership.NewBloomFilter(10000, 0.01)
	if err != nil {
		panic(err)
	}

	bf2.Add([]byte("hello"))
	bf2.Add([]byte("world"))
	bf2.Add([]byte("test"))
	bf2.Add([]byte("test2"))
	bf2.Add([]byte("test3"))

	bf3, err := bf.Merge(bf2)
	if err != nil {
		panic(err)
	}

	fmt.Println("Cardinality:", bf3.Cardinality())
	fmt.Println("False Positive Rate:", bf3.FalsePositiveRate())

	bf4, err := bf.Intersect(bf2)
	if err != nil {
		panic(err)
	}

	fmt.Println("Cardinality:", bf4.Cardinality())
	fmt.Println("False Positive Rate:", bf4.FalsePositiveRate())
}
