package main

import (
	"lngs"
	"tests/entities"
)

func main() {
	lngs.SetConfigFile("testserver.conf")
	lngs.RegisterEntityBehavior(entities.Boot{})
	lngs.RegisterEntityBehavior(entities.Avatar{})
	lngs.RegisterEntityBehavior(entities.OnlineManager{})

	lngs.SetBootEntityBehavior(entities.Boot{})

	lngs.CreateGlobalEntity("OnlineManager", "OnlineManager")

	lngs.Run("0.0.0.0:7777")

}
