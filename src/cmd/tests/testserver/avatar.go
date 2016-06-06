package main

import (
	"fmt"
	"log"
	"math/rand"

	. "lngs"
	"lngs/data"
	"lngs/utils"
	"strconv"
	"strings"
)

var (
	worldChatroom = NewChatroom("WorldChatroom")
)

type Avatar struct {
}

func (avatar *Avatar) Init(self *Entity) {
	// Initialize Avatar Attrs by default values
	self.Set("icon", self.GetInt("icon", 1))
	self.Attrs.Set("name", "") // test only
	self.Attrs.Set("exp", self.Attrs.GetInt("exp", 0))
	self.Attrs.Set("level", self.Attrs.GetInt("level", 1))
	self.Attrs.Set("gold", self.Attrs.GetInt("gold", 0))
	self.Attrs.Set("diamond", self.Attrs.GetInt("diamond", 0))

	self.Attrs.GetMapAttr("chests")
	self.Attrs.GetMapAttr("cards")

	self.Attrs.GetMapAttr("heroes")
	self.Attrs.GetMapAttr("items")
	// self.Attrs.Set("embattles", NewMapAttr())
	self.Attrs.GetMapAttr("embattles")

	self.Attrs.GetMapAttr("testAttr")
}

func (avatar *Avatar) Test(self *Entity) {
	log.Printf("!!!!!!!!!!!!!!!!!!!!!!! Avatar.Test !!!!!!!!!!!!!!!!!!!!!!!")
	testAttr := self.Attrs.GetMapAttr("testAttr")
	testAttr.Set("testAttr", testAttr.GetInt("testVal", 0)+1)
	self.NotifyAttrChange("testAttr")
}

func (avatar *Avatar) OnGetNewClient(self *Entity) {
	// 玩家上线
	onlineManager := GetGlobalEntity("OnlineManager")
	log.Printf("Entity %s get new client, OnlineManager %s", self, onlineManager)
	onlineManager.Call("NotifyAvatarLogin", self.Id())

	avatar.tryDailyRefresh(self)
}

func (avatar *Avatar) tryDailyRefresh(self *Entity) {
	// dailyTs := int64(self.GetInt("dailyTs", 0))
	ts := GetTime()
	// if IsSameDay(dailyTs, ts) {
	// 	return
	// }

	// 开始每日刷新逻辑
	self.Set("dailyTs", ts)
	self.NotifyAttrChange("dailyTs")

	avatar.dailyRefresh(self)
}

func (avatar *Avatar) dailyRefresh(self *Entity) {
	log.Printf("%v.dailyRefresh ...", self)
	avatar.refreshShop(self)
}

func (avatar *Avatar) refreshShop(self *Entity) {
	// 刷新商店
	shopInfo := self.GetMapAttr("shop") // 商店数据
	shopInfo.Set("sell1", fmt.Sprintf("H%d", RandHeroIndexOfClass(1)))
	shopInfo.Set("price1", 2)
	shopInfo.Set("sell2", fmt.Sprintf("H%d", RandHeroIndexOfClass(2)))
	shopInfo.Set("price2", 20)
	shopInfo.Set("sell3", fmt.Sprintf("H%d", RandHeroIndexOfClass(3)))
	shopInfo.Set("price3", 2000)

	self.NotifyAttrChange("shop")
}

func (avatar *Avatar) OnLoseClient(self *Entity, old_client *GameClient) {
	// 玩家下线
	onlineManager := GetGlobalEntity("OnlineManager")
	log.Printf("Entity %s lose client, OnlineManager %s", self, onlineManager)
	onlineManager.Call("NotifyAvatarLogout", self.Id())

	self.Destroy()
}

func (avatar *Avatar) Say(self *Entity, text string) {
	// player say something in the tribe
	if len(text) > 0 && text[0] == '$' {
		// process gm cmd
		gmcmd := text[1:]
		avatar.handleGmcmd(self, gmcmd)
		return
	}

	worldChatroom.Say(self, text)
}

func (avatar *Avatar) handleGmcmd(self *Entity, gmcmd string) {
	gmcmd = strings.TrimSpace(gmcmd)
	if gmcmd == "" {
		return
	}

	sp := strings.Split(gmcmd, " ")
	cmd := sp[0]
	args := sp[1:]

	log.Printf("%v: GM >>> %s %v", self, cmd, args)
	if cmd == "gold" {
		gold, _ := strconv.Atoi(args[0])
		self.Set("gold", gold)
		self.NotifyAttrChange("gold")
	} else if cmd == "chest" || cmd == "getChest" {
		chestID, _ := strconv.Atoi(args[0])
		var chestCount int = 1
		if len(args) >= 2 {
			chestCount, _ = strconv.Atoi(args[1])
		}
		avatar.addChest(self, chestID, chestCount)
	} else if cmd == "openChest" {
		chestID, _ := strconv.Atoi(args[0])
		avatar.openChest(self, chestID)
	} else if cmd == "clearCards" {
		avatar.clearCards(self)
	} else if cmd == "set" {
		attrName := args[0]
		val, _ := strconv.Atoi(args[1])
		self.Set(attrName, val)
	} else if cmd == "add" {
		attrName := args[0]
		addVal, _ := strconv.Atoi(args[1])
		self.Set(attrName, self.GetInt(attrName, 0)+addVal)
	} else {
		self.CallClient("Toast", "无法识别的GM指令："+gmcmd)
		return
	}

	self.CallClient("Toast", "GM指令执行成功："+gmcmd)
}

