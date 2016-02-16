package message

import (
	"time"
)

type LogMsg struct {
	UserName string
	LogOp    string
}

type ChatMsg struct {
	UserName string
	Message  string
}

type Message struct {
	MType   string //log msg or chat msg
	LogMSG  LogMsg
	ChatMSG ChatMsg
	Time    time.Time
}
