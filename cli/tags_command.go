package cli

import (
	"flag"
	"fmt"
	mcli "github.com/mitchellh/cli"
	"github.com/weynsee/go-phrase/phrase"
	"strings"
)

// TagsCommand will display all tags in the current project.
type TagsCommand struct {
	Ui     mcli.Ui
	Config *Config
	API    *phrase.Client
}

// Run executes the tags command.
func (c *TagsCommand) Run(args []string) int {
	list := false
	cmdFlags := flag.NewFlagSet("tags", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Config.Secret, "secret", c.Config.Secret, "")
	cmdFlags.BoolVar(&list, "list", false, "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	c.API.AuthToken = c.Config.Secret
	tags, err := c.API.Tags.ListAll()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error encountered while pulling tags from the API: %s", err.Error()))
		return 1
	}
	for _, tag := range tags {
		c.Ui.Output(tag.Name)
	}
	return 0
}

// Help displays available options for the tags command.
func (c *TagsCommand) Help() string {
	helpText := `
	Usage: phrase tags [options]

	  List all the tags in the current project.

	Options:

	  --list                    List all tags
	  --secret=YOUR_AUTH_TOKEN  The Auth Token to use for this operation instead of the saved one (optional)
	`
	return strings.TrimSpace(helpText)
}

// Synopsis displays a synopsis of the tags command.
func (c *TagsCommand) Synopsis() string {
	return "List all the tags in the current project"
}
