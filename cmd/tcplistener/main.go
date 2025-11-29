package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func main() {
	tcpListener, listenerErr := net.Listen("tcp", "127.0.0.1:42069")

	if listenerErr != nil {
		panic(listenerErr)
	}

	for {
		connection, connErr := tcpListener.Accept()

		if connErr != nil {
			panic(connErr)
		}

		fmt.Printf("TCP Connection Established: Remote Address: %s\n ", connection.RemoteAddr().String())
		msgs := getLinesFromReader(connection)

		for i := range msgs {
			fmt.Printf("read: %s\n", i)
		}
	}

}

func getLinesFromReader(connection io.ReadCloser) <-chan string {

	messages := make(chan string)

	currentline := ""

	go func() {
		for {
			readByte := make([]byte, 8)
			noOfBytes, err := connection.Read(readByte)
			if err != io.EOF {

				readByte = readByte[:noOfBytes]
				if i := bytes.IndexByte(readByte, '\n'); i != -1 {
					currentline += string(readByte[:i])
					messages <- currentline
					readByte = readByte[i+1:]
					currentline = ""
				}

				currentline += string(readByte)

			} else {
				break
			}
		}
		defer connection.Close()
		fmt.Printf("Channel for reading TCP connection messages from %s has been closed\n", connection.RemoteAddr().String())
		defer close(messages)
	}()

	return messages
}
