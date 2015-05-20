package phrase

import "fmt"

// TagsService provides access to the tags related functions
// in the PhraseApp API.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/tags/
type TagsService struct {
	client *Client
}

// Tag represents a tag.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Progress represents the progress of a translation.
type Progress struct {
	TranslationsCount int `json:"translations_count"`
	TranslatedCount   int `json:"translated_count"`
	UnverifiedCount   int `json:"unverified_count"`
	UntranslatedCount int `json:"untranslated_count"`
}

// LocaleProgress represents the translation progress with a Locale struct and a Progress struct.
type LocaleProgress struct {
	Locale   Locale   `json:"locale"`
	Progress Progress `json:"progress"`
}

// TagProgress represents the translation progress of a Tag and a map of locales to Progress.
type TagProgress struct {
	Tag      Tag                       `json:"tag"`
	Progress map[string]LocaleProgress `json:"progress"`
}

// ListAll returns a list of all tags in the current project.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/tags/#index
func (s *TagsService) ListAll() ([]Tag, error) {
	req, err := s.client.NewRequest("GET", "tags", nil)
	if err != nil {
		return nil, err
	}

	tags := new([]Tag)
	_, err = s.client.Do(req, tags)
	if err != nil {
		return nil, err
	}

	return *tags, err
}

// GetProgress returns detailed information and progress for a tag.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/tags/#show
func (s *TagsService) GetProgress(id int) (*TagProgress, error) {
	url := fmt.Sprintf("tags/%d", id)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	progress := new(TagProgress)
	_, err = s.client.Do(req, progress)
	if err != nil {
		return nil, err
	}

	return progress, err
}
