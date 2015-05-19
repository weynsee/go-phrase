package cli

import (
	"encoding/json"
	"github.com/weynsee/go-phrase/phrase"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewConfig_New(t *testing.T) {
	path := "testfile_that_does_not_exist"
	config, err := NewConfig(path)
	if err != nil {
		t.Errorf("NewConfig returned error: %v", err)
	}
	if config.path != path {
		t.Errorf("Path should be set to %+v", path)
	}
}

func TestNewConfig_exists(t *testing.T) {
	path := "config.phrase"

	ioutil.WriteFile(path, []byte(`{"secret":"token","default_locale":"ms"}`), 0644)
	defer os.Remove(path)

	config, err := NewConfig(path)
	if err != nil {
		t.Errorf("NewConfig returned error: %v", err)
	}

	if got := config.DefaultLocale; got != "ms" {
		t.Errorf("DefaultLocale got %+v, want \"ms\"", got)
	}

	if got := config.Secret; got != "token" {
		t.Errorf("Secret got %+v, want \"token\"", got)
	}
}

func TestNewConfig_error(t *testing.T) {
	path := "config.phrase"

	ioutil.WriteFile(path, []byte(`false`), 0644)
	defer os.Remove(path)

	_, err := NewConfig(path)
	if err == nil {
		t.Error("Invalid JSON will return an error")
	}
}

func TestNewConfig_ForLocale_parseConfig(t *testing.T) {
	c := Config{Format: "yml", Domain: "phrase", LocaleFilename: "<locale.name>.yml", LocaleDirectory: "<domain>/<locale.code>/"}
	l := &phrase.Locale{Name: "en", Code: "en-us"}
	lc := c.ForLocale(l)
	if want := "en.yml"; lc.LocaleFilename != want {
		t.Errorf("Parsing LocaleFilename expects %s, but got %s", want, lc.LocaleFilename)
	}
	if want := "phrase/en-us/"; lc.LocaleDirectory != want {
		t.Errorf("Parsing LocaleDirectory expects %s, but got %s", want, lc.LocaleDirectory)
	}
}

func TestNewConfig_ForLocale(t *testing.T) {
	c := Config{Format: "yml", Domain: "phrase"}
	l := &phrase.Locale{Name: "en", Code: "en-us"}
	lc := c.ForLocale(l)
	if want := "phrase.en.yml"; lc.LocaleFilename != want {
		t.Errorf("Parsing LocaleFilename expects %s, but got %s", want, lc.LocaleFilename)
	}
	if want := "./"; lc.LocaleDirectory != want {
		t.Errorf("Parsing LocaleDirectory expects %s, but got %s", want, lc.LocaleDirectory)
	}
}

func TestConfig_Valid_error(t *testing.T) {
	c := Config{Format: "blah"}
	if err := c.Valid(); err == nil {
		t.Error("Config with invalid format should not be valid")
	}
}

func TestConfig_Valid(t *testing.T) {
	c := Config{Format: "yml", Domain: "phrase", LocaleFilename: "<locale.name>.yml", LocaleDirectory: "<domain>/<locale.code>/"}
	if err := c.Valid(); err != nil {
		t.Error("Valid config returns no error")
	}
}

func TestNewConfig_Save(t *testing.T) {
	path := "newconfig.phrase"

	defer os.Remove(path)

	config, err := NewConfig(path)
	config.Secret = "newtoken"
	config.Domain = "somedomain"

	if err = config.Save(); err != nil {
		t.Errorf("Save returned error: %v", err)
	}

	newConfig := new(Config)
	f, _ := os.Open(path)
	json.NewDecoder(f).Decode(newConfig)

	if got := newConfig.Domain; got != "somedomain" {
		t.Errorf("Domain got %+v, want \"somedomain\"", got)
	}

	if got := newConfig.Secret; got != "newtoken" {
		t.Errorf("Secret got %+v, want \"newtoken\"", got)
	}
}
