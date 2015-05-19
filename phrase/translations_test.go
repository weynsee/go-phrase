package phrase

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestTranslationsService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"de":[%s]}`, testTranslationJSON()))
	})

	translations, err := client.Translations.ListAll()
	if err != nil {
		t.Errorf("Translations.ListAll returned error: %v", err)
	}

	want := []Translation{testTranslation()}
	if got := translations["de"]; !reflect.DeepEqual(got, want) {
		t.Errorf("Translations.ListAll returned %+v, want %+v", got, want)
	}
}

func TestTranslationsService_ListAll_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Translations.ListAll()
		return err
	})
}

func TestTranslationsService_Get(t *testing.T) {
	setup()
	defer teardown()

	updatedSince := time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

	mux.HandleFunc("/translations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryString(t, r, map[string]string{
			"locale_name":   "de",
			"updated_since": "20091117203458",
		})
		fmt.Fprint(w, "["+testTranslationJSON()+"]")
	})

	translations, err := client.Translations.Get("de", &updatedSince)
	if err != nil {
		t.Errorf("Translations.Get returned error: %v", err)
	}

	want := []Translation{testTranslation()}
	if !reflect.DeepEqual(translations, want) {
		t.Errorf("Translations.Get returned %+v, want %+v", translations, want)
	}
}

func TestTranslationsService_Get_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Translations.Get("fr", nil)
		return err
	})
}

func TestTranslationsService_GetByKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translations/fetch_list", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"locale": "de",
			"keys[]": "this.label",
		})
		fmt.Fprint(w, "["+testTranslationJSON()+"]")
	})

	translations, err := client.Translations.GetByKeys("de", []string{"this.label"})
	if err != nil {
		t.Errorf("Translations.GetByKeys returned error: %v", err)
	}

	want := []Translation{testTranslation()}
	if !reflect.DeepEqual(translations, want) {
		t.Errorf("Translations.GetByKeys returned %+v, want %+v", translations, want)
	}
}

func TestTranslationsService_GetByKeys_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Translations.GetByKeys("fr", []string{"blah"})
		return err
	})
}

func TestTranslationsService_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translations/store", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"locale":               "de",
			"key":                  "this.label",
			"skip_verification":    "1",
			"allow_update":         "0",
			"content":              "dummy content",
			"plural_suffix":        "suffix",
			"excluded_from_export": "1",
		})
		fmt.Fprint(w, testTranslationJSON())
	})

	translation := &Translation{PluralSuffix: "suffix", ExcludedFromExport: true, Content: "dummy content"}
	got, err := client.Translations.Update("de", "this.label", translation, true, true)
	if err != nil {
		t.Errorf("Translations.Update returned error: %v", err)
	}

	want := testTranslation()
	if !reflect.DeepEqual(*got, want) {
		t.Errorf("Translations.Update returned %+v, want %+v", got, want)
	}
}

func TestTranslationsService_Update_serverError(t *testing.T) {
	translation := &Translation{PluralSuffix: "suffix", ExcludedFromExport: true, Content: "dummy content"}
	testErrorHandling(t, func() error {
		_, err := client.Translations.Update("de", "this.label", translation, true, true)
		return err
	})
}

func TestTranslationsService_Download(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translations/download", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryString(t, r, map[string]string{
			"convert_emoji":              "1",
			"format":                     "yml",
			"include_empty_translations": "1",
			"keep_notranslate_tags":      "1",
			"locale":                     "ru",
			"tag":                        "dl",
			"updated_since":              "20091117203458",
		})
		w.Header().Add("X-Rate-Limit-Limit", "60")
		w.Header().Add("X-Rate-Limit-Remaining", "59")
		w.Header().Add("X-Rate-Limit-Reset", "1372700873")
	})

	updatedSince := time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	req := &DownloadRequest{
		Locale:       "ru",
		Format:       "yml",
		UpdatedSince: updatedSince,
		Tag:          "dl",
		IncludeEmptyTranslations: true,
		KeepNotranslateTags:      true,
		ConvertEmoji:             true,
	}
	writer := new(bytes.Buffer)
	limit, err := client.Translations.Download(req, writer)
	if err != nil {
		t.Errorf("Translations.Download returned error: %v", err)
	}
	if got, want := limit.Limit, 60; got != want {
		t.Errorf("Rate limit = %v, want %v", got, want)
	}
	if got, want := limit.Remaining, 59; got != want {
		t.Errorf("Rate remaining = %v, want %v", got, want)
	}
	reset := time.Date(2013, 7, 1, 17, 47, 53, 0, time.UTC)
	if limit.Reset.UTC() != reset {
		t.Errorf("Rate reset = %v, want %v", limit.Reset, reset)
	}
}

func TestTranslationsService_Download_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		req := &DownloadRequest{
			Locale: "ru",
			Format: "yml",
		}
		_, err := client.Translations.Download(req, new(bytes.Buffer))
		return err
	})
}

func testTranslationJSON() string {
	return `{"id":1,"content":"This is the help page","plural_suffix":"many","placeholders":[],"unverified":true,"excluded_from_export":true,"translation_key":{"id":23,"name":"page.help.title","description":"This explains what the help section is about","pluralized":true}}`
}

func testTranslation() Translation {
	return Translation{
		ID:                 1,
		Content:            "This is the help page",
		PluralSuffix:       "many",
		Unverified:         true,
		ExcludedFromExport: true,
		Placeholders:       []string{},
		Key: Key{
			ID:          23,
			Name:        "page.help.title",
			Description: "This explains what the help section is about",
			Pluralized:  true,
		},
	}
}
