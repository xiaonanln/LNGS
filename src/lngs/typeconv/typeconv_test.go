package typeconv

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

func TestJsonType(*testing.T) {
	var j []interface{}
	json.Unmarshal([]byte("[1, 2, 3]"), &j)
	log.Printf("Unmarshal: %v %T", j, j)
	val := reflect.ValueOf(j)
	elem := val.Index(0)
	log.Printf("first element: %v, type %v", elem, elem.Type())
	intValue, ok := elem.Interface().(float64)
	log.Println(intValue, ok, elem.InterfaceData(), (elem.Interface().(float64)))
}
