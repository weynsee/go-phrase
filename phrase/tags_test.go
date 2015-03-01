package phrase

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTagsService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id": 1,"name":"tag1"},{"id":2,"name":"tag2"}]`)
	})

	tags, err := client.Tags.ListAll()
	if err != nil {
		t.Errorf("Tags.ListAll returned error: %v", err)
	}

	want := []Tag{{ID: 1, Name: "tag1"}, {ID: 2, Name: "tag2"}}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("Tags.ListAll returned %+v, want %+v", tags, want)
	}
}

func TestTagsService_GetProgress(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tags/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		response := `{"tag":{"id":123,"name":"foo"},"progress":{"de":{"locale":{"id":123,"name":"de"},"progress":{"translations_count": 12,"translated_count":9,"unverified_count":11,"untranslated_count":3}}}}`
		fmt.Fprint(w, response)
	})

	progress, err := client.Tags.GetProgress(1)
	if err != nil {
		t.Errorf("Tags.GetProgress returned error: %v", err)
	}

	want := TagProgress{
		Tag: Tag{ID: 123, Name: "foo"},
		Progress: map[string]LocaleProgress{
			"de": LocaleProgress{
				Locale:   Locale{ID: 123, Name: "de"},
				Progress: Progress{TranslationsCount: 12, TranslatedCount: 9, UnverifiedCount: 11, UntranslatedCount: 3},
			},
		},
	}
	if !reflect.DeepEqual(progress.Progress, want.Progress) {
		t.Errorf("Tags.GetProgress.Progress returned %+v, want %+v", progress.Progress, want.Progress)
	}
	if !reflect.DeepEqual(progress.Tag, want.Tag) {
		t.Errorf("Tags.GetProgress.Tag returned %+v, want %+v", progress.Tag, want.Tag)
	}
}
