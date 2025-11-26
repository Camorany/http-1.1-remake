package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("message.txt")

	if err != nil {
		panic(err)
	}

	currentline := ""

	for {
		readByte := make([]byte, 8)
		noOfBytes, err := file.Read(readByte)
		if err != io.EOF {

			readByte = readByte[:noOfBytes]
			if i := bytes.IndexByte(readByte, '\n'); i != -1 {
				currentline += string(readByte[:i])
				readByte = readByte[i+1:]
				fmt.Printf("read: %s\n", currentline)
				currentline = ""
			}

			currentline += string(readByte)

			if len(currentline) != 0 {
				fmt.Printf("read: %s\n", currentline)
			}

		} else {
			break
		}
	}

	defer file.Close()
}
