package main

import (
	"lngs"
)

func main() {
	// lngs.SetConfigFile("testserver.conf")
	lngs.RegisterEntityBehavior(Boot{})
	lngs.RegisterEntityBehavior(Avatar{})
	lngs.RegisterEntityBehavior(OnlineManager{})
	lngs.RegisterEntityBehavior(Tribe{})
	lngs.RegisterEntityBehavior(TribeManager{})

	lngs.SetBootEntityBehavior(Boot{})

	lngs.CreateGlobalEntity("OnlineManager")
	lngs.CreateGlobalEntity("TribeManager")

	lngs.Run("0.0.0.0:7777")

}
