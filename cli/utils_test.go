package cli

import (
	"fmt"
	"net/http"
	"testing"
)

func TestUtils_replace(t *testing.T) {
	got := replace("hello <this>", map[string]string{"<this>": "world"})
	if want := "hello world"; got != want {
		t.Errorf("Utils replace expected %s, got %s", want, got)
	}
	got = replace("", map[string]string{"<this>": "world"})
	if want := ""; got != want {
		t.Errorf("Utils replace expected %s, got %s", want, got)
	}
}

func TestUtils_findDefaultLocaleNameError(t *testing.T) {
	setupAPI()
	defer shutdownAPI()

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `Service Unavailable`)
	})

	_, err := findDefaultLocaleName(client)
	if err == nil {
		t.Error("Utils findDefaultLocaleName must return server error")
	}
}

func TestUtils_findDefaultLocaleNameNoDefault(t *testing.T) {
	setupAPI()
	defer shutdownAPI()

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[]`)
	})

	name, err := findDefaultLocaleName(client)
	if err != nil {
		t.Errorf("Utils findDefaultLocaleName returned error %v", err.Error())
	}
	if name != "" {
		t.Error("Utils findDefaultLocaleName must return empty string if there is no default locale")
	}
}

var utf16bytes []byte = []byte{
	0xff, // BOM
	0xfe, // BOM
	'T',
	0x00,
	'E',
	0x00,
	'S',
	0x00,
	'T',
	0x00,
	0x6C,
	0x34,
	'\n',
	0x00,
}

func TestUtils_isUTF16(t *testing.T) {
	if !isUTF16(utf16bytes) {
		t.Error("UTF16 byte string not detected correctly")
	}
}

func TestUtils_decodeUTF16(t *testing.T) {
	_, err := decodeUTF16(utf16bytes)
	if err != nil {
		t.Errorf("decodeUTF16 encountered error %+v", err.Error())
	}
}

func TestUtils_decodeUTF16Error(t *testing.T) {
	_, err := decodeUTF16(utf16bytes[0:3])
	if err == nil {
		t.Error("decodeUTF16 should return error for invalid length byte slice")
	}
}
