package main

import (
	"math/rand"
)

func RandInt(a int, b int)  int {
	return a + rand.Intn(b-a + 1)
}