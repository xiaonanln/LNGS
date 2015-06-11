package lngs

import (
	"lngs/common"
)

func RegisterEntityBehavior(entityBehavior interface{}) {
	GetEntityManager().RegisterEntityBehavior(entityBehavior)
}

func SetBootEntityBehavior(entityBehavior interface{}) {
	GetEntityManager().SetBootEntityBehavior(entityBehavior)
}

func SetConfigFile(configPath string) {
	lngscommon.ReadConfigFile(configPath)
}
