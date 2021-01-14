package main

import (
	"fmt"
	"os"

	"github.com/mirkobriga/ufirst/client"
	"github.com/mirkobriga/ufirst/server"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("choice a function please")
		return
	}
	if len(args) < 3 {
		fmt.Println("choice an address please")
		return
	}
	if args[1] == "runserver" {
		server.RunServer(args[2])
	} else if args[1] == "connect" {
		client.RunClient(args[2])
	} else {
		fmt.Println("no function specified")
	}

}
