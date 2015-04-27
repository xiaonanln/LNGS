package main

import (
	"fmt"
)

func main() {
	m := make(map[int]int)
	m[1] = 1
	fmt.Printf("%v\n", m)
	x := m
	x[2] = 2
	fmt.Printf("%v %v\n", m, x)

	m = nil
	m[1] = 1
}
