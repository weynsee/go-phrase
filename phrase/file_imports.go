package phrase

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/url"
)

// FileImportsService provides access to the file upload service
// in the PhraseApp API.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/file_imports/
type FileImportsService struct {
	client *Client
}

// FileImportRequest represents a file upload request.
type FileImportRequest struct {
	Locale   string `url:"file_import[locale_code]"`
	Filename string `url:"file_import[filename],omitempty"`

	// http://docs.phraseapp.com/guides/formats/
	Format string   `url:"file_import[format],omitempty"`
	Tags   []string `url:"file_import[tag_names],comma,omitempty"`

	// Force the update of all translations with the file content.
	UpdateTranslations bool `url:"file_import[update_translations],int,omitempty"`

	// Do not initiate verification process for existing translations when updating translations.
	SkipUnverification bool `url:"file_import[skip_unverification],int,omitempty"`

	// Prevent creating an upload tag automatically.
	SkipUploadTags bool `url:"file_import[skip_upload_tags],int,omitempty"`

	// Enable Emoji conversion.
	ConvertEmoji bool `url:"file_import[convert_emoji],int,omitempty"`

	// Some formats can have additional options. Please check the format guide for more information about available options.
	FormatOptions formatOptions `url:"file_import[format_options],omitempty"`
}

type formatOptions map[string]interface{}

// Upload a localization file.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/file_imports/
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

func (o formatOptions) EncodeValues(key string, v *url.Values) error {
	for k, val := range o {
		(*v)[key+fmt.Sprintf("[%s]", k)] = []string{fmt.Sprintf("%v", val)}
	}
	return nil
}
