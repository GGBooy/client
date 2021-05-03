package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func Recv(ctx context.Context, c *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg := recvMsg(ctx, c)
			temp := msg
			switch temp["MessageType"].(string) {
			case "2":
				// 普通文字消息，直接输出
				printMsg(msg)
			case "3":
				// 对 向自己发送文件的请求 做回复
				RecverReply(ctx, c, msg)
			case "4":
				// 对方同意发送，开始发送
				SegSend(ctx, c, msg)
			case "5":
				// 对方拒绝发送/对方没有文件/文件传输完成，啥也不用干
				fmt.Println("refused or nothing or completed")
			case "6":
				// 接收到数据，写入本地
				SegRecv(ctx, c, msg)
			//case "0":
			//	// 服务端对切换请求做出回应
			//	if msg["State"].(bool) == false {
			//		// 切换失败
			//		fmt.Println(msg["Err"].(string))
			//		change <- 3
			//	} else {
			//		// 切换成功
			//		change <- 4
			//	}

			case "9":
				// 服务器命令本地退出
				ch <- 0
			}
		}
	}
}

func recvMsg(ctx context.Context, c *websocket.Conn) map[string]interface{} {
	var v interface{}
	err := wsjson.Read(ctx, c, &v)
	if err != nil {
		log.Println(err)
	}
	msg := v.(map[string]interface{})
	return msg
}

func printMsg(msg map[string]interface{}) {
	fmt.Println(msg["Sendername"].(string) + ": " + msg["Message"].(string))
	fmt.Println(time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
	fmt.Println()
}
