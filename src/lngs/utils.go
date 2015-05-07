package lngs

import (
	"lngs/common"
)

func RegisterEntityBehavior(entityBehavior interface{}) {
	lngscommon.GetEntityManager().RegisterEntityBehavior(entityBehavior)
}

func SetBootEntityBehavior(entityBehavior interface{}) {
	lngscommon.GetEntityManager().SetBootEntityBehavior(entityBehavior)
}
