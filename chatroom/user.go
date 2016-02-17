package chatroom

import (
	"github.com/golang/glog"
	"golang.org/x/net/websocket"
	"time"
	"wschat/message"
	. "wschat/util"
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
				break
			}

		case send := <-user.SendQueue:

			var content string
			if send.MType == CHATMSG {
				content = send.ChatMSG.UserName + "says: " + send.ChatMSG.Words
			} else if send.MType == USERLOGMSG {
				content = send.LogMSG.UserName + "  " + send.LogMSG.LogOp
			}

			_, err := user.Conn.Write([]byte(content))
			if err != nil {
				glog.Errorln("Can't send msg from server to user: ", user.Email, ".  reason: ", err)
				break
			}
		}
	}
	return

}

func (user *OnlineRoomUser) readuntildead() {
	for {
		var content string
		err := websocket.Message.Receive(user.Conn, &content)
		if err != nil {
			return
		}
		m := message.Message{
			MType: CHATMSG,
			Time:  time.Now(),
		}
		m.ChatMSG = message.ChatMsg{
			UserName: user.Email,
			Words:    content,
		}
		glog.Infoln(m)
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
