package main

import . "lngs"

type _ChatEntry []interface{}

type Chatroom struct {
	name        string
	avatars     map[*Entity]bool
	recentChats []_ChatEntry
}

func NewChatroom(name string) *Chatroom {
	return &Chatroom{
		name:        name,
		avatars:     make(map[*Entity]bool),
		recentChats: make([]_ChatEntry, 0, 10),
	}
}

func (self *Chatroom) Enter(avatar *Entity) {
	self.avatars[avatar] = true
	avatar.CallClient("OnEnterChatRoom", self.name, self.GetAvatarCount())
	myID := avatar.Id()
	for _, chatEntry := range self.recentChats[MaxInt(0, len(self.recentChats)-10):] {
		avatarID, avatarIcon, avatarName, text := chatEntry[0], chatEntry[1], chatEntry[2], chatEntry[3]
		isMe := avatarID == myID
		avatar.CallClient("OnSay", avatarIcon, avatarName, text, isMe)
	}
}

func (self *Chatroom) Leave(avatar *Entity) {
	delete(self.avatars, avatar)
}

func (self *Chatroom) Say(avatar *Entity, text string) {
	self.checkAvatarInChatroom(avatar)
	avatarID := avatar.Id()
	avatarIcon := avatar.GetInt("icon", 1)
	avatarName := avatar.GetStr("name", "")

	chatEntry := _ChatEntry{avatarID, avatarIcon, avatarName, text}
	self.recentChats = append(self.recentChats, chatEntry)
	self.recentChats = self.recentChats[MaxInt(0, len(self.recentChats)-10):] // only save most recent 10 chats

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
