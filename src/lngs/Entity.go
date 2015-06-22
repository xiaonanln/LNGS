package lngs

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	. "lngs/cmdque"
	. "lngs/common"

	"errors"

	. "lngs/db"
	. "lngs/rpc"
	"log"
	"reflect"
)

var (
	noneBehavior = reflect.ValueOf(struct{}{}) // entity with no behavior
)

func DocId2EntityId(id bson.ObjectId) string {
	return id.Hex()
}

func EntityId2DocId(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

func NewEntityId() string {
	return DocId2EntityId(bson.NewObjectId())
}

type Persistence interface {
	GetPersistentData() Doc
	InitWithPersistentData(data Doc)
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

func newEntity(id string) *Entity {
	if id == "" {
		id = NewEntityId()
	}
	return &Entity{id, nil, noneBehavior, GetCommandQueue(id)}
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

func (self *Entity) CallClient(method string, args ...interface{}) {
	Debug(self.id, "call client method %s %v", method, args)
	if self.client != nil {

		self.client.CallMethod(self.id, method, args...)
	}
}

func (self *Entity) FindDoc(collectionName string, query interface{}) (Doc, error) {
	cmd := Command{
		self.id,
		"find",
		[]interface{}{collectionName, query},
	}
	PostDbCommand(&cmd)

	for {
		cmd := <-self.commandQueue
		if cmd.Command == "find_cb" {
			// this is it
			doc := cmd.Data
			err, ok := doc.(error)
			if ok {
				// error
				return nil, err
			} else {
				return doc.(Doc), nil
			}
		} else {
			// this is wrong command
			Debug("Entity %s ignore command %v", self.id, cmd)
			continue
		}
	}
}

func (self *Entity) InsertDoc(collectionName string, doc Doc) error {
	cmd := Command{
		self.id,
		"insert",
		[]interface{}{collectionName, doc},
	}

	PostDbCommand(&cmd)

	for {
		cmd := <-self.commandQueue
		if cmd.Command == "insert_cb" {
			// this is it
			err, ok := cmd.Data.(error)
			if ok {
				// error
				return err
			} else {
				return nil
			}
		} else {
			// this is wrong command
			Debug("Entity %s ignore command %v", self.id, cmd)
			continue
		}
	}
}

func (self *Entity) UpdateDoc(collectionName string, query interface{}, doc Doc) error {
	cmd := Command{
		self.id,
		"update",
		[]interface{}{collectionName, query, doc},
	}

	PostDbCommand(&cmd)

	for {
		cmd := <-self.commandQueue
		if cmd.Command == "update_cb" {
			// this is it
			err, ok := cmd.Data.(error)
			if ok {
				// error
				return err
			} else {
				return nil
			}
		} else {
			// this is wrong command
			Debug("Entity %s ignore command %v", self.id, cmd)
			continue
		}
	}
}

func (self *Entity) PostCommand(cmd *Command) {
	self.commandQueue <- cmd
}

func (self *Entity) CreateEntity(behaviorName string, id string) (*Entity, error) {
	newEntity := entityManager.NewEntity(behaviorName, id)
	if !newEntity.IsPersistent() {
		Debug("Entity", "Non-persistent entity %s created successfully.", newEntity)
		return newEntity, nil
	}

	doc, err := self.FindDoc("entities", map[string]interface{}{"_id": EntityId2DocId(id)})
	if err != nil {
		Debug("Entity", "Create persistent entity failed: entity not found in entities collection: %s", id)
		return nil, err
	}

	if doc != nil && doc["_behavior"] != behaviorName {
		// entity behavior is wrong
		return nil, errors.New("wrong behavior")
	}

	if doc != nil {
		// no persistent data, just create entity
		newEntity.InitWithPersistentData(doc)
	} else {
		newEntity.InitWithPersistentData(Doc{})
	}
	return newEntity, nil
}

func (self *Entity) IsPersistent() bool {
	_, ok := self.behavior.Interface().(Persistence)
	return ok
}

func (self *Entity) GetPersistence() Persistence {
	p, ok := self.behavior.Interface().(Persistence)
	if ok {
		return p
	} else {
		return nil
	}
}

func (self *Entity) Save() error {
	p := self.GetPersistence()
	if p == nil {
		return nil
	}

	data := p.GetPersistentData()
	entityid := EntityId2DocId(self.id)
	data["_behavior"] = self.GetBehaviorName()

	query := map[string]interface{}{"_id": entityid}
	err := self.UpdateDoc("entities", query, Doc{"$set": data})
	Debug("Entity", "Entity %s saved, error = %v", self, err)
	return err
}

func (self *Entity) InitWithPersistentData(data Doc) {
	p := self.GetPersistence()
	if p == nil {
		return
	}

	Debug("Entity", "Entity %s init with data %v", self, data)
	p.InitWithPersistentData(data)
}

func (self *Entity) GiveClientTo(other *Entity) {
	if self == other || self.client == nil {
		return
	}

	client := self.client

	self.SetClient(nil)
	other.SetClient(client)
}

func (self *Entity) SendCommand(targetid string, cmdName string, cmdData interface{}) {
	cmd := Command{
		self.id,
		cmdName,
		cmdData,
	}
	PostCommandQueue(targetid, &cmd)
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
