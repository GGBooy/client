package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/GGBooy/message"
	"io"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
	"strconv"
)

func fileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 请求发送文件
func SenderRequest(ctx context.Context, c *websocket.Conn) {
	fmt.Println("please input the filename")
	var filename string
	fmt.Scan(&filename)

	// 检查文件存在性
	flag, err := fileExist(filename)
	if err != nil {
		log.Println(err)
	}
	if flag == false {
		fmt.Println("you don't have this file")
		return
	}

	// 发出发送请求
	fileseg := message.FileData{
		MessageType: "3",
		Filename:    filename,
		Sendername:  logData.Username,
	}
	err = wsjson.Write(ctx, c, fileseg)
	if err != nil {
		log.Println(err)
	}
}

func RecverReply(ctx context.Context, c *websocket.Conn, msg map[string]interface{}) {
	// 确认是否接收文件
	filename := msg["Filename"].(string)
	println("recvive the file ?(###y/n) " + filename)
	var temp string
	temp = <-chFile
	if temp == "n" {
		SendMsg(ctx, c, message.FileData{
			MessageType: "5",
			Filename:    filename,
			Sendername:  logData.Username,
		})
		return
	}
	FileRequest(ctx, c, filename)

}

func FileRequest(ctx context.Context, c *websocket.Conn, filename string) {
	flag, err := fileExist(filename)
	if err != nil {
		log.Println(err)
	}
	if flag == false {
		// 文件不存在，请求从头传输
		SendMsg(ctx, c, message.FileData{
			MessageType: "4",
			Filename:    filename,
			Offset:      "0",
			Sendername:  logData.Username,
		})
	} else {
		// 文件已经存在，请求断点续传
		fmt.Println("start from the position last time")
		f, err := os.Open(filename)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		offsetInt, _ := f.Seek(0, io.SeekEnd)
		offsetStr := strconv.FormatInt(offsetInt, 10)
		SendMsg(ctx, c, message.FileData{
			MessageType: "4",
			Filename:    filename,
			Offset:      offsetStr,
			Sendername:  logData.Username,
		})
	}
}

func SegSend(ctx context.Context, c *websocket.Conn, msg map[string]interface{}) {
	// 初始化各参数
	filename := msg["Filename"].(string)
	offsetStr := msg["Offset"].(string)
	offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	buffer := make([]byte, 4096)

	// 打开文件
	f, err := os.Open(filename)
	if err != nil {
		fileseg := message.FileData{
			MessageType: "5",
			Filename:    filename,
			Sendername:  logData.Username,
		}
		SendMsg(ctx, c, fileseg)
		fmt.Println("can't find this file")
		log.Println(err)
		return
	}
	defer f.Close()

	// 偏移至offset位置后发送一段数据
	_, _ = f.Seek(offsetInt, io.SeekStart)
	num, err := f.Read(buffer)
	if err == io.EOF {
		// 如果已经到达文件结尾，停止发送
		fileseg := message.FileData{
			MessageType: "5",
			Filename:    filename,
			Sendername:  logData.Username,
		}
		SendMsg(ctx, c, fileseg)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	fileseg := message.FileData{
		MessageType: "6",
		Filename:    filename,
		Offset:      offsetStr,
		Data:        buffer[:num],
		Sendername:  logData.Username,
	}
	SendMsg(ctx, c, fileseg)
}

func SegRecv(ctx context.Context, c *websocket.Conn, msg map[string]interface{}) {
	bufferStr := msg["Data"].(string)
	buffer, err := base64.StdEncoding.DecodeString(bufferStr)
	if err != nil {
		log.Println(err)
	}
	filename := msg["Filename"].(string)
	offsetStr := msg["Offset"].(string)
	offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	// 检查文件是否存在，不存在则创建一个
	flag, err := fileExist(filename)
	if flag == false {
		f, err := os.Create(filename)
		if err != nil {
			log.Println(err)
		}
		f.Close()
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0777)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	// 在Offset处(即末尾)写入数据
	_, _ = f.Seek(offsetInt, io.SeekStart)
	_, err = f.Write(buffer)
	if err != nil {
		log.Println(err)
	}

	// 获取写入数据后末尾偏移，请求该位置以后的数据
	posInt, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Println(err)
	}
	posStr := strconv.FormatInt(posInt, 10)
	fileseg := message.FileData{
		MessageType: "4",
		Filename:    filename,
		Offset:      posStr,
		Sendername:  logData.Username,
	}
	SendMsg(ctx, c, fileseg)
}
