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

type _UpgradeInfo struct {
	requireCount int
	requireGold  int
	getBaseExp   int
}

var (
	heroIndexes      []int
	cardUpgradeInfo  = map[string]_UpgradeInfo{}
	skillUpgradeInfo = map[string]_UpgradeInfo{}
	superUpgradeInfo = map[string]_UpgradeInfo{}
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
		cardUpgradeInfo[key] = _UpgradeInfo{
			requireCount: upgradeData.GetInt("UpgradeRequireCount"),
			requireGold:  upgradeData.GetInt("UpgradeRequireGold"),
			getBaseExp:   upgradeData.GetInt("UpgradeLevelGetBaseExp"),
		}

		skillKey := fmt.Sprintf("%d-%d", cardQuality, cardLevel)
		skillUpgradeInfo[skillKey] = _UpgradeInfo{
			requireCount: 0,
			requireGold:  upgradeData.GetInt("UpgradeSkillRequireGold"),
			getBaseExp:   upgradeData.GetInt("UpgradeSkillGetBaseExp"),
		}
		superUpgradeInfo[skillKey] = _UpgradeInfo{
			requireCount: 0,
			requireGold:  upgradeData.GetInt("UpgradeSuperRequireGold"),
			getBaseExp:   upgradeData.GetInt("UpgradeSuperGetBaseExp"),
		}
	}
	for cardQuality := 1; cardQuality <= 4; cardQuality++ {
		for cardLevel := 1; cardLevel < MAX_CARD_LEVEL; cardLevel++ {
			if _, ok := cardUpgradeInfo[fmt.Sprintf("H-%d-%d", cardQuality, cardLevel)]; !ok {
				log.Panicf("cardUpgradeInfo.H-1-1 not found")
			}
			if _, ok := skillUpgradeInfo[fmt.Sprintf("%d-%d", cardQuality, cardLevel)]; !ok {
				log.Panicf("skillUpgradeInfo.H-1-1 not found")
			}
			if _, ok := superUpgradeInfo[fmt.Sprintf("%d-%d", cardQuality, cardLevel)]; !ok {
				log.Panicf("superUpgradeInfo.H-1-1 not found")
			}
		}
	}
	log.Printf("Upgrade data checked OK!")
}

func GetCardUpgradeInfo(cardType string, cardQuality int, cardLevel int) _UpgradeInfo {
	key := fmt.Sprintf("%s-%d-%d", cardType, cardQuality, cardLevel)
	upgradeInfo, ok := cardUpgradeInfo[key]
	if !ok {
		log.Panicf("Card upgrade not found: %s", key)
	}
	return upgradeInfo
}

func GetSkillUpgradeInfo(cardQuality int, cardLevel int) _UpgradeInfo {
	key := fmt.Sprintf("%d-%d", cardQuality, cardLevel)
	upgradeInfo, ok := skillUpgradeInfo[key]
	if !ok {
		log.Panicf("Skill upgrade not found: %s", key)
	}
	return upgradeInfo
}

func GetSuperUpgradeInfo(cardQuality int, cardLevel int) _UpgradeInfo {
	key := fmt.Sprintf("%d-%d", cardQuality, cardLevel)
	upgradeInfo, ok := superUpgradeInfo[key]
	if !ok {
		log.Panicf("Skill upgrade not found: %s", key)
	}
	return upgradeInfo
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
