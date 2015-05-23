package entities

import (
	"log"
)

type Boot struct {
}

func (self *Boot) Test(a, b, c int) {
	log.Printf("Boot.Text called")
}

func (self *Boot) PlayGame(testInt int, testStr string, testMap map[string]interface{}, testList []interface{}) {
	log.Printf("Boot.PlayGame %v, %v, %v, %v", testInt, testStr, testMap, testList)
}

func (self *Boot) Login(username string, password string) {
	log.Println("Login", username, password)
	if username != "test" || password != "1234556" {
		return false
	}
	// 根据username找到对应的

}
