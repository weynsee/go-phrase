package cli

import (
	mcli "github.com/mitchellh/cli"
	"reflect"
	"testing"
)

func TestNewCLI(t *testing.T) {
	args := []string{"a", "b"}
	cli := NewCLI("0.0.1", args)
	if got := cli.cli.Name; got != "go-phrase" {
		t.Error("cli Name should be go-phrase, not %+v", got)
	}
	if got := cli.cli.Args; !reflect.DeepEqual(got, args) {
		t.Errorf("cli.Args returned %+v, want %+v", got, args)
	}
}

func TestNewCLI_Run(t *testing.T) {
	args := []string{"foo", "-bar", "-baz"}
	cli := NewCLI("0.0.1", args)
	command := new(mcli.MockCommand)
	cli.cli.Commands = map[string]mcli.CommandFactory{
		"foo": func() (mcli.Command, error) {
			return command, nil
		},
	}
	res, err := cli.Run()
	if err != nil {
		t.Fatalf("NewCLI Run returned error %s", err.Error())
	}
	if res != 0 {
		t.Fatal("NewCLI Run exit code should be 0")
	}
}
