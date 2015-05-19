package phrase

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestLocalesService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":122,"name":"german","code":"de-DE","country_code":"de","writing_direction":"ltr"},{"id":123,"name":"arabic","code":"ar","country_code":"sa","writing_direction":"rtl","pluralizations":{"zero":{"example":"0"},"other":{"gettext_mapping":"0","example":"5"}}}]`)
	})

	locales, err := client.Locales.ListAll()
	if err != nil {
		t.Errorf("Locales.ListAll returned error: %v", err)
	}

	want := []Locale{
		{
			ID:          122,
			Name:        "german",
			Code:        "de-DE",
			CountryCode: "de",
			Direction:   "ltr",
		},
		{
			ID:          123,
			Name:        "arabic",
			Code:        "ar",
			CountryCode: "sa",
			Direction:   "rtl",
			Pluralizations: map[string]map[string]string{
				"zero": map[string]string{"example": "0"},
				"other": map[string]string{
					"gettext_mapping": "0",
					"example":         "5",
				},
			},
		},
	}
	if !reflect.DeepEqual(locales, want) {
		t.Errorf("Locales.ListAll returned %+v, want %+v", locales, want)
	}
}

func TestLocalesService_ListAll_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Locales.ListAll()
		return err
	})
}

func TestLocalesService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/locales", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"locale[name]": "de",
		})
		fmt.Fprint(w, `{"id":123,"name":"de"}`)
	})

	locale, err := client.Locales.Create("de")
	if err != nil {
		t.Errorf("Locales.Create returned error: %v", err)
	}

	want := Locale{ID: 123, Name: "de"}
	if !reflect.DeepEqual(*locale, want) {
		t.Errorf("Locales.Create returned %+v, want %+v", locale, want)
	}
}

func TestLocalesService_Create_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Locales.Create("en")
		return err
	})
}

func TestLocalesService_MakeDefault(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/locales/de/make_default", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":123,"name":"de","code":"de","country_code":"de"}`)
	})

	locale, err := client.Locales.MakeDefault("de")
	if err != nil {
		t.Errorf("Locales.MakeDefault returned error: %v", err)
	}

	want := Locale{ID: 123, Name: "de", Code: "de", CountryCode: "de"}
	if !reflect.DeepEqual(*locale, want) {
		t.Errorf("Locales.MakeDefault returned %+v, want %+v", locale, want)
	}
}

func TestLocalesService_MakeDefault_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Locales.MakeDefault("en")
		return err
	})
}

func TestLocalesService_MakeDefault_invalidLocale(t *testing.T) {
	_, err := client.Locales.MakeDefault("%")
	testParseURLError(t, err)
}

func TestLocalesService_Download(t *testing.T) {
	setup()
	defer teardown()

	content := `{"key":"value"}`

	mux.HandleFunc("/locales/en.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, content)
	})

	output := new(bytes.Buffer)
	err := client.Locales.Download("en", "json", output)
	if err != nil {
		t.Errorf("Locales.Download returned error: %v", err)
	}

	if got := output.String(); !reflect.DeepEqual(content, got) {
		t.Errorf("Locales.Download returned %+v, want %+v", got, content)
	}
}

func TestLocalesService_Download_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		return client.Locales.Download("en", "json", new(bytes.Buffer))
	})
}

func TestLocalesService_Download_invalidLocale(t *testing.T) {
	err := client.Locales.Download("%", "yml", nil)
	testParseURLError(t, err)
}
