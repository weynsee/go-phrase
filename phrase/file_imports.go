package phrase

import (
	"github.com/google/go-querystring/query"
	"io"
)

// FileImportsService provides access to the file upload service
// in the Phrase API.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/file_imports/
type FileImportsService struct {
	client *Client
}

// FileImportRequest represents a file upload request
type FileImportRequest struct {
	Locale             string   `url:"file_import[locale_code]"`
	Filename           string   `url:"file_import[filename],omitempty"`
	Format             string   `url:"file_import[format],omitempty"`
	Tags               []string `url:"file_import[tag_names],comma,omitempty"`
	UpdateTranslations bool     `url:"file_import[update_translations],int,omitempty"`
	SkipUnverification bool     `url:"file_import[skip_unverification],int,omitempty"`
	SkipUploadTags     bool     `url:"file_import[skip_upload_tags],int,omitempty"`
	ConvertEmoji       bool     `url:"file_import[convert_emoji],int,omitempty"`
}

// Upload a localization file.
//
// Phrase API docs: http://docs.phraseapp.com/api/v1/file_imports/
func (s *FileImportsService) Upload(i *FileImportRequest, reader io.Reader) error {
	params, err := query.Values(i)
	if err != nil {
		return err
	}

	req, err := s.client.NewUploadRequest("file_imports", params, "file_import[file]", i.Filename, reader)

	if err != nil {
		return err
	}

	resp := new(successResponse)
	_, err = s.client.Do(req, resp)
	return err
}
