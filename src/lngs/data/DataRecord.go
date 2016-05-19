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
	return self.data[field].(int)
}


func (self *DataRecord) GetFloat(field string) float64 {
	return self.data[field].(float64)
}
