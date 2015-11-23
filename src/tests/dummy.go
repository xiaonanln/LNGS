package main

import (
	"fmt"
	"reflect"
)

type A struct {
	A int
	B int
	C map[int]int
	d *int
}

func main() {
	t := reflect.TypeOf(A{})
	fmt.Println(t)

}
