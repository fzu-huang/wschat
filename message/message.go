package message

type LogMsg struct {
	UserName string `json:"username"`
	LogOp    string `json:"logoperation"`
}

type ChatMsg struct {
	UserName string `json:"username"`
	Words    string `json:"words"`
}

type Message struct {
	MType   string  `json:"msgtype"` //log msg or chat msg
	LogMSG  LogMsg  `json:"logmsg,omitempty"`
	ChatMSG ChatMsg `json:"chatmsg,omitempty"`
	Time    string  `json:"time"`
}
