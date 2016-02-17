package handler

import (
	"github.com/golang/glog"
	"golang.org/x/net/websocket"
	//	"net/http"
	"sync"
	"time"
	room "wschat/chatroom"
	m "wschat/message"
	. "wschat/util"
)

var activeChatRooms map[string]room.ActiveChatRoom
var roomlistlock *sync.RWMutex

var defaultChatRoom *room.ActiveChatRoom

func init() {
	initChatRoom()
}

func initChatRoom() {
	defaultChatRoom = &room.ActiveChatRoom{
		RoomName:   "plaza",
		OnlineUser: make(map[string]room.OnlineRoomUser),
		Userlock:   new(sync.RWMutex),
		Broadcast:  make(chan m.Message, 1024),
		Exitsig:    make(chan bool),
	}
	go defaultChatRoom.Run()
	activeChatRooms = make(map[string]room.ActiveChatRoom)
	roomlistlock = new(sync.RWMutex)
}

//userroom handler
func WSEnterRoom(conn *websocket.Conn) {
	ret, reason := CheckAndAddUser(conn)
	if !ret {
		conn.Write([]byte(reason))
		return
	}
}

func WSExitRoom(conn *websocket.Conn) {
	ret, reason := CheckAndDelUser(conn)
	if !ret {
		conn.Write([]byte(reason))
		return
	}
	conn.Write([]byte(EXITSUS))
	return
}

func WSCreateRoom(conn *websocket.Conn) {
	roomname := conn.Request().URL.Query().Get("room")
	roomlistlock.Lock()
	if _, ok := activeChatRooms[roomname]; ok {
		roomlistlock.Unlock()
		conn.Write([]byte(ROOMRENAME))
		return
	}
	activeChatRooms[roomname] = room.ActiveChatRoom{
		RoomName:   roomname,
		OnlineUser: make(map[string]room.OnlineRoomUser),
		Userlock:   new(sync.RWMutex),
		Broadcast:  make(chan m.Message, 1024),
		Exitsig:    make(chan bool),
	}
	roomlistlock.Unlock()

	ret, reason := CheckAndAddUser(conn)
	if !ret {
		conn.Write([]byte(reason))
		return
	}
}

func WSDeleteRoom(conn *websocket.Conn) {

}

func CheckEmail(email string) bool {
	defaultChatRoom.Userlock.RLock()
	if _, ok := defaultChatRoom.OnlineUser[email]; ok {
		defaultChatRoom.Userlock.RUnlock()
		return false
	}
	defaultChatRoom.Userlock.RUnlock()
	return true
}

func CheckRoom(roomname string) *room.ActiveChatRoom {
	var selectroom *room.ActiveChatRoom
	if roomname == "" {
		selectroom = defaultChatRoom
	} else {
		roomlistlock.RLock()
		if _, ok := activeChatRooms[roomname]; ok {
			roomlistlock.RUnlock()
			*selectroom = activeChatRooms[roomname]
		} else {
			roomlistlock.RUnlock()
			selectroom = defaultChatRoom
		}
	}
	return selectroom
}

//need to optimize
func CheckAndDelUser(conn *websocket.Conn) (bool, string) {
	email := conn.Request().URL.Query().Get("email")
	roomname := conn.Request().URL.Query().Get("room")
	selectroom := CheckRoom(roomname)

	selectroom.Userlock.Lock()
	if _, ok := selectroom.OnlineUser[email]; ok {
		delete(selectroom.OnlineUser, email)
	}
	selectroom.Userlock.Unlock()
	msg := m.Message{
		MType: USERLOGMSG,
		Time:  time.Now(),
	}
	msg.LogMSG = m.LogMsg{email, EXIT}
	selectroom.Broadcast <- msg
	return true, EXITSUS
}

func CheckAndAddUser(conn *websocket.Conn) (bool, string) {

	email := conn.Request().URL.Query().Get("email")
	//conn.Request().ParseForm()
	//email := conn.Request().FormValue("email")
	//glog.Infoln("email:", conn.Request().URL.RequestURI())
	if email == "" {
		return false, "no email"
	}

	roomname := conn.Request().URL.Query().Get("room")
	selectroom := CheckRoom(roomname)
	glog.Infoln(email, " join the room:", selectroom)
	//	selectroom.Userlock.RLock()
	//	if _, ok := selectroom.OnlineUser[email]; ok {
	//		selectroom.Userlock.RUnlock()
	//		return false, EMAILREUSED
	//	}
	//	selectroom.Userlock.RUnlock()
	user := room.OnlineRoomUser{
		Conn:      conn,
		Email:     email,
		SendQueue: make(chan m.Message, 256),
		InRoom:    selectroom,
		ExitSig:   make(chan bool),
	}

	if AddUser(selectroom, user) {
		msg := m.Message{
			MType: USERLOGMSG,
			Time:  time.Now(),
		}
		msg.LogMSG = m.LogMsg{user.Email, JOIN}
		selectroom.Broadcast <- msg
	} else {
		return false, FINDNOROOM
	}
	user.Run()
	room.CloseUser(&user)
	return false, EXITSUS
}

func AddUser(room *room.ActiveChatRoom, user room.OnlineRoomUser) bool {
	if room == nil || room.OnlineUser == nil {
		return false
	}
	room.Userlock.Lock()
	room.OnlineUser[user.Email] = user
	room.Userlock.Unlock()
	return true
}
