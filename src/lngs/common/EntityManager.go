package lngscommon

import (
	"log"
	"reflect"
)

var (
	entityManager *EntityManager = new(EntityManager)
)

func init() {
	entityManager.init()
}

type EntityManager struct {
	entityBehaviorTypes        map[string]reflect.Type
	entities                   map[string]*Entity
	bootentityBehaviorTypeName string
}

func (self *EntityManager) init() {
	self.entityBehaviorTypes = map[string]reflect.Type{}
	self.entities = map[string]*Entity{}
}

func (self *EntityManager) RegisterEntityBehavior(entityBehavior interface{}) {
	entityBehaviorType := reflect.TypeOf(entityBehavior)
	behaviorName := entityBehaviorType.Name()

	log.Printf("Registering entity type %s => %v", behaviorName, entityBehaviorType)
	self.entityBehaviorTypes[behaviorName] = entityBehaviorType
}
func (self *EntityManager) NewEntityBehavior(behaviorName string) reflect.Value {
	var behaviorType reflect.Type = self.entityBehaviorTypes[behaviorName]
	if behaviorType != nil {
		return reflect.New(behaviorType)
	} else {
		log.Panicf("NewEntityBehavior: entity behavior not registered: %s\n", behaviorName)
		return noneBehavior
	}
}

func (self *EntityManager) SetBootEntityBehavior(entityBehavior interface{}) {
	self.RegisterEntityBehavior(entityBehavior)

	entityBehaviorType := reflect.TypeOf(entityBehavior)
	behaviorName := entityBehaviorType.Name()
	self.bootentityBehaviorTypeName = behaviorName
}

func (self *EntityManager) NewEntity(behaviorName string) *Entity {
	entityBehaviorType := self.entityBehaviorTypes[behaviorName]
	if entityBehaviorType == nil {
		log.Panicf("unknown behavior name: %s", behaviorName)
		return nil
	}

	// entityBehaviorValue := reflect.New(entityBehaviorType)
	var entity *Entity = NewEntity("")
	entity.SetBehavior(behaviorName)
	self.entities[behaviorName] = entity
	return entity
}

func (self *EntityManager) NewBootEntity() *Entity {
	if self.bootentityBehaviorTypeName == "" {
		log.Panicf("boot entity name is not set")
		return nil
	}
	return self.NewEntity(self.bootentityBehaviorTypeName)
}

func GetEntityManager() *EntityManager {
	return entityManager
}

func (self *EntityManager) GetEntity(id string) *Entity {
	for _, entity := range self.entities {
		return entity
	}
	return nil
}
