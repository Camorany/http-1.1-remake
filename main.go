package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("message.txt")

	if err != nil {
		panic(err)
	}

	for {
		readByte := make([]byte, 8)
		_, err := file.Read(readByte)
		if err != io.EOF {
			fmt.Printf("read: %s\n", readByte)
		} else {
			break
		}
	}

	defer file.Close()
}
