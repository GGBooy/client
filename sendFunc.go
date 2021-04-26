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
			if smsg[:2] == "###" {
				switch smsg[3] {
				case '0': // 退出
					ch <- 0
				case '1': // 重新连接
					ch <- 1
				case '2': // 发送文件
					//SendFile(ctx, c)
				case '3': // 主动接收（离线）文件、断点续传
					recvFile(ctx, c)

				}
				// 如果发出退出指令，确保Scan阻塞前收到信号
				time.Sleep(100 * time.Millisecond)
			} else {
				SendMsg(ctx, c, sendMessage{MessageType: 2, Message: smsg, Sendername: logData.Username})
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

//var offset int64
//func SendFile(ctx context.Context, c *websocket.Conn) {
//	b = make([]byte, 4096)
//
//	fi, err := os.Open(filename)
//	if err != nil {
//		log.Println(err)
//	}
//	defer fi.Close()
//
//	offset = 0
//	for {
//		num, err := fi.Read(b)
//		if err != io.EOF {
//			fileseg := fileData{Filename: filename, Offset: offset, Data: b[:num]}
//			e := wsjson.Write(ctx, c, fileseg)
//			if e != nil {
//				log.Println(e)
//			}
//		} else {
//			log.Println(err)
//			break
//		}
//	}
//}
