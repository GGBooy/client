package main

import (
	"context"
	"fmt"
	"github.com/GGBooy/message"
	"log"
	"nhooyr.io/websocket"
)

var serverAddr = "127.0.0.1:20229"
var ch = make(chan int, 64) // sendFunc发送信号到主函数
//var change = make(chan int, 64) // 切换连接信号
var chFile = make(chan string, 64)
var logData message.LoginMessage
var chatReq message.ChatRequest

//var state = 0 // 0尚未登录 1已经登录

func main() {
	//ch = make(chan int, 64)
	//change = make(chan int, 64)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := userLogin(ctx)
	if c == nil {
		log.Println("connect failed")
		return
	}
	for {
		var flag string
		fmt.Println("just login ? yes/no")
		_, _ = fmt.Scan(&flag)
		if flag == "yes" {
			break
		} else if flag == "no" {
			chatReqst(ctx, c)
			break
		}
	}
	go Recv(ctx, c)
	go Send(ctx, c)
	ct := control(cancel, ctx, c)
	switch ct {
	case 0:
		println("you have quit")
		return // 退出
	}

}

func userLogin(ctx context.Context) *websocket.Conn {
	c, _, err := websocket.Dial(ctx, "ws://"+serverAddr+"/login", nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	var uname, passwd string
	fmt.Println("please input your username, password")
	_, _ = fmt.Scan(&uname, &passwd)
	logData = message.LoginMessage{MessageType: "1", Username: uname, Password: passwd}

	SendMsg(ctx, c, logData)
	msg := recvMsg(ctx, c)
	if msg["State"] == true {
		println("login successful")
		return c
	} else {
		println(msg["Err"].(string))
		return nil
	}
}

func chatReqst(ctx context.Context, c *websocket.Conn) {
	for {
		var md, id string
		fmt.Println("please input your mode, ID")
		_, _ = fmt.Scan(&md, &id)
		chatReq = message.ChatRequest{MessageType: "7", Mode: md, ID: id}
		SendMsg(ctx, c, chatReq)
		//var temp int
		//temp = <-change
		//if temp == 3 {
		//	// nothing
		//} else if temp == 4 {
		//	fmt.Println("connect to " + chatReq.ID + " successfully")
		//	break
		//}
		break
	}
}

func control(cancel context.CancelFunc, ctx context.Context, c *websocket.Conn) int {
	for {
		var deal = 0
		deal = <-ch
		switch deal {
		case 0: //exit
			SendMsg(ctx, c, message.LogoutRequest{MessageType: "8"})
			closews(cancel, c)
			return 0
			//case 1: //change
			//	chatReqst(ctx, c)
		}
	}
}

func closews(cancel context.CancelFunc, c *websocket.Conn) {
	err := c.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		log.Println(err)
	}
	cancel()
}
