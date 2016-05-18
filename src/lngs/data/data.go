package lngsdata

import (
	"path/filepath"

	"log"
)

var (
	dataPath = ""
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
