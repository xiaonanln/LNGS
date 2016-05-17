package lngs

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
	persistentEntityBehaviors  map[string] bool
	bootentityBehaviorTypeName string
}

func (self *EntityManager) init() {
	self.entityBehaviorTypes = map[string]reflect.Type{}
	self.persistentEntityBehaviors = map[string]bool {}
	self.entities = map[string]*Entity{}
}

func (self *EntityManager) RegisterEntityBehavior(entityBehavior interface{}) {
	entityBehaviorType := reflect.TypeOf(entityBehavior)
	behaviorName := entityBehaviorType.Name()

	_, registered := self.entityBehaviorTypes[behaviorName] 
	if registered {
		return 
	}

	self.entityBehaviorTypes[behaviorName] = entityBehaviorType
	log.Printf("Registering entity type %s => %v", behaviorName, entityBehaviorType)
}
func (self *EntityManager) newEntityBehavior(entity *Entity, behaviorName string) reflect.Value {
	var behaviorType reflect.Type = self.entityBehaviorTypes[behaviorName]
	if behaviorType != nil {
		behavior := reflect.New(behaviorType)
		return behavior
	} else {
		log.Panicf("newEntityBehavior: entity behavior not registered: %s\n", behaviorName)
		return noneBehavior
	}
}

func (self *EntityManager) SetBootEntityBehavior(entityBehavior interface{}) {
	self.RegisterEntityBehavior(entityBehavior)

	entityBehaviorType := reflect.TypeOf(entityBehavior)
	behaviorName := entityBehaviorType.Name()
	self.bootentityBehaviorTypeName = behaviorName
}

func (self *EntityManager) newEntity(behaviorName string, id string) *Entity {
	entityBehaviorType := self.entityBehaviorTypes[behaviorName]
	if entityBehaviorType == nil {
		log.Panicf("unknown behavior name: %s", behaviorName)
		return nil
	}

	// entityBehaviorValue := reflect.New(entityBehaviorType)
	var entity *Entity = newEntity(id)
	entity.setBehavior(behaviorName)
	return entity
}

func (self *EntityManager) putEntity(entity *Entity) {
	existingEntity, exists := self.entities[entity.id]
	if exists {
		log.Panicf("Entity already exists: %v, duplicate entity: %v", existingEntity, entity)
		return 
	}

	self.entities[entity.id] = entity
}

func (self *EntityManager) delEntity(entity *Entity) {
	delete(self.entities, entity.id)
}

func (self *EntityManager) NewBootEntity() *Entity {
	if self.bootentityBehaviorTypeName == "" {
		log.Panicf("boot entity name is not set")
		return nil
	}
	boot, _ := CreateEntity(self.bootentityBehaviorTypeName, "")
	return boot 
}

func GetEntityManager() *EntityManager {
	return entityManager
}

func (self *EntityManager) GetEntity(id string) *Entity {
	return self.entities[id]
}

func GetEntity(id string) *Entity {
	return entityManager.GetEntity(id)
}

func Entities() map[string]*Entity {
	return entityManager.entities
}
