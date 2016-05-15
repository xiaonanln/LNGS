package main 

import (
	. "lngs"
)
type Chatroom struct {
	avatars map[*Entity] bool 
}

func NewChatroom() *Chatroom {
	return &Chatroom{
		avatars : make(map[*Entity] bool ), 
	}
}

func (self *Chatroom) Enter(avatar *Entity) {
	self.avatars[avatar] = true
}

func (self *Chatroom) Leave(avatar *Entity) {
	delete(self.avatars, avatar)
}

func (self *Chatroom) Say(avatar *Entity, text string) {
	avatarIcon := avatar.GetInt("icon", 1)
	avatarName := avatar.GetStr("name", "")
	for otherAvatar, _ := range self.avatars {
		otherAvatar.CallClient("OnSay", avatarIcon, avatarName, text)
	}
}