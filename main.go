package main

import (
	"github.com/weynsee/go-phrase/cli"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	c := cli.NewCLI("1.0.0", args)

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
