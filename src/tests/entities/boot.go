package entities

import (
	"log"

	. "lngs"
	. "lngs/common"
)

var ()

type Boot struct {
}

func (behavior *Boot) Test(self *Entity, a, b, c int) {
	log.Printf("Boot.Text called")
}

func (behavior *Boot) PlayGame(self *Entity, testInt int, testStr string, testMap map[string]interface{}, testList []interface{}) {
	log.Printf("Boot.PlayGame %v, %v, %v, %v", testInt, testStr, testMap, testList)
}

func (behavior *Boot) Login(self *Entity, username string, password string) {
	log.Println("Login", username, password)
	if username != "test" || password != "1234556" {
		log.Println("wrong username or password")
	}
	// 根据username找到对应的
}

func (behavior *Boot) Register(self *Entity, username string, password string) {
	log.Println("Register", username, password)
	// find the player before register
	doc, err := self.FindDoc("entities", map[string]string{"username": username})
	Debug("boot", "find player by username: %v, error=%v", doc, err)

	if doc != nil {
		self.CallClient("OnRegister", "fail")
		return
	}

	err = self.InsertDoc("entities", map[string]interface{}{"username": username, "password": password})
	if err != nil {
		self.CallClient("OnRegister", "fail")
		return
	}

	self.CallClient("OnRegister", "success")
}
