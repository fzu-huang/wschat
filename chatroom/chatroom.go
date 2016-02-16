package model

import (
	"wschat/user"
)

type ActiveChatRoom struct {
	RoomName   string
	onlineUser map[string]user.OnlineUser
}
