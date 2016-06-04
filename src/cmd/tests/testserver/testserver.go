package main

import (
	"lngs"
	"lngs/data"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// lngs.SetConfigFile("testserver.conf")
	exePath, _ := filepath.Abs(os.Args[0])
	binPath := filepath.Dir(exePath)
	dataPath := filepath.Join(binPath, "dota2_data")
	log.Printf("Data path: %s", dataPath)
	lngsdata.SetDataPath(dataPath)

	// log.Printf("buff[1] = %v\n", *lngsdata.GetDataRecord("buff", 1))
	// log.Printf("skill[1] = %v\n", *lngsdata.GetDataRecord("skill", 1))
	// log.Printf("hero[1] = %v\n", *lngsdata.GetDataRecord("hero", 1))
	// log.Printf("super[1] = %v\n", *lngsdata.GetDataRecord("super", 1))
	// log.Printf("chest[1] = %v\n", *lngsdata.GetDataRecord("chest", 1))

	cardUtilsInit()
	// lngsdata.Reload()

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
