package main

import (
	. "lngs"
)

type Chatroom struct {
	name    string
	avatars map[*Entity]bool
}

func NewChatroom(name string) *Chatroom {
	return &Chatroom{
		name:    name,
		avatars: make(map[*Entity]bool),
	}
}

func (self *Chatroom) Enter(avatar *Entity) {
	self.avatars[avatar] = true
	avatar.CallClient("OnEnterChatRoom", self.name, self.GetAvatarCount())
}

func (self *Chatroom) Leave(avatar *Entity) {
	delete(self.avatars, avatar)
}

func (self *Chatroom) Say(avatar *Entity, text string) {
	self.checkAvatarInChatroom(avatar)

	avatarIcon := avatar.GetInt("icon", 1)
	avatarName := avatar.GetStr("name", "")
	for otherAvatar, _ := range self.avatars {
		var isMe int
		if avatar == otherAvatar {
			isMe = 1
		} else {
			isMe = 0
		}
		otherAvatar.CallClient("OnSay", avatarIcon, avatarName, text, isMe)
	}
}

func (self *Chatroom) GetAvatarCount() int {
	return len(self.avatars)
}

func (self *Chatroom) checkAvatarInChatroom(avatar *Entity) {
	_, ok := self.avatars[avatar]
	if !ok {
		self.Enter(avatar)
	}
}
