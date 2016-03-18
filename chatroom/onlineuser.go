package chatroom

import (
	"github.com/golang/glog"
	"golang.org/x/net/websocket"
	//"time"
	"wschat/message"
	//. "wschat/util"
)

type OnlineRoomUser struct {
	InRoom    *ActiveChatRoom
	Conn      *websocket.Conn
	Email     string //USERINFO
	SendQueue chan message.Message
	ExitSig   chan bool //true when user exit the room
	
}

func (user *OnlineRoomUser) Run() {
	go user.readuntildead()
	for {
		select {
		case exit := <-user.ExitSig:
			if exit {
				return
			}

		case send := <-user.SendQueue:

			//			var content string
			//			if send.MType == CHATMSG {
			//				content = send.ChatMSG.UserName + " says: " + send.ChatMSG.Words
			//			} else if send.MType == USERLOGMSG {
			//				content = send.LogMSG.UserName + "  " + send.LogMSG.LogOp
			//			}

			err := websocket.JSON.Send(user.Conn, send)
			//_, err := user.Conn.Write([]byte(content))
			if err != nil {
				glog.Warningln("Can't send msg from server to user: ", user.Email, ".  reason: ", err)
				return
			}
		}
	}
	return

}

func (user *OnlineRoomUser) readuntildead() {
	for {
		//var content string
		var m message.Message
		err := websocket.JSON.Receive(user.Conn, &m)
		//err := websocket.Message.Receive(user.Conn, &content)
		if err != nil {
			glog.Warningln("read json msg wrong!,", err.Error())
			user.ExitSig <-true
			return
		}

		//glog.Infoln(m)
		user.InRoom.Broadcast <- m
	}
}

func CloseUser(user *OnlineRoomUser) {
	user.InRoom.Userlock.Lock()
	delete(user.InRoom.OnlineUser, user.Email)
	user.InRoom.Userlock.Unlock()

	user.Conn.Close()
	close(user.ExitSig)
	close(user.SendQueue)
	user = nil
}

func (user *OnlineRoomUser) waitforsend() {
	//var content string
	for b := range user.SendQueue {
		//content = buildMsgStr(b)
		err := websocket.JSON.Send(user.Conn, b)
		if err != nil {
			glog.Errorln("Can't send msg from server to user: ", user.Email)
			break
		}
	}
}
