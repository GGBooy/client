package main

type loginMessage struct {
	MessageType string // "1"
	Username    string
	Password    string
	Mode        string
	ID          string
}

type sendMessage struct {
	MessageType string // "2"
	Message     string
	Sendername  string
}

type fileData struct {
	MessageType string //3:请求发送 4:同意接收 5:拒绝接收 6:发送数据
	Filename    string
	Offset      string
	Data        []byte
}
