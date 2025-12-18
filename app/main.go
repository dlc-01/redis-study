package main

import (
	"errors"
	"fmt"
	"io"

	netinfra "github.com/codecrafters-io/redis-starter-go/internal/infrestructure/net"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

const (
	serverAddr = "localhost:6379"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	server := netinfra.NewTCPServer(serverAddr)
	server.Start(func(conn ports.Connection) {
		defer conn.Close()

		buf := make([]byte, 1024)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				fmt.Println(err)
				return
			}
			fmt.Printf("Received %d bytes\n", n)

			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})
	
}
