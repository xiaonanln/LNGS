package main

import (
	"lngs"
	"tests/entities"
)

func main() {
	lngs.SetConfigFile("testserver.conf")
	lngs.RegisterEntityBehavior(entities.Boot{})
	lngs.SetBootEntityBehavior(entities.Boot{})
	lngs.Run("0.0.0.0:7000")
}
