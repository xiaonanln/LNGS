package lngs

import (
	"fmt"
	. "lngs/cmdque"
	. "lngs/common"

	"gopkg.in/mgo.v2/bson"

	"errors"

	. "lngs/db"
	. "lngs/rpc"
	"log"
	"reflect"
	"sync"
	"lngs/typeconv"
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
	GetSaveInterval() int
}

type Entity struct {
	id           string
	client       *GameClient
	behavior     reflect.Value
	commandQueue CommandQueue
	lock         sync.Mutex 
	Attrs		 MapAttr 
} 

func (self *Entity) SetClient(client *GameClient) {
	if client == self.client {
		return // client not changed at all
	}

	old_client := self.client

	if self.client != nil {
		self.client.DestroyEntity(self.id)

		self.client.owner = nil
		self.client = nil
	}

	self.client = client
	if self.client != nil {
		self.client.owner = self
		self.client.CreateEntity(self.id, self.GetBehaviorName())
		self.client.BecomePlayer(self.id, self.GetPersistentData())
	}

	if old_client != nil && self.client == nil {
		self.onLoseClient(old_client)
	} else if old_client == nil && self.client != nil {
		self.onGetNewClient()
	} else {
		self.onChangeClient(old_client)
	}
}

func (self *Entity) Lock() {
	self.lock.Lock()
}

func (self *Entity) Unlock() {
	self.lock.Unlock()
}

func newEntity(id string) *Entity {
	if id == "" {
		id = NewEntityId()
	}
	return &Entity{
		id:           id,
		behavior:     noneBehavior,
		commandQueue: GetCommandQueue(id),
		Attrs : *NewMapAttr(), 
	}
}

func (self *Entity) Id() string {
	return self.id
}

func (self *Entity) GetBehaviorName() string {
	behaviorType := reflect.Indirect(self.behavior).Type()
	name := behaviorType.Name()
	return name
}

func (self *Entity) String() string {
	return fmt.Sprintf("%s<%s>", self.GetBehaviorName(), self.id)
}

func (self *Entity) setBehavior(behaviorName string) {
	behaviorValue := entityManager.newEntityBehavior(self, behaviorName)
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

	ARGS, ok := msg["ARGS"].([]interface{})
	if ok {
		targetEntity.Call(M, ARGS...)
	} else {
		targetEntity.Call(M)
	}
}

func (self *Entity) Call(methodname string, args ...interface{}) {
	log.Printf("Entity method: %s.%s", self, methodname)
	self.callBehaviorMethod(methodname, args...)
}

func (self *Entity) CallClient(method string, args ...interface{}) {
	Debug(self.id, "call client method %s %v", method, args)
	if self.client != nil {

		self.client.Call(self.id, method, args...)
	}
}

func (self *Entity) PostCommand(cmd *Command) {
	self.commandQueue <- cmd
}

func CreateEntity(behaviorName string, id string) (*Entity, error) {
	newEntity := entityManager.newEntity(behaviorName, id)

	if newEntity.IsPersistent() {
		var doc Doc
		var err error 

		if id != "" {
			doc, err = FindDoc("entities", map[string]interface{}{"_id": EntityId2DocId(id)})
			if err != nil {
				Debug("Entity", "Create persistent entity failed: entity not found in entities collection: %s", id)
				return nil, err
			}

			if doc != nil && doc["_behavior"] != behaviorName {
				// entity behavior is wrong
				return nil, errors.New("wrong behavior")
			}
		}
		
		entityManager.putEntity(newEntity) // after get data from DB successfully, put entity to entity manager

		if doc != nil {
			// no persistent data, just create entity
			newEntity.initWithPersistentData(doc)
		} else {
			newEntity.initWithPersistentData(Doc{})
		}

	} else {
		entityManager.putEntity(newEntity)
	}
	

	behavior := newEntity.behavior
	initMethod := behavior.MethodByName("Init")
	log.Println(behaviorName, "Init", initMethod)
	initMethod.Call([]reflect.Value{reflect.ValueOf(newEntity)})

	return newEntity, nil
}

