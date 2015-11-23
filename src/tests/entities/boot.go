package entities

import (
	"log"

	. "lngs"
	. "lngs/common"
	"lngs/db"
)

var ()

type Boot struct {
}

func (behavior *Boot) Init() {
}

func (behavior *Boot) Test(self *Entity, a, b, c int) {
	log.Printf("Boot.Text called")
}

func (behavior *Boot) PlayGame(self *Entity, testInt int, testStr string, testMap map[string]interface{}, testList []interface{}) {
	log.Printf("Boot.PlayGame %v, %v, %v, %v", testInt, testStr, testMap, testList)
}

func (behavior *Boot) Login(self *Entity, username string, password string) {
	log.Println("Login", username, password)

	doc, _ := lngsdb.FindDoc("entities", map[string]string{"username": username})
	if doc == nil {
		self.CallClient("OnLogin", "player_not_found", username)
		return
	}

	if doc["password"] != password {
		Debug("boot", "wrong password %s, correct is %s", password, doc["password"])
		self.CallClient("OnLogin", "wrong_password", username)
		return
	}

	// login success, create avatar now
	Debug("Boot", "found avatar %v", doc)
	entityid := doc.HexId()
	avatar, err := self.CreateEntity("Avatar", entityid)
	Debug("boot", "create entity %v, error %v", avatar, err)
	if err != nil {
		self.CallClient("OnLogin", "fail", username)
		return
	}

	self.CallClient("OnLogin", "success", username)
	Debug("boot", "%s login success", username)

	self.GiveClientTo(avatar)
}

func (behavior *Boot) Register(self *Entity, username string, password string) {
	log.Println("Register", username, password)
	// find the player before register
	doc, err := lngsdb.FindDoc("entities", map[string]string{"username": username})
	Debug("boot", "find player by username: %v, error=%v", doc, err)

	if doc != nil {
		self.CallClient("OnRegister", "fail")
		return
	}

	err = lngsdb.InsertDoc("entities", map[string]interface{}{"username": username, "password": password, "_behavior": "Avatar"})
	if err != nil {
		self.CallClient("OnRegister", "fail")
		return
	}

	self.CallClient("OnRegister", "success")
}
