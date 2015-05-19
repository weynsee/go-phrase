package phrase

// ProjectsService provides access to the projects related functions
// in the PhraseApp API.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/projects/
type ProjectsService struct {
	client *Client
}

// Project represents the project associated with the current auth token.
type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Get details for the current project.
//
// PhraseApp API docs: http://docs.phraseapp.com/api/v1/projects/
func (s *ProjectsService) Current() (*Project, error) {
	req, err := s.client.NewRequest("GET", "projects/current", nil)
	if err != nil {
		return nil, err
	}

	project := new(Project)
	_, err = s.client.Do(req, project)
	if err != nil {
		return nil, err
	}

	return project, err
}
