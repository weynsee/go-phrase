package main

import (
	"fmt"
	"github.com/weynsee/go-phrase/phrase"
)

func main() {
	client := phrase.NewClient()
	fmt.Println(client.Projects.List())
}
