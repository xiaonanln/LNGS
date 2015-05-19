package entities

import (
	"log"
)

type Boot struct {
}

func (self *Boot) Test(a, b, c int) {
	log.Printf("Boot.Text called")
}

func (self *Boot) PlayGame(testInt int, testStr string, testMap map[string]interface{}, testList []string) {
	log.Printf("Boot.PlayGame %v, %v, %v, %v", testInt, testStr, testMap, testList)
}