// GetSaveInterval returns the save interval of entities
func (avatar *Avatar) GetSaveInterval() int {
	return 10
}

// SetAvatarName : set avatar name from client
func (avatar *Avatar) SetAvatarName(self *Entity, name string) {
	if name == "" {
		self.CallClient("Toast", "名字不能为空")
		return
	}

	curName := self.GetStr("name", "")
	if curName != "" {
		// avatar already has name
		self.CallClient("OnSetAvatarName", curName)
		return
	}

	self.Set("name", name)
	self.NotifyAttrChange("name")
	self.CallClient("OnSetAvatarName", name)
}

// EnterWorldChatroom : enter world chat room request
func (avatar *Avatar) EnterWorldChatroom(self *Entity) {
	worldChatroom.Enter(self)
}

func (avatar *Avatar) LeaveWorldChatroom(self *Entity) {
	worldChatroom.Leave(self)
}

// FinishInstance : player finishs an instance
func (avatar *Avatar) FinishInstance(self *Entity, instanceID int, win bool) {
	log.Printf("%v.FinishBattle: instanceID=%v, win = %v", self, instanceID, win)
	instanceData := lngsdata.GetDataRecord("instance", instanceID)
	rewardChests := make([]int, 0, 4)
	if win {
		for chestID := 1; chestID <= 4; chestID++ {
			rewardChestProbKey := "RewardChest" + strconv.Itoa(chestID)
			prob := instanceData.GetFloat(rewardChestProbKey)
			if rand.Float64() < prob {
				// get the chest
				avatar.addChest(self, chestID, 1)
				rewardChests = append(rewardChests, chestID)
			}
		}
	}

	self.CallClient("OnFinishInstance", instanceID, win, rewardChests)
}

// EnterInstance : enter instance request
func (avatar *Avatar) EnterInstance(self *Entity, instanceID int) {
	instanceData := lngsdata.GetDataRecord("instance", instanceID)
	monsters := instanceData.GetList("Monsters")

	embattleCards := avatar.getEmbattleCards(self)
	log.Printf("%v.EnterInstance: instanceID=%v, monsters=%v, embattle=%v\n", self, instanceID, monsters, embattleCards)

	red := map[string]interface{}{
		"C1": embattleCards,
	}

	green := map[string]interface{}{
	//"monsters": monsters,
	}

	controlIndex := 1
	self.CallClient("OnEnterInstance", instanceID, red, green, controlIndex)
}

func (avatar *Avatar) getEmbattleCards(self *Entity) map[string]interface{} {
	cards := self.GetMapAttr("cards") // all cards
	embattleIndex := self.GetInt("embattleIndex", 0)
	embattles := self.GetMapAttr("embattles")
	embattle := embattles.GetMapAttr(strconv.Itoa(embattleIndex))

	embattleCards := make(map[string]interface{}, 5)
	for i := 0; i <= 5; i++ {
		cardID := embattle.GetStr(strconv.Itoa(i), "")

		if cardID != "" {
			cardInfo := *cards.GetMapAttr(cardID)
			log.Printf("Embattle ===> Card %s attrs %v", cardID, cardInfo.GetMap())
			embattleCards[cardID] = cardInfo.GetMap()
		}
	}
	return embattleCards
}

// OpenChest : open chest request
func (avatar *Avatar) OpenChest(self *Entity, chestID int) {
	avatar.openChest(self, chestID)
}

// BuyGold : buy gold with diamond
func (avatar *Avatar) BuyGold(self *Entity, gold int) {
	consumeDiamond := -1

	for _, buygoldData := range lngsdata.GetDataRecords("buygold") {
		if buygoldData.GetInt("Gold") == gold {
			consumeDiamond = buygoldData.GetInt("Diamond")
			break
		}
	}

	if consumeDiamond <= 0 {
		// data record not found
		self.CallClient("Toast", "金币数量错误")
		return
	}

	hasDiamond := self.GetInt("diamond", 0)
	if hasDiamond < consumeDiamond {
		self.CallClient("Toast", "宝石不足")
		return
	}

	self.Set("diamond", hasDiamond-consumeDiamond)
	self.Set("gold", self.GetInt("gold", 0)+gold)
	self.CallClient("OnBuyGold", gold, consumeDiamond)
}

