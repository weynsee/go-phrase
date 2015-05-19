package cli

import (
	mcli "github.com/mitchellh/cli"
	"github.com/weynsee/go-phrase/phrase"
	"os"
)

var commands map[string]mcli.CommandFactory

func init() {
	ui := &mcli.ConcurrentUi{
		Ui: &mcli.ColoredUi{
			Ui: &mcli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
			WarnColor:   mcli.UiColorYellow,
			ErrorColor:  mcli.UiColorRed,
			OutputColor: mcli.UiColorGreen,
		},
	}

	config, _ := NewConfig(".phrase")
	api := phrase.New(config.Secret)

	commands = map[string]mcli.CommandFactory{
		"init": func() (mcli.Command, error) {
			return &InitCommand{
				Ui:     ui,
				Config: config,
				API:    api,
			}, nil
		},
		"push": func() (mcli.Command, error) {
			return &PushCommand{
				Ui:     ui,
				Config: config,
				API:    api,
			}, nil
		},
		"pull": func() (mcli.Command, error) {
			return &PullCommand{
				Ui:     ui,
				Config: config,
				API:    api,
			}, nil
		},
		"tags": func() (mcli.Command, error) {
			return &TagsCommand{
				Ui:     ui,
				Config: config,
				API:    api,
			}, nil
		},
	}
}
