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
