package cli

import (
	"errors"
	"flag"
	"fmt"
	mcli "github.com/mitchellh/cli"
	"github.com/weynsee/go-phrase/phrase"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// PushCommand will upload locale files to PhraseApp.
type PushCommand struct {
	UI     mcli.Ui
	Config *Config
	API    *phrase.Client
}

var defaultLocaleFolder = filepath.Join("config", "locales")

var validTag = regexp.MustCompile(`\A[a-zA-Z0-9\_\-\.]+\z`)

// Run executes the push command.
func (c *PushCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("push", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	config := c.Config

	cmdFlags.StringVar(&config.Secret, "secret", config.Secret, "")
	cmdFlags.StringVar(&config.Format, "format", config.Format, "")

	req := new(phrase.UploadRequest)
	var recursive bool
	cmdFlags.BoolVar(&recursive, "recursive", false, "")
	var tags string
	cmdFlags.StringVar(&tags, "tags", "", "")
	cmdFlags.StringVar(&req.Locale, "locale", "", "")
	cmdFlags.BoolVar(&req.UpdateTranslations, "force-update-translations", false, "")
	cmdFlags.BoolVar(&req.SkipUnverification, "skip-unverification", false, "")
	cmdFlags.BoolVar(&req.SkipUploadTags, "skip-upload-tags", false, "")
	cmdFlags.BoolVar(&req.ConvertEmoji, "convert-emoji", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if tags != "" {
		req.Tags = strings.Split(tags, ",")
		for _, tag := range req.Tags {
			if !validTag.MatchString(tag) {
				c.UI.Error(fmt.Sprintf("Tag %s is invalid: Only letters, numbers, underscores and dashes are allowed", tag))
				return 1
			}
		}
	}

	c.API.AuthToken = config.Secret
	req.Format = config.Format

	return c.upload(req, cmdFlags.Args(), recursive)
}

func (c *PushCommand) upload(req *phrase.UploadRequest, args []string, recursive bool) int {
	selected, err := c.selectFiles(args, recursive)
	if err != nil {
		return 1
	}
	if len(selected) == 0 {
		c.UI.Error("Could not find any files to upload")
		return 1
	} else if len(selected) > 1 && req.Locale != "" {
		c.UI.Error("--locale should not be specified when multiple files are to be uploaded")
		return 1
	}

	supportedFormats := make(map[string]struct{})
	for _, format := range formats {
		for _, ext := range format.properties().extensions {
			supportedFormats[ext] = struct{}{}
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(selected))

	for _, file := range selected {
		ext := fileExtension(file)
		if _, ok := supportedFormats[ext]; ok || rendersLocaleAsExtension(req.Format) {
			go func(f string) {
				err := c.uploadFile(req, f)
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error uploading %s:\n\t%s", f, err.Error()))
				}
				wg.Done()
			}(file)
		} else {
			c.UI.Error(fmt.Sprintf("Could not upload %s (type not supported)", file))
			wg.Done()
		}
	}
	wg.Wait()
	return 0
}

func fileExtension(path string) string {
	return strings.ToLower(strings.Replace(filepath.Ext(path), ".", "", -1))
}

func rendersLocaleAsExtension(format string) bool {
	if format == "" {
		return false
	}
	if fmt, ok := formats[format]; ok {
		return fmt.properties().localeExtension
	}
	return false
}

func (c *PushCommand) uploadFile(req *phrase.UploadRequest, file string) error {
	var tagged string
	if len(req.Tags) > 0 {
		tagged = fmt.Sprintf(" (tagged: %s)", strings.Join(req.Tags, ", "))
	}
	c.UI.Output(fmt.Sprintf("Uploading %s%s...", file, tagged))
	if req.Locale == "" {
		var err error
		req.Locale, err = c.guessLocale(file, req.Format)
		if err != nil {
			return err
		}
	}
	return c.doUpload(*req, file)
}

func (c *PushCommand) doUpload(req phrase.UploadRequest, file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	var content string
	if isUTF16(bytes) {
		var err error
		content, err = decodeUTF16(bytes)
		if err != nil {
			return err
		}
	} else {
		content = string(bytes)
	}
	req.FileContent = content
	req.Filename = file
	return c.API.Keys.Upload(&req)
}

func (c *PushCommand) guessLocale(file, f string) (string, error) {
	if f == "" {
		f = guessFormatFromFileExtension(file)
	}
	validFormat, ok := formats[f]
	if ok && validFormat.properties().localeAware {
		return validFormat.extractLocaleFromPath(c.API, file)
	}
	return findDefaultLocaleName(c.API)
}

func guessFormatFromFileExtension(path string) string {
	extension := fileExtension(path)
	// if the extension is recognized as a format, use it
	if _, ok := formats[extension]; ok {
		return extension
	}
	var possibleFormat string
	for name, format := range formats {
		for _, ext := range format.properties().extensions {
			if ext == extension {
				possibleFormat = name
				break
			}
		}
	}
	return possibleFormat
}

func (c *PushCommand) selectFiles(filenames []string, recursive bool) ([]string, error) {
	if len(filenames) == 0 {
		if src, err := os.Stat(defaultLocaleFolder); err == nil && src.IsDir() {
			c.UI.Warn(fmt.Sprintf("No file or directory specified, using %s", defaultLocaleFolder))
			filenames = append(filenames, defaultLocaleFolder)
		} else {
			c.UI.Error("Need either a file or a directory:")
			c.UI.Error("go-phrase push FILE")
			c.UI.Error("go-phrase push DIRECTORY")
			return nil, errors.New("No files to push")
		}
	}
	var files []string
	for _, file := range filenames {
		if src, err := os.Stat(file); err == nil {
			if src.IsDir() {
				files = listFiles(file, files, recursive)
			} else {
				files = append(files, file)
			}
		}
	}
	return files, nil
}

func listFiles(dir string, list []string, recurse bool) []string {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		name := dir + string(os.PathSeparator) + file.Name()
		if file.IsDir() && recurse {
			list = listFiles(name, list, recurse)
		} else {
			list = append(list, name)
		}
	}
	return list
}

// Help displays available options for the push command.
func (c *PushCommand) Help() string {
	helpText := `
	Usage: phrase push [options] [FILE|DIRECTORY]

	  Upload the translation files in the current project to PhraseApp.

	Options:

        --tags=foo,bar                  List of tags for phrase push (separated by comma)
        --recursive                     Push files in subfolders as well (recursively)
        --locale=en                     Locale of the translations your file contain (required for formats that do not include the name of the locale in the file content)
        --format=yml                    See documentation for list of allowed formats
        --force-update-translations     Force update of existing translations with the file content
        --skip-unverification           When force updating translations, skip unverification of non-main locale translations
        --skip-upload-tags              Don't create upload tags automatically
        --convert-emoji                 Convert Emojis to store and display them correctlyin PhraseApp
        --secret=YOUR_AUTH_TOKEN        The Auth Token to use for this operation instead of the saved one (optional)
	`
	return strings.TrimSpace(helpText)
}

// Synopsis displays a synopsis of the push command.
func (c *PushCommand) Synopsis() string {
	return "Upload the translation files in the current project to PhraseApp"
}
