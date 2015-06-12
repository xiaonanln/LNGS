package lngs

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	. "lngs/common"
	"lngs/db"
	. "lngs/rpc"
	"log"
	"reflect"
)

var (
	noneBehavior = reflect.ValueOf(struct{}{}) // entity with no behavior
)

func NewEntityId() string {
	return bson.NewObjectId().Hex()
}

type Entity struct {
	id           string
	client       *GameClient
	behavior     reflect.Value
	commandQueue CommandQueue
}

func (self *Entity) SetClient(client *GameClient) {
	if self.client != nil {
		self.client.DestroyEntity(self.id)

		self.client.owner = nil
		self.client = nil
	}

	self.client = client
	if self.client != nil {
		self.client.owner = self
		self.client.CreateEntity(self.id, self.GetBehaviorName())
		self.client.BecomePlayer(self.id)
	}
}

func NewEntity(id string) *Entity {
	if id == "" {
		id = NewEntityId()
	}
	return &Entity{id, nil, noneBehavior, make(CommandQueue)}
}

func (self *Entity) GetBehaviorName() string {
	behaviorType := reflect.Indirect(self.behavior).Type()
	name := behaviorType.Name()
	log.Printf("behavior name %v %s", behaviorType, name)
	return name
}

func (self *Entity) String() string {
	return fmt.Sprintf("%s<%s>", self.GetBehaviorName(), self.id)
}

func (self *Entity) SetBehavior(behaviorName string) {
	behaviorValue := entityManager.NewEntityBehavior(behaviorName)
	self.behavior = behaviorValue
}

func (self *Entity) OnReceiveMessage(msg Message) {
	ID := msg["ID"].(string)
	M := msg["M"].(string)

	if ID != self.id {
		log.Printf("can not call entity %s from entity %s", ID, self.id)
		return
	}

	targetEntity := entityManager.GetEntity(ID)
	if targetEntity == nil {
		log.Printf("entity %s not found when calling method %s", ID, M)
		return
	}

	ARGS := msg["ARGS"].([]interface{})
	targetEntity.OnCallMethod(self, M, ARGS)
}

func (self *Entity) OnCallMethod(caller *Entity, methodname string, args []interface{}) {
	log.Printf("%s calling %s.%s", caller, self, methodname)

	method := self.behavior.MethodByName(methodname)
	log.Printf("method %v, total methods %d", method, self.behavior.NumMethod())
	in := make([]reflect.Value, len(args)+1)
	in[0] = reflect.ValueOf(self)
	for i, arg := range args {
		in[i+1] = reflect.ValueOf(arg)
	}
	// methodType := method.Type()
	// numArguments := methodType.NumIn()
	// for argIndex := 0; argIndex < numArguments; argIndex++ {
	// 	var argType reflect.Type = methodType.In(argIndex)
	// 	log.Println("arg type", argIndex, argType)
	// 	in[argIndex] = convertType(in[argIndex], argType)
	// }
	method.Call(in)
}

func (self *Entity) InsertDB(collectionName string, doc interface{}) {
	cmd := Command{
		self.id,
		"insertdb",
		[]interface{}{collectionName, doc},
	}

	lngsdb.PostDbCommand(&cmd)
	return
}

func (self *Entity) PostCommand(cmd *Command) {
	self.commandQueue <- cmd
}

// func convertType(val reflect.Value, targetType reflect.Type) reflect.Value {
// 	switch targetType.Kind() {
// 	case reflect.Slice:
// 		elemType := targetType.Elem()
// 		log.Println("element type", elemType)
// 		sliceLen := val.Len()
// 		newSlice := reflect.MakeSlice(targetType, sliceLen, sliceLen)
// 		for i := 0; i < sliceLen; i++ {
// 			elem := val.Index(i)
// 			log.Println("element", i, "is", elem.Interface().(string))
// 			newSlice.Index(i).Set(convertType(elem, elemType))
// 			// newSlice = reflect.Append(newSlice)
// 		}
// 		log.Println("slice len", sliceLen, "new slice", newSlice.Interface().([]string))
// 		return newSlice

// 	case reflect.Bool:
// 		return reflect.ValueOf(val.Interface().(bool))
// 	case reflect.String:
// 		return reflect.ValueOf(val.Interface().(string))
// 	}
// 	return val
// }
