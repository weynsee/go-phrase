package cli

import (
	"flag"
	"fmt"
	mcli "github.com/mitchellh/cli"
	"github.com/weynsee/go-phrase/phrase"
	"strings"
)

// InitCommand will initialize a PhraseApp config file named .phrase in the current directory.
type InitCommand struct {
	Ui     mcli.Ui
	Config *Config
	API    *phrase.Client
}

// Run executes the init command.
func (c *InitCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("init", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	config := c.Config
	cmdFlags.StringVar(&config.Secret, "secret", config.Secret, "")
	cmdFlags.StringVar(&config.DefaultLocale, "default-locale", config.DefaultLocale, "")
	cmdFlags.StringVar(&config.Format, "default-format", config.Format, "")
	cmdFlags.StringVar(&config.Domain, "domain", config.Domain, "")
	cmdFlags.StringVar(&config.LocaleDirectory, "locale-directory", config.LocaleDirectory, "")
	cmdFlags.StringVar(&config.LocaleFilename, "locale-filename", config.LocaleFilename, "")
	cmdFlags.StringVar(&config.TargetDirectory, "default-target", config.TargetDirectory, "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if config.Secret == "" {
		c.Ui.Error("No auth token was given")
		c.Ui.Error("Please provide the --secret=YOUR_SECRET parameter.")
		return 1
	}

	if err := config.Save(); err != nil {
		c.Ui.Error(fmt.Sprintf("Error encountered while saving the file: %s", err.Error()))
		return 1
	}

	c.Ui.Output("Updated config file .phrase")

	c.API.AuthToken = config.Secret
	if _, err := c.API.Locales.Create(config.DefaultLocale); err != nil {
		c.Ui.Warn(fmt.Sprintf("Notice: Locale \"%s\" could not be created (maybe it already exists)", config.DefaultLocale))
	}

	if _, err := c.API.Locales.MakeDefault(config.DefaultLocale); err != nil {
		c.Ui.Error(fmt.Sprintf("Error encountered while assigning locale %s as default: %s", config.DefaultLocale, err.Error()))
	}
	c.Ui.Output(fmt.Sprintf("Locale \"%s\" is now the default locale", config.DefaultLocale))

	return 0
}

// Help displays available options for the init command.
func (c *InitCommand) Help() string {
	helpText := `
	Usage: phrase init [options]

	  Initializes your project for use with phrase. It will create a default
	  locale if it does not exist yet. In addition, it will create a .phrase
	  config file that contains your client project configuration locally.

	Options:

	  --secret=YOUR_AUTH_TOKEN             Your auth token
	  --default-locale=en                  The default locale for your application
	  --default-format=json                The default format for locale files
	  --domain=phrase                      The default domain or app prefix for locale files
	  --locale-directory=./                The directory naming for locale files, e.g ./<locale.name>/ for subfolders with 'en' or 'de'
	  --locale-filename=<domain>.<format>  The filename for locale files
	  --default-target=phrase/locales/     The default target directory for locale files
	`
	return strings.TrimSpace(helpText)
}

// Synopsis displays a synopsis of the init command.
func (c *InitCommand) Synopsis() string {
	return "Initializes a phrase project"
}
