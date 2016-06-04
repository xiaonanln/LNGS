package main

import (
	"lngs/data"
	"lngs/utils"
	"log"
)

const (
	MIN_HERO_INDEX         = 1
	MAX_HERO_INDEX         = 999
	MIN_MONSTER_INDEX      = 1000
	MAX_MONSTER_INDEX      = 9999
	MIN_SPECIAL_UNIT_INDEX = 10000
	BASE_INDEX             = 10000
	TOWER_INDEX            = 10001
)

var (
	heroIndexes []int
)

func cardUtilsInit() {
	_heroIndexes := lngsdata.GetDataRecordIndexes("hero")
	heroIndexes = make([]int, 0, len(_heroIndexes))
	for _, index := range _heroIndexes {
		if index <= MAX_HERO_INDEX && index >= MIN_HERO_INDEX {
			heroIndexes = append(heroIndexes, index)
		}
	}
	log.Printf("Found %d hero cards", len(heroIndexes))
}

func RandHeroIndex() int {
	return lngsutils.ChooseInt(heroIndexes)
}

func RandHeroIndexOfClass(class int) int {
	classHeroIndexes := make([]int, 0, len(heroIndexes))
	for _, index := range heroIndexes {
		heroData := lngsdata.GetDataRecord("hero", index)
		if heroData.GetInt("Class") == class {
			classHeroIndexes = append(classHeroIndexes, index)
		}
	}
	return lngsutils.ChooseInt(classHeroIndexes)
}
