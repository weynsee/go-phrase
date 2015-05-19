package phrase

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestKeysService_ListAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1,"name":"helper.label","description":"Some description","pluralized":true,"data_type":"string","tag_list":["mobile","example-feature"]}]`)
	})

	keys, err := client.Keys.ListAll()
	if err != nil {
		t.Errorf("Keys.ListAll returned error: %v", err)
	}

	want := []Key{
		{
			ID: 1, Name: "helper.label",
			Description: "Some description",
			Pluralized:  true,
			DataType:    "string",
			Tags:        []string{"mobile", "example-feature"},
		},
	}
	if !reflect.DeepEqual(keys, want) {
		t.Errorf("Keys.ListAll returned %+v, want %+v", keys, want)
	}
}

func TestKeysService_ListAll_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Keys.ListAll()
		return err
	})
}

func TestKeysService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryString(t, r, map[string]string{
			"key_names[]": "my.key",
		})
		fmt.Fprint(w, `[{"id":1,"name":"helper.label","description":"Some description","pluralized":true,"data_type":"string","tag_list":["mobile","example-feature"]}]`)
	})

	keys, err := client.Keys.Get([]string{"my.key"})
	if err != nil {
		t.Errorf("Keys.Get returned error: %v", err)
	}

	want := []Key{
		{
			ID: 1, Name: "helper.label",
			Description: "Some description",
			Pluralized:  true,
			DataType:    "string",
			Tags:        []string{"mobile", "example-feature"},
		},
	}
	if !reflect.DeepEqual(keys, want) {
		t.Errorf("Keys.Get returned %+v, want %+v", keys, want)
	}
}

func TestKeysService_Get_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Keys.Get([]string{"my.key"})
		return err
	})
}

func TestKeysService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"translation_key[name]":                   "foo.key",
			"translation_key[name_plural]":            "foo.key.pluralized",
			"translation_key[pluralized]":             "1",
			"translation_key[description]":            "Some description",
			"translation_key[data_type]":              "string",
			"translation_key[tag_names]":              "foo,bar",
			"translation_key[unformatted]":            "1",
			"translation_key[xml_space_preserve]":     "1",
			"translation_key[max_characters_allowed]": "140",
		})
		fmt.Fprint(w, `{"id":123,"name":"foo.key","name_plural":"foo.key.pluralized","pluralized":true,"description":"Some description","data_type":"string","tag_list":["foo","bar"],"unformatted":true,"xml_space_preserve":true,"max_characters_allowed":140}`)
	})

	want := &Key{
		ID:               123,
		Name:             "foo.key",
		NamePlural:       "foo.key.pluralized",
		Pluralized:       true,
		Description:      "Some description",
		DataType:         "string",
		Tags:             []string{"foo", "bar"},
		Unformatted:      true,
		XMLSpacePreserve: true,
		MaxCharacters:    140,
	}
	key, err := client.Keys.Create(want)
	if err != nil {
		t.Errorf("Keys.Create returned error: %v", err)
	}

	if !reflect.DeepEqual(key, want) {
		t.Errorf("Keys.Create returned %+v, want %+v", key, want)
	}
}

func TestKeysService_Create_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Keys.Create(&Key{})
		return err
	})
}

func TestKeysService_Create_invalidParam(t *testing.T) {
	_, err := client.Keys.Create(nil)
	if err == nil {
		t.Fatal("Cannot do create key with a nil request")
	}
}

func TestKeysService_Destroy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"success":true}`)
	})

	err := client.Keys.Destroy(1)
	if err != nil {
		t.Errorf("Keys.Destroy returned error: %v", err)
	}
}

func TestKeysService_Destroy_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		err := client.Keys.Destroy(0)
		return err
	})
}

func TestKeysService_DestroyMultiple(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/destroy_multiple", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, map[string]string{
			"ids[]": "1",
		})
		fmt.Fprint(w, `{"success":true}`)
	})

	err := client.Keys.DestroyMultiple([]int{1})
	if err != nil {
		t.Errorf("Keys.DestroyMultiple returned error: %v", err)
	}
}

func TestKeysService_DestroyMultiple_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		err := client.Keys.DestroyMultiple([]int{1})
		return err
	})
}

func TestKeysService_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		testFormValues(t, r, map[string]string{
			"translation_key[name]":                   "foo.key",
			"translation_key[name_plural]":            "foo.key.pluralized",
			"translation_key[pluralized]":             "1",
			"translation_key[description]":            "Some description",
			"translation_key[data_type]":              "string",
			"translation_key[tag_names]":              "foo,bar",
			"translation_key[unformatted]":            "1",
			"translation_key[xml_space_preserve]":     "1",
			"translation_key[max_characters_allowed]": "140",
		})
		fmt.Fprint(w, `{"id":123,"name":"foo.key","name_plural":"foo.key.pluralized","pluralized":true,"description":"Some description","data_type":"string","tag_list":["foo","bar"],"unformatted":true,"xml_space_preserve":true,"max_characters_allowed":140}`)
	})

	want := &Key{
		ID:               123,
		Name:             "foo.key",
		NamePlural:       "foo.key.pluralized",
		Pluralized:       true,
		Description:      "Some description",
		DataType:         "string",
		Tags:             []string{"foo", "bar"},
		Unformatted:      true,
		XMLSpacePreserve: true,
		MaxCharacters:    140,
	}
	key, err := client.Keys.Update(want)
	if err != nil {
		t.Errorf("Keys.Update returned error: %v", err)
	}

	if !reflect.DeepEqual(key, want) {
		t.Errorf("Keys.Update returned %+v, want %+v", key, want)
	}
}

