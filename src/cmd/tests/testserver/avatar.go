package main

import (
	"fmt"
	"log"

	. "lngs"
	"lngs/data"
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
	} else if cmd == "chest" {
		chestID, _ := strconv.Atoi(args[0])
		behavior.addChest(self, chestID)
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

func (behavior *Avatar) FinishBattle(self *Entity, result int) {
	log.Printf("%v.FinishBattle: result = %v", self, result)
	if result == 1 {
		// win
		behavior.tryGetNewChest(self)
	}
}

// OpenChest : open chest request
func (behavior *Avatar) OpenChest(self *Entity, chestID int) {
	behavior.openChest(self, chestID)
}

func (behavior *Avatar) tryGetNewChest(self *Entity) {
	// get new chest according to avatar level
	chestID := RandInt(1, 4)
	behavior.addChest(self, chestID)
}

func (behavior *Avatar) addChest(self *Entity, chestID int) {
	// add a chest
	lngsdata.GetDataRecord("chest", chestID)

	chests := self.GetMapAttr("chests")

	chestIDStr := strconv.Itoa(chestID)
	chests.Set(chestIDStr, chests.GetInt(chestIDStr, 0)+1)
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

	behavior.addGold(self, RandInt(chestData.GetInt("GoldMin"), chestData.GetInt("GoldMax")))
	cardNum := chestData.GetInt("CardNum")
	cards := behavior.genRandomChestCards(chestID, cardNum)
	// put cards to avatar
	for cardID, num := range cards {
		behavior.addCard(self, cardID, num)
	}

	self.NotifyAttrChange("chests")
	self.NotifyAttrChange("cards")
}

func (behavior *Avatar) genRandomChestCards(chestID int, cardNum int) map[string]int {
	maxHeroIndex := lngsdata.GetMaxDataRecordIndex("hero")
	cards := map[string]int{}

	for i := 0; i < cardNum; i++ {
		heroIndex := RandInt(1, maxHeroIndex)
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
	cards.Set(cardID, cards.GetInt(cardID, 0)+num)
}
