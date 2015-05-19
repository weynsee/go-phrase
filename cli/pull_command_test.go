package cli

import (
	"fmt"
	mcli "github.com/mitchellh/cli"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
)

func TestPullCommand_Help(t *testing.T) {
	c := PullCommand{}
	if c.Help() == "" {
		t.Fatal("Help should not be empty")
	}
}

func TestPullCommand_Synopsis(t *testing.T) {
	c := PullCommand{}
	if c.Synopsis() == "" {
		t.Fatal("Synopsis should not be empty")
	}
}

func TestPullCommand_updatedSinceInvalidFormat(t *testing.T) {
	ui := new(mcli.MockUi)

	c := &PullCommand{Ui: ui, Config: new(Config), API: nil}
	code := c.Run([]string{"--updated-since=blah"})

	if code == 0 {
		t.Fatal("Pull command should return code != 0")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Error parsing updated-since") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPullCommand_unrecognizedFormat(t *testing.T) {
	ui := new(mcli.MockUi)

	c := &PullCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--format=blah"})

	if code == 0 {
		t.Fatal("Pull command should return code != 0")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Unrecognized format: blah") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPullCommand(t *testing.T) {
	setupAPI()
	defer tearDown()

	var counter int32 = 0
	mux.HandleFunc("/translations/download", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		w.Header().Add("X-Rate-Limit-Limit", "60")
		w.Header().Add("X-Rate-Limit-Remaining", "59")
		w.Header().Add("X-Rate-Limit-Reset", "1372700873")
		fmt.Fprint(w, "OK")
	})

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"name":"en","is_default":true},{"id":2,"name":"ms"}]`)
	})
	ui := new(mcli.MockUi)
	c := &PullCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--encoding=utf-8", "--secret=sikret", "--target=./test", "en", "ms", "unknown"})

	if code != 0 {
		t.Fatal("Pull command should return code == 0")
	}
	if token := c.API.AuthToken; token != "sikret" {
		t.Errorf("API token should be set to %s", "sikret")
	}
	if encoding := c.Config.Encoding; encoding != "utf-8" {
		t.Errorf("Config encoding should be set to %s", "utf-8")
	}
	if target := c.Config.TargetDirectory; target != "./test" {
		t.Errorf("Config target should be set to %s", "./test")
	}
	if format := c.Config.Format; format != "yml" {
		t.Errorf("Config default format should be set to %s", "yml")
	}
	if atomic.LoadInt32(&counter) != 2 {
		t.Errorf("Translations API should have been called %d times", counter)
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Skipping unknown locale unknown") == -1 {
		t.Error("Pull command should print warning for unknown locale")
	}
}

func TestPullCommand_allLocales(t *testing.T) {
	setupAPI()
	defer tearDown()

	var counter int32 = 0
	mux.HandleFunc("/translations/download", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		w.Header().Add("X-Rate-Limit-Limit", "60")
		w.Header().Add("X-Rate-Limit-Remaining", "59")
		w.Header().Add("X-Rate-Limit-Reset", "1372700873")
		fmt.Fprint(w, "OK")
	})

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"name":"en","is_default":true},{"id":2,"name":"ms"},{"id":3,"name":"ms"}]`)
	})
	ui := new(mcli.MockUi)
	c := &PullCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--target=./test"})

	if code != 0 {
		t.Fatal("Pull command should return code == 0")
	}
	if atomic.LoadInt32(&counter) != 3 {
		t.Errorf("Translations API should have been called %d times", counter)
	}
}

func TestPullCommand_listLocalesError(t *testing.T) {
	setupAPI()
	defer tearDown()

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `Service Unavailable`)
	})
	ui := new(mcli.MockUi)
	c := &PullCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--target=./test"})

	if code != 1 {
		t.Fatal("Pull command should return code == 1")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Error encountered fetching the locales") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPullCommand_rateLimited(t *testing.T) {
	setupAPI()
	defer tearDown()

	mux.HandleFunc("/translations/download", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Rate-Limit-Limit", "60")
		w.Header().Add("X-Rate-Limit-Remaining", "0")
		w.Header().Add("X-Rate-Limit-Reset", "1372700873")
	})

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"name":"en","is_default":true}]`)
	})
	ui := new(mcli.MockUi)
	c := &PullCommand{Ui: ui, Config: new(Config), API: client}
	c.Run([]string{"--target=./test"})

	if err := ui.ErrorWriter.String(); strings.Index(err, "Rate limit reached.") == -1 {
		t.Error("Pull command should print warning when rate limit has been reached.")
	}
}

func TestPullCommand_downloadError(t *testing.T) {
	setupAPI()
	defer tearDown()

	mux.HandleFunc("/translations/download", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `Service Unavailable`)
	})

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"name":"en","is_default":true}]`)
	})
	ui := new(mcli.MockUi)
	c := &PullCommand{Ui: ui, Config: new(Config), API: client}
	c.Run([]string{"--target=./test"})

	if err := ui.ErrorWriter.String(); strings.Index(err, "Error downloading locale en") == -1 {
		t.Error("Pull command should print warning when rate limit has been reached.")
	}
}
