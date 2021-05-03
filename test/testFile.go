package main

import (
	"io"
	"log"
	"os"
)

var b = make([]byte, 8192)

func main() {
	fi, err := os.Open("chuang.pdf")
	if err != nil {
		log.Println(err)
	}
	defer fi.Close()

	fo, err := os.Create("destination.pdf")
	if err != nil {
		log.Println(err)
	}
	defer fo.Close()

	for {
		num, err := fi.Read(b)
		//fmt.Println(num)
		if err != io.EOF {
			// fmt.Println(string(b)[:num])
			// offset += int64(num)
			offset, _ := fo.Seek(0, io.SeekEnd)
			fo.WriteAt(b[:num], offset)

		} else {
			log.Println(err)
			break
		}
	}

}
