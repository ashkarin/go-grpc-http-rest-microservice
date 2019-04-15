package main

import (
	"os"

	"github.com/ashkarin/go-grpc-http-rest-microservice/commands"
)

func main() {
	if cmd, err := commands.Execute(os.Args[1:]); err != nil {
		cmd.Println("")
		cmd.Println(cmd.UsageString())
	}
	os.Exit(-1)
}
