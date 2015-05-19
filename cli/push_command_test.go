package cli

import (
	"fmt"
	mcli "github.com/mitchellh/cli"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

func TestPushCommand_Help(t *testing.T) {
	c := PushCommand{}
	if c.Help() == "" {
		t.Fatal("Help should not be empty")
	}
}

func TestPushCommand_Synopsis(t *testing.T) {
	c := PushCommand{}
	if c.Synopsis() == "" {
		t.Fatal("Synopsis should not be empty")
	}
}

func TestPushCommand_tagsInvalidFormat(t *testing.T) {
	ui := new(mcli.MockUi)

	c := &PushCommand{Ui: ui, Config: new(Config), API: nil}
	code := c.Run([]string{"--tags=***"})

	if code == 0 {
		t.Fatal("Push command should return code != 0")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Tag *** is invalid") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPushCommand_noFilesError(t *testing.T) {
	ui := new(mcli.MockUi)

	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{""})

	if code == 0 {
		t.Fatal("Push command should return code != 0")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Could not find any files to upload") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPushCommand_noArgsError(t *testing.T) {
	setupAPI()
	defer tearDown()

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--secret=sikret"})

	if code == 0 {
		t.Fatal("Push command should return code != 0")
	}
	if err := ui.ErrorWriter.String(); strings.Index(err, "Need either a file or a directory") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func createTestFiles(files map[string][]byte) {
	prepareLocaleFiles(files, testFolder)
}

func prepareLocaleFiles(files map[string][]byte, folders ...string) {
	folder := filepath.Join(folders...)
	os.MkdirAll(folder, 0777)
	for filename, content := range files {
		file, _ := os.Create(filepath.Join(folder, filename))
		file.Write(content)
		file.Close()
	}
}

func testSuccessfulPush(t *testing.T, files map[string][]byte) {
	setupAPI()
	defer tearDown()

	createTestFiles(files)

	var counter int32 = 0
	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		fmt.Fprint(w, `{"success":true}`)
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{testFolder})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}
	if atomic.LoadInt32(&counter) != int32(len(files)) {
		t.Errorf("Translations API should have been called %d time(s), was called %d time(s)", len(files), counter)
	}
}

func TestPushCommand(t *testing.T) {
	setupAPI()
	defer tearDown()

	createTestFiles(map[string][]byte{
		"en.yml": []byte("testdata"),
	})

	var counter int32 = 0
	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		fmt.Fprint(w, `{"success":true}`)
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--tags=any", "--secret=sikret", "--format=yml", "--locale=en", testFolder})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}
	if token := c.API.AuthToken; token != "sikret" {
		t.Errorf("API token should be set to %s", "sikret")
	}
	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("Translations API should have been called 1 time, was called %d times", counter)
	}
}

func TestPushCommand_fileInvalidContent(t *testing.T) {
	setupAPI()
	defer tearDown()

	createTestFiles(map[string][]byte{
		"ms.yml": []byte{0xff, 0xfe, 'T'},
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	c.Run([]string{"--format=yml", "--locale=ms", testFolder})

	if err := ui.ErrorWriter.String(); strings.Index(err, "Must have even length byte slice") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPushCommand_defaultFolder(t *testing.T) {
	setupAPI()
	defer tearDown()

	prepareLocaleFiles(map[string][]byte{
		"en.yml": []byte("testdata"),
	}, "config", "locales",
	)
	defer os.RemoveAll("config")

	var counter int32 = 0
	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		fmt.Fprint(w, `{"success":true}`)
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--format=yml", "--locale=en"})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}
	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("Translations API should have been called 1 time, was called %d times", counter)
	}
}

func TestPushCommand_utf16(t *testing.T) {
	testSuccessfulPush(t, map[string][]byte{
		"en.yml": []byte{0xff, 0xfe},
	})
}

func TestPushCommand_recursive(t *testing.T) {
	setupAPI()
	defer tearDown()

	prepareLocaleFiles(map[string][]byte{
		"en.yml": []byte{0xff, 0xfe},
	}, "test", "locales",
	)

	var counter int32 = 0
	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		fmt.Fprint(w, `{"success":true}`)
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--format=yml", "--locale=en", "--recursive", testFolder})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}
	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("Translations API should have been called one time, was called %d times", counter)
	}
}

func TestPushCommand_fileInput(t *testing.T) {
	setupAPI()
	defer tearDown()

	createTestFiles(map[string][]byte{
		"en.yml": []byte{0xff, 0xfe},
	})

	var counter int32 = 0
	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		fmt.Fprint(w, `{"success":true}`)
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--format=yml", "--locale=en", filepath.Join(testFolder, "en.yml")})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}
	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("Translations API should have been called 1 time, was called %d times", counter)
	}
}

