package entities

import (
	"log"

	. "lngs"
	// . "lngs/common"
)

type OnlineManager struct {
	onlineAvatars map[string]bool
}

func (behavior *OnlineManager) Init() {
	behavior.onlineAvatars = make(map[string]bool)
}

func (behavior *OnlineManager) Test(self *Entity, args ...interface{}) {
	log.Printf("OnlineManager Test, args = %v\n", args)
}

//func (self *Avatar) GetPersistentData() Doc {
//	return Doc{
//		"exp": self.exp,
//	}
//}

//func (self *Avatar) InitWithPersistentData(data Doc) {
//	self.exp = data.Get("exp", 0).(int)
//}

func (behavior *OnlineManager) NotifyAvatarLogin(self *Entity, entityid string) {
	behavior.onlineAvatars[entityid] = true
	log.Printf("OnlineManager.NotifyAvatarLogin entityid = %s, total online %d", entityid, len(behavior.onlineAvatars))
}

func (behavior *OnlineManager) NotifyAvatarLogout(self *Entity, entityid string) {
	delete(behavior.onlineAvatars, entityid)

	log.Printf("OnlineManager.NotifyAvatarLogout entityid = %s, total online %d", entityid, len(behavior.onlineAvatars))
}
