package phrase

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFileImportsService_Upload(t *testing.T) {
	setup()
	defer teardown()
	filename := "en.csv"

	mux.HandleFunc("/file_imports", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(100000)
		if err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}
		testMethod(t, r, "POST")
		testParams(t, r.MultipartForm.Value, map[string]string{
			"file_import[skip_upload_tags]":                   "1",
			"file_import[skip_unverification]":                "1",
			"file_import[convert_emoji]":                      "1",
			"file_import[locale_code]":                        "en",
			"file_import[filename]":                           filename,
			"file_import[format]":                             "csv",
			"file_import[tag_names]":                          "new,upload",
			"file_import[update_translations]":                "1",
			"file_import[format_options][key_index]":          "1",
			"file_import[format_options][translation_index]":  "2",
			"file_import[format_options][comment_index]":      "3",
			"file_import[format_options][column_separator]":   ",",
			"file_import[format_options][quote_char]":         "\"",
			"file_import[format_options][header_content_row]": "false",
		})
		file := r.MultipartForm.File["file_import[file]"][0]
		if file.Filename != filename {
			t.Errorf("FileImports.Upload filename is %v, want %v", file.Filename, filename)
		}
		in, err := file.Open()
		if err != nil {
			t.Errorf("Failed to open multipart file: %v", err)
		}
		var out bytes.Buffer
		io.Copy(&out, in)
		if out.String() != "this is a test" {
			t.Errorf("Failed to read multipart file")
		}
		in.Close()
		fmt.Fprint(w, `{"success":true}`)
	})

	upload := &FileImportRequest{
		Locale:             "en",
		Filename:           filename,
		Format:             "csv",
		Tags:               []string{"new", "upload"},
		UpdateTranslations: true,
		SkipUnverification: true,
		SkipUploadTags:     true,
		ConvertEmoji:       true,
		FormatOptions: map[string]interface{}{
			"key_index":          1,
			"translation_index":  2,
			"comment_index":      3,
			"column_separator":   ",",
			"quote_char":         "\"",
			"header_content_row": false,
		},
	}
	err := client.FileImports.Upload(upload, strings.NewReader("this is a test"))
	if err != nil {
		t.Errorf("FileImports.Upload returned error: %v", err)
	}
}

func TestFileImportsService_serverError(t *testing.T) {
	testErrorHandling(t, func() error {
		upload := &FileImportRequest{
			Locale:   "en",
			Filename: "test.json",
			Format:   "json",
		}
		return client.FileImports.Upload(upload, strings.NewReader("this is a test"))
	})
}
