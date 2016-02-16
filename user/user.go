package user

import (
	"golang.org/x/net/websocket"
	///"wschat/chatroom"
	"wschat/message"
)

type OnlineUser struct {
	//InRoom    *chatroom.ActiveChatRoom
	Conn      *websocket.Conn
	Email     string //USERINFO
	SendQueue chan message.Message
}
