package cli

import (
	"encoding/json"
	"errors"
	"github.com/weynsee/go-phrase/phrase"
	"os"
)

// Config stores the default values for some of the properties of the PhraseApp API client
type Config struct {
	path string

	// Project auth token. You can find the auth token in your project overview or project settings form.
	Secret string `json:"secret"`
	// Default locale  your PhraseApp project (default is en).
	DefaultLocale string `json:"default_locale"`
	// Set a domain for use with Gettext translation files (default is phrase).
	Domain string `json:"domain,omitempty"`
	// Specify a format that should be used as the default format when downloading files (default is yml).
	Format string `json:"format,omitempty"`
	// Set the target directly to store your localization files retrieved by phrase pull. Allows placeholders: https://github.com/phrase/phrase#allowed-placeholders-for-advanced-configuration
	TargetDirectory string `json:"target_directory,omitempty"`
	// Set the directory that contains your source locales (used by the go-phrase push command). Allows placeholders: https://github.com/phrase/phrase#allowed-placeholders-for-advanced-configuration
	LocaleDirectory string `json:"locale_directory,omitempty"`
	// Set the filename for files you download from PhraseApp via phrase pull. Allows placeholders: https://github.com/phrase/phrase#allowed-placeholders-for-advanced-configuration
	LocaleFilename string `json:"locale_filename,omitempty"`
	// Set the encoding for your localization files to UTF-8, UTF-16 or Latin-1. Please note that the encodings only work for a handful of formats like IOS .strings or Java .properties. The default will be UTF-8. If none is provided the default encoding of the formats is used.
	Encoding string `json:"encoding,omitempty"`
}

// Config stores locale specific configuration options
type LocaleConfig struct {
	Config
	format
}

// NewConfig returns a Config instance. Its properties will be
// populated from the file found in the path argument
// (assumed to be in a certain JSON format), and some properties
// will have default values assigned.
func NewConfig(path string) (*Config, error) {
	config := new(Config)
	config.path = path
	if _, err := os.Stat(path); err == nil {
		// if file exists, initialize the object with its contents
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		if err = json.NewDecoder(f).Decode(config); err != nil {
			return nil, err
		}
	}

	// defaults
	if config.Domain == "" {
		config.Domain = "phrase"
	}
	if config.DefaultLocale == "" {
		config.DefaultLocale = "en"
	}

	return config, nil
}

// ForLocale returns a LocaleConfig for the given Locale.
func (c *Config) ForLocale(l *phrase.Locale) *LocaleConfig {
	return newLocaleConfig(c, l)
}

// Valid validates the current config to check whether all values are valid.
func (c *Config) Valid() error {
	_, found := formats[c.Format]
	if !found {
		return errors.New("Unrecognized format: " + c.Format)
	}
	return nil
}

func newLocaleConfig(c *Config, l *phrase.Locale) *LocaleConfig {
	format := formats[c.Format]
	lc := &LocaleConfig{*c, format}
	if lc.LocaleDirectory = replacePlaceholders(c.LocaleDirectory, c, l); lc.LocaleDirectory == "" {
		lc.LocaleDirectory = format.directoryForLocale(c, l)
	}
	if lc.LocaleFilename = replacePlaceholders(c.LocaleFilename, c, l); lc.LocaleFilename == "" {
		lc.LocaleFilename = format.filenameForLocale(c, l)
	}
	if lc.TargetDirectory = c.TargetDirectory; lc.TargetDirectory == "" {
		lc.TargetDirectory = format.properties().targetDirectory
		if lc.TargetDirectory == "" {
			lc.TargetDirectory = "phrase/locales/"
		}
	}
	return lc
}

// Save saves current values of the properties of the Config
// instance to disk.
func (c *Config) Save() error {
	f, err := os.OpenFile(c.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return f.Close()
}
