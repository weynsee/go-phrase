package cli

import (
	"flag"
	"fmt"
	mcli "github.com/mitchellh/cli"
	"github.com/weynsee/go-phrase/phrase"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// PullCommand will download locale files from PhraseApp.
type PullCommand struct {
	UI     mcli.Ui
	Config *Config
	API    *phrase.Client
}

const (
	timeFormat            = "20060102150405"
	concurrency           = 2
	defaultDownloadFormat = "yml"
)

// Run executes the pull command.
func (c *PullCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("pull", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	config := c.Config

	cmdFlags.StringVar(&config.Secret, "secret", config.Secret, "")
	cmdFlags.StringVar(&config.TargetDirectory, "target", config.TargetDirectory, "")
	cmdFlags.StringVar(&config.Encoding, "encoding", config.Encoding, "")
	cmdFlags.StringVar(&config.Format, "format", config.Format, "")

	req := new(phrase.DownloadRequest)
	cmdFlags.StringVar(&req.Tag, "tag", "", "")
	var updatedSince string
	cmdFlags.StringVar(&updatedSince, "updated-since", "", "")
	cmdFlags.BoolVar(&req.ConvertEmoji, "convert-emoji", false, "")
	cmdFlags.BoolVar(&req.SkipUnverifiedTranslations, "skip-unverified-translations", false, "")
	cmdFlags.BoolVar(&req.IncludeEmptyTranslations, "include-empty-translations", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if updatedSince != "" {
		var err error
		req.UpdatedSince, err = time.Parse(timeFormat, updatedSince)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing updated-since (%s), format should be YYYYMMDDHHMMSS", updatedSince))
			return 1
		}
	}

	if config.Format == "" {
		config.Format = defaultDownloadFormat
	}

	c.API.AuthToken = config.Secret
	req.Encoding = config.Encoding
	req.Format = config.Format

	if err := config.Valid(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	err := c.fetch(req, cmdFlags.Args())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error encountered fetching the locales:\n\t%s", err.Error()))
		return 1
	}
	return 0
}

func (c *PullCommand) fetch(req *phrase.DownloadRequest, locales []string) error {
	selected, err := c.selectLocales(locales)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(selected))
	gates := make(chan struct{}, concurrency)

	for _, locale := range selected {
		go func(l phrase.Locale) {
			<-gates

			c.fetchLocale(*req, l)

			// start other locales that might still be waiting
			gates <- struct{}{}

			wg.Done()
		}(locale)
	}
	for i := 0; i < concurrency; i++ {
		gates <- struct{}{}
	}

	wg.Wait()
	return nil
}

func (c *PullCommand) fetchLocale(req phrase.DownloadRequest, locale phrase.Locale) {
	lc := c.Config.ForLocale(&locale)
	folder := filepath.Join(lc.TargetDirectory, lc.LocaleDirectory)
	err := os.MkdirAll(folder, 0777)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating folder %s:\n\t%s", folder, err.Error()))
		return
	}
	path := filepath.Join(folder, lc.LocaleFilename)
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating file %s:\n\t%s", path, err.Error()))
		return
	}

	req.Locale = locale.Name
	limit, err := c.API.Translations.Download(&req, file)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error downloading locale %s:\n\t%s", req.Locale, err.Error()))
		return
	}
	if limit.Remaining == 0 {
		c.UI.Error(fmt.Sprintf("Rate limit reached. Please try again at %v", limit.Reset))
		return
	}

	c.UI.Output(fmt.Sprintf("Downloaded %s", path))
}

func (c *PullCommand) selectLocales(locales []string) ([]phrase.Locale, error) {
	all, err := c.API.Locales.ListAll()
	if err != nil {
		return nil, err
	}
	if len(locales) == 0 {
		return all, nil
	}
	localeMap := make(map[string]phrase.Locale)
	for _, locale := range all {
		localeMap[locale.Name] = locale
	}
	var selected = make([]phrase.Locale, 0, len(locales))
	for _, locale := range locales {
		if l, ok := localeMap[locale]; ok {
			selected = append(selected, l)
		} else {
			c.UI.Warn(fmt.Sprintf("Skipping unknown locale %s", locale))
		}
	}
	return selected, nil
}

// Help displays available options for the pull command.
func (c *PullCommand) Help() string {
	helpText := `
	Usage: phrase pull [options] [LOCALE]

	  Download the translation files in the current project.

	Options:

        --format=yml                    See documentation for list of allowed formats
        --target=./phrase/locales       Target folder to store locale files
        --tag=foo                       Limit results to a given tag instead of all translations
        --updated-since=YYYYMMDDHHMMSS  Limit results to translations updated after the given date (UTC)
        --include-empty-translations    Include empty translations in the result
        --convert-emoji                 Convert Emoji symbols
        --encoding=utf-8                Convert .strings or .properties with alternate encoding
        --skip-unverified-translations  Skip unverified translations in the result
        --secret=YOUR_AUTH_TOKEN        The Auth Token to use for this operation instead of the saved one (optional)
	`
	return strings.TrimSpace(helpText)
}

// Synopsis displays a synopsis of the pull command.
func (c *PullCommand) Synopsis() string {
	return "Download the translation files in the current project"
}
