package main

import (
	"lngs"
	"os"
	"log"
	"path/filepath"
)

func main() {
	// lngs.SetConfigFile("testserver.conf")
	exePath, _ := filepath.Abs(os.Args[0])
	binPath := filepath.Dir(exePath)
	dataPath := filepath.Join(binPath, "dota2_data")
	log.Printf("Data path: %s", dataPath)
	
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
