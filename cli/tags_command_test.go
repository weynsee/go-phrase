package cli

import (
	"fmt"
	mcli "github.com/mitchellh/cli"
	"net/http"
	"strings"
	"testing"
)

func TestTagsCommand_Run(t *testing.T) {
	setupAPI()
	defer shutdownAPI()

	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id": 1,"name":"tag1"},{"id":2,"name":"tag2"}]`)
	})

	ui := new(mcli.MockUi)

	c := &TagsCommand{UI: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--secret=secrettoken"})

	if code != 0 {
		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
	}
}

func TestTagsCommand_error(t *testing.T) {
	setupAPI()
	defer shutdownAPI()

	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	})

	ui := new(mcli.MockUi)

	c := &TagsCommand{UI: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--secret=secrettoken"})

	if code == 0 {
		t.Fatal("API error should return code != 0")
	}

	if err := ui.ErrorWriter.String(); strings.Index(err, "Error encountered") == -1 {
		t.Fatal("UI should display error message")
	}
}

func TestTestCommand_Help(t *testing.T) {
	c := TagsCommand{}
	if c.Help() == "" {
		t.Fatal("Help should not be empty")
	}
}

func TestTagsCommand_Synopsis(t *testing.T) {
	c := TagsCommand{}
	if c.Synopsis() == "" {
		t.Fatal("Synopsis should not be empty")
	}
}
