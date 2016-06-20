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
	cardUpgradeRequireGold  = make(map[string]int)
	skillUpgradeRequireGold = make(map[string]int)
	superUpgradeRequireGold = make(map[string]int)
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

		var cardTypeS string
		if cardType == 1 {
			cardTypeS = "H"
		} else if cardType == 2 {
			cardTypeS = "I"
		}

		key := fmt.Sprintf("%s-%d-%d", cardTypeS, cardQuality, cardLevel)
		cardUpgradeRequireCount[key] = upgradeData.GetInt("UpgradeRequireCount")
		cardUpgradeRequireGold[key] = upgradeData.GetInt("UpgradeRequireGold")

		skillKey := fmt.Sprintf("%d-%d", cardQuality, cardLevel)
		skillUpgradeRequireGold[skillKey] = upgradeData.GetInt("UpgradeSkillRequireGold")
		superUpgradeRequireGold[skillKey] = upgradeData.GetInt("UpgradeSuperRequireGold")
	}
	for cardQuality := 1; cardQuality <= 4; cardQuality++ {
		for cardLevel := 1; cardLevel < MAX_CARD_LEVEL; cardLevel++ {
			if _, ok := cardUpgradeRequireCount[fmt.Sprintf("H-%d-%d", cardQuality, cardLevel)]; !ok {
				log.Panicf("cardUpgradeRequireCount.H-1-1 not found")
			}
			if _, ok := cardUpgradeRequireGold[fmt.Sprintf("H-%d-%d", cardQuality, cardLevel)]; !ok {
				log.Panicf("cardUpgradeRequireGold.H-1-1 not found")
			}
			if _, ok := skillUpgradeRequireGold[fmt.Sprintf("%d-%d", cardQuality, cardLevel)]; !ok {
				log.Panicf("cardUpgradeRequireCount.H-1-1 not found")
			}
			if _, ok := superUpgradeRequireGold[fmt.Sprintf("%d-%d", cardQuality, cardLevel)]; !ok {
				log.Panicf("cardUpgradeRequireCount.H-1-1 not found")
			}
		}
	}
	log.Printf("Upgrade data checked OK!")
}

func GetCardUpgradeRequireCountGold(cardType string, cardQuality int, cardLevel int) (int, int) {
	key := fmt.Sprintf("%s-%d-%d", cardType, cardQuality, cardLevel)
	count, ok := cardUpgradeRequireCount[key]
	if !ok {
		log.Panicf("Card upgrade require count not found: %s", key)
	}

	gold, ok := cardUpgradeRequireGold[key]
	if !ok {
		log.Panicf("Card upgrade require gold not found: %s", key)
	}

	return count, gold
}

func GetSkillUpgradeRequireGold(cardQuality int, cardLevel int) int {
	key := fmt.Sprintf("%d-%d", cardQuality, cardLevel)
	gold, ok := skillUpgradeRequireGold[key]
	if !ok {
		log.Panicf("Skill upgrade not found: %s", key)
	}
	return gold
}

func GetSuperUpgradeRequireGold(cardQuality int, cardLevel int) int {
	key := fmt.Sprintf("%d-%d", cardQuality, cardLevel)
	gold, ok := superUpgradeRequireGold[key]
	if !ok {
		log.Panicf("Skill upgrade not found: %s", key)
	}
	return gold
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
