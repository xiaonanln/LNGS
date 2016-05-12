package main

import (
	// "log"

	. "lngs"
	// . "lngs/common"
)

type Tribe struct {
}

func (behavior *Tribe) Init(self *Entity) {
}

func (behavior *Tribe) GetSaveInterval() int {
	return 10
}
