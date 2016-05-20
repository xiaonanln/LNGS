package lngsutils

import "math/rand"

func RandInt(a int, b int) int {
	return a + rand.Intn(b-a+1)
}

func ChooseInt(values []int) int {
	i := rand.Intn(len(values))
	return values[i]
}
