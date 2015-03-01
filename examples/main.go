package main

import (
	"fmt"
	"github.com/weynsee/go-phrase/phrase"
	"os"
)

func main() {
	token := os.Getenv("PHRASE_TOKEN")
	if token == "" {
		fmt.Println("Please supply a PHRASE_TOKEN env variable")
		return
	}

	// test out the apis here
	client := phrase.New(token)

	locales, e := client.Locales.ListAll()
	if e != nil {
		panic(e)
	} else {
		for _, locale := range locales {
			fmt.Println(locale)
		}
	}

	f, e := os.Create("en.yml")
	defer f.Close()
	if e != nil {
		panic(e)
	}
	e = client.Locales.Download("en", "yml", f)
	if e != nil {
		panic(e)
	}

	keys, e := client.Keys.ListAll()
	if e != nil {
		panic(e)
	} else {
		for _, key := range keys {
			fmt.Println(key)
		}
	}

	tags, e := client.Tags.ListAll()
	if e != nil {
		panic(e)
	} else {
		for _, tag := range tags {
			fmt.Println(tag)
		}
	}
}
