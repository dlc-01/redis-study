package main

import (
	"bufio"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	"github.com/codecrafters-io/redis-starter-go/internal/infrastructure/codec"
	netinfra "github.com/codecrafters-io/redis-starter-go/internal/infrastructure/net"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

const serverAddr = "localhost:6379"

func main() {
	server := netinfra.NewTCPServer(serverAddr)
	processor := application.NewProcessor()
	builder := codec.BuilderResponse{}

	server.Start(func(conn ports.Connection) {
		defer conn.Close()

		reader := bufio.NewReader(conn)
		parser := codec.NewRespParser(reader)

		for {
			cmd, err := parser.ReadCommand()
			if err != nil {
				fmt.Println("read error:", err)
				return
			}

			res, err := processor.Process(cmd)
			if err != nil {
				fmt.Println("process error:", err)
				return
			}

			raw, err := builder.Build(res)
			if err != nil {
				fmt.Println("build error:", err)
				return
			}

			if _, err := conn.Write(raw); err != nil {
				return
			}
		}
	})
}