func (self *Entity) CreateEntity(behaviorName string, id string) (*Entity, error) {
	return CreateEntity(behaviorName, id)
}

func (self *Entity) IsPersistent() bool {
	_, ok := self.behavior.Interface().(Persistence)
	return ok
}

func (self *Entity) getPersistence() Persistence {
	p, ok := self.behavior.Interface().(Persistence)
	if ok {
		return p
	} else {
		return nil
	}
}

func (self *Entity) Save() error {
	p := self.getPersistence()
	if p == nil {
		return nil
	}

	data := self.GetPersistentData()
	entityid := EntityId2DocId(self.id)
	data["_behavior"] = self.GetBehaviorName()

	query := map[string]interface{}{"_id": entityid}
	err := UpdateDoc("entities", query, Doc{"$set": data})
	for err != nil {
		log.Println("Save entity %s failed: %v, retry...", self, err)
		err = UpdateDoc("entities", query, Doc{"$set": data})
	}

	Debug("Entity", "Entity %s saved successfuly", self)
	return err
}

func (self *Entity) Destroy() error {
	err := self.Save()
	if err != nil {
		return err 
	}

	entityManager.delEntity(self)
	return nil 
}

func (self *Entity) GetPersistentData()  Doc {
	return self.Attrs.ToDoc()
}

func (self *Entity) initWithPersistentData(data Doc) {
	self.Attrs.AssignDoc(data)
	Debug("Entity", "Entity %s init with data %v, Attrs = %v", self, data, self.Attrs)
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

func (self *Entity) callBehaviorMethod(methodname string, args ...interface{}) bool {
	method := self.behavior.MethodByName(methodname)
	if !method.IsValid() {
		log.Printf("Entity %s: method %s not found\n", self, methodname)
		return false
	}

	methodType := method.Type()

	in := make([]reflect.Value, len(args)+1)
	in[0] = reflect.ValueOf(self)

	for i, arg := range args {
		argType := methodType.In(i+1)
		argVal := reflect.ValueOf(arg)
		in[i+1] = typeconv.Convert(argVal, argType)
	}
	method.Call(in)
	return true
}

func (self *Entity) onGetNewClient() {
	// get new client: self.client
	self.callBehaviorMethod("OnGetNewClient")
}

func (self *Entity) onLoseClient(old_client *GameClient) {
	// lose client: old_client
	self.callBehaviorMethod("OnLoseClient", old_client)
}

func (self *Entity) onChangeClient(old_client *GameClient) {
	// client changed from old_client to self.client
	self.callBehaviorMethod("OnChangeClient", old_client)
}

func (self *Entity) NotifyAttrChange(attrName string) {
	attrVal := self.Attrs.attrs[attrName]
	mapAttrVal, ok := attrVal.(*MapAttr)
	if ok {
		attrVal = mapAttrVal.ToDoc()
	} else {
		attrVal = attrVal
	}
	self.CallClient("OnAttrChange", attrName, attrVal)
}

func (self *Entity) Set(key string, val interface{}) {
	self.Attrs.Set(key, val)
}


func (self *Entity) GetInt(key string, defaultVal int) int {
	return self.Attrs.GetInt(key, defaultVal)
}

func (self *Entity) GetStr(key string, defaultVal string) string {
	return self.Attrs.GetStr(key, defaultVal)
}

func (self *Entity) GetMapAttr(key string) *MapAttr {
	return self.Attrs.GetMapAttr(key)
}

func (self *Entity) GetFloat(key string, defaultVal float64) float64 {
	return self.Attrs.GetFloat(key, defaultVal)
}

func (self *Entity) GetBool(key string, defaultVal bool) bool {
	return self.Attrs.GetBool(key, defaultVal)
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
