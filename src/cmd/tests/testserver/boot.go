package main

import (
	. "lngs"
	. "lngs/common"
	"lngs/db"
	"log"
)

var ()

type Boot struct {
}

func (behavior *Boot) Init(self *Entity) {
}

func (behavior *Boot) Test(self *Entity, a, b, c int) {
	log.Printf("Boot.Test called")
}

func (behavior *Boot) PlayGame(self *Entity, testInt int, testStr string, testMap map[string]interface{}, testList []interface{}) {
	log.Printf("Boot.PlayGame %v, %v, %v, %v", testInt, testStr, testMap, testList)
}

// Login : login request
func (behavior *Boot) Login(self *Entity, accountId string) {
	var err error
	log.Println("Login", accountId)

	doc, _ := lngsdb.FindDoc("entities", map[string]string{"account": accountId})
	if doc == nil {
		err = lngsdb.InsertDoc("entities", map[string]interface{}{"account": accountId, "_behavior": "Avatar"})
		if err != nil {
			self.CallClient("Login", "fail")
			return
		}
		doc, _ := lngsdb.FindDoc("entities", map[string]string{"account": accountId})
		if doc == nil {
			self.CallClient("Login", "fail")
			return
		}
	}

	// login success, create avatar now
	Debug("Boot", "found avatar %v", doc)
	entityid := doc.HexId()

	avatar := GetEntity(entityid)
	if avatar == nil {
		// avatar already exists
		avatar, err = self.CreateEntity("Avatar", entityid)
		Debug("boot", "create entity %v, error %v", avatar, err)
		if err != nil {
			self.CallClient("OnLogin", "fail")
			return
		}
	}

	self.CallClient("OnLogin", "success")
	Debug("boot", "%s login success", accountId)

	self.GiveClientTo(avatar)
}

func (behavior *Boot) OnLoseClient(self *Entity, old_client *GameClient) {
	self.Destroy()
}
