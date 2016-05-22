package lngsutils

import (
	"fmt"
	"log"
	"math/rand"
)

func RandInt(a int, b int) int {
	return a + rand.Intn(b-a+1)
}

func ChooseInt(values []int) int {
	i := rand.Intn(len(values))
	return values[i]
}

func Debug(category string, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	log.Printf("= %s = %s\n", category, msg)
}
