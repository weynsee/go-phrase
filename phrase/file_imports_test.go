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
	filename := "en.yml"

	mux.HandleFunc("/file_imports", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(100000)
		if err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}
		testMethod(t, r, "POST")
		testParams(t, r.MultipartForm.Value, map[string]string{
			"file_import[skip_upload_tags]":    "1",
			"file_import[skip_unverification]": "1",
			"file_import[convert_emoji]":       "1",
			"file_import[locale_code]":         "en",
			"file_import[filename]":            filename,
			"file_import[format]":              "yml",
			"file_import[tag_names]":           "new,upload",
			"file_import[update_translations]": "1",
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
		Format:             "yml",
		Tags:               []string{"new", "upload"},
		UpdateTranslations: true,
		SkipUnverification: true,
		SkipUploadTags:     true,
		ConvertEmoji:       true,
	}
	err := client.FileImports.Upload(upload, strings.NewReader("this is a test"))
	if err != nil {
		t.Errorf("FileImports.Upload returned error: %v", err)
	}
}