func TestPushCommand_findDefaultLocaleFromAPI(t *testing.T) {
	setupAPI()
	defer tearDown()

	var localeCounter int32 = 0
	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&localeCounter, 1)
		fmt.Fprint(w, `[{"id":1,"name":"en","is_default":true}]`)
	})

	createTestFiles(map[string][]byte{
		"locale.ini": []byte("test"),
	})

	var counter int32 = 0
	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		fmt.Fprint(w, `{"success":true}`)
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{testFolder})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}
	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("Translations API should have been called 1 time, was called %d times", counter)
	}
	if atomic.LoadInt32(&localeCounter) != 1 {
		t.Errorf("Locales API should have been called 1 time, was called %d times", localeCounter)
	}
}

func TestPushCommand_findDefaultLocaleFromAPIError(t *testing.T) {
	setupAPI()
	defer tearDown()

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `Service Unavailable`)
	})

	createTestFiles(map[string][]byte{
		"locale.ini": []byte("test"),
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	c.Run([]string{testFolder})

	if err := ui.ErrorWriter.String(); strings.Index(err, "Error uploading") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPushCommand_guessFormatFromFileExtension(t *testing.T) {
	testSuccessfulPush(t, map[string][]byte{
		"en.yml": []byte("test"),
	})
}

func TestPushCommand_guessFormatFromSupportedExtensions(t *testing.T) {
	testSuccessfulPush(t, map[string][]byte{
		"en.po": []byte("test"),
	})
}

func TestPushCommand_formatNotSupported(t *testing.T) {
	defer tearDown()

	createTestFiles(map[string][]byte{
		"en.xxx": []byte("testdata"),
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--locale=en", testFolder})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}
}

func TestPushCommand_formatAndExtensionNotSupported(t *testing.T) {
	defer tearDown()

	createTestFiles(map[string][]byte{
		"en.xxx": []byte("testdata"),
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	c.Run([]string{"--format=blah", testFolder})

	if err := ui.ErrorWriter.String(); strings.Index(err, "(type not supported)") == -1 {
		t.Fatal("Ui should display error message")
	}
}

func TestPushCommand_playProperties(t *testing.T) {
	setupAPI()
	defer tearDown()

	createTestFiles(map[string][]byte{
		"messages.en": []byte("testdata"),
	})

	var counter int32 = 0
	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&counter, 1)
		fmt.Fprint(w, `{"success":true}`)
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--locale=en", "--format=play_properties", testFolder})

	if code != 0 {
		t.Fatal("Push command should return code == 0")
	}

	if atomic.LoadInt32(&counter) != 1 {
		t.Errorf("Translations API should have been called 1 time, was called %d times", counter)
	}
}

func TestPushCommand_localeSpecifiedMultipleFilesError(t *testing.T) {
	setupAPI()
	defer tearDown()

	createTestFiles(map[string][]byte{
		"fr.yml": []byte("test"),
		"ru.yml": []byte("test"),
	})

	ui := new(mcli.MockUi)
	c := &PushCommand{Ui: ui, Config: new(Config), API: client}
	code := c.Run([]string{"--format=yml", "--locale=fr", testFolder})

	if code == 0 {
		t.Fatal("Push command should return code != 0")
	}

	if err := ui.ErrorWriter.String(); strings.Index(err, "--locale should not be specified") == -1 {
		t.Fatal("Ui should display error message")
	}
}