func (avatar *Avatar) Embattle(self *Entity, embattleIndex int, cardID string, embattlePos int) {
	cards := self.GetMapAttr("cards")
	if !cards.HasKey(cardID) {
		// card not found
		self.CallClient("Toast", "找不到英雄卡牌："+cardID)
		return
	}

	embattles := self.GetMapAttr("embattles")
	embattle := embattles.GetMapAttr(strconv.Itoa(embattleIndex))
	embattle.Set(strconv.Itoa(embattlePos), cardID)
	self.NotifyAttrChange("embattles")

	self.CallClient("OnEmbattle", embattleIndex)
}

func (avatar *Avatar) SetEmbattleIndex(self *Entity, embattleIndex int) {
	if embattleIndex != 1 && embattleIndex != 2 && embattleIndex != 3 {
		return
	}

	self.Set("embattleIndex", embattleIndex)
	self.NotifyAttrChange("embattleIndex")
}

func (avatar *Avatar) UpgradeCard(self *Entity, cardID string) {
	cardType, _, cardData := InterpretCardID(cardID)
	cardQuality := cardData.GetInt("Class")

	cards := self.GetMapAttr("cards")
	if !cards.HasKey(cardID) {
		log.Printf("Card %s not found", cardID)
		return
	}

	cardInfo := cards.GetMapAttr(cardID)
	cardLevel := cardInfo.GetInt("lv", 0)
	cardCount := cardInfo.GetInt("num", 0)

	requireCount := GetCardUpgradeRequireCount(cardType, cardQuality, cardLevel)

	if cardLevel >= MAX_CARD_LEVEL || cardCount < requireCount {
		return
	}
	cardInfo.Set("num", cardCount-requireCount)
	cardLevel = cardLevel + 1
	cardInfo.Set("lv", cardLevel)

	self.NotifyAttrChange("cards")
	self.CallClient("OnUpgradeCard", cardID, cardLevel)
}

// func (avatar *Avatar) getCardLevel(self *Entity, cardID string) int {
// 	cards := self.GetMapAttr("cards")
// 	if !cards.HasKey(cardID) {
// 		return 0
// 	}

// 	return
// }

func (avatar *Avatar) addChest(self *Entity, chestID int, count int) {
	// add a chest
	lngsdata.GetDataRecord("chest", chestID)

	chests := self.GetMapAttr("chests")

	chestIDStr := strconv.Itoa(chestID)
	chests.Set(chestIDStr, chests.GetInt(chestIDStr, 0)+count)
	self.NotifyAttrChange("chests")
}

func (avatar *Avatar) openChest(self *Entity, chestID int) {
	// open a chest
	chests := self.GetMapAttr("chests")
	chestData := lngsdata.GetDataRecord("chest", chestID)

	chestIDStr := strconv.Itoa(chestID)
	chestCount := chests.GetInt(chestIDStr, 0)
	if chestCount <= 0 {
		// got no chest
		return
	}

	chests.Set(chestIDStr, chestCount-1) // reduce chest count first

	addGold := lngsutils.RandInt(chestData.GetInt("GoldMin"), chestData.GetInt("GoldMax"))
	avatar.addGold(self, addGold)
	cardNum := chestData.GetInt("CardNum")
	cards := avatar.genRandomChestCards(chestID, cardNum)
	// put cards to avatar
	for cardID, num := range cards {
		avatar.addCard(self, cardID, num)
	}

	self.NotifyAttrChange("chests")
	self.NotifyAttrChange("cards")
	self.NotifyAttrChange("gold")

	self.CallClient("OnOpenChest", addGold, cards)
}

func (avatar *Avatar) clearCards(self *Entity) {
	self.Set("cards", NewMapAttr())
	self.Set("embattles", NewMapAttr())

	self.NotifyAttrChange("cards", "embattles")
}

func (avatar *Avatar) genRandomChestCards(chestID int, cardNum int) map[string]int {
	cards := map[string]int{}

	for i := 0; i < cardNum; i++ {
		heroIndex := RandHeroIndex()
		cardID := fmt.Sprintf("H%d", heroIndex)
		cards[cardID] = cards[cardID] + 1
	}

	return cards
}

func (avatar *Avatar) addGold(self *Entity, gold int) {
	if gold < 0 {
		log.Panicf("addGold: negative gold %d", gold)
		return
	}

	self.Set("gold", self.GetInt("gold", 0)+gold)
}

func (avatar *Avatar) addCard(self *Entity, cardID string, num int) {
	if num < 0 {
		log.Panicf("addCard %s: negative num %d", cardID, num)
		return
	}

	cards := self.GetMapAttr("cards")
	cardInfo := cards.GetMapAttr(cardID)

	cardLv := cardInfo.GetInt("lv", 0)
	if cardLv == 0 {
		cardInfo.Set("lv", 1) // 第一次获得卡牌的时候自动设置为等级1
	}

	cardInfo.Set("num", cardInfo.GetInt("num", 0)+num)
}
