package main

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func Recv(ctx context.Context, c *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg := recvMsg(ctx, c)
			switch msg["MessageType"] {
			case 2:
				printMsg(msg)
			case 3:
				recvFile(ctx, c)
			case 4:
				ch <- 0

			}
		}
	}
}

func recvMsg(ctx context.Context, c *websocket.Conn) map[string]interface{} {
	var v interface{}
	err := wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}
	msg := v.(map[string]interface{})
	return msg
}

func printMsg(msg map[string]interface{}) {
	fmt.Println(msg["Username"])
	fmt.Println(": ")
	fmt.Println(msg["Message"])
}

func recvFile(cxt context.Context, c *websocket.Conn) {

}
