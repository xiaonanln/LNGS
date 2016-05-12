package main

import (
	"log"

	. "lngs"
	// . "lngs/common"
)

type Avatar struct {
	exp int
}

func (behavior *Avatar) Init(self *Entity) {
	// Initialize Avatar Attrs by default values
	self.Attrs.Set("exp", self.Attrs.GetInt("exp", 0))
	self.Attrs.Set("cups", self.Attrs.GetInt("cups", 0))
	self.Attrs.Set("gold", self.Attrs.GetInt("gold", 0))
	self.Attrs.Set("diamond", self.Attrs.GetInt("diamond", 0))

	self.Attrs.GetMapAttr("chests")
	self.Attrs.GetMapAttr("heroes") 
	self.Attrs.GetMapAttr("items")
	self.Attrs.GetMapAttr("embattles")

	self.Attrs.GetMapAttr("testAttr")
}

func (behavior *Avatar) Test(self *Entity) {
	log.Printf("!!!!!!!!!!!!!!!!!!!!!!! Avatar.Test !!!!!!!!!!!!!!!!!!!!!!!")
	testAttr := self.Attrs.GetMapAttr("testAttr")
	testAttr.Set("testAttr", testAttr.GetInt("testVal", 0) + 1)
	self.NotifyAttrChange("testAttr")
}

func (behavior *Avatar) OnGetNewClient(self *Entity) {
	onlineManager := GetGlobalEntity("OnlineManager")
	log.Printf("Entity %s get new client, OnlineManager %s", self, onlineManager)
	onlineManager.CallMethod("NotifyAvatarLogin", self.Id())
}

func (behavior *Avatar) OnLoseClient(self *Entity, old_client *GameClient) {
	onlineManager := GetGlobalEntity("OnlineManager")
	log.Printf("Entity %s lose client, OnlineManager %s", self, onlineManager)
	onlineManager.CallMethod("NotifyAvatarLogout", self.Id())
}

func (behavior *Avatar) Say(self *Entity, text string) {
	// player say something in the tribe
	for _, entity := range Entities() {
		if entity.GetBehaviorName() == "Avatar" {
			entity.CallClient("OnSay", self.Id(), text)
		}
	}
}

func (behavior *Avatar) GetSaveInterval() int {
	return 10
}

