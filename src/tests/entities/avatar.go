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
