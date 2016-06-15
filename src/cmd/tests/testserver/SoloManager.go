package main

import (
	. "lngs"
	"log"
	"math"
	"sort"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	MIN_SCORE_DIVERSE            = 100 // 一开始，只能匹配到分数差距为10分的人
	ADD_SCORE_DIVERSE_PER_SECOND = 100 // 每秒增加对分数差距的容忍
	MAX_SOLO_MATCH_TIME          = 15  // 最长匹配时间为15s，超出则必然匹配到机器人

	SOLO_INSTANCE_ID = 10000
)

var (
	theSoloManager *SoloManager
)

func init() {
	theSoloManager = newSoloManager()
}

type SoloManager struct {
	avatars map[*Entity]int64
}

func newSoloManager() *SoloManager {
	soloMgr := &SoloManager{
		avatars: make(map[*Entity]int64),
	}
	soloMatchTicker := time.Tick(3 * time.Second)
	go func() {
		for {
			<-soloMatchTicker
			soloMgr.tryMatchAvatars()
		}
	}()
	return soloMgr
}

// OnAvatarLogout 通知玩家登出
func (self *SoloManager) OnAvatarLogout(avatar *Entity) {
	delete(self.avatars, avatar)
}

func (self *SoloManager) StartSolo(avatar *Entity) {
	if _, ok := self.avatars[avatar]; ok {
		return // already soloing, do not do it again ...
	}

	self.avatars[avatar] = GetTime()
}

type _AvatarSST struct {
	avatar *Entity
	sst    int64
}

type SortAvatarSST []_AvatarSST

func (avatars SortAvatarSST) Len() int {
	return len(avatars)
}

func (avatars SortAvatarSST) Less(i, j int) bool {
	asst1 := avatars[i]
	asst2 := avatars[j]

	return asst1.sst < asst2.sst
}

func (avatars SortAvatarSST) Swap(i, j int) {
	tmp := avatars[i]
	avatars[i] = avatars[j]
	avatars[j] = tmp
}

func (self *SoloManager) newBattleID() string {
	return bson.NewObjectId().Hex()
}

func (self *SoloManager) tryMatchAvatars() {
	// 找出合适的玩家并进行匹配，如果实在找不到合适的就让玩家打电脑
	log.Printf("Matching %d avatars...", len(self.avatars))
	if len(self.avatars) <= 0 {
		return
	}

	avatars := make([]_AvatarSST, 0, len(self.avatars))

	for avatar, sst := range self.avatars {
		avatars = append(avatars, _AvatarSST{avatar, sst})
	}

	sort.Sort(SortAvatarSST(avatars)) // 将所有玩家按照匹配时间排列
	matchedAvatars := make(map[*Entity]bool, len(avatars))
	matchResult := make([]struct {
		avatar1 *Entity
		avatar2 *Entity
	}, 0, len(avatars)/2+1)

	for i, avatarSST := range avatars {
		avatar := avatarSST.avatar

		if matchedAvatars[avatar] {
			continue
		}

		avatarCups := avatar.GetInt("cups", 0)

		minIndex := -1
		minScoreDist := 0

		for j := i + 1; j < len(avatars); j++ {
			otherSST := avatars[j]
			other := otherSST.avatar

			if matchedAvatars[other] {
				// avatar already matched, ignore...
				continue
			}

			otherCups := other.GetInt("cups", 0)
			cupsDist := int(math.Abs(float64(avatarCups - otherCups)))
			if minIndex == -1 || cupsDist < minScoreDist {
				minIndex = j
				minScoreDist = cupsDist
			}
		}

		if minIndex == -1 {
			// 没有其他玩家可以匹配了...
			continue
		}

		// 找出了最小分数差距的对付玩家j
		matchedAvatars[avatar] = true
		other := avatars[minIndex].avatar
		matchedAvatars[other] = true
		matchResult = append(matchResult, struct {
			avatar1 *Entity
			avatar2 *Entity
		}{avatar, other})
	}

	log.Printf("Matched %d couples", len(matchResult))
	for _, match := range matchResult {
		avatar1, avatar2 := match.avatar1, match.avatar2
		delete(self.avatars, avatar1)
		delete(self.avatars, avatar2)

		OnSoloMatched(self.newBattleID(), avatar1, avatar2)
	}
}
