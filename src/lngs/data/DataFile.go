package lngsdata

import (
	"io/ioutil"
	"encoding/json"
	"strconv"
	"path/filepath"
	"log"
)

type DataFile struct {
	name string
	records map[int] *DataRecord
}

func openDataFile(dataName string) *DataFile {
	dataFilePath := filepath.Join(dataPath, dataName + ".json")
	log.Printf("Reading data [%s] from file %s ...\n", dataName, dataFilePath)

	data, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		panic(err)
	}

	var j map[string]interface{} 
	err = json.Unmarshal(data, &j)
	if err != nil{
		panic(err)
	}

	records := make(map[int] *DataRecord, len(j))

	for indexStr, data := range j {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			panic(err)
		}
		record := interpretDataRecord(data)
		records[index] = record
	}

	return &DataFile{
		name: dataName, 
		records: records, 
	}
}