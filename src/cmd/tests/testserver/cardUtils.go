package main

import (
	"fmt"
	"lngs/data"
	"lngs/utils"
	"log"
	"strconv"
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
	heroIndexes             []int
	cardUpgradeRequireCount = make(map[string]int)
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

	// 预处理卡牌升级数据
	upgradeIndexes := lngsdata.GetDataRecordIndexes("upgrade")
	for _, index := range upgradeIndexes {
		upgradeData := lngsdata.GetDataRecord("upgrade", index)

		// "2": {
		// 	"CardLevel": 2,
		// 	"CardQuality": 4,
		// 	"CardType": 1,
		// 	"Index": 2,
		// 	"UpgradeRequireCount": 4
		// },
		cardLevel := upgradeData.GetInt("CardLevel")
		cardQuality := upgradeData.GetInt("CardQuality")
		cardType := upgradeData.GetInt("CardType")
		upgradeRequireCount := upgradeData.GetInt("UpgradeRequireCount")

		var cardTypeS string
		if cardType == 1 {
			cardTypeS = "H"
		} else if cardType == 2 {
			cardTypeS = "I"
		}

		key := fmt.Sprintf("%s-%d-%d", cardTypeS, cardQuality, cardLevel)
		cardUpgradeRequireCount[key] = upgradeRequireCount
	}
	for k, v := range cardUpgradeRequireCount {
		log.Printf("Card upgrade %s = %d", k, v)
	}

}

func GetCardUpgradeRequireCount(cardType string, cardQuality int, cardLevel int) int {
	key := fmt.Sprintf("%s-%d-%d", cardType, cardQuality, cardLevel)
	count, ok := cardUpgradeRequireCount[key]
	if !ok {
		log.Panicf("Card upgrade not found: %s", key)
	}
	return count
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

func InterpretCardID(cardID string) (cardType string, cardIndex int, cardData *lngsdata.DataRecord) {
	cardType = cardID[:1]
	cardIndex, err := strconv.Atoi(cardID[1:])
	if err != nil {
		panic(err)
	}

	if cardType == "H" {
		cardData = lngsdata.GetDataRecord("hero", cardIndex)
	} else if cardType == "I" {
		cardData = lngsdata.GetDataRecord("item", cardIndex)
	} else {
		log.Panicf("Invalid card ID: %s", cardID)
	}

	return
}
