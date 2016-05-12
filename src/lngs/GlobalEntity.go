package lngs

import (
	"log"
)

var (
	globalEntities = map[string]*Entity{}
)

func CreateGlobalEntity(behaviorName string) *Entity {
	if entity, ok := globalEntities[behaviorName]; ok {
		// global entity exists
		if entity.GetBehaviorName() == behaviorName {
			log.Printf("Global entity %s already exists\n", behaviorName)
			return entity
		} else {
			log.Panicf("Global entity %s already exists, expected behavior is %s\n", entity.GetBehaviorName(), behaviorName)
		}
	}

	entity, err := CreateEntity(behaviorName, "")
	if err != nil {
		log.Panicf("Create global entity (behavior %s) failed", behaviorName)
	}

	globalEntities[behaviorName] = entity
	return entity
}

func GetGlobalEntity(behaviorName string) *Entity {
	return CreateGlobalEntity(behaviorName)
}
