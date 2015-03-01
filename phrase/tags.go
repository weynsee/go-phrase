package phrase

import "fmt"

type TagsService struct {
	client *Client
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Progress struct {
	TranslationsCount int `json:"translations_count"`
	TranslatedCount   int `json:"translated_count"`
	UnverifiedCount   int `json:"unverified_count"`
	UntranslatedCount int `json:"untranslated_count"`
}

type LocaleProgress struct {
	Locale   Locale   `json:"locale"`
	Progress Progress `json:"progress"`
}

type TagProgress struct {
	Tag      Tag                       `json:"tag"`
	Progress map[string]LocaleProgress `json:"progress"`
}

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

func (t Tag) String() string {
	return fmt.Sprintf("Tag ID: %d Name: %s",
		t.ID, t.Name)
}

func (t TagProgress) String() string {
	return fmt.Sprintf("TagProgress Tag: %s Progress: %v",
		t.Tag, t.Progress)
}
