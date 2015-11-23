package lngs

import (
	"log"
)

var (
	globalEntities = map[string]*Entity{}
)

func CreateGlobalEntity(behaviorName string, entityName string) *Entity {
	if entity, ok := globalEntities[entityName]; ok {
		// global entity exists
		if entity.GetBehaviorName() == behaviorName {
			log.Printf("Global entity %s (behavior %s) already exists\n")
			return entity
		} else {
			log.Panicf("Global entity %s (behavior %s) already exists, expected behavior is %s\n", entityName, entity.GetBehaviorName(), behaviorName)
		}
	}

	entity := entityManager.newEntity(behaviorName, "")
	globalEntities[entityName] = entity
	return entity
}

func GetGlobalEntity(behaviorName string, entityName string) *Entity {
	return CreateGlobalEntity(behaviorName, entityName)
}
