package hander

import (
	"golang.org/x/net/websocket"
	room "wschat/chatroom"
	m "wschat/message"
	u "wschat/user"
)

var activeChatRooms []room.ActiveChatRoom

func init() {
	initChatRoom()
}

func initChatRoom() {
	activeChatRooms = []room.ActiveChatRoom{}
}

func WSLogin(conn *websocket.Conn) {
	email := conn.Request().URL.Query().Get("email")
	user := u.OnlineUser{
		Conn:      conn,
		Email:     email,
		SendQueue: make(chan m.Message, 128),
	}

}

func WSLogout(conn *websocket.Conn) {

}

func WSEnterRoom(conn *websocket.Conn) {

}

func WSExitRoom(conn *websocket.Conn) {

}

func WSCreateRoom(conn *websocket.Conn) {

}

func WSDeleteRoom(conn *websocket.Conn) {

}

func WSRandW(conn *websocket.Conn) {
	var msg string
	websocket.Message.Receive(conn, msg)
}
