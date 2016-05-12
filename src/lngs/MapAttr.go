package lngs

import "lngs/db"

type MapAttr struct {
	attrs map[string]interface{}
	dirtySet map[string] bool
}

func (self *MapAttr) Set(key string, val interface{}) {
	self.attrs[key] = val
	self.dirtySet[key] = true
}

func (self *MapAttr) GetInt(key string, defaultVal int) int {
	val, ok := self.attrs[key]
	if !ok {
		return defaultVal
	}
	return val.(int)
}

func (self *MapAttr) GetStr(key string, defaultVal string) string {
	val, ok := self.attrs[key]
	if !ok {
		return defaultVal
	}
	return val.(string)
}

func (self *MapAttr) GetMapAttr(key string) *MapAttr {
	val, ok := self.attrs[key]
	if !ok {
		attrs := NewMapAttr()
		self.attrs[key] = attrs
		return attrs
	}
	return val.(*MapAttr)
}

func (self *MapAttr) GetFloat(key string, defaultVal float64) float64 {
	val, ok := self.attrs[key]
	if !ok {
		return defaultVal
	}
	return val.(float64)
}

func (self *MapAttr) GetBool(key string, defaultVal bool) bool {
	val, ok := self.attrs[key]
	if !ok {
		return defaultVal
	}
	return val.(bool)
}

func (self *MapAttr) ToDoc() lngsdb.Doc {
	doc := lngsdb.Doc{}
	for k, v := range self.attrs {
		innerMapAttr, isInnerMapAttr := v.(*MapAttr)
		if isInnerMapAttr {
			doc[k] = innerMapAttr.ToDoc()
		} else {
			doc[k] = v
		}
	}
	return doc 
}


func (self *MapAttr) AssignDoc(doc lngsdb.Doc) {
	for k, v := range doc {
		innerMap, ok := v.(lngsdb.Doc)
		if ok {
			innerMapAttr := NewMapAttr()
			innerMapAttr.AssignDoc(innerMap)
			self.attrs[k] = innerMapAttr
		} else {
			self.attrs[k] = v
		}
		
	}
}

func NewMapAttr() *MapAttr {
	return &MapAttr{
		attrs: make(map[string]interface{}), 
		dirtySet: make(map[string] bool), 
	}
}