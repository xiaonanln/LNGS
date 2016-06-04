package lngsdata

import (
	"log"
)

type DataRecord struct {
	data map[string]interface{}
}

func interpretDataRecord(data interface{}) *DataRecord {
	recordData, ok := data.(map[string]interface{})
	if !ok {
		log.Panicf("record should be type map[string]interface{}: %v", data)
	}

	return &DataRecord{
	data: recordData,
	}
}

func (self *DataRecord) Index() int {
	return self.data["Index"].(int)
}

func (self *DataRecord) GetInt(field string) int {
	if self.data[field] == nil {
		return 0
	}

	v := self.data[field].(float64)

	if v != float64(int(v)) {
		// v is not int
		log.Panicf("Field %s is not int: %v", field, v)
	}
	return int(v)
}

func (self *DataRecord) GetFloat(field string) float64 {
	if self.data[field] == nil {
		return 0.0
	}
	return self.data[field].(float64)
}

func (self *DataRecord) GetList(field string) []interface{} {
	return self.data[field].([]interface{})
}
