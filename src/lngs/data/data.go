package lngsdata

import (
	"path/filepath"

	"log"
)

var (
	dataPath  = ""
	dataCache = map[string]*DataFile{}
)

func SetDataPath(path string) error {
	_dataPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dataPath = _dataPath
	log.Printf("SetDataPath: %s", dataPath)
	return nil
}

func getDataFile(dataName string) *DataFile {
	data, ok := dataCache[dataName]
	if ok {
		return data
	}

	return openDataFile(dataName)
}

// GetDataRecordIndexe : return all index of data
func GetDataRecordIndexes(dataName string) []int {
	datafile := getDataFile(dataName)
	return datafile.getIndexes()
}

func GetMaxDataRecordIndex(dataName string) int {
	return getDataFile(dataName).getMaxIndex()
}

func GetMinDataRecordIndex(dataName string) int {
	return getDataFile(dataName).getMinIndex()
}

func GetDataRecords(dataName string) map[int]*DataRecord {
	return getDataFile(dataName).records
}

func GetDataRecord(dataName string, recordId int) *DataRecord {
	dataFile := getDataFile(dataName)
	record, ok := dataFile.records[recordId]
	if !ok {
		log.Panicf("record %d not exists", recordId)
	}

	return record
}

func Reload() {
	dataCache = map[string]*DataFile{} // clear
}
