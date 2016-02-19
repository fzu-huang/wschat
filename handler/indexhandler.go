package handler

import (
	"github.com/Unknwon/macaron"
	"github.com/golang/glog"
	room "wschat/chatroom"
	. "wschat/util"
)

func Index(ctx *macaron.Context) {
	ctx.HTML(200, "index")
}

//userlog handler
func WCLogin(context *macaron.Context) {
	//email := context.Req.Request.URL.Query().Get("email")
	email := context.Req.FormValue("email")
	glog.Infoln(email)
	if !CheckEmail(email) {
		context.Write([]byte(email + "," + EMAILREUSED))
		return
	} else {
		//context.Write([]byte(email + "," + LOGINSUS))
		context.Data["WSSERVER"] = WSSERVER
		context.Data["email"] = email
		//context.Write([]byte(email + "," + LOGINSUS))
		context.HTML(200, "room")
		return
	}
}

func WCLogout(context *macaron.Context) {

}

func CreateRoom(ctx *macaron.Context) {
	roomname := ctx.Req.FormValue("roomname")
	if roomname == "" || roomname == PLAZA {
		ctx.Error(404, "invalid name")
		return
	}
	roomlistlock.Lock()
	if _, ok := activeChatRooms[roomname]; ok {
		roomlistlock.Unlock()
		ctx.Error(404, "invalid name")
		return
	}
	tmproom := room.NewRoom(roomname)
	activeChatRooms[roomname] = tmproom
	roomlistlock.Unlock()
	ctx.JSON(200, "create success!")
	return
}

func ListRoom(ctx *macaron.Context) {
	//	hasuser := ctx.Req.FormValue("hasuser")
	//	if hasuser == "" {
	//		rooms := GetAllRooms()
	//		ctx.JSON(200, rooms)
	//	}
	var rooms []string
	roomlistlock.RLock()
	rooms = make([]string, len(activeChatRooms)+1)
	i := 1
	rooms[0] = PLAZA
	for key, _ := range activeChatRooms {
		rooms[i] = key
		i++
	}
	roomlistlock.RUnlock()
	ctx.JSON(200, rooms)
	return
}

func ListusersinRoom(ctx *macaron.Context) {
	roomname := ctx.Req.FormValue("room")
	if roomname == "" {
		ctx.Write([]byte("room name can not be empty!"))
		return
	}
	var result []string
	if roomname == PLAZA {
		defaultChatRoom.Userlock.RLock()
		result = make([]string, len(defaultChatRoom.OnlineUser))
		i := 0
		for name, _ := range defaultChatRoom.OnlineUser {
			result[i] = name
			i++
		}
		defaultChatRoom.Userlock.RUnlock()
	} else {
		var selectroom room.ActiveChatRoom
		roomlistlock.RLock()
		selectroom, ok := activeChatRooms[roomname]
		if !ok {
			roomlistlock.RUnlock()
			ctx.Write([]byte("room name can not be found!"))
			return
		}
		roomlistlock.RUnlock()

		selectroom.Userlock.RLock()
		result = make([]string, len(selectroom.OnlineUser))
		i := 0
		for name, _ := range selectroom.OnlineUser {
			result[i] = name
			i++
		}
		selectroom.Userlock.RUnlock()

	}

	ctx.JSON(200, result)
	return

}
