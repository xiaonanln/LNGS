package main

import (
	"lngs"
)

func main() {
	// lngs.SetConfigFile("testserver.conf")
	lngs.RegisterEntityBehavior(Boot{})
	lngs.RegisterEntityBehavior(Avatar{})
	lngs.RegisterEntityBehavior(OnlineManager{})

	lngs.SetBootEntityBehavior(Boot{})

	lngs.CreateGlobalEntity("OnlineManager", "OnlineManager")

	lngs.Run("0.0.0.0:7777")

}
