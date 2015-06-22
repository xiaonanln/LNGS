package lngsdb

import (
	"gopkg.in/mgo.v2/bson"
)

type Doc map[string]interface{} // Type of MongoDB Document

func (self Doc) Id() bson.ObjectId {
	id := self["_id"]
	return id.(bson.ObjectId)
}

func (self Doc) HexId() string {
	return self.Id().Hex()
}

func (self Doc) Get(key string, defaultValue interface{}) interface{} {
	v, exists := self[key]
	if exists {
		return v
	} else {
		return defaultValue
	}
}
