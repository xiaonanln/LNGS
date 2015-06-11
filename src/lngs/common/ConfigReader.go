package lngscommon

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	config = make(map[string]interface{})
)

func ReadConfigFile(configPath string) {
	configFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}
}

func GetConfig() map[string]interface{} {
	return config
}
