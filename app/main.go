package main

import (
	"bufio"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/application"
	"github.com/codecrafters-io/redis-starter-go/internal/infrastructure/codec"
	netinfra "github.com/codecrafters-io/redis-starter-go/internal/infrastructure/net"
	"github.com/codecrafters-io/redis-starter-go/internal/infrastructure/storage"
	"github.com/codecrafters-io/redis-starter-go/internal/ports"
)

func main() {
	server := netinfra.NewTCPServer("localhost:6379")
	storageMem := storage.NewMemoryStorage()
	processor := application.NewProcessor(storageMem)
	builder := codec.BuilderResponse{}

	server.Start(func(conn ports.Connection) {
		defer conn.Close()

		reader := bufio.NewReader(conn)

		for {
			parser := codec.NewRespParser(reader)

			cmd, err := parser.ReadCommand()
			if err != nil {
				fmt.Println(err)
				return
			}

			resp, err := processor.Handle(cmd)
			if err != nil {
				fmt.Println(err)
				return
			}

			raw, err := builder.Build(resp)
			if err != nil {
				fmt.Println(err)
				return
			}

			conn.Write(raw)
		}
	})
}
