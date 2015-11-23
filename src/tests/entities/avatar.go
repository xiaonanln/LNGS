package entities

import (
	"log"

	. "lngs"
	. "lngs/db"
	// . "lngs/common"
)

type Avatar struct {
	exp int
}

func (behavior *Avatar) Init() {
}

func (behavior *Avatar) Test(self *Entity, a, b, c int) {
	log.Printf("Boot.Text called")
}

func (self *Avatar) GetPersistentData() Doc {
	return Doc{
		"exp": self.exp,
	}
}

func (self *Avatar) InitWithPersistentData(data Doc) {
	self.exp = data.Get("exp", 0).(int)
}

func (behavior *Avatar) AddExp(self *Entity, exp int) {
	log.Printf("Avatar.AddExp %v -> %v", exp, behavior.exp+exp)
	behavior.exp += exp
	self.Save()
}

func (behavior *Avatar) OnGetNewClient(self *Entity) {
	onlineManager := GetGlobalEntity("OnlineManager", "OnlineManager")
	log.Printf("Entity %s get new client, OnlineManager %s", self, onlineManager)
	onlineManager.CallMethod("NotifyAvatarLogin", self.Id())
}

func (behavior *Avatar) OnLoseClient(self *Entity, old_client *GameClient) {
	onlineManager := GetGlobalEntity("OnlineManager", "OnlineManager")
	log.Printf("Entity %s lose client, OnlineManager %s", self, onlineManager)
	onlineManager.CallMethod("NotifyAvatarLogout", self.Id())
}
