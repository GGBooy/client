package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
)

var serverAddr = "192.168.3.16:20229"
var ch chan int // sendFunc发送信号到主函数
var logData loginMessage
var state int // 0尚未登录 1已经登录

func main() {
	for {
		ch = make(chan int, 64)
		ctx, cancel := context.WithCancel(context.Background())
		c := userLogin(ctx)
		if c == nil {
			log.Println("connect failed")
			return
		}
		go Recv(ctx, c)
		go Send(ctx, c)
		ct := control(cancel, c)
		switch ct {
		case 0:
			println("you have quit")
			return // 退出
		case 1:
			// 再次for循环重置连接
		}
	}

}

func userLogin(ctx context.Context) *websocket.Conn {
	//ctx, cancel := context.WithTimeout(ctx, time.Minute)
	//defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://"+serverAddr+"/login", nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")

	var uname, passwd, md, id string
	if state == 0 {
		fmt.Println("please input your username, password, mode, ID")
		_, _ = fmt.Scan(&uname, &passwd, &md, &id)
		logData = loginMessage{MessageType: "1", Username: uname, Password: passwd, Mode: md, ID: id}
	} else if state == 1 {
		fmt.Println("please input your mode, ID")
		_, _ = fmt.Scan(&md, &id)
		logData.Mode = md
		logData.ID = id
	}
	SendMsg(ctx, c, logData)

	msg := recvMsg(ctx, c)
	if msg["State"] == true {
		println("login successful")
		state = 1
		return c
	} else {
		println(msg["Err"].(string))
		ch <- 0
		return nil
	}
}

func control(cancel context.CancelFunc, c *websocket.Conn) int {
	for {
		var deal = 0
		deal = <-ch
		switch deal {
		case 0: //exit
			closews(cancel, c)
			return 0
		case 1: //change
			closews(cancel, c)
			return 1
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
