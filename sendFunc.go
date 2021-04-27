package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func Send(ctx context.Context, c *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var smsg string
			_, _ = fmt.Scan(&smsg)
			if smsg[:3] == "###" {
				switch smsg[3] {
				case '0':
					// 退出
					ch <- 0
				case '1':
					// 重新连接
					ch <- 1
				case '2':
					// 发送文件
					SenderRequest(ctx, c)
				case '3':
					// 主动接收（离线）文件、断点续传
					FileReq(ctx, c)

				}
				// 如果发出退出指令，确保Scan阻塞前收到信号
				time.Sleep(100 * time.Millisecond)
			} else {
				SendMsg(ctx, c, sendMessage{MessageType: "2", Message: smsg, Sendername: logData.Username})
			}
		}
	}
}

func SendMsg(ctx context.Context, c *websocket.Conn, sendData interface{}) {
	err := wsjson.Write(ctx, c, sendData)
	if err != nil {
		log.Println(err)
	}
}

func FileReq(ctx context.Context, c *websocket.Conn) {
	var filename string
	fmt.Println("input the filename you want")
	_, _ = fmt.Scan(&filename)
	FileRequest(ctx, c, filename)
}