func TestKeysService_Update_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Keys.Update(&Key{})
		return err
	})
}

func TestKeysService_Update_invalidParam(t *testing.T) {
	_, err := client.Keys.Update(nil)
	if err == nil {
		t.Fatal("Cannot do key update with a nil request")
	}
}

func TestKeysService_ListUntranslated(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/untranslated", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryString(t, r, map[string]string{"locale_name": "fr"})
		fmt.Fprint(w, `[{"id":1,"name":"helper.label","description":"Some description","pluralized":true,"data_type":"string"}]`)
	})

	keys, err := client.Keys.ListUntranslated("fr")
	if err != nil {
		t.Errorf("Keys.ListUntranslated returned error: %v", err)
	}

	want := []Key{
		{
			ID: 1, Name: "helper.label",
			Description: "Some description",
			Pluralized:  true,
			DataType:    "string",
		},
	}
	if !reflect.DeepEqual(keys, want) {
		t.Errorf("Keys.ListUntranslated returned %+v, want %+v", keys, want)
	}
}

func TestKeysService_ListUntranslated_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Keys.ListUntranslated("fr")
		return err
	})
}

func TestKeysService_Tag(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/tag", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"ids[]":  "123",
			"tags[]": "mobile",
		})
		fmt.Fprint(w, `{"success":true}`)
	})

	err := client.Keys.Tag([]int{123}, []string{"mobile"})
	if err != nil {
		t.Errorf("Keys.Tag returned error: %v", err)
	}
}

func TestKeysService_Tag_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		return client.Keys.Tag([]int{123}, []string{"mobile"})
	})
}

func TestKeysService_Translate_Map(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/translate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryString(t, r, map[string]string{"key": "mykey"})
		fmt.Fprint(w, `{"success":true,"translate":{"title":"A new title","body":"The help text body"}}`)
	})

	translation, err := client.Keys.Translate("mykey")
	if err != nil {
		t.Errorf("Keys.Translate returned error: %v", err)
	}
	want := map[string]string{
		"title": "A new title",
		"body":  "The help text body",
	}

	if !reflect.DeepEqual(translation.Map, want) {
		t.Errorf("Keys.Translate returned %+v, want %+v", translation.Map, want)
	}
}

func TestKeysService_Translate_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		_, err := client.Keys.Translate("mykey")
		return err
	})
}

func TestKeysService_Translate_String(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/translate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testQueryString(t, r, map[string]string{"key": "mykey"})
		fmt.Fprint(w, `{"success":true,"translate":"translation"}`)
	})

	translation, err := client.Keys.Translate("mykey")
	if err != nil {
		t.Errorf("Keys.Translate returned error: %v", err)
	}
	want := "translation"

	if !reflect.DeepEqual(translation.String, want) {
		t.Errorf("Keys.Translate returned %+v, want %+v", translation.String, want)
	}
}

func TestKeysService_Translate_jsonError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/translate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `1`)
	})

	_, err := client.Keys.Translate("mykey")
	if err == nil {
		t.Error("Keys.Translate will return marshal error if response cannot be parsed as JSON")
	}
}

func TestKeysService_Translate_failed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/translate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":false,"translate":false}`)
	})

	key, _ := client.Keys.Translate("mykey")
	if key != nil {
		t.Error("Keys.Translate will return nil if translate failed")
	}
}

func TestKeysService_Translate_parseError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/translate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"success":true,"translate":false}`)
	})

	_, err := client.Keys.Translate("mykey")
	if err == nil {
		t.Error("Keys.Translate will return error if translate field cannot be parsed")
	}
}

func TestKeysService_Upload(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/translation_keys/upload", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testFormValues(t, r, map[string]string{
			"filename":            "de.yml",
			"locale_name":         "de",
			"file_content":        "blah",
			"tags[]":              "tag",
			"file_format":         "yml",
			"update_translations": "1",
			"skip_unverification": "1",
			"skip_upload_tags":    "1",
		})
		fmt.Fprint(w, `{"success":true}`)
	})

	req := &UploadRequest{
		Filename:           "de.yml",
		Format:             "yml",
		Locale:             "de",
		FileContent:        "blah",
		Tags:               []string{"tag"},
		UpdateTranslations: true,
		SkipUnverification: true,
		SkipUploadTags:     true,
	}
	err := client.Keys.Upload(req)
	if err != nil {
		t.Errorf("Keys.Upload returned error: %v", err)
	}
}

func TestKeysService_Upload_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		return client.Keys.Upload(&UploadRequest{})
	})
}

func TestKeysService_Upload_invalidParam(t *testing.T) {
	err := client.Keys.Upload(nil)
	if err == nil {
		t.Fatal("Cannot do key upload with a nil request")
	}
}
