package main

import (
	"os"
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

//func negotiate(ctx context.Context, c *websocket.Conn) int64 {
//	fmt.Println("please input the filename")
//	var filename string
//	fmt.Scan(&filename)
//	fileseg := fileData{MessageType: 3, Filename: filename}
//}
