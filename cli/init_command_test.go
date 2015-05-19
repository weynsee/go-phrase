package cli

import (
	mcli "github.com/mitchellh/cli"
	"os"
	"strings"
	"testing"
)

func TestInitCommand_Run(t *testing.T) {
	setupAPI()
	defer shutdownAPI()

	ui := new(mcli.MockUi)
	path := ".testing"
	config, _ := NewConfig(path)
	defer os.Remove(path)

	c := &InitCommand{Ui: ui, Config: config, API: client}
	code := c.Run([]string{"--secret=secrettoken", "--default-format=json"})

	if code != 0 {
		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
	}

	if got := config.Secret; got != "secrettoken" {
		t.Errorf("Config.Secret should be set to %s, was %s", "secrettoken", got)
	}

	if got := config.Format; got != "json" {
		t.Errorf("Config.Format should be set to %s, was %s", "json", got)
	}

	if got := client.AuthToken; got != "secrettoken" {
		t.Errorf("API.AuthToken should be set to %s, was %s", "secrettoken", got)
	}
}

func TestInitCommand_Run_noSecretError(t *testing.T) {
	ui := new(mcli.MockUi)
	c := &InitCommand{Ui: ui, Config: &Config{}}
	code := c.Run([]string{})

	if code != 1 {
		t.Fatal("Run should fail")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "No auth token was given") == -1 {
		t.Fatalf("Error message should indicate that secret is required")
	}
}

func TestInitCommand_Run_saveError(t *testing.T) {
	ui := new(mcli.MockUi)
	c := &InitCommand{Ui: ui, Config: &Config{}}
	code := c.Run([]string{"--secret=this"})

	if code != 1 {
		t.Fatal("Run should fail")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Error encountered while saving the file") == -1 {
		t.Fatalf("Error message should indicate there was an error while saving")
	}
}

func TestInitCommand_Help(t *testing.T) {
	c := InitCommand{}
	if c.Help() == "" {
		t.Fatal("Help should not be empty")
	}
}

func TestInitCommand_Synopsis(t *testing.T) {
	c := InitCommand{}
	if c.Synopsis() == "" {
		t.Fatal("Synopsis should not be empty")
	}
}
