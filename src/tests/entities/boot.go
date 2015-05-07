package entities

import (
	"log"
)

type Boot struct {
}

func (self *Boot) Test(a, b, c int) {
	log.Printf("Boot.Text called")
}
