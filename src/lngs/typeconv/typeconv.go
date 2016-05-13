package typeconv

import (
	"fmt"
	"reflect"
)

var (
	emptyInterface     interface{} = 1
	emptyInterfaceType             = reflect.TypeOf(emptyInterface)
)

// try to convert value to target type, panic if fail
func Convert(val reflect.Value, targetType reflect.Type) reflect.Value {
	valType := val.Type()
	if valType.ConvertibleTo(targetType) {
		return val.Convert(targetType)
	}

	fmt.Printf("Value type is %v, emptyInterfaceType is %v, equals %v\n", valType, emptyInterfaceType, valType == emptyInterfaceType)
	interfaceVal := val.Interface()

	switch realVal := interfaceVal.(type) {
	case float64:
		return reflect.ValueOf(realVal).Convert(targetType)
	case []interface{}:
		// val is of type []interface{}, try to convert to typ
		sliceSize := len(realVal)
		targetSlice := reflect.MakeSlice(targetType, 0, sliceSize)
		elemType := targetType.Elem()
		for i := 0; i < sliceSize; i++ {
			targetSlice = reflect.Append(targetSlice, Convert(val.Index(i), elemType))
		}
		return targetSlice
	}

	panic(fmt.Errorf("convert from type %v to %v failed: %v", valType, targetType, val))
}
