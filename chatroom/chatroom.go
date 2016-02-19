package chatroom

import (
	"sync"
	"wschat/message"
)

type ActiveChatRoom struct {
	RoomName   string
	OnlineUser map[string]OnlineRoomUser
	Userlock   *sync.RWMutex
	Broadcast  chan message.Message
	Exitsig    chan bool //true when chatroom was destroied
}

func NewRoom(name string) ActiveChatRoom {
	room := ActiveChatRoom{name, make(map[string]OnlineRoomUser), new(sync.RWMutex), make(chan message.Message, 1024), make(chan bool)}
	return room
}

func (room *ActiveChatRoom) Run() {
	for {
		select {
		case msg := <-room.Broadcast:
			for _, user := range room.OnlineUser {
				user.SendQueue <- msg
			}
		case exit := <-room.Exitsig:
			if exit {
				close(room.Broadcast)
				close(room.Exitsig)
				//do close room  ,delete room from map
				return
			}
		}
	}
}
