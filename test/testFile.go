package test

import (
	"fmt"
	"io"
	"log"
	"os"
)

var b = make([]byte, 1)

func main() {
	fi, err := os.Open("source.txt")
	if err != nil {
		log.Println(err)
	}
	defer fi.Close()

	fo, err := os.Create("destination.txt")
	if err != nil {
		log.Println(err)
	}
	defer fo.Close()

	for {
		num, err := fi.Read(b)
		if err != io.EOF {
			fmt.Println(string(b)[:num])
			// offset += int64(num)
			offset, _ := fo.Seek(0, io.SeekEnd)
			fo.WriteAt(b[:num], offset)

		} else {
			log.Println(err)
			break
		}
	}

}
