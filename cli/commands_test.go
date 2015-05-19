package cli

import (
	"testing"
)

func TestCommands(t *testing.T) {
	keys := []string{"push", "pull", "tags", "init"}
	for _, command := range keys {
		_, err := commands[command]()
		if err != nil {
			t.Errorf("%s command returned error %s", command, err.Error())
		}
	}
}
