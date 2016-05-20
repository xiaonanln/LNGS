package lngsdata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
)

type DataFile struct {
	name    string
	records map[int]*DataRecord
	indexes []int
}

func openDataFile(dataName string) *DataFile {
	dataFilePath := filepath.Join(dataPath, dataName+".json")
	log.Printf("Reading data [%s] from file %s ...\n", dataName, dataFilePath)

	data, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		panic(err)
	}

	var j map[string]interface{}
	err = json.Unmarshal(data, &j)
	if err != nil {
		panic(err)
	}

	records := make(map[int]*DataRecord, len(j))
	indexes := make([]int, 0, len(j))

	for indexStr, data := range j {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			panic(err)
		}
		record := interpretDataRecord(data)
		records[index] = record
		indexes = append(indexes, index)
	}

	sort.Ints(indexes) // sort indexes

	return &DataFile{
		name:    dataName,
		records: records,
		indexes: indexes,
	}
}

func (self *DataFile) getMaxIndex() int {
	if len(self.indexes) > 0 {
		return self.indexes[len(self.indexes)-1]
	} else {
		return -1
	}
}

func (self *DataFile) getMinIndex() int {
	if len(self.indexes) > 0 {
		return self.indexes[0]
	} else {
		return -1
	}
}

func (self *DataFile) getIndexes() []int {
	return self.indexes
}
