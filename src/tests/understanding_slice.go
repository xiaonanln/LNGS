package main

import (
	"fmt"
)

func main() {
	a := []byte{1, 2, 3}
	debugSlice(a)

	debugSlice(a[0:1])
	debugSlice(a[1:])

	b := a[0:1]
	b[0] = 100
	debugSlice(b)
	debugSlice(a)
	b = append(b, 9)
	debugSlice(b)
	debugSlice(a)

	b = append(b, 9)
	debugSlice(b)
	debugSlice(a)

	b = append(b, 9)
	debugSlice(b)
	debugSlice(a)

	a = append(a, 10)
	debugSlice(b)
	debugSlice(a)
}

func debugSlice(slice []byte) {
	fmt.Printf("%v len %d cap %d\n", slice, len(slice), cap(slice))
}
