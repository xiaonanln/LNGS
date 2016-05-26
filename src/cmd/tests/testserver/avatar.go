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

func (behavior *Avatar) Init(self *Entity) {
	// Initialize Avatar Attrs by default values
	self.Set("icon", self.GetInt("icon", 1))
	self.Attrs.Set("name", "") // test only
	self.Attrs.Set("exp", self.Attrs.GetInt("exp", 0))
	self.Attrs.Set("gold", self.Attrs.GetInt("gold", 0))
	self.Attrs.Set("diamond", self.Attrs.GetInt("diamond", 0))

	self.Attrs.GetMapAttr("chests")
	self.Attrs.GetMapAttr("cards")

	self.Attrs.GetMapAttr("heroes")
	self.Attrs.GetMapAttr("items")
	self.Attrs.GetMapAttr("embattles")

	self.Attrs.GetMapAttr("testAttr")
}

func (behavior *Avatar) Test(self *Entity) {
	log.Printf("!!!!!!!!!!!!!!!!!!!!!!! Avatar.Test !!!!!!!!!!!!!!!!!!!!!!!")
	testAttr := self.Attrs.GetMapAttr("testAttr")
	testAttr.Set("testAttr", testAttr.GetInt("testVal", 0)+1)
	self.NotifyAttrChange("testAttr")
}

func (behavior *Avatar) OnGetNewClient(self *Entity) {
	onlineManager := GetGlobalEntity("OnlineManager")
	log.Printf("Entity %s get new client, OnlineManager %s", self, onlineManager)
	onlineManager.Call("NotifyAvatarLogin", self.Id())
}

func (behavior *Avatar) OnLoseClient(self *Entity, old_client *GameClient) {
	onlineManager := GetGlobalEntity("OnlineManager")
	log.Printf("Entity %s lose client, OnlineManager %s", self, onlineManager)
	onlineManager.Call("NotifyAvatarLogout", self.Id())

	self.Destroy()
}

func (behavior *Avatar) Say(self *Entity, text string) {
	// player say something in the tribe
	if len(text) > 0 && text[0] == '$' {
		// process gm cmd
		gmcmd := text[1:]
		behavior.handleGmcmd(self, gmcmd)
		return
	}

	worldChatroom.Say(self, text)
}

func (behavior *Avatar) handleGmcmd(self *Entity, gmcmd string) {
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
		behavior.addChest(self, chestID, chestCount)
	} else if cmd == "openChest" {
		chestID, _ := strconv.Atoi(args[0])
		behavior.openChest(self, chestID)
	} else if cmd == "clearCards" {
		behavior.clearCards(self)
	} else {
		self.CallClient("Toast", "无法识别的GM指令："+gmcmd)
		return
	}

	self.CallClient("Toast", "GM指令执行成功："+gmcmd)
}

// GetSaveInterval returns the save interval of entities
func (behavior *Avatar) GetSaveInterval() int {
	return 10
}

// SetAvatarName : set avatar name from client
func (behavior *Avatar) SetAvatarName(self *Entity, name string) {
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
func (behavior *Avatar) EnterWorldChatroom(self *Entity) {
	worldChatroom.Enter(self)
}

func (behavior *Avatar) LeaveWorldChatroom(self *Entity) {
	worldChatroom.Leave(self)
}

// FinishInstance : player finishs an instance
func (behavior *Avatar) FinishInstance(self *Entity, instanceID int, win bool) {
	log.Printf("%v.FinishBattle: instanceID=%v, win = %v", self, instanceID, win)
	instanceData := lngsdata.GetDataRecord("instance", instanceID)
	rewardChests := make([]int, 0, 4)
	if win {
		for chestID := 1; chestID <= 4; chestID++ {
			rewardChestProbKey := "RewardChest" + strconv.Itoa(chestID)
			prob := instanceData.GetFloat(rewardChestProbKey, 0.0)
			if rand.Float64() < prob {
				// get the chest
				behavior.addChest(self, chestID, 1)
				rewardChests = append(rewardChests, chestID)
			}
		}
	}

	self.CallClient("OnFinishInstance", instanceID, win, rewardChests)
}

// EnterInstance : enter instance request
func (behavior *Avatar) EnterInstance(self *Entity, instanceID int) {
	instanceData := lngsdata.GetDataRecord("instance", instanceID)
	monsters := instanceData.GetList("Monsters")
	log.Printf("%v.EnterInstance: instanceID=%v, monsters=%v\n", self, instanceID, monsters)

	cards := self.GetMapAttr("cards") // all cards
	heros := [][]interface{}{
		[]interface{}{[]int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}},
		[]interface{}{[]int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}},
		[]interface{}{[]int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}, []int{0, 0, 0, 0}},
	}

	pos := 0
	for cardID, _cardInfo := range cards.GetAttrs() {
		if cardID[0] != 'H' {
			continue
		}
		heroKind, _ := strconv.Atoi(cardID[1:])
		cardInfo := _cardInfo.(*MapAttr)
		cardLv := cardInfo.GetInt("lv", 1)
		row := pos / 5
		col := pos % 5
		heros[row][col] = []int{heroKind, cardLv, cardLv, cardLv}
		pos = pos + 1

		if pos >= 5 {
			break
		}
	}

	red := map[string]interface{}{
		"heros": heros,
	}

	green := map[string]interface{}{
		"monsters": monsters,
	}

	self.CallClient("OnEnterInstance", instanceID, red, green)
}

// OpenChest : open chest request
func (behavior *Avatar) OpenChest(self *Entity, chestID int) {
	behavior.openChest(self, chestID)
}

func (behavior *Avatar) addChest(self *Entity, chestID int, count int) {
	// add a chest
	lngsdata.GetDataRecord("chest", chestID)

	chests := self.GetMapAttr("chests")

	chestIDStr := strconv.Itoa(chestID)
	chests.Set(chestIDStr, chests.GetInt(chestIDStr, 0)+count)
	self.NotifyAttrChange("chests")
}

func (behavior *Avatar) openChest(self *Entity, chestID int) {
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
	behavior.addGold(self, addGold)
	cardNum := chestData.GetInt("CardNum")
	cards := behavior.genRandomChestCards(chestID, cardNum)
	// put cards to avatar
	for cardID, num := range cards {
		behavior.addCard(self, cardID, num)
	}

	self.NotifyAttrChange("chests")
	self.NotifyAttrChange("cards")
	self.NotifyAttrChange("gold")

	self.CallClient("OnOpenChest", addGold, cards)
}

func (behavior *Avatar) clearCards(self *Entity) {
	self.Set("cards", NewMapAttr())
	self.NotifyAttrChange("cards")
}

func (behavior *Avatar) genRandomChestCards(chestID int, cardNum int) map[string]int {
	cards := map[string]int{}
	heroIndexes := lngsdata.GetDataRecordIndexes("hero")

	for i := 0; i < cardNum; i++ {
		heroIndex := lngsutils.ChooseInt(heroIndexes)
		cardID := fmt.Sprintf("H%d", heroIndex)
		cards[cardID] = cards[cardID] + 1
	}

	return cards
}

func (behavior *Avatar) addGold(self *Entity, gold int) {
	if gold < 0 {
		log.Panicf("addGold: negative gold %d", gold)
		return
	}

	self.Set("gold", self.GetInt("gold", 0)+gold)
}

func (behavior *Avatar) addCard(self *Entity, cardID string, num int) {
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
