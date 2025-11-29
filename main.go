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

	msgs := getLinesChannel(file)

	for i := range msgs {
		fmt.Printf("read: %s\n", i)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	messages := make(chan string)

	currentline := ""

	go func() {
		for {
			readByte := make([]byte, 8)
			noOfBytes, err := f.Read(readByte)
			if err != io.EOF {

				readByte = readByte[:noOfBytes]
				if i := bytes.IndexByte(readByte, '\n'); i != -1 {
					currentline += string(readByte[:i])
					messages <- currentline
					readByte = readByte[i+1:]
					currentline = ""
				}

				currentline += string(readByte)
				messages <- currentline

			} else {
				break
			}
		}
		defer f.Close()
		defer close(messages)
	}()

	return messages
}
