package main

import (
	"context"
	"fmt"
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
	fileseg := fileData{MessageType: "3", Filename: filename}
	err = wsjson.Write(ctx, c, fileseg)
	if err != nil {
		log.Println(err)
	}
}

func RecverReply(ctx context.Context, c *websocket.Conn, msg map[string]interface{}) {
	// 确认是否接收文件
	filename := msg["Filename"].(string)
	println("recvive the file ?(y/n) " + filename)
	var temp string
	fmt.Scan(&temp)
	if temp == "n" {
		SendMsg(ctx, c, fileData{MessageType: "5", Filename: filename})
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
		SendMsg(ctx, c, fileData{MessageType: "4", Filename: filename, Offset: "0"})
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
		SendMsg(ctx, c, fileData{MessageType: "4", Filename: filename, Offset: offsetStr})
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

	f, err := os.Open(filename)
	if err != nil {
		fileseg := fileData{MessageType: "5", Filename: filename}
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
		fileseg := fileData{MessageType: "5", Filename: filename}
		SendMsg(ctx, c, fileseg)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	fileseg := fileData{MessageType: "6", Filename: filename, Offset: offsetStr, Data: buffer[:num]}
	SendMsg(ctx, c, fileseg)
}

func SegRecv(ctx context.Context, c *websocket.Conn, msg map[string]interface{}) {
	buffer := msg["Data"].([]byte)
	filename := msg["Filename"].(string)
	offsetStr := msg["Offset"].(string)
	offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	// 在末尾(即offset)写入数据
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
	fileseg := fileData{MessageType: "4", Filename: filename, Offset: posStr}
	SendMsg(ctx, c, fileseg)
}
