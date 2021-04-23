package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

type login struct {
	MessageType int // 1
	Username string
	Password string
	Mode int
	ID string
}

type sendMessage struct {
	Message string // 2
}




var serverAddr = "192.168.3.16:20229"

func main() {
	ctx, cancel, c := userLogin()
	if c == nil {
		log.Println("connect failed")
		return
	}
	defer cancel()
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	c.Ping(ctx)
}

func userLogin() (context.Context, context.CancelFunc, *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	//defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://" + serverAddr + "/login", nil)
	if err != nil {
		log.Println(err)
		return nil, nil, nil
	}
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")

	var uname, passwd, id string
	var md int
	fmt.Println("please input your username, password, mode, ID")
	fmt.Scan(&uname, &passwd, &md, &id)
	loginData := login{MessageType: 1, Username: uname, Password: passwd, Mode: md, ID: id}

	//fmt.Println(string(b))

	err = wsjson.Write(ctx, c, loginData)
	if err != nil {
		log.Println(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}
	msg := v.(map[string]interface{})
	if msg["State"] == true {
		return ctx, cancel, c
	} else {
		println(msg["Err"])
		close(c)
		return nil, nil, nil
	}
}

func close(c *websocket.Conn)  {
	c.Close(websocket.StatusNormalClosure, "")
}