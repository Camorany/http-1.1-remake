package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddr, addrErr := net.ResolveUDPAddr("udp", "localhost:42069")

	if addrErr != nil {
		panic(udpAddr)
	}

	udpConn, connErr := net.DialUDP("udp", nil, udpAddr)

	if connErr != nil {
		panic(udpAddr)
	}

	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">>>")
		content, readErr := reader.ReadString('\n')

		if readErr != nil {
			panic(readErr)
		}

		udpConn.Write([]byte(content))
	}
}
